// src/lib/parseqr.ts
// Vue + Vite friendly TypeScript port of your Go QR EMV parser

export type TagMap = Record<string, string>;

type TLV = { id: string; len: number; value: string };

/** Read a single TLV at position i; throws on malformed input */
function readTLV(s: string, i: number): { tlv: TLV; next: number } {
    if (i + 4 > s.length) {
        throw new Error(`TLV header truncated at index ${i}`);
    }
    const id = s.slice(i, i + 2);
    const lenStr = s.slice(i + 2, i + 4);
    const len = Number.parseInt(lenStr, 10);
    if (Number.isNaN(len)) {
        throw new Error(`Invalid length '${lenStr}' at index ${i + 2}`);
    }
    const start = i + 4;
    const end = start + len;
    if (end > s.length) {
        throw new Error(`TLV value out of bounds: need ${len} chars from ${start}, have ${s.length - start}`);
    }
    const value = s.slice(start, end);
    return { tlv: { id, len, value }, next: end };
}

/** Generic TLV loop */
function parseTLVString(s: string): TLV[] {
    const out: TLV[] = [];
    let i = 0;
    while (i < s.length) {
        const { tlv, next } = readTLV(s, i);
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
        if (dbl) {
            n *= 2;
            if (n > 9) n -= 9;
        }
        sum += n;
        dbl = !dbl;
    }
    return (10 - (sum % 10)) % 10;
}

/** Sub-parsers (26/27/28/40/51/62.99) */

function parse26(value: string, out: TagMap) {
    for (const tlv of parseTLVString(value)) {
        switch (tlv.id) {
            case '00':
                out['data26_00'] = tlv.value; break
            case '01':
                out['data26_01'] = tlv.value + String(luhnCheckDigit(tlv.value)); break
            case '02':
                out['data26_02'] = tlv.value; break
            case '03':
                out['data26_03'] = tlv.value; break
            default:
                out[`data26_${tlv.id}`] = tlv.value
        }
    }
}
function parse27(value: string, out: TagMap) {
    for (const tlv of parseTLVString(value)) {
        switch (tlv.id) {
            case '00':
                out['data27_00'] = tlv.value; break
            case '01':
                out['data27_01'] = tlv.value + String(luhnCheckDigit(tlv.value)); break
            case '02':
                out['data27_02'] = tlv.value; break
            case '03':
                out['data27_03'] = tlv.value; break
            default:
                out[`data27_${tlv.id}`] = tlv.value
        }
    }
}
function parse28(value: string, out: TagMap) {
    for (const tlv of parseTLVString(value)) {
        switch (tlv.id) {
            case '00':
                out['data28_00'] = tlv.value; break
            case '01':
                out['data28_01'] = tlv.value + String(luhnCheckDigit(tlv.value)); break
            case '02':
                out['data28_02'] = tlv.value; break
            case '03':
                out['data28_03'] = tlv.value; break
            default:
                out[`data28_${tlv.id}`] = tlv.value
        }
    }
}
function parse40(value: string, out: TagMap) {
    for (const tlv of parseTLVString(value)) {
        switch (tlv.id) {
            case '00':
                out['data40_00'] = tlv.value; break
            case '01':
                out['data40_01'] = tlv.value; break
            case '02':
                out['data40_02'] = tlv.value; break
            default:
                out[`data40_${tlv.id}`] = tlv.value
        }
    }
}
function parse51(value: string, out: TagMap) {
    for (const tlv of parseTLVString(value)) {
        switch (tlv.id) {
            case '00':
                out['data51_00'] = tlv.value; break
            case '01':
                out['data51_01'] = tlv.value; break
            case '02':
                out['data51_02'] = tlv.value; break
            case '03':
                out['data51_03'] = tlv.value; break
            default:
                out[`data51_${tlv.id}`] = tlv.value
        }
    }
}
function parse6299(value: string, out: TagMap) {
    for (const tlv of parseTLVString(value)) {
        if (tlv.id === '01') out['data62_99_01'] = tlv.value
        else out[`data62_99_${tlv.id}`] = tlv.value
    }
}
function parse62(value: string, out: TagMap) {
    for (const tlv of parseTLVString(value)) {
        switch (tlv.id) {
            case '02':
                out['data62_02'] = tlv.value; break
            case '07':
                out['data62_07'] = tlv.value; break
            case '99':
                parse6299(tlv.value, out); break
            default:
                out[`data62_${tlv.id}`] = tlv.value
        }
    }
}

/** Optional: metadata map (mirrors your Go emvTags) */
export const emvTags: Record<
    string,
    { id: string; name: string; format: string; presence: 'M' | 'C' }
> = {
    '00': { id: '00', name: 'payloadFormatIndocator', format: 'ans', presence: 'C' },
    '01': { id: '01', name: 'pointOfInitiationMethod', format: 'ans', presence: 'C' },
    '02': { id: '02', name: 'merchantAccountInformationVisa', format: 'ans', presence: 'C' },
    '03': { id: '03', name: 'merchantAccountInformationVisa', format: 'ans', presence: 'C' },
    '04': { id: '04', name: 'merchantAccountInformationMasterCard', format: 'ans', presence: 'C' },
    '05': { id: '05', name: 'merchantAccountInformationMasterCard', format: 'ans', presence: 'C' },
    '06': { id: '06', name: 'merchantAccountInformationEMVCo', format: 'ans', presence: 'C' },
    '07': { id: '07', name: 'merchantAccountInformationEMVCo', format: 'ans', presence: 'C' },
    '08': { id: '08', name: 'merchantAccountInformationEMVCo', format: 'ans', presence: 'C' },
    '09': { id: '09', name: 'merchantAccountInformationDiscover', format: 'ans', presence: 'C' },
    '10': { id: '10', name: 'merchantAccountInformationDiscover', format: 'ans', presence: 'C' },
    '11': { id: '11', name: 'merchantAccountInformationAmex', format: 'ans', presence: 'C' },
    '12': { id: '12', name: 'merchantAccountInformationAmex', format: 'ans', presence: 'C' },
    '13': { id: '13', name: 'merchantAccountInformationJCB', format: 'ans', presence: 'C' },
    '14': { id: '14', name: 'merchantAccountInformationJCB', format: 'ans', presence: 'C' },
    '15': { id: '15', name: 'merchantAccountInformationUnionPay', format: 'ans', presence: 'C' },
    '16': { id: '16', name: 'merchantAccountInformationUnionPay', format: 'ans', presence: 'C' },
    '17': { id: '17', name: 'merchantAccountInformationEMVCo', format: 'ans', presence: 'C' },
    '18': { id: '18', name: 'merchantAccountInformationEMVCo', format: 'ans', presence: 'C' },
    '19': { id: '19', name: 'merchantAccountInformationEMVCo', format: 'ans', presence: 'C' },
    '20': { id: '20', name: 'merchantAccountInformationEMVCo', format: 'ans', presence: 'C' },
    '21': { id: '21', name: 'merchantAccountInformationEMVCo', format: 'ans', presence: 'C' },
    '22': { id: '22', name: 'merchantAccountInformationEMVCo', format: 'ans', presence: 'C' },
    '23': { id: '23', name: 'merchantAccountInformationEMVCo', format: 'ans', presence: 'C' },
    '24': { id: '24', name: 'merchantAccountInformationEMVCo', format: 'ans', presence: 'C' },
    '25': { id: '25', name: 'merchantAccountInformationEMVCo', format: 'ans', presence: 'C' },
    '26': { id: '26', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '27': { id: '27', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '28': { id: '28', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '29': { id: '29', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '30': { id: '30', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '31': { id: '31', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '32': { id: '32', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '33': { id: '33', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '34': { id: '34', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '35': { id: '35', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '36': { id: '36', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '37': { id: '37', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '38': { id: '38', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '39': { id: '39', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '40': { id: '40', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '41': { id: '41', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '42': { id: '42', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '43': { id: '43', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '44': { id: '44', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '45': { id: '45', name: 'merchantAccountInformationDomestic', format: 'ans', presence: 'C' },
    '46': { id: '46', name: 'merchantAccountInformationReservedDomesticId', format: 'ans', presence: 'C' },
    '47': { id: '47', name: 'merchantAccountInformationReservedDomesticId', format: 'ans', presence: 'C' },
    '48': { id: '48', name: 'merchantAccountInformationReservedDomesticId', format: 'ans', presence: 'C' },
    '49': { id: '49', name: 'merchantAccountInformationReservedDomesticId', format: 'ans', presence: 'C' },
    '50': { id: '50', name: 'merchantAccountInformationReservedDomesticId', format: 'ans', presence: 'C' },
    '51': { id: '51', name: 'merchantAccountInformationDomesticCentralRepository', format: 'ans', presence: 'C' },
    '52': { id: '52', name: 'merchantCategoryCode', format: 'ans', presence: 'M' },
    '53': { id: '53', name: 'transactionCurrency', format: 'ans', presence: 'M' },
    '54': { id: '54', name: 'transactionAmount', format: 'ans', presence: 'C' },
    '55': { id: '55', name: 'tipIndicator', format: 'ans', presence: 'C' },
    '56': { id: '56', name: 'tipValueOfFixed', format: 'ans', presence: 'C' },
    '57': { id: '57', name: 'tipValueOfPercentage', format: 'ans', presence: 'C' },
    '58': { id: '58', name: 'countryCode', format: 'ans', presence: 'C' },
    '59': { id: '59', name: 'merchant_name', format: 'ans', presence: 'M' },
    '60': { id: '60', name: 'merchant_city', format: 'ans', presence: 'M' },
    '61': { id: '61', name: 'postalCode', format: 'ans', presence: 'M' },
    '62': { id: '62', name: 'addionalData', format: 'ans', presence: 'C' },
    '63': { id: '63', name: 'crc', format: 'ans', presence: 'M' },
    '64': { id: '64', name: 'merchantInfoLanguage', format: 'ans', presence: 'M' },
    '65': { id: '65', name: 'rfu', format: 'ans', presence: 'M' },
    '80': { id: '80', name: 'unreservedTemplates', format: 'ans', presence: 'M' },
};

/**
 * Main entry: replicates ParseTag from Go.
 * Returns keys like data00, data59, data26_01, data62_99_01, etc.
 */
export function parseTag(emvData: string): TagMap {
    const tags: TagMap = {};
    const tlvs = parseTLVString(emvData);

    for (const tlv of tlvs) {
        const { id, value } = tlv;

        switch (id) {
            case '26': parse26(value, tags); break
            case '27': parse27(value, tags); break
            case '28': parse28(value, tags); break
            case '40': parse40(value, tags); break
            case '51': parse51(value, tags); break
            case '62': parse62(value, tags); break
            default:
                tags[`data${id}`] = value;
                if (emvTags[id]) {
                    const mappedId = emvTags[id].id;
                    tags[`data${mappedId}`] = value;
                }
        }
    }
    return tags;
}

/** Convenience: safe parse that returns {} on error (optional) */
export function tryParseTag(emvData: string): TagMap {
    try {
        return parseTag(emvData);
    } catch {
        return {};
    }
}
