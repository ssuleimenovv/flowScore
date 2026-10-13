# Design

The mockup lives in `mockup/`: one `.dc.html` per screen and `canvas.json`, the canvas that lays them out.

- Original mockup: https://claude.ai/artifact/HKt9MPjJvKDFysSLXWUfjz
- extended-v, the screens the original did not have: https://claude.ai/artifact/4neAe9Jim8EJQ8Sj2gd4N1
  - `MatchExtended`: lineups on the pitch and as a list, bench, form and head-to-head, odds against our model, news
  - `MatchPhases`: the match before kick-off, live and after the whistle, the shot map with our xG, players by their Flow
  - `CompetitionsExtended`: European and domestic cups, the Champions League table, two-legged ties, extra time and shoot-outs
  - `HomeDays`: the home screen by day, matches grouped by competition
  - `LeagueStats`: top scorers, assists, our xG and Flow leaders, fixtures and results
  - `PlayerSeason`: a player's season in numbers, goals against xG, the match log
  - `Onboarding`: first launch, picking competitions, teams and notifications
  - `About`: data sources and their terms, our models and how well they do, version, legal
  - `StatesExtended`: empty, loading and error states of the screens above

Each comes as desktop and mobile, dark and light (`…Light`, `Mobile…`). On the canvas the extended-v boards sit below the others, under the title "extended-v".
