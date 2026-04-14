/// <reference types="vite/client" />

declare module 'plantuml-encoder' {
  export function encode(input: string): string
  export function decode(input: string): string
  const plantumlEncoder: { encode: typeof encode; decode: typeof decode }
  export default plantumlEncoder
}
