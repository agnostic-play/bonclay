import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import Index from "@/pages/mockApi/Index.vue";
import ShowCollection from "@/pages/mockApi/ShowCollection.vue";
import History from "@/pages/mockApi/History.vue";
import MermaidDiagramToolsIndex from "@/pages/mermaidDiagram/Index.vue";
import MermaidDiagramToolsShow from "@/pages/mermaidDiagram/Show.vue";
import EncryptionToolsIndex from "@/pages/encryptionTools/index.vue";
import EmvcoParser from "@/pages/parserTools/Emvco.vue";
import IsoParser from "@/pages/parserTools/IsoParser.vue";

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        redirect: '/tools/mock-api',
    },
    {
        path: '/tools',
        name: 'tools',
        children: [
            {
                path: 'mock-api',
                name: 'MockApiTools',
                children: [
                    {
                        // Collection show as a sibling route (full page)
                        path: ':slug/collection/:collectionSlug',
                        name: 'MockApiTools-CollectionShow',
                        component: ShowCollection,
                    },
                    {
                        path: ':slug/collection/:collectionSlug/history',
                        name: 'MockApiTools-CollectionHistory',
                        component: History,
                    },
                    {
                        path: ':slug?',
                        name: 'MockApiTools-Index',
                        component: Index,
                    },
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