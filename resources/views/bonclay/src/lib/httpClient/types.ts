// ================================
// files: types.ts or httpClient/types.ts
// ================================

// Base HTTP Client Types

export interface RequestOptions {
    headers?: Record<string, string>
    params?: Record<string, any>
    timeout?: number
    signal?: AbortSignal
    [key: string]: any
}

// Notification Types
export type NotificationType = 'success' | 'error' | 'warning' | 'info'

export interface ApiNotification {
    type: NotificationType
    title: string
    message?: string
    duration?: number
}

export interface NotificationConfig {
    showSuccess?: boolean
    showError?: boolean
    successMessage?: string
    errorMessage?: string
}

// Pagination Types
export interface PaginationParams {
    page_number?: number
    page_size?: number
    keyword?: string
    sort_by?: string
    order?: 'asc' | 'desc'
}

export interface PaginatedResponse<T> {
    list: T[]
    page_number: number
    page_size: number
    total_items: number
    total_pages: number
    next_page: number | null
    previous_page: number | null
    keyword?: string | null
}

// API Response Types
export interface ApiResponse<T = any> {
    data?: T
    message?: string
    notification?: ApiNotification
    status?: number
    success?: boolean
}

export interface ApiErrorResponse {
    data?: Record<string, string>
    message?: string
    status?: number
}

// Extended Request Options with Notification
export interface ExtendedRequestOptions extends RequestOptions {
    notification?: NotificationConfig
}

// HTTP Method Types
export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

// HTTP Client Configuration
export interface HttpClientConfig {
    baseURL?: string
    timeout?: number
    headers?: Record<string, string>
    withCredentials?: boolean
}

// Query Builder Types (if you want to add advanced querying)
export interface QueryFilter {
    field: string
    operator: 'eq' | 'ne' | 'gt' | 'gte' | 'lt' | 'lte' | 'in' | 'like'
    value: any
}

export interface QueryOptions extends PaginationParams {
    filters?: QueryFilter[]
    include?: string[] // For relationship inclusion
    fields?: string[] // For field selection
}

// Batch Operations Types
export interface BatchRequest<T> {
    method: HttpMethod
    url: string
    data?: Partial<T>
}

export interface BatchResponse<T> {
    success: boolean
    data?: T
    error?: string
}

// File Upload Types
export interface FileUploadOptions extends RequestOptions {
    onProgress?: (progress: number) => void
    maxSize?: number // in bytes
    allowedTypes?: string[]
}

export interface UploadedFile {
    id: string
    name: string
    size: number
    type: string
    url: string
    created_at: string
}

// Webhook/Event Types (for real-time notifications)
export interface WebhookEvent<T = any> {
    event: string
    data: T
    timestamp: string
    id: string
}

// Cache Types
export interface CacheOptions {
    ttl?: number // time to live in seconds
    key?: string
    enabled?: boolean
}

// Retry Types
export interface RetryOptions {
    maxRetries?: number
    retryDelay?: number
    retryOn?: number[] // HTTP status codes to retry on
}

// Complete Request Options with all features
export interface CompleteRequestOptions extends RequestOptions {
    notification?: NotificationConfig
    cache?: CacheOptions
    retry?: RetryOptions
}

// ================================
// Type Guards
// ================================

export function isApiResponse<T>(response: any): response is ApiResponse<T> {
    return response && typeof response === 'object' && 'data' in response
}

export function isApiErrorResponse(error: any): error is ApiErrorResponse {
    return error && typeof error === 'object' && ('error' in error || 'message' in error)
}

export function isPaginatedResponse<T>(response: any): response is PaginatedResponse<T> {
    return (
        response &&
        typeof response === 'object' &&
        'data' in response &&
        Array.isArray(response.data) &&
        'page_number' in response &&
        'total_items' in response
    )
}

// ================================
// Utility Types
// ================================

// Make all properties of T nullable
export type Nullable<T> = {
    [P in keyof T]: T[P] | null
}

// Make specific properties optional
export type Optional<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>

// Make specific properties required
export type Required<T, K extends keyof T> = T & {
    [P in K]-?: T[P]
}

// Extract data type from PaginatedResponse
export type ExtractData<T> = T extends PaginatedResponse<infer U> ? U : never

// ================================
// Domain Model Base Types
// ================================

export interface BaseEntity {
    id: string | number
    created_at?: string
    updated_at?: string
    deleted_at?: string | null
}

export interface TimestampedEntity extends BaseEntity {
    created_at: string
    updated_at: string
}

export interface SoftDeletableEntity extends TimestampedEntity {
    deleted_at: string | null
}

// ================================
// Example Domain Types
// ================================

export interface User extends TimestampedEntity {
    name: string
    email: string
    avatar?: string
    role: 'admin' | 'user' | 'guest'
    is_active: boolean
}

export interface Post extends TimestampedEntity {
    title: string
    content: string
    author_id: number
    status: 'draft' | 'published' | 'archived'
    tags?: string[]
}

export interface Comment extends TimestampedEntity {
    post_id: number
    user_id: number
    content: string
    parent_id?: number | null
}

export interface ApiError {
    error?: string
    message?: string
    code?: string
    status?: number
    data?: any
}



// ================================
// Export all types
// ================================

export type {
    // Add any additional exports here
}

// ================================
// Constants
// ================================

export const HTTP_STATUS = {
    OK: 200,
    CREATED: 201,
    NO_CONTENT: 204,
    BAD_REQUEST: 400,
    UNAUTHORIZED: 401,
    FORBIDDEN: 403,
    NOT_FOUND: 404,
    UNPROCESSABLE_ENTITY: 422,
    INTERNAL_SERVER_ERROR: 500,
} as const

export const DEFAULT_PAGE_SIZE = 10
export const DEFAULT_TIMEOUT = 30000
export const MAX_RETRIES = 3
export const RETRY_DELAY = 1000

