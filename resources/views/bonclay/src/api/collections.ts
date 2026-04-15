/**
 * Collection API
 * All collection-related HTTP calls live here.
 * Endpoint base: /api/v2/collections
 */
import client from '@/api/client'
import type { CollectionDetail } from '@/types/api.types'

// ── Response shapes ───────────────────────────────────────────────────────────

export interface CollectionInfo {
  id: string
  name: string
  slug: string
  desc: string
  docs: string
  squad_id: string
  forward_proxy_url: string
  is_proxy_enable: boolean
}

// ── Request shapes ────────────────────────────────────────────────────────────

export interface CreateCollectionPayload {
  name: string
  desc?: string
  squad_id: string
}

export interface UpdateCollectionPayload {
  name?: string
  desc?: string
  docs?: string
  forward_proxy_url?: string
  is_proxy_enable?: boolean
}

// ── API functions ─────────────────────────────────────────────────────────────

/**
 * POST /api/v2/collections/create
 * Creates a new collection under the specified squad.
 */
export const createCollection = (payload: CreateCollectionPayload): Promise<void> =>
  client.post<void>('/collections/create', payload)

/**
 * GET /api/v2/collections/:slug/endpoint_scenario
 * Returns collection metadata and endpoints grouped by category with their scenarios.
 */
export const getCollectionDetail = (slug: string): Promise<CollectionDetail> =>
  client.get<CollectionDetail>(`/collections/${slug}/endpoint_scenario`)

/**
 * GET /api/v2/collections/:id/show
 * Returns the collection metadata (name, desc, slug, etc).
 */
export const getCollectionInfo = (id: string): Promise<CollectionInfo> =>
  client.get<CollectionInfo>(`/collections/${id}/show`)

/**
 * PATCH /api/v2/collections/:id/update
 * Updates collection metadata (name, desc, docs, forward_proxy_url).
 */
export const updateCollection = (id: string, payload: UpdateCollectionPayload): Promise<void> =>
  client.patch<void>(`/collections/${id}/update`, payload)
