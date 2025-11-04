// src/api/diagramCollectionAPI.ts
import type { PaginatedResponse, PaginationParams, RequestOptions} from "@/lib/httpClient/types";
import {BaseHttpClient} from "@/lib/httpClient/http.ts";
import {BaseCRUDServices} from "@/api/BaseCRUDServices.ts";
import type {DiagramCollection, DiagramDetail, DiagramSummary} from "@/types/entities.ts"; // your base from the snippet


class DiagramCollectionServices extends BaseCRUDServices<DiagramCollection> {
    constructor() {
        super("diagram_collections");
    }

    async getDiagram(
        collectionID:string
    ): Promise<PaginatedResponse<DiagramSummary>> {
        return this.get<PaginatedResponse<DiagramSummary>>(`/${this.resource}/${collectionID}/diagrams`, )
    }
}


const diagramCollectionAPI = new DiagramCollectionServices();
export default diagramCollectionAPI;
