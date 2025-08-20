import axios, { AxiosError } from 'axios';
import type {
    AxiosInstance,
    AxiosRequestConfig,
    AxiosResponse,
    InternalAxiosRequestConfig,
    AxiosRequestHeaders,
} from 'axios';
import type { ApiError, ApiResponse, RequestOptions } from './types';

export class BaseHttpClient {
    protected readonly instance: AxiosInstance;
    protected readonly baseURL: string;

    constructor(baseURL: string, config: AxiosRequestConfig = {}) {
        this.baseURL = baseURL;

        const defaultConfig: AxiosRequestConfig = {
            baseURL,
            timeout: 10_000,
            headers: { 'Content-Type': 'application/json' },
            ...config,
        };

        this.instance = axios.create(defaultConfig);
        this.setupInterceptors();
    }

    private setupInterceptors(): void {
        // Request interceptor
        this.instance.interceptors.request.use(
            (config: InternalAxiosRequestConfig) => {
                const headers = (config.headers ?? {}) as AxiosRequestHeaders;

                // Add auth token if available
                const token = this.getAuthToken();
                if (token) {
                    headers.Authorization = `Bearer ${token}`;
                }

                // Add request ID for tracing
                headers['X-Request-ID'] = this.generateRequestId();

                // Add timestamp
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

        // Response interceptor
        this.instance.interceptors.response.use(
            (response: AxiosResponse) => {
                console.log(`✅ ${response.config.method?.toUpperCase()} ${response.config.baseURL ?? ''}${response.config.url}`, {
                    status: response.status,
                    data: response.data,
                });
                // Keep the axios wrapper so downstream methods can standardize unwrapping
                return response;
            },
            async (error: AxiosError) => {
                console.error(`❌ ${error.config?.method?.toUpperCase()} ${error.config?.url}`, {
                    status: error.response?.status,
                    data: error.response?.data,
                });

                // Handle specific error cases
                if (error.response?.status === 401) {
                    await this.handleUnauthorized();
                } else if (error.response?.status === 403) {
                    this.handleForbidden();
                } else if (error.response?.status === undefined || error.response?.status >= 500) {
                    this.handleServerError(error);
                }

                return Promise.reject(this.normalizeError(error));
            }
        );
    }

    // Protected methods for subclasses to override
    protected getAuthToken(): string | null {
        // Override in subclasses to provide token logic. Safe for SSR.
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
        // Example: router.push('/login')
    }

    protected handleForbidden(): void {
        console.warn('Forbidden request');
    }

    protected handleServerError(error: AxiosError): void {
        console.error('Server error:', error);
    }

    private generateRequestId(): string {
        // Prefer crypto.randomUUID if available
        try {
            // @ts-ignore - node/web crypto
            if (typeof crypto !== 'undefined' && crypto.randomUUID) {
                // @ts-ignore
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
            details: axErr?.response?.data,
        };
        return apiError;
    }

    // ---- Helpers -------------------------------------------------------------

    private isEmptyBody(body: unknown): boolean {
        if (body == null) return true; // null or undefined

        // String payload
        if (typeof body === 'string') return body.trim().length === 0;

        // FormData payload
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore: FormData may not exist in all runtimes
        if (typeof FormData !== 'undefined' && body instanceof FormData) {
            // FormData doesn't expose size directly; iterate once.
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            for (const _ of body.entries()) {
                return false;
            }
            return true;
        }

        // ArrayBuffer / Blob considered non-empty if byteLength/size > 0
        if (typeof ArrayBuffer !== 'undefined' && body instanceof ArrayBuffer) {
            return body.byteLength === 0;
        }
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        if (typeof Blob !== 'undefined' && body instanceof Blob) {
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            return body.size === 0;
        }

        // Plain object or array
        if (typeof body === 'object') {
            try {
                if (Array.isArray(body)) return body.length === 0;
                return Object.keys(body as Record<string, unknown>).length === 0;
            } catch {
                return false;
            }
        }

        // Numbers/booleans/etc. are considered "present"
        return false;
    }

    /** Unwraps various response shapes into T:
     * - AxiosResponse<ApiResponse<T>>
     * - AxiosResponse<T>
     * - ApiResponse<T>
     * - T
     */
    private unwrap<T>(response: unknown): T {
        // AxiosResponse
        if (response && typeof response === 'object' && 'data' in (response as any)) {
            const r = response as AxiosResponse<ApiResponse<T> | T>;
            const payload = r.data as any;

            // ApiResponse<T> shape with data field
            if (payload && typeof payload === 'object' && 'data' in payload) {
                return (payload as ApiResponse<T>).data as T;
            }
            return payload as T;
        }

        // ApiResponse<T> or T passed directly
        const payload = response as any;
        if (payload && typeof payload === 'object' && 'data' in payload) {
            return (payload as ApiResponse<T>).data as T;
        }
        return payload as T;
    }

    // ---- Core HTTP methods with proper typing & error handling ---------------

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
                const err: ApiError = {
                    message: 'Request body is required and cannot be empty.',
                    status: 0,
                };
                throw err;
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
                const err: ApiError = {
                    message: 'Request body is required and cannot be empty.',
                    status: 0,
                };
                throw err;
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
                const err: ApiError = {
                    message: 'Request body is required and cannot be empty.',
                    status: 0,
                };
                throw err;
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

    // Utility methods
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
