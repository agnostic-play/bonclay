import { BaseHttpClient } from '@/lib/httpClient/http.ts'
import type { PaginatedResponse, PaginationParams, RequestOptions } from '@/lib/httpClient/types.ts'
import {BonclayHttpClient} from "@/api/BonclayHttpClient.ts";

export abstract class BaseCRUDServices<T> extends BonclayHttpClient {
    // --- CRUD (no notifications / no UI error handling) -----------------------
    async getAll(
        params?: PaginationParams,
        options?: RequestOptions,
    ): Promise<PaginatedResponse<T>> {
        return this.get<PaginatedResponse<T>>(`/${this.resource}/all`, { ...options, params })
    }

    async getList(
        params?: PaginationParams,
        options?: RequestOptions,
    ): Promise<PaginatedResponse<T>> {
        return this.get<PaginatedResponse<T>>(`/${this.resource}/list`, { ...options, params })
    }

    async getDetail(
        id: string | number,
        options?: RequestOptions,
    ): Promise<T> {
        return this.get<T>(`/${this.resource}/${id}/show`, options)
    }

    async actCreate(
        data: Partial<T>,
        options?: RequestOptions,
    ): Promise<T> {
        return this.post<T>(`/${this.resource}/create`, data, options)
    }

    async actUpdate(
        id: string | number,
        data: Partial<T>,
        options?: RequestOptions,
    ): Promise<T> {
        return this.patch<T>(`/${this.resource}/${id}/update`, data, options)
    }

    async actRemove(
        id: string | number,
        options?: RequestOptions,
    ): Promise<void> {
        await this.delete(`/${this.resource}/${id}/remove`, options)
    }

    hasNextPage(r: PaginatedResponse<T>): boolean {
        return r.next_page !== null
    }

    hasPreviousPage(r: PaginatedResponse<T>): boolean {
        return r.previous_page !== null
    }

    async getNextPage(
        current: PaginatedResponse<T>,
        options?: RequestOptions,
    ): Promise<PaginatedResponse<T> | null> {
        if (!this.hasNextPage(current)) return null
        return this.getList(
            {
                page_number: current.next_page!,
                page_size: current.page_size,
                keyword: current.keyword || undefined,
            },
            options,
        )
    }

    async getPreviousPage(
        current: PaginatedResponse<T>,
        options?: RequestOptions,
    ): Promise<PaginatedResponse<T> | null> {
        if (!this.hasPreviousPage(current)) return null
        return this.getList(
            {
                page_number: current.previous_page!,
                page_size: current.page_size,
                keyword: current.keyword || undefined,
            },
            options,
        )
    }
}