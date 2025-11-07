import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import MermaidDiagramToolsIndex from "@/pages/mermaidDiagram/Index.vue";
import MermaidDiagramToolsShow from "@/pages/mermaidDiagram/Show.vue";
import EncryptionToolsIndex from "@/pages/encryptionTools/Index.vue";
import EmvcoParser from "@/pages/parserTools/Emvco.vue";
import IsoParser from "@/pages/parserTools/IsoParser.vue";
import MockAPIToolsIndex from "@/pages/mockApi/Index.vue";
import MockAPIToolsIndexShowCollection from "@/pages/mockApi/ShowCollection.vue";

const routes: RouteRecordRaw[] = [
    {
        path: '/tools',
        name: 'tools',
        children: [
            {
                path: 'mock-api',
                name: 'MockApiTools',
                children: [
                    {
                        path: '',
                        name: 'MockApiTools-Index',
                        component: MockAPIToolsIndex,
                    },
                    {
                        path: 'collection/:id',
                        name: 'MockApiTools-CollectionShow',
                        component: MockAPIToolsIndexShowCollection,
                    }
                ]
            },
            {
                path: 'encryption-tools',
                name: 'EncryptionTools',
                children: [
                    {
                        path: '',
                        name: 'EncryptionTools-Index',
                        component: EncryptionToolsIndex,
                    }
                ]
            },
            {
                path: 'mermaid-diagram',
                name: 'MermaidDiagramTools',
                children: [
                    {
                        path: '',
                        name: 'MermaidDiagramTools-Index',
                        component: MermaidDiagramToolsIndex,
                    },
                    {
                        path: 'collection/:id',
                        name: 'MermaidDiagramTools-ProjectShow',
                        component: MermaidDiagramToolsShow,
                    }
                ]
            },
            {
                path: 'parser',
                name: 'ParserTools',
                children: [
                    {
                        path: 'emvco',
                        name: 'ParserTools-Emvco',
                        component: EmvcoParser,
                    },
                    {
                        path: 'iso',
                        name: 'ParserTools-ISO',
                        component: IsoParser,
                    },
                ]
            },
            // You can add more children here
        ],
    }
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes
})

export default router