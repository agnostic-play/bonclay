// src/lib/iso8583.ts
// Minimal, pragmatic ISO8583 (1987/1993-ish) parser for HEX or ASCII payloads.
//
// Assumptions:
// - Message starts with 4-char MTI (ASCII: "0200", etc.).
// - Bitmap(s) immediately follow MTI, expressed as 16/32 hex chars (primary(+secondary)).
// - Field values for n/an/ans are ASCII when inputMode='ascii'. When inputMode='hex', they are ASCII-encoded bytes represented as hex pairs.
// - VAR length prefixes (LL/LLL) are ASCII digits when ascii; when hex, the digits are ASCII bytes encoded as hex (e.g., '31 32' => "12").
//
// You can extend FIELD_DEFS to your host’s spec easily.

export type InputMode = 'hex' | 'ascii';

export type FieldDef = {
    id: number;
    name: string;
    type: 'n' | 'an' | 'ans' | 'b' | 'z';
    // fixed length (chars for ascii, bytes for binary), or max length for var
    length: number;
    var?: 2 | 3; // LLVAR or LLLVAR
    // store value as: 'ascii' (string) or 'hex' (preserve hex)
    store?: 'ascii' | 'hex';
};

export type ParsedIso8583 = {
    mti: string;
    bitmapPrimaryHex: string;
    bitmapSecondaryHex?: string;
    fields: Record<string, string>;
};

function stripSpaces(s: string) {
    return s.replace(/\s+/g, '');
}

function asciiToHex(s: string): string {
    // "0200" -> "30323030"
    let out = '';
    for (let i = 0; i < s.length; i++) {
        out += s.charCodeAt(i).toString(16).padStart(2, '0');
    }
    return out;
}

function hexToAscii(hex: string): string {
    let out = '';
    for (let i = 0; i < hex.length; i += 2) {
        out += String.fromCharCode(parseInt(hex.slice(i, i + 2), 16));
    }
    return out;
}

function hexToBits(hex: string): string {
    return hex
        .match(/.{1,2}/g)!
        .map(b => parseInt(b, 16).toString(2).padStart(8, '0'))
        .join('');
}

function takeAscii(src: string, pos: number, len: number) {
    // ASCII string slice
    return { val: src.slice(pos, pos + len), next: pos + len };
}

function takeHex(srcHex: string, pos: number, lenChars: number) {
    // lenChars = number of ASCII chars to read, each is one byte => 2 hex chars per ASCII char
    const hexLen = lenChars * 2;
    const chunk = srcHex.slice(pos, pos + hexLen);
    return { valHex: chunk, next: pos + hexLen };
}

function readVarLen(
    src: string,
    srcHex: string,
    pos: number,
    mode: InputMode,
    varDigits: 2 | 3
): { length: number; next: number } {
    if (mode === 'ascii') {
        const { val, next } = takeAscii(src, pos, varDigits);
        const n = parseInt(val, 10);
        if (Number.isNaN(n)) throw new Error(`Bad VAR length '${val}' @${pos}`);
        return { length: n, next };
    } else {
        // hex mode: LL/LLL are ASCII digits encoded as hex
        const { valHex, next } = takeHex(srcHex, pos, varDigits);
        const digits = hexToAscii(valHex);
        const n = parseInt(digits, 10);
        if (Number.isNaN(n)) throw new Error(`Bad VAR length '${digits}' @${pos}`);
        return { length: n, next };
    }
}

function readValue(
    src: string,
    srcHex: string,
    pos: number,
    mode: InputMode,
    def: FieldDef,
    dataLen: number // characters for ascii, characters for n/an/ans; bytes for 'b' if store='hex'
): { value: string; next: number } {
    if (def.store === 'hex' || (def.type === 'b' && mode === 'hex')) {
        // consume 'dataLen' ASCII chars => 2 * dataLen hex chars
        const { valHex, next } = takeHex(srcHex, pos, dataLen);
        return { value: valHex.toLowerCase(), next };
    }

    if (mode === 'ascii') {
        const { val, next } = takeAscii(src, pos, dataLen);
        return { value: val, next };
    } else {
        const { valHex, next } = takeHex(srcHex, pos, dataLen);
        return { value: hexToAscii(valHex), next };
    }
}

/** Default field dictionary — adjust to your host spec as needed */
export const FIELD_DEFS: Record<number, FieldDef> = {
    2:  { id: 2,  name: 'Primary Account Number', type: 'n',   length: 19, var: 2 },
    3:  { id: 3,  name: 'Processing Code',        type: 'n',   length: 6 },
    4:  { id: 4,  name: 'Amount, Transaction',    type: 'n',   length: 12 },
    7:  { id: 7,  name: 'Transmission Date & Time', type: 'n', length: 10 },
    11: { id: 11, name: 'STAN',                   type: 'n',   length: 6 },
    12: { id: 12, name: 'Time, Local Txn',        type: 'n',   length: 6 },
    13: { id: 13, name: 'Date, Local Txn',        type: 'n',   length: 4 },
    14: { id: 14, name: 'Expiration Date',        type: 'n',   length: 4 },
    15: { id: 15, name: 'Settlement Date',        type: 'n',   length: 4 },
    18: { id: 18, name: 'Merchant Category Code', type: 'n',   length: 4 },
    22: { id: 22, name: 'POS Entry Mode',         type: 'n',   length: 3 },
    23: { id: 23, name: 'PAN Sequence',           type: 'n',   length: 3 },
    25: { id: 25, name: 'POS Condition Code',     type: 'n',   length: 2 },
    32: { id: 32, name: 'Acquiring Institution ID', type: 'n', length: 11, var: 2 },
    35: { id: 35, name: 'Track 2 Data',           type: 'z',   length: 37, var: 2 },
    37: { id: 37, name: 'Retrieval Reference No', type: 'an',  length: 12 },
    38: { id: 38, name: 'Auth ID Response',       type: 'an',  length: 6 },
    39: { id: 39, name: 'Response Code',          type: 'an',  length: 2 },
    41: { id: 41, name: 'Terminal ID',            type: 'an',  length: 8 },
    42: { id: 42, name: 'Merchant ID',            type: 'an',  length: 15 },
    43: { id: 43, name: 'Card Acceptor Name/Loc', type: 'ans', length: 40 },
    48: { id: 48, name: 'Additional Data (Priv)', type: 'ans', length: 999, var: 3 },
    49: { id: 49, name: 'Currency, Txn',          type: 'n',   length: 3 },
    52: { id: 52, name: 'PIN Data',               type: 'b',   length: 8,  store: 'hex' }, // 8 bytes -> 16 hex
    55: { id: 55, name: 'EMV ICC Data',           type: 'b',   length: 255, var: 3, store: 'hex' },
    60: { id: 60, name: 'Advice/Reserved',        type: 'ans', length: 999, var: 3 },
    61: { id: 61, name: 'POS Data',               type: 'ans', length: 999, var: 3 },
    62: { id: 62, name: 'Invoice/Priv',           type: 'ans', length: 999, var: 3 },
    63: { id: 63, name: 'Network Data',           type: 'ans', length: 999, var: 3 },
    64: { id: 64, name: 'MAC',                    type: 'b',   length: 8,  store: 'hex' },
};

function ensureHexBitmap(hex: string) {
    if (!/^[0-9a-fA-F]+$/.test(hex) || (hex.length !== 16 && hex.length !== 32)) {
        throw new Error(`Bitmap must be 16 or 32 hex chars, got length=${hex.length}`);
    }
}

/**
 * Parse ISO8583 message.
 * @param raw Input string (hex or ascii)
 * @param mode 'hex' (default) treats the input as hex pairs. 'ascii' treats input as plain text.
 */
export function parseIso8583(raw: string, mode: InputMode = 'hex', defs: Record<number, FieldDef> = FIELD_DEFS): ParsedIso8583 {
    if (!raw) throw new Error('Empty message');
    let src = mode === 'ascii' ? raw : '';
    let srcHex = mode === 'hex' ? stripSpaces(raw) : asciiToHex(raw);

    // Read MTI (4 ASCII chars)
    let pos = 0;
    let mti = '';
    if (mode === 'ascii') {
        const t = takeAscii(src, pos, 4);
        mti = t.val;
        pos = t.next;
    } else {
        const t = takeHex(srcHex, pos, 4);
        mti = hexToAscii(t.valHex);
        pos = t.next;
    }
    if (!/^\d{4}$/.test(mti)) throw new Error(`Bad MTI: ${mti}`);

    // Read Primary Bitmap (16 hex chars => 8 bytes)
    let bmp1Hex = '';
    if (mode === 'ascii') {
        // Next 16 hex chars are bitmap (as ASCII hex)
        const t = takeAscii(src, pos, 16);
        bmp1Hex = t.val;
        pos = t.next;
    } else {
        const t = takeHex(srcHex, pos, 16);
        bmp1Hex = hexToAscii(t.valHex); // ascii view of hex, e.g., "7238A0..."
        pos = t.next;
    }
    bmp1Hex = stripSpaces(bmp1Hex);
    ensureHexBitmap(bmp1Hex);

    const bmp1Bits = hexToBits(bmp1Hex);
    const hasSecondary = bmp1Bits[0] === '1';

    let bmp2Hex: string | undefined;
    let bmpBits = bmp1Bits;

    if (hasSecondary) {
        if (mode === 'ascii') {
            const t = takeAscii(src, pos, 16);
            bmp2Hex = t.val;
            pos = t.next;
        } else {
            const t = takeHex(srcHex, pos, 16);
            bmp2Hex = hexToAscii(t.valHex);
            pos = t.next;
        }
        bmp2Hex = stripSpaces(bmp2Hex);
        ensureHexBitmap(bmp2Hex);
        bmpBits += hexToBits(bmp2Hex);
    }

    const fields: Record<string, string> = {};
    // Iterate DE 2..128
    for (let de = 2; de <= bmpBits.length; de++) {
        if (bmpBits[de - 1] !== '1') continue;

        const def = defs[de];
        if (!def) {
            // Unknown field: we can't parse without spec; abort gracefully.
            throw new Error(`Bitmap indicates field ${de} but no definition provided.`);
        }

        // Read value
        let value = '';
        if (def.var) {
            // Read LL/LLL length
            const { length: vlen, next } = readVarLen(src, srcHex, pos, mode, def.var);
            pos = next;

            if (def.type === 'b' && (def.store === 'hex' || mode === 'hex')) {
                // binary-ish var field: length is bytes; in hex we need 2 chars per byte
                const { value: val, next: n } = readValue(src, srcHex, pos, mode, def, vlen);
                value = val;
                pos = n;
            } else {
                // ascii-ish var field: vlen chars
                const { value: val, next: n } = readValue(src, srcHex, pos, mode, def, vlen);
                value = val;
                pos = n;
            }
        } else {
            // fixed
            const fixedLen = def.type === 'b' && (def.store === 'hex' || mode === 'hex')
                ? def.length // bytes
                : def.length; // chars
            const { value: val, next: n } = readValue(src, srcHex, pos, mode, def, fixedLen);
            value = val;
            pos = n;
        }

        fields[String(de).padStart(3, '0')] = value;
    }

    return {
        mti,
        bitmapPrimaryHex: bmp1Hex.toUpperCase(),
        bitmapSecondaryHex: bmp2Hex ? bmp2Hex.toUpperCase() : undefined,
        fields,
    };
}
