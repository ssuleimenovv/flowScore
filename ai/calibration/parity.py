"""Check that both Python Flows match the Go Flow Engine on the whole season.

Usage (from the repo root, after `go run ./cmd/export` in services/):
    python ai/calibration/parity.py

Compares every minute of every match with ai/data/flow.csv, which Go wrote
with the same params: flow.py, the readable port, and fast.py, the vectorized
one the calibration runs on. A calibration is only worth something if the params it
finds mean the same thing in Go, and this is the proof.
"""

import sys
import time

import numpy as np
import pandas as pd

from fast import Season
from flow import DATA, Params, load_events, trace_season

# Float64 in two languages: the same operations in the same order give the
# same bits, and the rest is a last-digit difference in minutes of decay
TOLERANCE = 1e-9


def check(name: str, minutes: pd.DataFrame, home: np.ndarray, away: np.ndarray, ref: pd.DataFrame, elapsed: float) -> None:
    keys = ["match_id", "period", "second"]
    if not minutes[keys].reset_index(drop=True).equals(ref[keys]):
        sys.exit(f"{name}: sampled minutes differ: {len(minutes)} in Python, {len(ref)} in Go")

    diff = max(np.abs(home - ref.home.to_numpy()).max(), np.abs(away - ref.away.to_numpy()).max())
    print(f"{name:<8} {len(ref)} minutes of {ref.match_id.nunique()} matches in {elapsed:.2f}s, max difference {diff:.1e}")
    if diff > TOLERANCE:
        sys.exit(f"{name}: Flow differs from Go by up to {diff}")


def main() -> None:
    ref = pd.read_csv(DATA / "flow.csv", dtype={"match_id": str})
    matches = load_events()

    start = time.perf_counter()
    ours = trace_season(matches, Params())
    check("flow.py", ours, ours.home.to_numpy(), ours.away.to_numpy(), ref, time.perf_counter() - start)

    season = Season(matches)
    start = time.perf_counter()
    home, away = season.trace(Params())
    check("fast.py", season.samples(), home, away, ref, time.perf_counter() - start)


if __name__ == "__main__":
    main()
