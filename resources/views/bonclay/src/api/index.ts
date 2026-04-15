/**
 * API barrel — import all API functions from here.
 *
 * Usage:
 *   import { getSquads, createSquad } from '@/api'
 *   import { createCollection, getCollectionDetail } from '@/api'
 *   import { createEndpoint } from '@/api'
 *   import { setActiveScenario, createScenarioResponse } from '@/api'
 */
export * from './squads'
export * from './collections'
export * from './endpoints'
export * from './scenarios'
export * from './history'
export { default as client } from './client'
