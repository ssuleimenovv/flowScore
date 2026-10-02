"""Flow Momentum in Python: a line-by-line port of services/internal/flow.

Calibration fits the params here and hands them back to Go, so this port must
compute exactly what the Flow Engine does. parity.py checks it on every minute
of the season against the numbers Go wrote to ai/data/flow.csv.
"""

import math
from dataclasses import dataclass, field
from pathlib import Path

import pandas as pd

DATA = Path(__file__).resolve().parents[1] / "data"

SHOTS = {"shot_on_target", "shot_off_target", "shot_blocked"}
OPPONENT = {"home": "away", "away": "home"}

# Match second each period kicks off at, as periodStart in trace.go
PERIOD_START = {1: 0, 2: 45 * 60, 3: 90 * 60, 4: 105 * 60, 5: 120 * 60}


@dataclass(frozen=True)
class Params:
    """DefaultParams() from params.go; docs/FLOW.md, section 4."""

    tau: float = 5.0  # impulse memory, match minutes
    k: float = 18.7  # scale of the 0–100 curve
    base: float = 3.0  # S₀, the Flow floor
    halftime_keep: float = 0.7
    shot_base: float = 6.0
    shot_per_xg: float = 7.3
    default_xg: float = 0.1
    possession: float = 9.0
    weights: dict[str, float] = field(
        default_factory=lambda: {
            "goal": 23.0,
            "yellow_card": 5.0,
            "corner": 4.0,
            "substitution": 5.0,
            "key_pass": 10.0,
            "red_card": -18.0,
        }
    )


    def weight(self, type_: str, side: str, xg: float) -> tuple[str, float]:
        """Which team the event benefits and by how much, as Params.Weight."""
        if type_ in SHOTS:
            if math.isnan(xg):
                xg = self.default_xg
            return side, self.shot_base + self.shot_per_xg * xg
        if type_ == "yellow_card":
            return OPPONENT[side], self.weights[type_]
        return side, self.weights.get(type_, 0.0)


class State:
    """The impulse of both teams, as State in state.go. Time is in match seconds."""

    def __init__(self, p: Params):
        self.p = p
        self.impulse = {"home": 0.0, "away": 0.0}
        self.at = 0
        self.period = 0

    def start_period(self, period: int, at: int) -> None:
        if self.period != 0 and period > self.period:
            for side in self.impulse:
                self.impulse[side] *= self.p.halftime_keep
            self.at = at
        self.period = period

    def apply(self, period: int, second: int, type_: str, side: str, xg: float, home_share: float) -> None:
        self.start_period(period, second)
        self.advance(second)

        if type_ == "possession":
            if not math.isnan(home_share):
                lead = home_share - 0.5
                self._add("home", self.p.possession * lead)
                self._add("away", -self.p.possession * lead)
            return

        side, w = self.p.weight(type_, side, xg)
        self._add(side, w)

    def _add(self, side: str, w: float) -> None:
        self.impulse[side] = max(0.0, self.impulse[side] + w)

    def advance(self, t: int) -> None:
        if t <= self.at:
            return
        decay = math.exp(-((t - self.at) / 60) / self.p.tau)
        for side in self.impulse:
            self.impulse[side] *= decay
        self.at = t

    def flow(self) -> tuple[float, float]:
        return self._flow("home"), self._flow("away")

    def _flow(self, side: str) -> float:
        return 100 * (1 - math.exp(-(self.p.base + self.impulse[side]) / self.p.k))


def trace(events: list[tuple], p: Params) -> list[tuple[int, int, float, float]]:
    """Flow at every whole minute of one match, by the rule of Trace in trace.go.

    events are (period, second, type, side, xg, home_share) in match order;
    the result is (period, second, home, away).
    """
    s = State(p)
    out = []
    i = 0
    while i < len(events):
        period = events[i][0]
        end = i
        while end < len(events) and events[end][0] == period:
            end += 1
        last = events[end - 1][1]

        s.start_period(period, PERIOD_START[period])
        t = PERIOD_START[period] + 60
        while t <= last:
            while i < end and events[i][1] <= t:
                s.apply(*events[i])
                i += 1
            s.advance(t)
            out.append((period, t, *s.flow()))
            t += 60
        while i < end:  # events after the last whole minute
            s.apply(*events[i])
            i += 1
    return out


def load_events(path: Path = DATA / "events.csv") -> dict[str, list[tuple]]:
    """events.csv from `go run ./cmd/export`, as match id → events in match order."""
    # An empty side (possession) stays "", an empty number becomes NaN
    df = pd.read_csv(
        path,
        dtype={"match_id": str, "side": str},
        keep_default_na=False,
        na_values={"xg": [""], "home_share": [""]},
    )
    columns = ["period", "second", "type", "side", "xg", "home_share"]
    return {
        match_id: list(group[columns].itertuples(index=False, name=None))
        for match_id, group in df.groupby("match_id", sort=False)
    }


def trace_season(matches: dict[str, list[tuple]], p: Params) -> pd.DataFrame:
    rows = [(match_id, *sample) for match_id, events in matches.items() for sample in trace(events, p)]
    return pd.DataFrame(rows, columns=["match_id", "period", "second", "home", "away"])
