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