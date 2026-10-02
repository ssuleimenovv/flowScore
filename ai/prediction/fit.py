"""Fit and validate the outcome model: who wins from here (docs/PREDICTION.md).

Usage (from the repo root, after dataset.py):
    python ai/prediction/fit.py

The model counts goals. From any minute on, each team scores a Poisson number
of further goals with mean

    rate · minutes left,   log rate = a + b · features of that team

and the outcome is the current score plus those goals: home win, draw or away
win. The features say who is stronger and who has a reason to attack:

- home:   1 for the home team
- rating: this team's rating minus the other's, before the match
- red:    the other team's red cards minus this team's
- lead:   this team's goal lead, capped at 2 either way

The fit is a Poisson regression of the goals each team still scored, per
minute left: target goals / minutes, weighted by minutes.

Two independent Poisson counts give too few draws: a team level late in a
match settles for the point, and the other often lets it. One more number
fixes that: the draw is weighted by e^draw and the three renormalized, with
draw fitted on the same training matches after the regression.

Validation is the same as Flow's: 5 blocks of the season by date, each
predicted by a model fitted on the other 4, intervals by resampling whole
matches. The model without the draw fix is compared too, and two challengers
add what the Flow screen shows, Flow and xG: if they do not beat the plain
model, the plain model ships.

Writes ai/prediction/model.json (what Go loads) and ai/prediction/report.json.
"""

import json
from pathlib import Path

import numpy as np
import pandas as pd
from scipy.optimize import minimize_scalar
from scipy.stats import poisson
from sklearn.linear_model import PoissonRegressor

from dataset import DATA, PRIOR_GAMES

HERE = Path(__file__).resolve().parent
FOLDS = 5
BOOTSTRAP = 2000
LEAD_CAP = 2
MAX_GOALS = 10  # further goals per team; the chance of more is below 1e-6
VERSION = "Poisson v1"

FEATURES = ["home", "rating", "red", "lead"]
# name: (features, fix draws)
MODELS = {
    "outcome": (FEATURES, True),
    "no draw fix": (FEATURES, False),
    "with flow": (FEATURES + ["flow"], True),
    "with xg": (FEATURES + ["xg"], True),
}
OUTCOMES = ["home", "draw", "away"]


def features(df: pd.DataFrame, side: str) -> pd.DataFrame:
    """The features of one team in every snapshot."""
    other = "away" if side == "home" else "home"
    sign = 1.0 if side == "home" else -1.0
    return pd.DataFrame(
        {
            "home": 1.0 if side == "home" else 0.0,
            "rating": sign * df.rating,
            "red": df[f"{other}_reds"] - df[f"{side}_reds"],
            "lead": (df[f"{side}_goals"] - df[f"{other}_goals"]).clip(-LEAD_CAP, LEAD_CAP),
            "flow": (df[f"{side}_flow"] - df[f"{other}_flow"]) / 100,
            "xg": df[f"{side}_xg"] - df[f"{other}_xg"],
        }
    )


def result(df: pd.DataFrame) -> np.ndarray:
    """0 home win, 1 draw, 2 away win."""
    return np.sign(df.final_away - df.final_home).to_numpy().astype(int) + 1


def log_loss(p: np.ndarray, y: np.ndarray) -> np.ndarray:
    return -np.log(np.clip(p[np.arange(len(y)), y], 1e-12, None))


def outcome(home_mean: np.ndarray, away_mean: np.ndarray, lead: np.ndarray) -> np.ndarray:
    """P(home win, draw, away win) for further goals with these means and the current home lead."""
    k = np.arange(MAX_GOALS + 1)
    home = poisson.pmf(k, home_mean[:, None])
    away = poisson.pmf(k, away_mean[:, None])
    joint = home[:, :, None] * away[:, None, :]  # [snapshot, home goals, away goals]
    final = lead[:, None, None] + k[None, :, None] - k[None, None, :]
    p = np.stack(
        [(joint * (final > 0)).sum((1, 2)), (joint * (final == 0)).sum((1, 2)), (joint * (final < 0)).sum((1, 2))],
        axis=1,
    )
    return p / p.sum(1, keepdims=True)


def with_draw(p: np.ndarray, draw: float) -> np.ndarray:
    """The draw weighted by e^draw, the three renormalized."""
    q = p * np.array([1.0, np.exp(draw), 1.0])
    return q / q.sum(1, keepdims=True)


class Model:
    """The goal rates and the draw fix, fitted on some snapshots."""

    def __init__(self, df: pd.DataFrame, columns: list[str], fix_draws: bool = True):
        self.columns = columns
        # Both teams of every snapshot as rows: the goals each still scored per minute left
        x = pd.concat([features(df, "home")[columns], features(df, "away")[columns]])
        goals = np.concatenate([df.final_home - df.home_goals, df.final_away - df.away_goals])
        minutes = np.concatenate([df.remaining, df.remaining])
        self.rates = PoissonRegressor(alpha=0, max_iter=1000).fit(x, goals / minutes, sample_weight=minutes)

        self.draw = 0.0
        if fix_draws:
            p, y = self.goals(df), result(df)
            best = minimize_scalar(lambda d: log_loss(with_draw(p, d), y).mean(), bounds=(-1, 1), method="bounded")
            self.draw = float(best.x)

    def goals(self, df: pd.DataFrame) -> np.ndarray:
        """The outcome from the goal counts alone, before the draw fix."""
        home = self.rates.predict(features(df, "home")[self.columns]) * df.remaining.to_numpy()
        away = self.rates.predict(features(df, "away")[self.columns]) * df.remaining.to_numpy()
        return outcome(home, away, (df.home_goals - df.away_goals).to_numpy())

    def predict(self, df: pd.DataFrame) -> np.ndarray:
        return with_draw(self.goals(df), self.draw)


def gain_interval(better: np.ndarray, worse: np.ndarray, match: np.ndarray) -> tuple[float, float]:
    """95% interval of how much lower the log loss of `better` is, resampling whole matches."""
    per_match = pd.Series(worse - better).groupby(match).agg(["sum", "count"])
    sums, counts = per_match["sum"].to_numpy(), per_match["count"].to_numpy()
    picks = np.random.default_rng(7).integers(0, len(sums), size=(BOOTSTRAP, len(sums)))
    means = sums[picks].sum(axis=1) / counts[picks].sum(axis=1)
    return float(np.percentile(means, 2.5)), float(np.percentile(means, 97.5))


def calibration(p: np.ndarray, hit: np.ndarray, bins: int = 10) -> tuple[float, list]:
    """How far stated chances are from what happened: a 70% that comes true 70%
    of the time is calibrated. Returns the mean gap and the table by decile."""
    edges = np.minimum((p * bins).astype(int), bins - 1)
    table, gap = [], 0.0
    for b in range(bins):
        inside = edges == b
        if inside.any():
            said, came = p[inside].mean(), hit[inside].mean()
            table.append([round(float(said), 3), round(float(came), 3), int(inside.sum())])
            gap += inside.mean() * abs(said - came)
    return float(gap), table


def season_ratings(df: pd.DataFrame) -> dict[str, float]:
    """Every team's rating after the whole season, the one a new match starts from."""
    events = pd.read_csv(DATA / "events.csv", dtype={"match_id": str})
    xg = events.groupby(["match_id", "side"]).xg.sum().unstack(fill_value=0.0)
    games: dict[str, list[float]] = {}
    for m in df.drop_duplicates("match_id").itertuples():
        diff = xg.at[m.match_id, "home"] - xg.at[m.match_id, "away"]
        games.setdefault(m.home, []).append(diff)
        games.setdefault(m.away, []).append(-diff)
    return {team: round(sum(g) / (len(g) + PRIOR_GAMES), 4) for team, g in sorted(games.items())}


def main() -> None:
    df = pd.read_csv(DATA / "snapshots.csv", dtype={"match_id": str})
    by_date = df.drop_duplicates("match_id").sort_values("date", kind="stable").match_id.to_numpy()
    block = {m: k for k, part in enumerate(np.array_split(by_date, FOLDS)) for m in part}
    df["block"] = df.match_id.map(block)
    y = result(df)
    first = ~df.match_id.duplicated().to_numpy()  # the first minute: the pre-match view
    print(f"{len(by_date)} matches, {len(df)} snapshots, {FOLDS} blocks by date")

    predictions = {name: np.zeros((len(df), 3)) for name in ["constant", *MODELS]}
    draws = []
    for k in range(FOLDS):
        train, test = df[df.block != k], (df.block == k).to_numpy()
        shares = np.bincount(result(train.drop_duplicates("match_id")), minlength=3)
        predictions["constant"][test] = shares / shares.sum()
        for name, (columns, fix_draws) in MODELS.items():
            model = Model(train, columns, fix_draws)
            predictions[name][test] = model.predict(df[test])
            if name == "outcome":
                draws.append(round(model.draw, 3))

    report = {"model": VERSION, "matches": len(by_date), "snapshots": len(df), "models": {}, "gain95": {}}
    ll = {}
    print("\nout of fold, every match predicted by a model that never saw it:")
    for name, p in predictions.items():
        ll[name] = log_loss(p, y)
        brier = ((p - np.eye(3)[y]) ** 2).sum(1).mean()
        report["models"][name] = {
            "logLoss": round(float(ll[name].mean()), 5),
            "preMatchLogLoss": round(float(ll[name][first].mean()), 5),
            "brier": round(float(brier), 5),
        }
        print(f"  {name:<11} log loss {ll[name].mean():.4f}   before kickoff {ll[name][first].mean():.4f}   Brier {brier:.4f}")

    pairs = [("outcome", "constant"), ("outcome", "no draw fix"), ("with flow", "outcome"), ("with xg", "outcome")]
    for better, worse in pairs:
        low, high = gain_interval(ll[better], ll[worse], df.match_id.to_numpy())
        verdict = "better" if low > 0 else "worse" if high < 0 else "no clear difference"
        report["gain95"][f"{better} over {worse}"] = [round(low, 5), round(high, 5)]
        print(f"  {better} over {worse}: log loss gain 95% [{low:.4f}, {high:.4f}], {verdict}")

    # A fix worth trusting lands in the same place on every block
    print(f"\ndraw fix fitted on each block: {draws}")
    report["drawByBlock"] = draws

    print("calibration, mean gap between stated chance and how often it came true:")
    report["calibration"] = {}
    for i, name in enumerate(OUTCOMES):
        gap, table = calibration(predictions["outcome"][:, i], y == i)
        plain, _ = calibration(predictions["no draw fix"][:, i], y == i)
        report["calibration"][name] = {"meanGap": round(gap, 4), "noDrawFix": round(plain, 4), "deciles": table}
        print(f"  {name:<5} {gap:.3f}   without the draw fix {plain:.3f}")

    # The model that ships is fitted on the whole season
    model = Model(df, FEATURES)
    weights = dict(zip(FEATURES, (round(float(w), 5) for w in model.rates.coef_)))
    intercept, draw = round(float(model.rates.intercept_), 5), round(model.draw, 5)
    print(f"\nthe whole season: weights {weights}, intercept {intercept}, draw {draw}")
    out = {
        "model": VERSION,
        "intercept": intercept,
        "weights": weights,
        "draw": draw,
        "leadCap": LEAD_CAP,
        "maxGoals": MAX_GOALS,
        "ratings": season_ratings(df),
    }
    report["final"] = weights | {"intercept": intercept, "draw": draw}
    (HERE / "model.json").write_text(json.dumps(out, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")
    (HERE / "report.json").write_text(json.dumps(report, indent=2) + "\n")
    print(f"\n-> {HERE / 'model.json'}\n-> {HERE / 'report.json'}")


if __name__ == "__main__":
    main()
