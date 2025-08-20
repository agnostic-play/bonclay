export interface ApiResponse<T = any> {
    data?: T;
    message?: string;
    status: number;
}

export interface ApiError {
    message: string;
    code?: string;
    status: number;
    details?: unknown;
}

export interface RequestOptions {
    timeout?: number;
    headers?: Record<string, string>;
    params?: Record<string, unknown>;
    retries?: number;
}
