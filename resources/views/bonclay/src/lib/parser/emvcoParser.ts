// Vue + Vite friendly TypeScript port of your Go QR EMV parser
// — with absolute index positions for each tag value.

export type TagMap = Record<string, string>;
export type TagPositions = Record<string, { start: number; end: number }>;

type TLV = {
    id: string;
    len: number;
    value: string;
    // absolute positions in the ORIGINAL string (0-based)
    start: number;       // header start (ID)
    end: number;         // end of TLV (exclusive)
    valueStart: number;  // start of the value (after 4-char header)
    valueEnd: number;    // end of value (exclusive)
};

/** Read a single TLV at position i; throws on malformed input */
function readTLV(s: string, i: number, base = 0): { tlv: TLV; next: number } {
    if (i + 4 > s.length) throw new Error(`TLV header truncated at index ${base + i}`);
    const id = s.slice(i, i + 2);
    const lenStr = s.slice(i + 2, i + 4);
    const len = Number.parseInt(lenStr, 10);
    if (Number.isNaN(len)) throw new Error(`Invalid length '${lenStr}' at index ${base + i + 2}`);
    const start = i + 4;
    const end = start + len;
    if (end > s.length) {
        throw new Error(
            `TLV value out of bounds: need ${len} chars from ${base + start}, have ${s.length - start}`
        );
    }
    const value = s.slice(start, end);
    return {
        tlv: {
            id,
            len,
            value,
            start: base + i,
            end: base + end,
            valueStart: base + start,
            valueEnd: base + end,
        },
        next: end,
    };
}

/** Generic TLV loop (supports absolute base offset) */
function parseTLVString(s: string, base = 0): TLV[] {
    const out: TLV[] = [];
    let i = 0;
    while (i < s.length) {
        const { tlv, next } = readTLV(s, i, base);
        out.push(tlv);
        i = next;
    }
    return out;
}

/** Luhn mod-10 check digit (returns the single digit to append) */
function luhnCheckDigit(num: string): number {
    const digits = num.replace(/\D+/g, '').split('').map(d => Number(d));
    let sum = 0;
    let dbl = true;
    for (let i = digits.length - 1; i >= 0; i--) {
        let n = digits[i];
        if (dbl) { n *= 2; if (n > 9) n -= 9; }
        sum += n;
        dbl = !dbl;
    }
    return (10 - (sum % 10)) % 10;
}

function setKV(
    map: TagMap,
    pos: TagPositions,
    key: string,
    value: string,
    start: number,
    end: number
) {
    map[key] = value;
    pos[key] = { start, end };
}

/* -------- Sub-parsers (26/27/28/40/51/62.99) -------- */

function parse26(value: string, out: TagMap, positions: TagPositions, base: number) {
    for (const tlv of parseTLVString(value, base)) {
        const k =
            tlv.id === '00' ? 'data26_00' :
                tlv.id === '01' ? 'data26_01' :
                    tlv.id === '02' ? 'data26_02' :
                        tlv.id === '03' ? 'data26_03' : `data26_${tlv.id}`;
        const v = tlv.id === '01' ? tlv.value + String(luhnCheckDigit(tlv.value)) : tlv.value;
        setKV(out, positions, k, v, tlv.valueStart, tlv.valueEnd);
    }
}

function parse27(value: string, out: TagMap, positions: TagPositions, base: number) {
    for (const tlv of parseTLVString(value, base)) {
        const k =
            tlv.id === '00' ? 'data27_00' :
                tlv.id === '01' ? 'data27_01' :
                    tlv.id === '02' ? 'data27_02' :
                        tlv.id === '03' ? 'data27_03' : `data27_${tlv.id}`;
        const v = tlv.id === '01' ? tlv.value + String(luhnCheckDigit(tlv.value)) : tlv.value;
        setKV(out, positions, k, v, tlv.valueStart, tlv.valueEnd);
    }
}

function parse28(value: string, out: TagMap, positions: TagPositions, base: number) {
    for (const tlv of parseTLVString(value, base)) {
        const k =
            tlv.id === '00' ? 'data28_00' :
                tlv.id === '01' ? 'data28_01' :
                    tlv.id === '02' ? 'data28_02' :
                        tlv.id === '03' ? 'data28_03' : `data28_${tlv.id}`;
        const v = tlv.id === '01' ? tlv.value + String(luhnCheckDigit(tlv.value)) : tlv.value;
        setKV(out, positions, k, v, tlv.valueStart, tlv.valueEnd);
    }
}

function parse40(value: string, out: TagMap, positions: TagPositions, base: number) {
    for (const tlv of parseTLVString(value, base)) {
        const k =
            tlv.id === '00' ? 'data40_00' :
                tlv.id === '01' ? 'data40_01' :
                    tlv.id === '02' ? 'data40_02' : `data40_${tlv.id}`;
        setKV(out, positions, k, tlv.value, tlv.valueStart, tlv.valueEnd);
    }
}

function parse51(value: string, out: TagMap, positions: TagPositions, base: number) {
    for (const tlv of parseTLVString(value, base)) {
        const k =
            tlv.id === '00' ? 'data51_00' :
                tlv.id === '01' ? 'data51_01' :
                    tlv.id === '02' ? 'data51_02' :
                        tlv.id === '03' ? 'data51_03' : `data51_${tlv.id}`;
        setKV(out, positions, k, tlv.value, tlv.valueStart, tlv.valueEnd);
    }
}

function parse6299(value: string, out: TagMap, positions: TagPositions, base: number) {
    for (const tlv of parseTLVString(value, base)) {
        const k = tlv.id === '01' ? 'data62_99_01' : `data62_99_${tlv.id}`;
        setKV(out, positions, k, tlv.value, tlv.valueStart, tlv.valueEnd);
    }
}

function parse62(value: string, out: TagMap, positions: TagPositions, base: number) {
    for (const tlv of parseTLVString(value, base)) {
        switch (tlv.id) {
            case '02': setKV(out, positions, 'data62_02', tlv.value, tlv.valueStart, tlv.valueEnd); break;
            case '07': setKV(out, positions, 'data62_07', tlv.value, tlv.valueStart, tlv.valueEnd); break;
            case '99': parse6299(tlv.value, out, positions, tlv.valueStart); break;
            default:   setKV(out, positions, `data62_${tlv.id}`, tlv.value, tlv.valueStart, tlv.valueEnd);
        }
    }
}

/** Metadata map aligned closer to EMVCo MPM 1.1 */
export const emvTags: Record<
    string,
    { id: string; name: string; format: string; presence: 'M' | 'C' | 'O' }
> = {
    '00': { id: '00', name: 'payloadFormatIndicator', format: 'N2', presence: 'M' },
    '01': { id: '01', name: 'pointOfInitiationMethod', format: 'N2', presence: 'O' },

    // 02–25: scheme/network MAI (primitive or templates depending on scheme)
    '02': { id: '02', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '03': { id: '03', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '04': { id: '04', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '05': { id: '05', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '06': { id: '06', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '07': { id: '07', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '08': { id: '08', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '09': { id: '09', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '10': { id: '10', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '11': { id: '11', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '12': { id: '12', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '13': { id: '13', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '14': { id: '14', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '15': { id: '15', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '16': { id: '16', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '17': { id: '17', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '18': { id: '18', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '19': { id: '19', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '20': { id: '20', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '21': { id: '21', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '22': { id: '22', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '23': { id: '23', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '24': { id: '24', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },
    '25': { id: '25', name: 'merchantAccountInformation', format: 'ans', presence: 'C' },

    // 26–51: Domestic/Additional MAI templates
    '26': { id: '26', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '27': { id: '27', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '28': { id: '28', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '29': { id: '29', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '30': { id: '30', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '31': { id: '31', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '32': { id: '32', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '33': { id: '33', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '34': { id: '34', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '35': { id: '35', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '36': { id: '36', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '37': { id: '37', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '38': { id: '38', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '39': { id: '39', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '40': { id: '40', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '41': { id: '41', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '42': { id: '42', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '43': { id: '43', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '44': { id: '44', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '45': { id: '45', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '46': { id: '46', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '47': { id: '47', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '48': { id: '48', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '49': { id: '49', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '50': { id: '50', name: 'merchantAccountInformationTemplate', format: 'S', presence: 'C' },
    '51': { id: '51', name: 'merchantAccountInformationTemplateCentralRepository', format: 'S', presence: 'C' },

    '52': { id: '52', name: 'merchantCategoryCode', format: 'N4', presence: 'M' },
    '53': { id: '53', name: 'transactionCurrency', format: 'N3', presence: 'M' },
    '54': { id: '54', name: 'transactionAmount', format: 'ans', presence: 'C' },
    '55': { id: '55', name: 'tipOrConvenienceIndicator', format: 'N2', presence: 'O' },
    '56': { id: '56', name: 'valueOfConvenienceFeeFixed', format: 'ans', presence: 'C' },
    '57': { id: '57', name: 'valueOfConvenienceFeePercentage', format: 'ans', presence: 'C' },
    '58': { id: '58', name: 'countryCode', format: 'ans2', presence: 'M' },
    '59': { id: '59', name: 'merchantName', format: 'ans', presence: 'M' },
    '60': { id: '60', name: 'merchantCity', format: 'ans', presence: 'M' },
    '61': { id: '61', name: 'postalCode', format: 'ans', presence: 'O' },
    '62': { id: '62', name: 'additionalDataFieldTemplate', format: 'S', presence: 'O' },
    '63': { id: '63', name: 'crc', format: 'ans4', presence: 'M' },
    '64': { id: '64', name: 'merchantInformationLanguageTemplate', format: 'S', presence: 'O' },
    '65': { id: '65', name: 'rfuForEMVCo', format: 'S', presence: 'O' },
    '80': { id: '80', name: 'unreservedTemplates', format: 'S', presence: 'O' },
};

/**
 * Main entry: Parse a root TLV string and flatten to keys like:
 * data00, data59, data26_00, data62_99_01, etc., with absolute value positions.
 */
export function parseTagWithPos(emvData: string): { tags: TagMap; positions: TagPositions } {
    const tags: TagMap = {};
    const positions: TagPositions = {};
    const tlvs = parseTLVString(emvData, 0);

    for (const tlv of tlvs) {
        const { id, value } = tlv;

        switch (id) {
            case '26': parse26(value, tags, positions, tlv.valueStart); break;
            case '27': parse27(value, tags, positions, tlv.valueStart); break;
            case '28': parse28(value, tags, positions, tlv.valueStart); break;
            case '40': parse40(value, tags, positions, tlv.valueStart); break;
            case '51': parse51(value, tags, positions, tlv.valueStart); break;
            case '62': parse62(value, tags, positions, tlv.valueStart); break;
            default: {
                const key = `data${id}`;
                setKV(tags, positions, key, value, tlv.valueStart, tlv.valueEnd);
                if (emvTags[id]) {
                    const mappedId = emvTags[id].id;
                    setKV(tags, positions, `data${mappedId}`, value, tlv.valueStart, tlv.valueEnd);
                }
            }
        }
    }
    return { tags, positions };
}

// Back-compat API (no positions)
export function parseTag(emvData: string): TagMap {
    return parseTagWithPos(emvData).tags;
}

/** Convenience: safe parse that returns {} on error */
export function tryParseTag(emvData: string): TagMap {
    try { return parseTag(emvData); } catch { return {}; }
}
