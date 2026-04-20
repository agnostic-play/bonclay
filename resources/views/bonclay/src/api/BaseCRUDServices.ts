import { BonClayHttpClient } from '@/api/bonClayHttpClient'
import type { RequestOptions } from '@/lib/httpClient/types'

export interface PaginatedResponse<T> {
    list: T[]
    total: number
    page: number
    page_size: number
    next_page: number | null
    previous_page: number | null
    keyword?: string
}

export interface PaginationParams {
    page_number?: number
    page_size?: number
    keyword?: string
}

export abstract class BaseCRUDServices<T> extends BonClayHttpClient {
    protected readonly resource: string

    constructor(resource: string) {
        super()
        this.resource = resource
    }

    async getAll(params?: PaginationParams, options?: RequestOptions): Promise<PaginatedResponse<T>> {
        return this.get<PaginatedResponse<T>>(`/${this.resource}/all`, { ...options, params })
    }

    async getList(params?: PaginationParams, options?: RequestOptions): Promise<PaginatedResponse<T>> {
        return this.get<PaginatedResponse<T>>(`/${this.resource}/list`, { ...options, params })
    }

    async getDetail(id: string | number, options?: RequestOptions): Promise<T> {
        return this.get<T>(`/${this.resource}/${id}/show`, options)
    }

    async actCreate(data: Partial<T>, options?: RequestOptions): Promise<T> {
        return this.post<T>(`/${this.resource}/create`, data, options)
    }

    async actUpdate(id: string | number, data: Partial<T>, options?: RequestOptions): Promise<T> {
        return this.patch<T>(`/${this.resource}/${id}/update`, data, options)
    }

    async actRemove(id: string | number, options?: RequestOptions): Promise<void> {
        await this.delete(`/${this.resource}/${id}/remove`, options)
    }
}
