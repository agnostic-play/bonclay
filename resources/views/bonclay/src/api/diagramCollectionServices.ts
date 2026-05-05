import { BaseCRUDServices, type PaginatedResponse } from '@/api/BaseCRUDServices'
import type { DiagramCollection, DiagramSummary } from '@/types/entities'

class DiagramCollectionServices extends BaseCRUDServices<DiagramCollection> {
    constructor() {
        super('diagram_collections')
    }

    async getDiagram(collectionID: string, params?: { search?: string }): Promise<PaginatedResponse<DiagramSummary>> {
        return this.get<PaginatedResponse<DiagramSummary>>(`/${this.resource}/${collectionID}/diagrams`, { params })
    }
}

const diagramCollectionAPI = new DiagramCollectionServices()
export default diagramCollectionAPI
