import { BaseCRUDServices } from '@/api/BaseCRUDServices'
import type { DiagramSummary } from '@/types/entities'

class DiagramServices extends BaseCRUDServices<DiagramSummary> {
    constructor() {
        super('diagrams')
    }
}

const diagramServices = new DiagramServices()
export default diagramServices
