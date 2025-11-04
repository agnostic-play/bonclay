import axios, { AxiosError } from 'axios';
import type {
    AxiosInstance,
    AxiosRequestConfig,
    AxiosResponse,
    InternalAxiosRequestConfig,
    AxiosRequestHeaders,
} from 'axios';
import type { ApiError, ApiErrorResponse, ApiResponse, RequestOptions } from './types';

export class BaseHttpClient {
    protected readonly instance: AxiosInstance;
    protected readonly host: string;
    protected readonly pathPrefix: string;
    protected readonly baseURL: string;

    constructor(host: string, pathPrefix: string, config: AxiosRequestConfig = {}) {
        this.host = host;
        this.pathPrefix = pathPrefix;
        this.baseURL = host + "/" + pathPrefix;

        const defaultConfig: AxiosRequestConfig = {
            baseURL: this.baseURL,
            timeout: 10_000,
            headers: { 'Content-Type': 'application/json' },
            ...config,
        };

        this.instance = axios.create(defaultConfig);
        this.setupInterceptors();
    }

    public getHost(): string {
        return this.host;
    }

    public getPathPrefix(): string {
        return this.pathPrefix;
    }

    private setupInterceptors(): void {
        // --- REQUEST INTERCEPTOR ---
        this.instance.interceptors.request.use(
            (config: InternalAxiosRequestConfig) => {
                const headers = (config.headers ?? {}) as AxiosRequestHeaders;

                const token = this.getAuthToken();
                if (token) headers.Authorization = `Bearer ${token}`;

                headers['X-Request-ID'] = this.generateRequestId();
                headers['X-Request-Time'] = new Date().toISOString();

                config.headers = headers;

                console.log(`🚀 ${config.method?.toUpperCase()} ${config.baseURL ?? ''}${config.url}`, {
                    params: config.params,
                    data: config.data,
                });

                return config;
            },
            (error: AxiosError) => {
                console.error('Request interceptor error:', error);
                return Promise.reject(this.normalizeError(error));
            }
        );

        // --- RESPONSE INTERCEPTOR ---
        this.instance.interceptors.response.use(
            (response: AxiosResponse) => {
                console.log(
                    `✅ ${response.config.method?.toUpperCase()} ${response.config.baseURL ?? ''}${response.config.url}`,
                    {
                        status: response.status,
                        data: response.data,
                    }
                );

                // Handle non-200 as structured errors
                if (response.status !== 200) {
                    const payload = response.data as ApiErrorResponse | any;
                    const apiError: ApiError = {
                        message: payload?.message || `Unexpected response status: ${response.status}`,
                        status: response.status,
                        data: payload?.data as Record<string, string> | undefined,
                    };
                    throw apiError;
                }

                return response;
            },
            async (error: AxiosError) => {
                console.error(
                    `❌ ${error.config?.method?.toUpperCase()} ${error.config?.url}`,
                    {
                        status: error.response?.status,
                        data: error.response?.data,
                    }
                );

                const status = error.response?.status;

                if (status === 401) {
                    await this.handleUnauthorized();
                } else if (status === 403) {
                    this.handleForbidden();
                } else if (status === undefined || status >= 500) {
                    this.handleServerError(error);
                }

                // --- normalize API error structure ---
                const payload = error.response?.data as ApiErrorResponse | any;
                if (payload && typeof payload === 'object') {
                    const structuredError: ApiError = {
                        message: payload?.message || `HTTP ${status ?? 'Error'}`,
                        status: status ?? 0,
                        data: payload?.data as Record<string, string> | undefined,
                    };
                    return Promise.reject(structuredError);
                }

                return Promise.reject(this.normalizeError(error));
            }
        );
    }

    // --- Auth / Error Helpers -------------------------------------------------

    protected getAuthToken(): string | null {
        try {
            if (typeof window !== 'undefined' && 'localStorage' in window) {
                return window.localStorage.getItem('auth_token');
            }
        } catch {
            // ignore
        }
        return null;
    }

    protected async handleUnauthorized(): Promise<void> {
        console.warn('Unauthorized request - consider redirecting to login');
    }

    protected handleForbidden(): void {
        console.warn('Forbidden request');
    }

    protected handleServerError(error: AxiosError): void {
        console.error('Server error:', error);
    }

    private generateRequestId(): string {
        try {
            if (typeof crypto !== 'undefined' && crypto.randomUUID) {
                return crypto.randomUUID();
            }
        } catch {
            // ignore
        }
        return `${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
    }

    private normalizeError(error: AxiosError | unknown): ApiError {
        const axErr = error as AxiosError;
        const apiError: ApiError = {
            message: axErr?.message || 'An unexpected error occurred',
            status: axErr?.response?.status ?? 0,
            code: axErr?.code,
            data: axErr?.response?.data,
        };
        return apiError;
    }

    // --- Internal Helpers ----------------------------------------------------

    private isEmptyBody(body: unknown): boolean {
        if (body == null) return true;
        if (typeof body === 'string') return body.trim().length === 0;

        if (typeof FormData !== 'undefined' && body instanceof FormData) {
            for (const _ of body.entries()) return false;
            return true;
        }

        if (typeof ArrayBuffer !== 'undefined' && body instanceof ArrayBuffer)
            return body.byteLength === 0;
        // @ts-ignore
        if (typeof Blob !== 'undefined' && body instanceof Blob)
            // @ts-ignore
            return body.size === 0;

        if (typeof body === 'object') {
            try {
                if (Array.isArray(body)) return body.length === 0;
                return Object.keys(body as Record<string, unknown>).length === 0;
            } catch {
                return false;
            }
        }

        return false;
    }

    private unwrap<T>(response: AxiosResponse): T {
        if (response && typeof response === 'object' && 'data' in (response as any)) {
            const r = response as AxiosResponse<ApiResponse<T> | T>;
            const payload = r.data as any;
            if (payload && typeof payload === 'object' && 'data' in payload) {
                return (payload as ApiResponse<T>).data as T;
            }
            return payload as T;
        }
        return (response as any) as T;
    }

    // --- Core HTTP methods ---------------------------------------------------

    protected async get<T = unknown>(
        url: string,
        options: RequestOptions = {}
    ): Promise<T> {
        try {
            const res = await this.instance.get<ApiResponse<T> | T>(url, {
                params: options.params,
                timeout: options.timeout,
                headers: options.headers,
            });
            return this.unwrap<T>(res);
        } catch (err) {
            const norm = this.normalizeError(err);
            console.error('Error in GET request:', norm);
            throw norm;
        }
    }

    protected async post<T = unknown>(
        url: string,
        data?: unknown,
        options: RequestOptions = {}
    ): Promise<T> {
        try {
            if (this.isEmptyBody(data)) {
                throw {
                    message: 'Request body is required and cannot be empty.',
                    status: 0,
                } as ApiError;
            }

            const res = await this.instance.post<ApiResponse<T> | T>(url, data, {
                timeout: options.timeout,
                headers: options.headers,
            });
            return this.unwrap<T>(res);
        } catch (err) {
            const norm = this.normalizeError(err);
            console.error('Error in POST request:', norm);
            throw norm;
        }
    }

    protected async put<T = unknown>(
        url: string,
        data?: unknown,
        options: RequestOptions = {}
    ): Promise<T> {
        try {
            if (this.isEmptyBody(data)) {
                throw {
                    message: 'Request body is required and cannot be empty.',
                    status: 0,
                } as ApiError;
            }

            const res = await this.instance.put<ApiResponse<T> | T>(url, data, {
                timeout: options.timeout,
                headers: options.headers,
            });
            return this.unwrap<T>(res);
        } catch (err) {
            const norm = this.normalizeError(err);
            console.error('Error in PUT request:', norm);
            throw norm;
        }
    }

    protected async patch<T = unknown>(
        url: string,
        data?: unknown,
        options: RequestOptions = {}
    ): Promise<T> {
        try {
            if (this.isEmptyBody(data)) {
                throw {
                    message: 'Request body is required and cannot be empty.',
                    status: 0,
                } as ApiError;
            }

            const res = await this.instance.patch<ApiResponse<T> | T>(url, data, {
                timeout: options.timeout,
                headers: options.headers,
            });
            return this.unwrap<T>(res);
        } catch (err) {
            const norm = this.normalizeError(err);
            console.error('Error in PATCH request:', norm);
            throw norm;
        }
    }

    protected async delete<T = unknown>(
        url: string,
        options: RequestOptions = {}
    ): Promise<T> {
        try {
            const res = await this.instance.delete<ApiResponse<T> | T>(url, {
                timeout: options.timeout,
                headers: options.headers,
            });
            return this.unwrap<T>(res);
        } catch (err) {
            const norm = this.normalizeError(err);
            console.error('Error in DELETE request:', norm);
            throw norm;
        }
    }

    // --- Utilities -----------------------------------------------------------

    public getBaseURL(): string {
        return this.baseURL;
    }

    public setDefaultHeader(key: string, value: string): void {
        this.instance.defaults.headers.common[key] = value;
    }

    public removeDefaultHeader(key: string): void {
        delete this.instance.defaults.headers.common[key];
    }
}
