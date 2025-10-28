import {BonclayHttpClient} from "@/api/bonclayHttpClient.ts";
import type {MermaidDiagramTools} from "@/types/api.types.ts";


export class MermaidDiagramAPI extends BonclayHttpClient<MermaidDiagramTools> {
    constructor() {
        super('users');
    }
}