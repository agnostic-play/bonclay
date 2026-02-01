export interface Parameter {
    name: string
    type: string
    required: boolean
    description: string
}

export interface Response {
    code: string
    description: string
    example?: string
}

export interface Endpoint {
    method: string
    path: string
    name: string
    description: string
    parameters: Parameter[]
    responses: Response[]
}

export interface ApiCategory {
    id: string
    name: string
    endpoints: Endpoint[]
}

export interface Scenario {
    id: string
    endpointPath: string
    httpStatus: number
    name: string
    description: string
    responseDelay: number
    isActive: boolean
    createdAt: string
    response?: any
}

export interface ScenarioResponse {
    id: string
    endpoint_id: string
    desc: string
    header: string          // JSON string
    body: string            // JSON string
    status_header: number
    delay: number
}

export interface CollectionEndpoint {
    id: string
    path: string
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | string
    collection_id: string
    category: string
    active_scenario: string
    desc: string
    script: string
    scenario_response: ScenarioResponse[]
}

export type EndpointByCategory = Record<string, CollectionEndpoint[]>