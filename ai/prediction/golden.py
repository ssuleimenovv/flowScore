"""Cases for the Go outcome model to match, the way parity.py checks Flow.

Usage (from the repo root, after fit.py):
    python ai/prediction/golden.py

Takes minutes of the season plus a few corner cases (red cards, a big lead,
the final whistle), predicts them in Python with the numbers of model.json
and writes inputs and chances to services/internal/predict/testdata. The Go
test predicts the same inputs with the same model.json and must agree.
"""

import json
from pathlib import Path

import numpy as np
import pandas as pd

from dataset import DATA, ROOT
from fit import outcome, with_draw

MODEL = Path(__file__).resolve().parent / "model.json"
OUT = ROOT / "services/internal/predict/testdata/golden.json"
SAMPLE = 300


def predict(model: dict, df: pd.DataFrame) -> np.ndarray:
    """The chances from model.json alone, without sklearn: what Go computes."""
    w = model["weights"]
    cap = model["leadCap"]
    lead = (df.home_goals - df.away_goals).to_numpy()

    def rate(home: float, rating, red, lead):
        return np.exp(
            model["intercept"] + w["home"] * home + w["rating"] * rating + w["red"] * red + w["lead"] * np.clip(lead, -cap, cap)
        )

    left = df.remaining.to_numpy()
    home = rate(1.0, df.rating, df.away_reds - df.home_reds, lead) * left
    away = rate(0.0, -df.rating, df.home_reds - df.away_reds, -lead) * left
    return with_draw(outcome(home.to_numpy(), away.to_numpy(), lead), model["draw"])


def main() -> None:
    model = json.loads(MODEL.read_text(encoding="utf-8"))
    df = pd.read_csv(DATA / "snapshots.csv", dtype={"match_id": str})
    df = df.sample(SAMPLE, random_state=7)
    df["finished"] = False

    corners = pd.DataFrame(
        [
            # period, minute, remaining, goals, reds, rating, finished
            (1, 0.0, 95.0, 0, 0, 0, 0, 0.0, False),
            (1, 46.5, 48.5, 1, 0, 0, 1, 0.4, False),
            (2, 60.0, 33.0, 0, 3, 2, 0, -0.9, False),
            (2, 95.0, 0.5, 2, 2, 0, 0, 1.2, False),
            (2, 94.0, 0.0, 1, 0, 0, 0, -0.3, True),
            (2, 94.0, 0.0, 2, 2, 0, 1, 0.0, True),
        ],
        columns=["period", "minute", "remaining", "home_goals", "away_goals", "home_reds", "away_reds", "rating", "finished"],
    )
    df = pd.concat([df, corners], ignore_index=True)
    p = predict(model, df)

    cases = [
        {
            "period": int(r.period),
            "seconds": round(r.minute * 60),
            "homeGoals": int(r.home_goals),
            "awayGoals": int(r.away_goals),
            "homeReds": int(r.home_reds),
            "awayReds": int(r.away_reds),
            "rating": float(r.rating),
            "finished": bool(r.finished),
            "home": float(p[i, 0]),
            "draw": float(p[i, 1]),
            "away": float(p[i, 2]),
        }
        for i, r in enumerate(df.itertuples())
    ]
    OUT.parent.mkdir(parents=True, exist_ok=True)
    OUT.write_text(json.dumps(cases, indent=1) + "\n")
    print(f"{len(cases)} cases -> {OUT}")


if __name__ == "__main__":
    main()
