"""Every shot of the season as one row, for the xG model (docs/FLOW.md, section 5).

Usage (from the repo root):
    python ai/xg/shots.py

A row says where the shot was taken from and how, whether it went in, and
what StatsBomb's own xG gave it, to compare with later. Coordinates stay in
StatsBomb's yards: the pitch is 120 × 80 and the goal the shot aims at is
on the line x = 120, between y = 36 and y = 44.

Writes ai/data/shots.csv.
"""

import json
from pathlib import Path

import pandas as pd

ROOT = Path(__file__).resolve().parents[2]
DATA = ROOT / "ai/data"
EVENTS = ROOT / "services/data/statsbomb"
MATCHES = EVENTS / "matches-2-27.json"


def situation(event: dict) -> str:
    """What led to the shot, in the words a live feed also has."""
    kind = event["shot"]["type"]["name"]
    if kind == "Penalty":
        return "penalty"
    if kind == "Free Kick":
        return "free_kick"  # a shot straight from a free kick
    if event["play_pattern"]["name"] in ("From Corner", "From Free Kick"):
        return "set_piece"  # play that started with a corner or a free kick
    return "open_play"


def shots_of(match_id: int, date: str) -> list[dict]:
    events = json.loads((EVENTS / f"events-{match_id}.json").read_text(encoding="utf-8"))
    rows = []
    for e in events:
        if e["type"]["name"] != "Shot":
            continue
        shot = e["shot"]
        rows.append(
            {
                "match_id": match_id,
                "date": date,
                "x": e["location"][0],
                "y": e["location"][1],
                "header": int(shot["body_part"]["name"] == "Head"),
                "situation": situation(e),
                "goal": int(shot["outcome"]["name"] == "Goal"),
                "statsbomb_xg": shot["statsbomb_xg"],
            }
        )
    return rows


def main() -> None:
    matches = json.loads(MATCHES.read_text(encoding="utf-8"))
    rows = []
    for m in matches:
        rows += shots_of(m["match_id"], m["match_date"])
    shots = pd.DataFrame(rows).sort_values(["date", "match_id"], kind="stable")

    DATA.mkdir(exist_ok=True)
    shots.to_csv(DATA / "shots.csv", index=False)
    print(f"{len(shots)} shots in {shots.match_id.nunique()} matches, {shots.goal.sum()} goals")
    print(shots.situation.value_counts().to_string())


if __name__ == "__main__":
    main()
