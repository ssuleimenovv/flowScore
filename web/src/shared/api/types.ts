import type { components } from './schema'

type Schemas = components['schemas']

export type WsMessage = Schemas['WsMessage']
export type MatchEvent = Schemas['MatchEvent']
export type FlowValues = Schemas['FlowValues']
export type FlowPoint = Schemas['FlowPoint']
export type Match = Schemas['Match']
export type FlowSeries = Schemas['FlowSeries']
export type EventList = Schemas['EventList']
export type Score = Schemas['Score']
export type Side = Schemas['Side']
export type Problem = Schemas['Problem']
