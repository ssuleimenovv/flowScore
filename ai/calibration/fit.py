"""Calibrate and validate Flow on a season (docs/FLOW.md, section 8).

Usage (from the repo root, after `go run ./cmd/export` in services/):
    python ai/calibration/fit.py

The question Flow has to answer: will this team shoot in the next 10 minutes?
For every team and every minute we take Flow at that minute and whether the
team had a shot in the following 10 minutes, and link them by

    P(shot) = sigmoid(a + b · Flow / 100)

Validation is time-blocked cross-validation: the season is cut by date into 5
blocks, each block is predicted by a model fitted on the other 4. Every match
is predicted once, by a model that never saw it: 380 test matches instead of
the 95 of a single split, which is what a signal this small needs.

Three models are compared on the same predictions:
- constant: the training shot rate, whatever Flow says
- current:  Flow with the params in flow.py (the ones Go runs)
- fitted:   Flow with the event weights refitted on every training block

Memory (tau) and the break are never fitted: they define what "momentum"
means, and left free tau runs to its bound, because a long memory predicts
shots best by measuring which team is stronger. The slope b is fitted with
the current params and kept, or the fit shrinks every weight and steepens b:
the same predictions from a Flow that barely moves on screen.

Writes ai/calibration/report.json.
"""

import argparse
import json
import time
from concurrent.futures import ProcessPoolExecutor
from dataclasses import replace
from pathlib import Path

import numpy as np
import pandas as pd
from scipy.optimize import minimize
from sklearn.metrics import roc_auc_score

from fast import Season
from flow import DATA, Params, load_events

HORIZON = 10 * 60  # seconds ahead a shot counts
FOLDS = 5
BOOTSTRAP = 2000
OUT = Path(__file__).resolve().parent / "report.json"

FITTED = ["shot_base", "shot_per_xg", "possession"]
WEIGHTS = ["goal", "corner", "yellow_card", "red_card", "substitution", "key_pass", "tackle", "take_on", "foul"]
BOUNDS = [(0.0, 60.0)] * len(FITTED) + [(-30.0, 60.0)] * len(WEIGHTS)


def labels(df: pd.DataFrame, samples: pd.DataFrame) -> tuple[np.ndarray, np.ndarray, np.ndarray]:
    """For each sampled minute: did home and away shoot in the next 10 minutes,
    and does the period last long enough to tell. A shot is any event with xG,
    so an own goal, which has none, is not one."""
    shots = df[df.xg.notna()]
    times = {key: g.second.to_numpy() for key, g in shots.groupby(["match_id", "period", "side"])}
    last = df.groupby(["match_id", "period"]).second.max()
    second = samples.second.to_numpy()

    home = np.zeros(len(samples), dtype=bool)
    away = np.zeros(len(samples), dtype=bool)
    for (match_id, period), idx in samples.groupby(["match_id", "period"]).indices.items():
        t = second[idx]
        for side, out in (("home", home), ("away", away)):
            s = times.get((match_id, period, side), np.empty(0))
            # Shots in (t, t + HORIZON]: two binary searches instead of a loop
            out[idx] = np.searchsorted(s, t + HORIZON, side="right") > np.searchsorted(s, t, side="right")

    ends = last.loc[list(zip(samples.match_id, samples.period))].to_numpy()
    return home, away, second + HORIZON <= ends


class Block:
    """The team-minutes of some matches, with what happened after each."""

    def __init__(self, matches: dict, df: pd.DataFrame, ids: list[str]):
        self.season = Season({i: matches[i] for i in ids})
        samples = self.season.samples()
        home, away, self.valid = labels(df, samples)
        self.y = np.concatenate([home[self.valid], away[self.valid]]).astype(float)
        match = samples.match_id.to_numpy()[self.valid]
        self.match = np.concatenate([match, match])

    def flow(self, p: Params) -> np.ndarray:
        home, away = self.season.trace(p)
        return np.concatenate([home[self.valid], away[self.valid]])


def losses(y: np.ndarray, flow: np.ndarray, a: float, b: float) -> np.ndarray:
    """Log loss of every team-minute, written to stay finite for any z."""
    z = a + b * flow / 100
    return np.logaddexp(0, z) - y * z


def fit_link(y: np.ndarray, flow: np.ndarray) -> tuple[float, float]:
    """The best a and b for fixed Flow values: a two-parameter logistic regression."""
    res = minimize(lambda ab: losses(y, flow, *ab).mean(), [0.0, 1.0], method="L-BFGS-B")
    return float(res.x[0]), float(res.x[1])


def to_params(theta: np.ndarray, base: Params) -> Params:
    fitted = dict(zip(FITTED, theta[: len(FITTED)]))
    weights = dict(base.weights) | dict(zip(WEIGHTS, theta[len(FITTED) :]))
    return replace(base, **fitted, weights=weights)


def refit(blocks: list[Block], base: Params, b: float) -> tuple[Params, float]:
    """Event weights and the intercept a that minimize the log loss on the blocks."""
    y = np.concatenate([bl.y for bl in blocks])

    def loss(x: np.ndarray) -> float:
        p = to_params(x[:-1], base)
        return float(losses(y, np.concatenate([bl.flow(p) for bl in blocks]), x[-1], b).mean())

    x0 = [getattr(base, name) for name in FITTED] + [base.weights.get(w, 0.0) for w in WEIGHTS] + [0.0]
    res = minimize(loss, x0, method="L-BFGS-B", bounds=BOUNDS + [(None, None)], options={"maxiter": 500})
    return to_params(res.x[:-1], base), float(res.x[-1])


def gain_interval(better: np.ndarray, worse: np.ndarray, match: np.ndarray) -> tuple[float, float]:
    """95% interval of how much lower the log loss of `better` is. Whole matches
    are resampled: the minutes of one match are not independent, matches are."""
    per_match = pd.Series(worse - better).groupby(match).agg(["sum", "count"])
    sums, counts = per_match["sum"].to_numpy(), per_match["count"].to_numpy()
    picks = np.random.default_rng(7).integers(0, len(sums), size=(BOOTSTRAP, len(sums)))
    means = sums[picks].sum(axis=1) / counts[picks].sum(axis=1)
    return float(np.percentile(means, 2.5)), float(np.percentile(means, 97.5))


def validate(blocks: list[Block], k: int, current: Params) -> dict:
    """Predict block k with the three models fitted on the other blocks."""
    test = blocks[k]
    train = [bl for j, bl in enumerate(blocks) if j != k]
    train_y = np.concatenate([bl.y for bl in train])

    rate = train_y.mean()
    a, b = fit_link(train_y, np.concatenate([bl.flow(current) for bl in train]))
    fitted, a_fit = refit(train, current, b)
    return {
        "constant": losses(test.y, np.zeros(len(test.y)), np.log(rate / (1 - rate)), 0.0),
        "current": losses(test.y, test.flow(current), a, b),
        "fitted": losses(test.y, test.flow(fitted), a_fit, b),
        "flow": {"current": test.flow(current), "fitted": test.flow(fitted)},
        "weights": {name: getattr(fitted, name) for name in FITTED} | fitted.weights,
    }


def main() -> None:
    parser = argparse.ArgumentParser(description="Calibrate and validate Flow on a season.")
    # Each process holds numpy, pandas and scipy, a few hundred MB: more jobs
    # finish sooner on a machine with memory to spare
    parser.add_argument("--jobs", type=int, default=2, help="blocks validated at once")
    args = parser.parse_args()

    matches = load_events()
    df = pd.read_csv(DATA / "events.csv", dtype={"match_id": str})
    by_date = list(df.groupby("match_id", sort=False).date.first().sort_values(kind="stable").index)
    blocks = [Block(matches, df, [by_date[i] for i in part]) for part in np.array_split(range(len(by_date)), FOLDS)]
    y = np.concatenate([bl.y for bl in blocks])
    match = np.concatenate([bl.match for bl in blocks])
    print(f"{len(by_date)} matches, {len(y)} team-minutes, shot rate {y.mean():.3f}, {FOLDS} blocks by date")

    current = Params()
    started = time.perf_counter()
    if args.jobs == 1:
        # One block at a time in this process: the least memory
        folds = [validate(blocks, k, current) for k in range(FOLDS)]
    else:
        # The blocks are independent, so they can run in separate processes
        with ProcessPoolExecutor(max_workers=args.jobs) as pool:
            folds = list(pool.map(validate, [blocks] * FOLDS, range(FOLDS), [current] * FOLDS))
    print(f"  {FOLDS} blocks validated in {time.perf_counter() - started:.0f}s")

    ll = {model: np.concatenate([f[model] for f in folds]) for model in ("constant", "current", "fitted")}
    flows = {model: [f["flow"][model] for f in folds] for model in ("current", "fitted")}
    fold_weights = [f["weights"] for f in folds]
    report = {"matches": len(by_date), "teamMinutes": int(len(y)), "shotRate": round(float(y.mean()), 4), "models": {}}
    print("\nout of fold, every match predicted by a model that never saw it:")
    for model, values in ll.items():
        auc = roc_auc_score(y, np.concatenate(flows[model])) if model in flows else 0.5
        report["models"][model] = {"logLoss": round(float(values.mean()), 5), "auc": round(float(auc), 4)}
        print(f"  {model:<9} log loss {values.mean():.4f}   AUC {auc:.3f}")

    report["gain95"] = {}
    for better, worse in (("current", "constant"), ("fitted", "current")):
        low, high = gain_interval(ll[better], ll[worse], match)
        verdict = "better" if low > 0 else "worse" if high < 0 else "no clear difference"
        report["gain95"][f"{better} over {worse}"] = [round(low, 5), round(high, 5)]
        print(f"  {better} over {worse}: log loss gain 95% [{low:.4f}, {high:.4f}], {verdict}")

    # A weight worth trusting lands in the same place on every block
    table = pd.DataFrame(fold_weights).round(1)
    print("\nweights refitted on each block:")
    print(table.T.assign(current=[getattr(current, n, current.weights.get(n, 0.0)) for n in table.columns]).to_string())
    report["refittedWeights"] = table.to_dict(orient="list")

    OUT.write_text(json.dumps(report, indent=2) + "\n")
    print(f"\n-> {OUT}")


if __name__ == "__main__":
    main()
