/**
 * Endpoint API
 * All endpoint-related HTTP calls live here.
 * Endpoint base: /api/v2/endpoints
 */
import client from '@/api/client'

// ── Request shapes ────────────────────────────────────────────────────────────

export interface CreateEndpointPayload {
  collection_id: string
  path: string
  method: string
  category?: string
  desc?: string
}

// ── API functions ─────────────────────────────────────────────────────────────

export interface UpdateEndpointPayload {
  path?: string
  method?: string
  category?: string
  desc?: string
  delay?: number
}

/**
 * POST /api/v2/endpoints/create
 * Creates a new endpoint under a collection.
 */
export const createEndpoint = (payload: CreateEndpointPayload): Promise<void> =>
  client.post<void>('/endpoints/create', payload)

/**
 * PATCH /api/v2/endpoints/:id/update
 * Updates an endpoint's path, method, category, or desc.
 */
export const updateEndpoint = (id: string, payload: UpdateEndpointPayload): Promise<void> =>
  client.patch<void>(`/endpoints/${id}/update`, payload)

/**
 * PATCH /api/v2/endpoints/:id/script
 * Saves (or clears) the JavaScript snippet run on every mock request for this endpoint.
 * The script can read and mutate the collection's custom variables via the `env` object.
 */
export const updateEndpointScript = (id: string, script: string): Promise<void> =>
  client.patch<void>(`/endpoints/${id}/script`, { script })

/**
 * DELETE /api/v2/endpoints/:id/remove
 * Permanently deletes an endpoint and all its scenarios.
 */
export const deleteEndpoint = (id: string): Promise<void> =>
  client.delete<void>(`/endpoints/${id}/remove`)
