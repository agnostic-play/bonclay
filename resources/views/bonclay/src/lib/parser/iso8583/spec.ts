
// =========================
// File: src/lib/iso8583/specs.ts
// =========================

export type LengthType = 'FIXED' | 'LL' | 'LLL'

export interface FieldSpec {
    id: number
    name: string
    length: number // max length for LL/LLL, exact length for FIXED
    type: 'n' | 'ans' | 'b' | 'z' // numeric / alphanumeric+special / binary / track2
    lenType: LengthType
}

export interface MessageSpec {
    name: string
    bitmap: 'HEX' // (Hex-ASCII encoded bitmap for 1987 ASCII spec)
    fields: Record<number, FieldSpec>
}

// Helper to build field quickly
function f(
    id: number,
    name: string,
    length: number,
    lenType: LengthType,
    type: FieldSpec['type'] = 'ans'
): FieldSpec {
    return { id, name, length, lenType, type }
}

// --- ISO 8583:1987 ASCII (common subset) ---
// Notes: This is a pragmatic subset matching your Go spec where relevant.
// You can extend/trim fields as needed.
export const ISO87_ASCII: MessageSpec = {
    name: 'ISO 8583 v1987 ASCII',
    bitmap: 'HEX',
    fields: {
        2:  f(2,  'Primary Account Number', 19, 'LL', 'n'),
        3:  f(3,  'Processing Code', 6, 'FIXED', 'n'),
        4:  f(4,  'Transaction Amount', 12, 'FIXED', 'n'),
        5:  f(5,  'Settlement Amount', 12, 'FIXED', 'n'),
        6:  f(6,  'Billing Amount', 12, 'FIXED', 'n'),
        7:  f(7,  'Transmission Date & Time', 10, 'FIXED', 'n'),
        8:  f(8,  'Billing Fee Amount', 8, 'FIXED', 'n'),
        9:  f(9,  'Settlement Conversion Rate', 8, 'FIXED', 'n'),
        10: f(10, 'Cardholder Billing Conversion Rate', 8, 'FIXED', 'n'),
        11: f(11, 'STAN', 6, 'FIXED', 'n'),
        12: f(12, 'Local Transaction Time', 6, 'FIXED', 'n'),
        13: f(13, 'Local Transaction Date', 4, 'FIXED', 'n'),
        14: f(14, 'Expiration Date', 4, 'FIXED', 'n'),
        15: f(15, 'Settlement Date', 4, 'FIXED', 'n'),
        16: f(16, 'Currency Conversion Date', 4, 'FIXED', 'n'),
        17: f(17, 'Capture Date', 4, 'FIXED', 'n'),
        18: f(18, 'Merchant Type', 4, 'FIXED', 'n'),
        19: f(19, 'Acquiring Inst. Country Code', 3, 'FIXED', 'n'),
        20: f(20, 'PAN Extended Country Code', 3, 'FIXED', 'n'),
        21: f(21, 'Forwarding Inst. Country Code', 3, 'FIXED', 'n'),
        22: f(22, 'POS Entry Mode', 3, 'FIXED', 'n'),
        23: f(23, 'Card Sequence Number', 3, 'FIXED', 'n'),
        24: f(24, 'Function Code', 3, 'FIXED', 'n'),
        25: f(25, 'POS Condition Code', 2, 'FIXED', 'n'),
        26: f(26, 'POS PIN Capture Code', 2, 'FIXED', 'n'),
        27: f(27, 'Auth ID Response Length', 1, 'FIXED', 'n'),
        28: f(28, 'Amount Convenience Fee', 9, 'FIXED', 'ans'),
        29: f(29, 'Settlement Fee Amount', 9, 'FIXED', 'ans'),
        30: f(30, 'Transaction Processing Fee', 9, 'FIXED', 'ans'),
        31: f(31, 'Settlement Processing Fee', 9, 'FIXED', 'ans'),
        32: f(32, 'Acquiring Inst. ID Code', 11, 'LL', 'n'),
        33: f(33, 'Forwarding Inst. ID Code', 11, 'LL', 'n'),
        34: f(34, 'Extended PAN', 28, 'LL', 'n'),
        35: f(35, 'Track 2 Data', 37, 'LL', 'z'),
        36: f(36, 'Track 3 Data', 104, 'LLL', 'z'),
        37: f(37, 'Retrieval Reference Number', 12, 'FIXED', 'ans'),
        38: f(38, 'Approval Code', 6, 'FIXED', 'ans'),
        39: f(39, 'Response Code', 2, 'FIXED', 'ans'),
        40: f(40, 'Service Restriction Code', 3, 'FIXED', 'ans'),
        41: f(41, 'Terminal ID', 16, 'FIXED', 'ans'),
        42: f(42, 'Card Acceptor ID Code', 15, 'FIXED', 'ans'),
        43: f(43, 'Card Acceptor Name/Location', 40, 'FIXED', 'ans'),
        44: f(44, 'Additional Data', 99, 'LL', 'ans'),
        45: f(45, 'Track 1 Data', 76, 'LL', 'ans'),
        46: f(46, 'Additional Data (ISO)', 999, 'LLL', 'ans'),
        47: f(47, 'Additional Data (National)', 999, 'LLL', 'ans'),
        48: f(48, 'Additional Data', 999, 'LLL', 'ans'),
        49: f(49, 'Transaction Currency Code', 3, 'FIXED', 'n'),
        50: f(50, 'Settlement Currency Code', 3, 'FIXED', 'n'),
        51: f(51, 'Cardholder Billing Currency Code', 3, 'FIXED', 'n'),
        52: f(52, 'PIN Data', 8, 'FIXED', 'b'),
        53: f(53, 'Security Related Control Info', 16, 'FIXED', 'ans'),
        54: f(54, 'Additional Amounts', 120, 'LLL', 'ans'),
        55: f(55, 'ICC Data – EMV TLV', 765, 'LLL', 'ans'),
        56: f(56, 'Reserved (ISO)', 999, 'LLL', 'ans'),
        57: f(57, 'Additional Data - National', 999, 'LLL', 'ans'),
        58: f(58, 'Reserved (National)', 999, 'LLL', 'ans'),
        59: f(59, 'Reserved (National)', 999, 'LLL', 'ans'),
        60: f(60, 'Reserved (National)', 999, 'LLL', 'ans'),
        61: f(61, 'Reserved (Private)', 999, 'LLL', 'ans'),
        62: f(62, 'Reserved (Private)', 999, 'LLL', 'ans'),
        63: f(63, 'Reserved (Private)', 999, 'LLL', 'ans'),
        64: f(64, 'MAC', 8, 'FIXED', 'b'),
        70: f(70, 'Network Management Code', 3, 'FIXED', 'n'),
        90: f(90, 'Original Data Elements', 42, 'FIXED', 'ans'),
        100: f(100, 'Issuer ID', 11, 'LL', 'n'),
        102: f(102, 'Account Identification 1', 19, 'LL', 'ans'),
        103: f(103, 'Beneficiary Account', 28, 'LL', 'ans'),
        123: f(123, 'Invoice Number', 20, 'LLL', 'ans'),
        125: f(125, 'Transaction Indicator', 255, 'LLL', 'ans'),
    },
}

// --- Your custom Rintis spec ---
// Mirrors your Go definitions; field names aligned for UI clarity.
export const RINTIS_ASCII: MessageSpec = {
    name: 'ISO 8583 v1987 ASCII (Rintis custom)',
    bitmap: 'HEX',
    fields: ISO87_ASCII.fields, // identical layout per your provided Go spec
}

// --- Free-form specs -------------------------------------------------
function makeFreeSpec(name: string, lenType: LengthType): MessageSpec {
    const fields: Record<number, FieldSpec> = {}
    const maxLen = lenType === 'LL' ? 99 : 999
    for (let id = 2; id <= 128; id++) {
        fields[id] = f(id, `DE ${id}`, maxLen, lenType, 'ans')
    }
    return { name, bitmap: 'HEX', fields }
}

export const FREE_LL: MessageSpec = makeFreeSpec('Free (LL everywhere, max 99)', 'LL')
export const FREE_LLL: MessageSpec = makeFreeSpec('Free (LLL everywhere, max 999)', 'LLL')

export const AVAILABLE_SPECS: Record<string, MessageSpec> = {
    'iso87': ISO87_ASCII,
    'rintis': RINTIS_ASCII,
    'free_ll': FREE_LL,
    'free_lll': FREE_LLL,
}
