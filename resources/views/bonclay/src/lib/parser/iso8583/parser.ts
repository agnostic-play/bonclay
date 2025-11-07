
// =========================
// File: src/lib/iso8583/parser.ts
// =========================

import type { FieldSpec, MessageSpec } from '@/lib/parser/iso8583/spec'

export interface ParsedField {
    id: number
    name: string
    lenType: 'FIXED' | 'LL' | 'LLL'
    declaredLength?: number // for LL/LLL
    actualLength: number
    value: string
}

export interface ParsedISOMessage {
    mti: string
    bitmapHex: string
    bits: number[]
    fields: ParsedField[]
    map: Record<string, string>
}

function hexToBits(hex: string): number[] {
    const bits: number[] = []
    for (let i = 0; i < hex.length; i++) {
        const nibble = parseInt(hex[i], 16)
        if (Number.isNaN(nibble)) throw new Error('Bitmap contains non-hex characters')
        for (let b = 3; b >= 0; b--) bits.push((nibble >> b) & 1)
    }
    return bits
}

function readLen(src: string, pos: number, lenDigits: number): { len: number; pos: number } {
    const slice = src.slice(pos, pos + lenDigits)
    if (slice.length !== lenDigits) throw new Error('Unexpected end while reading length prefix')
    const len = parseInt(slice, 10)
    if (Number.isNaN(len)) throw new Error('Length prefix is not numeric')
    return { len, pos: pos + lenDigits }
}

function readFixed(src: string, pos: number, n: number): { v: string; pos: number } {
    const end = pos + n
    if (end > src.length) throw new Error('Unexpected end while reading fixed field')
    return { v: src.slice(pos, end), pos: end }
}

export function parseISO8583(raw: string, spec: MessageSpec): ParsedISOMessage {
    const s = (raw || '').trim()
    if (s.length < 20) throw new Error('Message too short (need MTI + bitmap)')

    // MTI (ASCII 4)
    let pos = 0
    const mti = s.slice(pos, pos + 4)
    if (mti.length !== 4) throw new Error('Invalid MTI')
    pos += 4

    // Primary bitmap (16 hex ASCII)
    const primaryBitmap = s.slice(pos, pos + 16)
    if (primaryBitmap.length !== 16) throw new Error('Invalid primary bitmap length')
    pos += 16

    // If bit 1 of primary bitmap is set, read secondary (next 16 hex)
    const primaryBits = hexToBits(primaryBitmap)
    let bitmapHex = primaryBitmap
    if (primaryBits[0] === 1) {
        const secondary = s.slice(pos, pos + 16)
        if (secondary.length !== 16) throw new Error('Secondary bitmap expected but missing')
        pos += 16
        bitmapHex += secondary
    }

    const bits = hexToBits(bitmapHex)
    const present: number[] = []
    for (let i = 1; i < bits.length; i++) if (bits[i] === 1) present.push(i + 1)

    const fields: ParsedField[] = []
    const map: Record<string, string> = {}

    for (const id of present) {
        const def: FieldSpec | undefined = spec.fields[id]
        if (!def) {
            // Skip gracefully: if a field is not defined, stop with error
            throw new Error(`Field ${id} is present in bitmap but not defined in spec`)
        }

        let value = ''
        let declaredLength: number | undefined

        if (def.lenType === 'FIXED') {
            const { v, pos: p2 } = readFixed(s, pos, def.length)
            pos = p2
            value = v
        } else if (def.lenType === 'LL') {
            const r = readLen(s, pos, 2)
            declaredLength = r.len
            pos = r.pos
            if (declaredLength > def.length) throw new Error(`Field ${id} length ${declaredLength} exceeds max ${def.length}`)
            const { v, pos: p2 } = readFixed(s, pos, declaredLength)
            pos = p2
            value = v
        } else if (def.lenType === 'LLL') {
            const r = readLen(s, pos, 3)
            declaredLength = r.len
            pos = r.pos
            if (declaredLength > def.length) throw new Error(`Field ${id} length ${declaredLength} exceeds max ${def.length}`)
            const { v, pos: p2 } = readFixed(s, pos, declaredLength)
            pos = p2
            value = v
        }

        fields.push({
            id,
            name: def.name,
            lenType: def.lenType,
            declaredLength,
            actualLength: value.length,
            value,
        })
        map[id.toString().padStart(3, '0')] = value
    }

    return { mti, bitmapHex, bits, fields, map }
}
