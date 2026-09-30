"""Flow for the whole season at once: flow.trace, vectorized across matches.

flow.py walks one match at a time, event by event; the optimizer needs the
season hundreds of times, and a Python loop over 130 000 steps is too slow for
that. Here the matches move in lockstep instead: step k is the k-th event or
minute mark of every match at the same time, done by numpy on a column of 380.
Within a match the operations and their order stay those of flow.py, so the
numbers stay those of Go (parity.py checks both).
"""

import numpy as np
import pandas as pd

from flow import PERIOD_START, SHOTS, Params

# What a step of a match does
PAD, EVENT, SAMPLE, KICKOFF = 0, 1, 2, 3


class Season:
    """The matches laid out as steps, once; Flow for any params is then cheap."""

    def __init__(self, matches: dict[str, list[tuple]]):
        self.ids = list(matches)
        rows = [self._steps(events) for events in matches.values()]
        width = max(len(r) for r in rows)

        def column(i, fill):
            out = np.full((len(rows), width), fill, dtype=float)
            for m, r in enumerate(rows):
                out[m, : len(r)] = [step[i] for step in r]
            return out

        self.kind = column(0, PAD).astype(int)
        self.period = column(1, 0).astype(int)
        self.second = column(2, 0)
        self.xg = column(5, np.nan)
        self.home_share = column(6, np.nan)

        types = [[step[3] for step in r] + [""] * (width - len(r)) for r in rows]
        sides = [[step[4] for step in r] + [""] * (width - len(r)) for r in rows]
        self.type = np.array(types, dtype=object)
        self.side = np.array(sides, dtype=object)
        # Each type as a number once, so a weight table indexes them per params
        self.names, self.codes = np.unique(self.type, return_inverse=True)
        self.shot = np.isin(self.type, list(SHOTS))
        self.possession = (self.type == "possession") & ~np.isnan(self.home_share)
        # A card helps the other team (Params.Weight), so its weight goes there
        card = self.type == "yellow_card"
        self.home = ((self.side == "home") & ~card) | ((self.side == "away") & card)
        self.away = ((self.side == "away") & ~card) | ((self.side == "home") & card)

    @staticmethod
    def _steps(events: list[tuple]) -> list[tuple]:
        """The steps of flow.trace for one match, as (kind, period, second, type, side, xg, home_share)."""
        steps = []
        i = 0
        while i < len(events):
            period = events[i][0]
            end = i
            while end < len(events) and events[end][0] == period:
                end += 1
            last = events[end - 1][1]

            steps.append((KICKOFF, period, PERIOD_START[period], "", "", np.nan, np.nan))
            t = PERIOD_START[period] + 60
            while t <= last:
                while i < end and events[i][1] <= t:
                    steps.append((EVENT, *events[i]))
                    i += 1
                steps.append((SAMPLE, period, t, "", "", np.nan, np.nan))
                t += 60
            while i < end:
                steps.append((EVENT, *events[i]))
                i += 1
        return steps

    def samples(self) -> pd.DataFrame:
        """Which minute each Flow value of trace is, as flow.trace_season lists them."""
        rows, cols = np.nonzero(self.kind == SAMPLE)
        return pd.DataFrame(
            {
                "match_id": np.array(self.ids)[rows],
                "period": self.period[rows, cols],
                "second": self.second[rows, cols].astype(int),
            }
        )

    def weights(self, p: Params) -> tuple[np.ndarray, np.ndarray]:
        """What each event adds to the home and the away impulse under p."""
        w = np.array([p.weights.get(t, 0.0) for t in self.names])[self.codes]
        xg = np.where(np.isnan(self.xg), p.default_xg, self.xg)
        w = np.where(self.shot, p.shot_base + p.shot_per_xg * xg, w)

        lead = p.possession * (self.home_share - 0.5)
        home = np.where(self.possession, lead, np.where(self.home, w, 0.0))
        away = np.where(self.possession, -lead, np.where(self.away, w, 0.0))
        return home, away

    def trace(self, p: Params) -> tuple[np.ndarray, np.ndarray]:
        """Home and away Flow at every sample, in the order of flow.trace_season."""
        add_home, add_away = self.weights(p)
        n, width = self.kind.shape
        home, away = np.zeros(n), np.zeros(n)
        at, period = np.zeros(n), np.zeros(n, dtype=int)
        out_home, out_away = np.zeros((n, width)), np.zeros((n, width))

        for k in range(width):
            kind, second = self.kind[:, k], self.second[:, k]

            # State.start_period: the break keeps part of the impulse
            brk = (kind == KICKOFF) & (period != 0) & (self.period[:, k] > period)
            home[brk] *= p.halftime_keep
            away[brk] *= p.halftime_keep
            at[brk] = second[brk]
            period = np.where(kind == KICKOFF, self.period[:, k], period)

            # State.advance: decay up to the event or the minute mark
            moving = ((kind == EVENT) | (kind == SAMPLE)) & (second > at)
            decay = np.exp(-((second[moving] - at[moving]) / 60) / p.tau)
            home[moving] *= decay
            away[moving] *= decay
            at[moving] = second[moving]

            # State._add: adding 0 to a non-negative impulse changes nothing, so
            # every match can take its step's weights, most of them zero
            home = np.maximum(0.0, home + add_home[:, k])
            away = np.maximum(0.0, away + add_away[:, k])

            out_home[:, k] = home
            out_away[:, k] = away

        sample = self.kind == SAMPLE
        flow = lambda s: 100 * (1 - np.exp(-(p.base + s[sample]) / p.k))
        return flow(out_home), flow(out_away)
