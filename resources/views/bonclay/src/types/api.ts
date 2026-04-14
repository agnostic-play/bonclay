export interface EncryptionToolReq  {
    encryptionType: string;
    type: string;
    key: string;
    outputFormat: string;
    inputFormat: string;
    keyDeriveMethod: string;
    rawValue: string;
    iv?: string;
}

export interface EncryptionToolResp  {
    output: string;
}