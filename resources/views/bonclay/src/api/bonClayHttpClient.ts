import {BaseHttpClient} from "@/lib/httpClient/http.ts";


export abstract class BonclayHttpClient extends BaseHttpClient {
    protected readonly resource: string


    protected constructor(resource: string) {
        // http://api-mockup.alexandria.allobank.local
        // http://localhost:6106
        super("http://api-mockup.alexandria.allobank.local","api/v2", {
            timeout: 10_000,
            headers: { 'Content-Type': 'application/json' },
        })
        this.resource = resource.replace(/^\/+|\/+$/g, '')
    }

    getBaseURL(): string {
        return super.getHost()
    }
}
