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

export interface MermaidDiagramTools {
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


export interface DiagramCollection {
    id: string | number;
    name: string;
    description?: string;
}

export interface DiagramSummary {
    id: string | number;
    title: string;
    syntax: string;
    description?: string;
    latestUpdated?: string;
}

export interface DiagramDetail extends DiagramSummary {
    id: string | number;
    title: string;
    syntax: string;
    description?: string;
    latestUpdated?: string; // ISO string if backend returns one
}

export interface DiagramCollection {
    id: string | number;
    name: string;
    description?: string;
}
