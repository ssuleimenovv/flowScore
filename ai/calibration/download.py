"""Download a season of StatsBomb Open Data for calibrating Flow.

Usage (from the repo root):
    python ai/calibration/download.py

Saves the season's matches file plus events and lineups of every match into
services/data/statsbomb, named the way the Go replay expects. Files that are
already there are skipped, so an interrupted run can simply be started again.
"""

import argparse
import json
import sys
import urllib.request
from concurrent.futures import ThreadPoolExecutor, as_completed
from pathlib import Path

BASE = "https://raw.githubusercontent.com/statsbomb/open-data/master/data"
ROOT = Path(__file__).resolve().parents[2]


def fetch(url: str, dest: Path, attempts: int = 3) -> None:
    if dest.exists():
        return
    for attempt in range(1, attempts + 1):
        try:
            with urllib.request.urlopen(url, timeout=60) as res:
                data = res.read()
            # Write to a temporary name first: a half-written file must not
            # look finished to the next run
            tmp = dest.with_suffix(".part")
            tmp.write_bytes(data)
            tmp.replace(dest)
            return
        except OSError as err:
            if attempt == attempts:
                raise RuntimeError(f"{url}: {err}") from err


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument("--competition", type=int, default=2, help="2 = Premier League")
    parser.add_argument("--season", type=int, default=27, help="27 = 2015/16")
    parser.add_argument("--out", type=Path, default=ROOT / "services/data/statsbomb")
    parser.add_argument("--workers", type=int, default=8)
    args = parser.parse_args()

    args.out.mkdir(parents=True, exist_ok=True)
    matches_path = args.out / f"matches-{args.competition}-{args.season}.json"
    fetch(f"{BASE}/matches/{args.competition}/{args.season}.json", matches_path)
    ids = [m["match_id"] for m in json.loads(matches_path.read_text(encoding="utf-8"))]

    jobs = []
    for match_id in ids:
        jobs.append((f"{BASE}/events/{match_id}.json", args.out / f"events-{match_id}.json"))
        jobs.append((f"{BASE}/lineups/{match_id}.json", args.out / f"lineups-{match_id}.json"))

    failed = []
    with ThreadPoolExecutor(max_workers=args.workers) as pool:
        futures = {pool.submit(fetch, url, dest): dest for url, dest in jobs}
        for done, future in enumerate(as_completed(futures), start=1):
            if future.exception():
                failed.append(str(future.exception()))
            print(f"\r{done}/{len(jobs)} files", end="", flush=True)
    print()

    if failed:
        print(f"{len(failed)} files failed; run again to retry:", *failed[:5], sep="\n  ")
        sys.exit(1)
    print(f"{len(ids)} matches in {args.out}")


if __name__ == "__main__":
    main()
