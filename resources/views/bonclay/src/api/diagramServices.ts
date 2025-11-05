// src/api/diagramCollectionAPI.ts
import type { PaginatedResponse, PaginationParams, RequestOptions} from "@/lib/httpClient/types";
import {BaseHttpClient} from "@/lib/httpClient/http.ts";
import {BaseCRUDServices} from "@/api/BaseCRUDServices.ts";
import type {DiagramCollection,  DiagramSummary} from "@/types/entities.ts"; // your base from the snippet


class DiagramServices extends BaseCRUDServices<DiagramSummary> {
    constructor() {
        super("diagrams");
    }
}


const diagramServices = new DiagramServices();
export default diagramServices;
