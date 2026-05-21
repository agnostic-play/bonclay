/**
 * Kong API
 * Calls the Go backend, which proxies the Kong Admin API.
 * Endpoint base: /kong
 */
import client from '@/api/client'

export interface KongService {
  id: string
  name: string
  protocol: string
  host: string
  port: number
  path: string | null
  tags: string[] | null
  enabled: boolean
  retries: number
  connect_timeout: number
  read_timeout: number
  write_timeout: number
  client_certificate: string | null
  ca_certificates: string[] | null
  tls_verify: boolean | null
  tls_verify_depth: number | null
  created_at: number
  updated_at: number
}

export interface ListKongServicesResponse {
  data: KongService[]
  next: string | null
}

export interface ListKongServicesParams {
  size?: number
  offset?: string
  tag?: string
}

/**
 * GET /api/v2/kong/services
 * Returns Kong services proxied from the Kong Admin API.
 */
export const listKongServices = (
  params?: ListKongServicesParams,
): Promise<ListKongServicesResponse> =>
  client.get<ListKongServicesResponse>('/kong/services', { params })

export interface UpdateKongServicePayload {
  enabled?: boolean
}

/**
 * PATCH /api/v2/kong/services/:id
 * Updates a Kong service. Currently used for enable/disable toggle.
 */
export const updateKongService = (
  id: string,
  payload: UpdateKongServicePayload,
): Promise<KongService> =>
  client.patch<KongService>(`/kong/services/${id}`, payload)
