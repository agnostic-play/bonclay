/**
 * Squad API
 * All squad-related HTTP calls live here.
 * Endpoint base: /squads
 */
import client from '@/api/client'

// ── Response shapes ───────────────────────────────────────────────────────────

export interface SquadListItem {
  id: string
  name: string
  slug: string
  desc?: string
}

/** Paginated list wrapper returned by GET /squads/all */
export interface SquadListResponse {
  list: SquadListItem[]
  total: number
  page: number
  page_size: number
  total_pages: number
  has_next: boolean
  has_prev: boolean
}

export interface SquadDetail {
  id: string
  name: string
  slug: string
  desc: string
  collections: Collection[]
}

export interface Collection {
  id: string
  name: string
  docs: string
  slug: string
  desc: string
}

// ── Request shapes ────────────────────────────────────────────────────────────

export interface CreateSquadPayload {
  name: string
  desc?: string
}

// ── API functions ─────────────────────────────────────────────────────────────

/**
 * GET /squads/all
 * Returns a paginated list of squads.
 */
export const getSquads = (): Promise<SquadListResponse> =>
  client.get<SquadListResponse>('/squads/all')

/**
 * GET /api/squad/detail/:slug
 * Returns full squad detail including collections.
 */
export const getSquadDetail = (slug: string): Promise<SquadDetail> =>
  client.get<SquadDetail>(`/squad/detail/${slug}`)

/**
 * POST /api/v2/squads/create
 * Creates a new squad.
 */
export const createSquad = (payload: CreateSquadPayload): Promise<void> =>
  client.post<void>('/squads/create', payload)
