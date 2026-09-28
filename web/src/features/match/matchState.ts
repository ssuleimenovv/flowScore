import type {
  EventList,
  FlowPoint,
  FlowSeries,
  FlowValues,
  Match,
  MatchEvent,
  Score,
  Side,
  WsMessage,
} from '@/shared/api/types'

// How many timeline events to keep; older ones are dropped.
const MAX_EVENTS = 100

// MatchLive is the match screen's data: a REST snapshot plus the stream on top.
export interface MatchLive {
  match: Match
  score: Score
  flow: FlowValues
  delta10: FlowValues
  points: FlowPoint[] // by minute, ascending
  events: MatchEvent[] // newest first
  // Seq of each part of the snapshot. The three requests are answered at
  // slightly different moments, so each part skips the messages it already has.
  seq: { match: number; flow: number; events: number }
}

export function fromSnapshot(match: Match, flow: FlowSeries, events: EventList): MatchLive {
  return {
    match,
    score: match.score,
    flow: flow.current,
    delta10: flow.delta10,
    points: [...flow.points].sort((a, b) => a.minute - b.minute),
    events: events.items.slice(0, MAX_EVENTS),
    seq: { match: match.seq, flow: flow.seq, events: events.seq },
  }
}

// applyMessage returns the state with one stream message applied.
// It never changes the old state, so Vue sees a new object and re-renders.
export function applyMessage(state: MatchLive, message: WsMessage): MatchLive {
  switch (message.type) {
    case 'flow.update': {
      if (message.seq <= state.seq.flow) return state
      const { current, delta10, point, clock } = message.data
      return {
        ...state,
        // The clock resyncs on every update: the client ticks in real time,
        // while a replay runs faster and a real match pauses at half-time.
        // The match part may be newer than the flow part, so only move it forward
        match: message.seq > state.seq.match ? { ...state.match, clock } : state.match,
        flow: current,
        delta10,
        points: upsertPoint(state.points, point),
        seq: { ...state.seq, flow: message.seq },
      }
    }

    case 'match.event': {
      const event = message.data
      let next = state
      if (message.seq > state.seq.events) {
        next = {
          ...next,
          events: [event, ...next.events].slice(0, MAX_EVENTS),
          seq: { ...next.seq, events: message.seq },
        }
      }
      if (event.type === 'goal' && message.seq > state.seq.match) {
        next = {
          ...next,
          score: addGoal(next.score, event.side),
          seq: { ...next.seq, match: message.seq },
        }
      }
      return next
    }

    default:
      return state // prediction.update and insight.update come with the AI step
  }
}

// upsertPoint replaces the point of the same minute or inserts it in order.
function upsertPoint(points: FlowPoint[], point: FlowPoint): FlowPoint[] {
  const rest = points.filter((p) => p.minute !== point.minute)
  const at = rest.findIndex((p) => p.minute > point.minute)
  return at === -1 ? [...rest, point] : [...rest.slice(0, at), point, ...rest.slice(at)]
}

function addGoal(score: Score, side: Side | null): Score {
  if (side === 'home') return { ...score, home: score.home + 1 }
  if (side === 'away') return { ...score, away: score.away + 1 }
  return score
}
