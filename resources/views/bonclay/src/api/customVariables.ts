import client from '@/api/client'
import type { CustomVariable } from '@/types/api.types'

export interface CreateVariablePayload {
  key: string
  value: string
}

export const listCollectionVariables = (collectionSlug: string): Promise<CustomVariable[]> =>
  client.get<CustomVariable[]>(`/collections/${collectionSlug}/variables`)

export const createCollectionVariable = (
  collectionId: string,
  collectionSlug: string,
  payload: CreateVariablePayload,
): Promise<CustomVariable> =>
  client.post<CustomVariable>(
    `/collections/${collectionSlug}/variables?collection_id=${collectionId}`,
    payload,
  )

export const updateCollectionVariable = (
  id: string,
  payload: CreateVariablePayload,
): Promise<CustomVariable> =>
  client.patch<CustomVariable>(`/collections/variables/${id}`, payload)

export const deleteCollectionVariable = (id: string): Promise<void> =>
  client.delete<void>(`/collections/variables/${id}`)
