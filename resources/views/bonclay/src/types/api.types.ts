export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

export interface CustomVariable {
  id: string
  collection_id: string
  key: string
  value: string
}

export interface ScenarioResponse {
  id: string
  endpoint_id: string
  desc: string
  header: string   // JSON string
  body: string     // JSON string
  status_header: number
  delay: number
}

export interface CollectionEndpoint {
  id: string
  path: string
  method: HttpMethod
  collection_id: string
  category: string
  active_scenario: string
  desc: string
  script: string
  delay?: number
  scenario_response: ScenarioResponse[]
}

export type EndpointByCategory = Record<string, CollectionEndpoint[]>

export interface CollectionDetail {
  id: string
  name: string
  slug: string
  desc: string
  docs: string
  forward_proxy_url: string
  is_proxy_enable: boolean
  endpoints: EndpointByCategory
}
