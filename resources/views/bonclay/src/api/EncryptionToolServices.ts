import { BonClayHttpClient } from '@/api/bonClayHttpClient'
import type { EncryptionToolReq, EncryptionToolResp } from '@/types/api'

class EncryptionToolService extends BonClayHttpClient {
    private readonly resource = 'tools/config-encryption'

    async actEncryptDecrypt(
        phaseType: 'encrypt' | 'decrypt',
        req: EncryptionToolReq
    ): Promise<EncryptionToolResp> {
        return this.post<EncryptionToolResp>(`/${this.resource}/${phaseType}`, req)
    }
}

const encryptionToolService = new EncryptionToolService()
export default encryptionToolService
