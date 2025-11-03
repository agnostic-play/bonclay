// src/api/diagramCollectionAPI.ts
import type { PaginatedResponse, PaginationParams, RequestOptions} from "@/lib/httpClient/types";
import {BaseHttpClient} from "@/lib/httpClient/http.ts";
import {BonClayHttpClient} from "@/api/bonClayHttpClient.ts";
import type {DiagramCollection,  DiagramSummary} from "@/types/api.types.ts"; // your base from the snippet


class DiagramServices extends BonClayHttpClient<DiagramSummary> {
    constructor() {
        super("diagrams");
    }
}


const diagramServices = new DiagramServices();
export default diagramServices;
