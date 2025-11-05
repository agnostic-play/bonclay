// src/api/diagramCollectionAPI.ts
import type { PaginatedResponse, PaginationParams, RequestOptions} from "@/lib/httpClient/types";
import {BaseHttpClient} from "@/lib/httpClient/http.ts";
import {BaseCRUDServices} from "@/api/BaseCRUDServices.ts";
import type {DiagramCollection,  DiagramSummary} from "@/types/entities.ts";
import {BonclayHttpClient} from "@/api/BonclayHttpClient.ts";
import type {EncryptionToolReq, EncryptionToolResp} from "@/types/api.ts"; // your base from the snippet


class EncryptionToolService extends BonclayHttpClient {
    constructor() {
        super("tools/config-encryption");
    }

    async actEncryptDecrypt(
        phaseType: 'encrypt' | 'decrypt',
        req: EncryptionToolReq
    ): Promise<EncryptionToolResp> {
        return this.post<EncryptionToolResp>(`/${this.resource}/${phaseType}`, req);
    }
}

const encryptionToolService = new EncryptionToolService();
export default encryptionToolService;
