"""Check that the Python Flow matches the Go Flow Engine on the whole season.

Usage (from the repo root, after `go run ./cmd/export` in services/):
    python ai/calibration/parity.py

Compares every minute of every match with ai/data/flow.csv, which Go wrote
with the same params. A calibration is only worth something if the params it
finds mean the same thing in Go, and this is the proof.
"""

import sys
import time

import pandas as pd

from flow import DATA, Params, load_events, trace_season

# Float64 in two languages: the same operations in the same order give the
# same bits, and the rest is a last-digit difference in minutes of decay
TOLERANCE = 1e-9


def main() -> None:
    start = time.perf_counter()
    ours = trace_season(load_events(), Params())
    elapsed = time.perf_counter() - start
    ref = pd.read_csv(DATA / "flow.csv", dtype={"match_id": str})

    keys = ["match_id", "period", "second"]
    if not ours[keys].equals(ref[keys]):
        sys.exit(f"sampled minutes differ: {len(ours)} in Python, {len(ref)} in Go")

    diff = (ours[["home", "away"]] - ref[["home", "away"]]).abs().to_numpy().max()
    print(f"{len(ref)} minutes of {ref.match_id.nunique()} matches in {elapsed:.1f}s, max difference {diff:.1e}")
    if diff > TOLERANCE:
        sys.exit(f"Python Flow differs from Go by up to {diff}")


if __name__ == "__main__":
    main()
