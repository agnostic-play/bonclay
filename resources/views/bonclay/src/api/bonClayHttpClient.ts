import { BaseHttpClient } from "@/lib/httpClient/http.ts";

export class BonClayHttpClient extends BaseHttpClient {
    constructor() {
        const base = (import.meta.env as any)?.VITE_API_BASE_URL ?? "";
        super(base);
    }

    public async get<T = unknown>(url: string, options: any = {}): Promise<T> {
        return super.get<T>(url, options);
    }

    public async post<T = unknown>(url: string, data?: unknown, options: any = {}): Promise<T> {
        return super.post<T>(url, data, options);
    }

    public async put<T = unknown>(url: string, data?: unknown, options: any = {}): Promise<T> {
        return super.put<T>(url, data, options);
    }

    public async patch<T = unknown>(url: string, data?: unknown, options: any = {}): Promise<T> {
        return super.patch<T>(url, data, options);
    }

    public async delete<T = unknown>(url: string, options: any = {}): Promise<T> {
        return super.delete<T>(url, options);
    }
}

const bonClayHttpClient = new BonClayHttpClient();
export default bonClayHttpClient;