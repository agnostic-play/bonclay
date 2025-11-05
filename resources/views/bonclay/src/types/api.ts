import {BaseHttpClient} from "@/lib/httpClient/http.ts";


export interface EncryptionToolReq  {
    encryptionType: string;
    type: string;
    key: string;
    outputFormat: string;
    keyDeriveMethod: string;
    rawValue: string;
    iv: string;
}

export interface EncryptionToolResp  {
    output: string;
}