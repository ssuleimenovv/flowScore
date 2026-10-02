"""The season as minute-by-minute snapshots for the outcome model.

Usage (from the repo root, after `go run ./cmd/export` in services/):
    python ai/prediction/dataset.py

Each row is one minute of one match: what a viewer knows at that moment
(score, red cards, xG, Flow, time left) and how the match really ended. The
model learns from that how many more goals each team will score.

Team strength before the match comes from the matches it played earlier in
the season, never from later ones: the rating is the team's average xG
difference, shrunk towards zero while there are only a few games behind it.

Writes ai/data/snapshots.csv.
"""

import json
from pathlib import Path

import numpy as np
import pandas as pd

ROOT = Path(__file__).resolve().parents[2]
DATA = ROOT / "ai/data"
MATCHES = ROOT / "services/data/statsbomb/matches-2-27.json"

# A team's rating counts this many average games before its own: after 5
# matches it is half the team's record and half "an average team"
PRIOR_GAMES = 5
# The time a match is expected to last, stoppage included. A viewer does not
# know the real stoppage time, so neither does the model
HALF_END = {1: 47, 2: 93}
SECOND_HALF = 48  # minutes the second half is expected to last


def ratings(matches: pd.DataFrame, xg: pd.DataFrame) -> pd.Series:
    """Home rating minus away rating before each match, from earlier matches only."""
    history: dict[str, list[float]] = {}

    def rating(team: str) -> float:
        games = history.get(team, [])
        return sum(games) / (len(games) + PRIOR_GAMES)

    out = {}
    # Matches on the same day do not see each other
    for _, day in matches.sort_values("date").groupby("date"):
        for m in day.itertuples():
            out[m.match_id] = rating(m.home) - rating(m.away)
        for m in day.itertuples():
            diff = xg.at[m.match_id, "home"] - xg.at[m.match_id, "away"]
            history.setdefault(m.home, []).append(diff)
            history.setdefault(m.away, []).append(-diff)
    return pd.Series(out, name="rating")


def snapshots(events: pd.DataFrame, flow: pd.DataFrame) -> pd.DataFrame:
    """For every minute of flow.csv: the score, cards and xG up to that minute."""
    events = events.sort_values(["match_id", "period", "second"], kind="stable")
    parts = []
    for match_id, rows in flow.groupby("match_id", sort=False):
        e = events[events.match_id == match_id]
        # An event at the very second of a sample has already happened, as in flow.trace
        at = np.searchsorted(e.period * 10_000 + e.second, rows.period * 10_000 + rows.second, side="right")
        out = rows.copy()
        for side in ("home", "away"):
            mine = e.side == side
            for name, hit in (("goals", e.type == "goal"), ("reds", e.type == "red_card")):
                count = np.concatenate([[0], np.cumsum(mine & hit)])
                out[f"{side}_{name}"] = count[at]
            xg = np.concatenate([[0.0], np.cumsum(np.where(mine, e.xg.fillna(0), 0))])
            out[f"{side}_xg"] = xg[at]
        parts.append(out)
    return pd.concat(parts, ignore_index=True)


def main() -> None:
    events = pd.read_csv(DATA / "events.csv", dtype={"match_id": str})
    flow = pd.read_csv(DATA / "flow.csv", dtype={"match_id": str})
    raw = json.loads(MATCHES.read_text(encoding="utf-8"))
    matches = pd.DataFrame(
        {
            "match_id": str(m["match_id"]),
            "date": m["match_date"],
            "home": m["home_team"]["home_team_name"],
            "away": m["away_team"]["away_team_name"],
            "final_home": m["home_score"],
            "final_away": m["away_score"],
        }
        for m in raw
    )

    xg = events.groupby(["match_id", "side"]).xg.sum().unstack(fill_value=0.0)
    matches = matches.join(ratings(matches, xg), on="match_id")

    df = snapshots(events, flow).rename(columns={"home": "home_flow", "away": "away_flow"})
    df["minute"] = df.second / 60
    left = (HALF_END[1] - df.minute).clip(lower=0) + SECOND_HALF
    df["remaining"] = left.where(df.period == 1, (HALF_END[2] - df.minute).clip(lower=0.5))
    df = df.merge(matches, on="match_id")

    # The goals in the events must be the final score, or the labels lie
    last = df.groupby("match_id").tail(1)
    wrong = last[(last.home_goals != last.final_home) | (last.away_goals != last.final_away)]
    if len(wrong):
        raise SystemExit(f"{len(wrong)} matches where events disagree with the final score: {wrong.match_id.tolist()[:5]}")

    columns = [
        "match_id", "date", "home", "away", "period", "minute", "remaining",
        "home_goals", "away_goals", "home_reds", "away_reds", "home_xg", "away_xg",
        "home_flow", "away_flow", "rating", "final_home", "final_away",
    ]
    df[columns].to_csv(DATA / "snapshots.csv", index=False)
    print(f"{len(df)} snapshots of {df.match_id.nunique()} matches -> {DATA / 'snapshots.csv'}")


if __name__ == "__main__":
    main()
