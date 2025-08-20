import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import Index from "@/pages/mockApi/Index.vue";
import ShowCollection from "@/pages/mockApi/ShowCollection.vue";

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
                        component: Index,
                    },
                    {
                        path: 'collection/:id',
                        name: 'MockApiTools-CollectionShow',
                        component: ShowCollection,
                    }
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