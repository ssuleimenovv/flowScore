import type { Match, Side } from '@/shared/api/types'

// Crests and league logos from ESPN's image server, cut to the size they are
// shown at. A real match's team ids are ESPN's, so they lead straight to the
// crests; a demo match's are StatsBomb's and have none.
const IMAGES = 'https://a.espncdn.com/combiner/i'

// League logos have numbers of their own: 23 is the Premier League
const LEAGUE_LOGOS: Record<string, number> = {
  'eng.1': 23,
  'esp.1': 15,
  'ger.1': 10,
  'ita.1': 12,
  'fra.1': 9,
}

// Twice the size on screen, so the edges stay sharp on a retina display
function image(path: string, px: number): string {
  const size = px * 2
  return `${IMAGES}?img=${path}&w=${size}&h=${size}`
}

export function teamLogo(teamId: string, px: number): string {
  return image(`/i/teamlogos/soccer/500/${encodeURIComponent(teamId)}.png`, px)
}

// null for a league whose logo we do not know the number of. ESPN draws each
// logo twice: for a light background, and in lighter colors for a dark one.
export function leagueLogo(
  leagueId: string,
  px: number,
  on: 'light' | 'dark' = 'light',
): string | null {
  const n = LEAGUE_LOGOS[leagueId]
  const folder = on === 'dark' ? '500-dark' : '500'
  return n === undefined ? null : image(`/i/leaguelogos/soccer/${folder}/${n}.png`, px)
}

// The demo matches are the Premier League's of 2015/16, whose teams have
// StatsBomb ids: their crests are found by name. Norwich and Leicester are
// not in the league now, but ESPN keeps their crests.
const DEMO_TEAMS: Record<string, string> = {
  Arsenal: '359',
  'Leicester City': '375',
  Liverpool: '364',
  'Manchester City': '382',
  'Manchester United': '360',
  'Norwich City': '381',
}

// A team of a match: a real match's id is ESPN's, a demo team is found by name
export function matchTeamLogo(
  match: Pick<Match, 'source' | 'home' | 'away'>,
  side: Side,
  px: number,
): string | null {
  const team = match[side]
  const id = match.source === 'espn' ? team.id : DEMO_TEAMS[team.name]
  return id ? teamLogo(id, px) : null
}
