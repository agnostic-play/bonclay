/**
 * Scenarios / Scenario Response API
 * All scenario-related HTTP calls live here.
 * Endpoint base: /api/v2/scenario_response
 */
import client from '@/api/client'
import type { ScenarioResponse } from '@/types/api.types'

// ── Request shapes ────────────────────────────────────────────────────────────

export interface SetActiveScenarioPayload {
  endpoint_id: string
  scenario_id: string
}

export interface CreateScenarioPayload {
  endpoint_id: string
  desc: string
  status_header: number
  header?: string
  body?: string
  delay?: number
}

export interface UpdateScenarioPayload {
  endpoint_id: string
  desc?: string
  status_header?: number
  header?: string
  body?: string
  delay?: number
}

// ── API functions ─────────────────────────────────────────────────────────────

/**
 * POST /api/v2/scenario_response/set_active
 * Sets the active scenario for a given endpoint.
 */
export const setActiveScenario = (payload: SetActiveScenarioPayload): Promise<void> =>
  client.post<void>('/scenario_response/set_active', payload)

/**
 * POST /api/v2/scenario_response/:id/remove_active
 * Clears the active scenario for the given endpoint.
 */
export const removeActiveScenario = (endpointId: string): Promise<void> =>
  client.post<void>(`/scenario_response/${endpointId}/remove_active`, { endpoint_id: endpointId })

/**
 * POST /api/v2/scenario_response/create
 * Creates a new scenario response for an endpoint.
 */
export const createScenarioResponse = (payload: CreateScenarioPayload): Promise<ScenarioResponse> =>
  client.post<ScenarioResponse>('/scenario_response/create', payload)

/**
 * PATCH /api/v2/scenario_response/:id/update
 * Updates an existing scenario response. endpoint_id is required for backend validation.
 */
export const updateScenarioResponse = (id: string, payload: UpdateScenarioPayload): Promise<void> =>
  client.patch<void>(`/scenario_response/${id}/update`, payload)

/**
 * DELETE /api/v2/scenario_response/:id/remove
 * Permanently deletes a scenario response.
 */
export const deleteScenarioResponse = (id: string): Promise<void> =>
  client.delete<void>(`/scenario_response/${id}/remove`)
