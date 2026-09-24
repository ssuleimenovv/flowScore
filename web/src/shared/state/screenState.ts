export type RequestStatus = 'idle' | 'pending' | 'success' | 'error'
export type SocketStatus = 'open' | 'connecting' | 'closed'

export type ScreenView = 'loading' | 'empty' | 'error' | 'offline' | 'content'
export type ScreenBanner = 'offline' | 'reconnecting' | null

export interface ScreenInput {
  status: RequestStatus
  hasData: boolean
  isEmpty: boolean
  online: boolean
  socket?: SocketStatus
}

export interface ScreenState {
  view: ScreenView
  banner: ScreenBanner
}

export function resolveScreenState(input: ScreenInput): ScreenState {
  return {
    view: resolveView(input),
    banner: resolveBanner(input),
  }
}

function resolveView({ status, hasData, isEmpty, online }: ScreenInput): ScreenView {
  if (hasData) return isEmpty ? 'empty' : 'content'
  if (!online) return 'offline'
  if (status === 'error') return 'error'
  return 'loading'
}

function resolveBanner({ hasData, online, socket }: ScreenInput): ScreenBanner {
  if (!hasData) return null
  if (!online) return 'offline'
  if (socket && socket !== 'open') return 'reconnecting'
  return null
}
