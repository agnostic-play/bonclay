/// <reference types="vite/client" />



interface ImportMetaEnv {
    readonly BONCLAY_URL: string
    // add more as needed...
    // readonly VITE_API_KEY: string
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}