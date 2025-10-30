// src/api/diagramCollectionAPI.ts
import type { PaginatedResponse, PaginationParams, RequestOptions} from "@/lib/httpClient/types";
import {BaseHttpClient} from "@/lib/httpClient/http.ts";
import {BonClayHttpClient} from "@/api/bonClayHttpClient.ts";
import type {DiagramCollection, DiagramDetail, DiagramSummary} from "@/types/api.types.ts"; // your base from the snippet


class DiagramCollectionAPI extends BonClayHttpClient<DiagramCollection> {
    constructor() {
        super("diagram_collections");
    }


    async getDiagram(
        collectionID:string
    ): Promise<PaginatedResponse<DiagramSummary>> {
        return this.get<PaginatedResponse<DiagramSummary>>(`/${this.resource}/${collectionID}/diagrams`, )
    }


    // ----- Diagrams under a collection -----
    /**
     * GET /diagram_collections/:id/diagrams
     * Returns { list: DiagramSummary[], ...pagination }
     */
    // async getDiagrams(
    //     collectionId: string | number,
    //     options?: RequestOptions
    // ): Promise<ListResp<DiagramSummary>> {
    //     return this.get<ListResp<DiagramSummary>>(
    //         `/${this.resource}/${collectionId}/diagrams`,
    //         options
    //     );
    // }
    //
    // async getDiagrams(
    //     params?: PaginationParams,
    //     options?: RequestOptions,
    // ): Promise<PaginatedResponse<T>> {
    //     return this.get<PaginatedResponse<T>>(`/${this.resource}/list`, { ...options, params })
    // }
    //
    //
    // /**
    //  * GET /diagram_collections/:id/diagrams/:diagramId
    //  */
    // async getDiagram(
    //     collectionId: string | number,
    //     diagramId: string | number,
    //     options?: RequestOptions
    // ): Promise<DiagramDetail> {
    //     return this.get<DiagramDetail>(
    //         `/${this.resource}/${collectionId}/diagrams/${diagramId}`,
    //         options
    //     );
    // }
    //
    // /**
    //  * PATCH /diagram_collections/:id/diagrams/:diagramId
    //  * payload e.g. { syntax, title?, description? }
    //  */
    // async saveDiagram(
    //     collectionId: string | number,
    //     diagramId: string | number,
    //     payload: Partial<Pick<DiagramDetail, "syntax" | "title" | "description">>,
    //     options?: RequestOptions
    // ) {
    //     return this.patch<DiagramDetail>(
    //         `/${this.resource}/${collectionId}/diagrams/${diagramId}`,
    //         payload,
    //         options
    //     );
    // }
    //
    // /**
    //  * POST /diagram_collections/:id/diagrams
    //  * payload e.g. { title, description?, syntax? }
    //  */
    // async createDiagram(
    //     collectionId: string | number,
    //     payload: { title: string; description?: string; syntax?: string },
    //     options?: RequestOptions
    // ) {
    //     return this.post<DiagramDetail>(
    //         `/${this.resource}/${collectionId}/diagrams`,
    //         payload,
    //         options
    //     );
    // }
}


const diagramCollectionAPI = new DiagramCollectionAPI();
export default diagramCollectionAPI;
