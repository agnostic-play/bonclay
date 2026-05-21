import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import Index from "@/pages/mockApi/Index.vue";
import ShowCollection from "@/pages/mockApi/ShowCollection.vue";
import History from "@/pages/mockApi/History.vue";
import MermaidDiagramToolsIndex from "@/pages/mermaidDiagram/Index.vue";
import MermaidDiagramToolsShow from "@/pages/mermaidDiagram/Show.vue";
import KongIndex from "@/pages/kong/Index.vue";

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
                path: 'kong',
                name: 'KongTools',
                children: [
                    {
                        path: 'services',
                        name: 'KongTools-Services',
                        component: KongIndex,
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