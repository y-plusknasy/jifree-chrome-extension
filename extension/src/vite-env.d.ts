/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_ENDPOINT: string
  readonly VITE_SHARED_SECRET: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
