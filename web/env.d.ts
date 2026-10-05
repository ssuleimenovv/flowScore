/// <reference types="vite/client" />

interface ImportMetaEnv {
    // the gateway's address in production; unset in development (see shared/api/origin.ts)
    readonly VITE_API_URL?: string
}
// Fontsource ships no type declarations for this package (Onest and JetBrains Mono have them)
declare module '@fontsource-variable/big-shoulders-display'
