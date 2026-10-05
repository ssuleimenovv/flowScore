// Where the gateway lives. In development VITE_API_URL is not set: Vite serves 
// the page and proxies /api and /ws to the gateway on the same host.
// In production the page is on Vercel and the gateway on Render, so the build 
// gets its address: VITE_API_URL=https://flowscore-api.onrender.com
const API_URL = (import.meta.env.VITE_API_URL ?? '').replace(/\/+$/, '')

// "api/v1/matches/1" -> "https://flowscore-api-onrender.com/api/v1/matches/1"
export function apiUrl(path: string): string {
    return API_URL + path
}

// the same for the stream, with ws:// for http:// and wss:// for https://
export function streamUrl(path: string): string {
    return (API_URL || location.origin).replace(/^http/, 'ws') + path
}
