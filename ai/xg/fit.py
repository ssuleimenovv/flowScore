"""Fit and validate our own xG model (docs/FLOW.md, section 5).

Usage (from the repo root, after shots.py):
    python ai/xg/fit.py

xG is the chance that a shot goes in. The model is a logistic regression:

    xG = 1 / (1 + e^-(a + b · features))

on what any live feed knows about a shot:

- distance:  yards from the shot to the middle of the goal
- angle:     how wide the goal looks from there, in radians
- header:    1 for a header
- set_piece, free_kick, penalty: what led to the shot (open play is 0 · 0 · 0)

Validation is the same as for the other models: 5 blocks of the season by
date, each predicted by a model fitted on the other 4. It is scored against
two references: a constant (every shot gets the season's goal rate) and
StatsBomb's own xG, which also sees the keeper and the defenders. The
interval of the gap to StatsBomb comes from resampling whole matches.

Writes ai/xg/model.json (what Go loads) and ai/xg/report.json.
"""

import json
from pathlib import Path

import numpy as np
import pandas as pd
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import brier_score_loss, log_loss, roc_auc_score

ROOT = Path(__file__).resolve().parents[2]
SHOTS = ROOT / "ai/data/shots.csv"
OUT = Path(__file__).resolve().parent

GOAL_X, GOAL_Y, HALF_GOAL = 120.0, 40.0, 4.0  # the goal is 8 yards wide
FEATURES = ["distance", "angle", "header", "set_piece", "free_kick", "penalty"]
BLOCKS = 5
RESAMPLES = 1000


def features(shots: pd.DataFrame) -> pd.DataFrame:
    dx = GOAL_X - shots.x
    dy = shots.y - GOAL_Y
    out = pd.DataFrame(index=shots.index)
    out["distance"] = np.hypot(dx, dy)
    # The angle between the lines to the two posts; atan2 keeps it right
    # even beside the goal, where the plain formula's denominator is negative
    out["angle"] = np.arctan2(2 * HALF_GOAL * dx, dx**2 + dy**2 - HALF_GOAL**2)
    out["header"] = shots.header
    for s in ["set_piece", "free_kick", "penalty"]:
        out[s] = (shots.situation == s).astype(int)
    return out


def fit(x: pd.DataFrame, y: pd.Series) -> LogisticRegression:
    return LogisticRegression(penalty=None, max_iter=1000).fit(x, y)


def scores(y: np.ndarray, p: np.ndarray) -> dict:
    p = np.clip(p, 1e-6, 1 - 1e-6)
    return {
        "log_loss": round(log_loss(y, p), 4),
        "brier": round(brier_score_loss(y, p), 4),
        "auc": round(roc_auc_score(y, p), 4),
    }


def cross_validate(shots: pd.DataFrame, x: pd.DataFrame) -> pd.Series:
    """Each shot predicted by a model that never saw its block of the season."""
    days = np.sort(shots.date.unique())
    block_of_day = {d: i * BLOCKS // len(days) for i, d in enumerate(days)}
    blocks = shots.date.map(block_of_day)
    out = pd.Series(np.nan, index=shots.index)
    for b in range(BLOCKS):
        test = blocks == b
        model = fit(x[~test], shots.goal[~test])
        out[test] = model.predict_proba(x[test])[:, 1]
    return out


def gap_interval(shots: pd.DataFrame, ours: pd.Series) -> list[float]:
    """95% interval of our log loss minus StatsBomb's, resampling matches."""
    p = np.clip(ours, 1e-6, 1 - 1e-6)
    q = np.clip(shots.statsbomb_xg, 1e-6, 1 - 1e-6)
    y = shots.goal
    loss = -(y * np.log(p) + (1 - y) * np.log(1 - p)) + (y * np.log(q) + (1 - y) * np.log(1 - q))
    per_match = pd.DataFrame({"gap": loss, "match": shots.match_id}).groupby("match").gap.agg(["sum", "count"])
    rng = np.random.default_rng(7)
    gaps = []
    for _ in range(RESAMPLES):
        pick = per_match.iloc[rng.integers(0, len(per_match), len(per_match))]
        gaps.append(pick["sum"].sum() / pick["count"].sum())
    return [round(float(v), 4) for v in np.percentile(gaps, [2.5, 97.5])]


def calibration(y: pd.Series, p: pd.Series, bins: int = 10) -> list[dict]:
    """Shots grouped by predicted xG: the mean prediction next to the goal rate."""
    groups = pd.qcut(p, bins, duplicates="drop")
    table = pd.DataFrame({"y": y, "p": p}).groupby(groups, observed=True)
    return [
        {"predicted": round(g.p.mean(), 3), "scored": round(g.y.mean(), 3), "shots": len(g)}
        for _, g in table
    ]


def main() -> None:
    shots = pd.read_csv(SHOTS)
    x = features(shots)
    y = shots.goal

    ours = cross_validate(shots, x)
    constant = np.full(len(y), y.mean())
    report = {
        "shots": len(shots),
        "goals": int(y.sum()),
        "constant": scores(y, constant),
        "ours": scores(y, ours),
        "statsbomb": scores(y, shots.statsbomb_xg),
        "ours_minus_statsbomb_log_loss_95": gap_interval(shots, ours),
        "calibration": calibration(y, ours),
    }

    model = fit(x, y)
    weights = dict(zip(FEATURES, (round(float(w), 6) for w in model.coef_[0])))
    out = {"name": "xG v1", "intercept": round(float(model.intercept_[0]), 6), "weights": weights}
    report["model"] = out

    (OUT / "model.json").write_text(json.dumps(out, indent=2) + "\n")
    (OUT / "report.json").write_text(json.dumps(report, indent=2, ensure_ascii=False) + "\n")
    print(json.dumps({k: v for k, v in report.items() if k != "calibration"}, indent=2))
    print("calibration (predicted -> scored):")
    for row in report["calibration"]:
        print(f"  {row['predicted']:.3f} -> {row['scored']:.3f}  ({row['shots']} shots)")


if __name__ == "__main__":
    main()
