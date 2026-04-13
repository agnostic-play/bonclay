/**
 * History API
 * Request/response history stored in Redis per collection.
 * Endpoint base: /api/v2/collections/:slug/history
 */
import client from '@/api/client'

// ── Response shapes ───────────────────────────────────────────────────────────

export interface RequestHistory {
  id: string
  collection_slug: string
  method: string
  path: string
  request_headers: Record<string, string>
  request_body: string
  status_code: number
  response_headers?: Record<string, string>
  response_body: string
  duration_ms?: number
  proxy_used: boolean
  scenario_id?: string
  timestamp: string
}

// ── API functions ─────────────────────────────────────────────────────────────

/**
 * GET /api/v2/collections/:slug/history
 * Returns the last 100 request history entries for a collection.
 */
export const getHistory = (slug: string): Promise<RequestHistory[]> =>
  client.get<RequestHistory[]>(`/collections/${slug}/history`)

/**
 * DELETE /api/v2/collections/:slug/history
 * Clears all history for a collection.
 */
export const clearHistory = (slug: string): Promise<void> =>
  client.delete<void>(`/collections/${slug}/history`)
