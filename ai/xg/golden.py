"""Cases for the Go xG model to match, the way golden.py does for the outcome model.

Usage (from the repo root, after fit.py):
    python ai/xg/golden.py

Takes shots of the season plus a few corner cases (a penalty, a shot from
the goal line, one from the halfway line), computes their xG in Python from
the numbers of model.json and writes inputs and xG to
services/internal/xg/testdata. The Go test computes the same shots with the
same model.json and must agree.
"""

import json
from pathlib import Path

import numpy as np
import pandas as pd

from fit import FEATURES, ROOT, SHOTS, features

MODEL = Path(__file__).resolve().parent / "model.json"
OUT = ROOT / "services/internal/xg/testdata/golden.json"
SAMPLE = 200


def predict(model: dict, shots: pd.DataFrame) -> np.ndarray:
    """xG from model.json alone, without sklearn: what Go computes."""
    x = features(shots)[FEATURES].to_numpy()
    w = np.array([model["weights"][f] for f in FEATURES])
    return 1 / (1 + np.exp(-(model["intercept"] + x @ w)))


def main() -> None:
    model = json.loads(MODEL.read_text(encoding="utf-8"))
    shots = pd.read_csv(SHOTS).sample(SAMPLE, random_state=7)
    corners = pd.DataFrame(
        [
            # x, y, header, situation
            (108.0, 40.0, 0, "penalty"),
            (119.5, 36.0, 0, "open_play"),  # beside the post, almost on the line
            (120.0, 30.0, 0, "open_play"),  # on the goal line itself
            (60.0, 40.0, 0, "open_play"),  # from the halfway line
            (114.0, 42.0, 1, "set_piece"),
            (95.0, 25.0, 0, "free_kick"),
        ],
        columns=["x", "y", "header", "situation"],
    )
    shots = pd.concat([shots[corners.columns], corners], ignore_index=True)
    xg = predict(model, shots)

    cases = [
        {"x": r.x, "y": r.y, "header": bool(r.header), "situation": r.situation, "xg": float(xg[i])}
        for i, r in enumerate(shots.itertuples())
    ]
    OUT.parent.mkdir(parents=True, exist_ok=True)
    OUT.write_text(json.dumps(cases, indent=1) + "\n")
    print(f"{len(cases)} cases -> {OUT}")


if __name__ == "__main__":
    main()
