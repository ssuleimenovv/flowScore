"""The Go tests of state.go and trace.go, ported: the same inputs, the same numbers.

Run from the repo root:
    python -m unittest discover -s ai/calibration
"""

import math
import unittest

from flow import Params, State, trace

NAN = float("nan")
FLOOR = 14.82  # 100 · (1 − e^(−3/18.7))


def event(type_, side, period, minute, xg=NAN, home_share=NAN):
    return (period, minute * 60, type_, side, xg, home_share)


class StateTest(unittest.TestCase):
    def flow_after(self, *events):
        s = State(Params())
        for e in events:
            s.apply(*e)
        return s

    def test_idle_teams_stay_at_floor(self):
        home, away = State(Params()).flow()
        self.assertAlmostEqual(home, FLOOR, places=2)
        self.assertAlmostEqual(away, FLOOR, places=2)

    def test_goal(self):
        home, away = self.flow_after(event("goal", "home", 1, 10)).flow()
        self.assertAlmostEqual(home, 75.10, places=2)
        self.assertAlmostEqual(away, FLOOR, places=2)

    def test_decay_after_tau(self):
        s = self.flow_after(event("goal", "home", 1, 10))
        s.advance(15 * 60)
        self.assertAlmostEqual(s.flow()[0], 45.82, places=2)

    def test_yellow_card_helps_opponent(self):
        home, away = self.flow_after(event("yellow_card", "away", 1, 20)).flow()
        self.assertAlmostEqual(home, 34.81, places=2)
        self.assertAlmostEqual(away, FLOOR, places=2)

    def test_halftime_keeps_part_of_impulse(self):
        s = self.flow_after(event("goal", "home", 1, 45), event("foul", "away", 2, 45))
        self.assertAlmostEqual(s.flow()[0], 63.99, places=2)

    def test_possession_lead_helps_holder(self):
        home, away = self.flow_after(event("possession", "", 1, 1, home_share=0.7)).flow()
        self.assertAlmostEqual(home, 22.64, places=2)
        self.assertAlmostEqual(away, FLOOR, places=2)

    def test_no_debt_after_long_spell_without_ball(self):
        spell = [event("possession", "", 1, m, home_share=0.7) for m in range(1, 21)]
        shot = event("shot_on_target", "away", 1, 20)
        self.assertAlmostEqual(self.flow_after(*spell, shot).flow()[1], self.flow_after(shot).flow()[1])


class TraceTest(unittest.TestCase):
    def test_samples_every_minute(self):
        p = Params()
        samples = trace([event("goal", "home", 1, 44), event("foul", "away", 2, 46)], p)

        self.assertEqual(len(samples), 45)
        self.assertAlmostEqual(samples[42][2], FLOOR, places=2)
        self.assertAlmostEqual(samples[43][2], 75.10, places=2)

        period, second, home, away = samples[44]
        self.assertEqual((period, second), (2, 46 * 60))
        impulse = 23 * p.halftime_keep * math.exp(-1 / p.tau)
        self.assertAlmostEqual(home, 100 * (1 - math.exp(-(p.base + impulse) / p.k)))
        self.assertAlmostEqual(away, FLOOR, places=2)


if __name__ == "__main__":
    unittest.main()
