import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
    // {
    //     path: '/',
    //     name: 'home',
    //     component: () => import('@/App.vue'),
    //     meta: {
    //         title: 'Home',
    //         icon: 'Home'
    //     }
    // },
    {
        path: '/test',
        name: 'tesst',
        component: () => import('@/pages/common/Dashboard.vue'),
        meta: {
            title: 'tesst',
            icon: 'tesst'
        }
    },
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes
})

export default router