import type { components } from './schema'

type Schemas = components['schemas']

export type WsMessage = Schemas['WsMessage']
export type MatchEvent = Schemas['MatchEvent']
export type FlowValues = Schemas['FlowValues']
export type FlowPoint = Schemas['FlowPoint']
