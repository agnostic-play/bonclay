import {createRouter, createWebHistory} from 'vue-router'
import type {RouteRecordRaw} from 'vue-router'

// Layout
import AppLayout from '@/components/Layout/AppLayout.vue'

// Pages
import Dashboard from '@/pages/Dashboard.vue'

// Auth pages (without layout)
import Login from '@/pages/auth/Login.vue'
import EncryptionTools from "@/pages/tools/EncryptionTools.vue";

const routes: RouteRecordRaw[] = [
    // Auth routes (no layout)
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: {requiresAuth: false}
    },
    // Main app routes (with layout)
    {
        path: '/',
        component: AppLayout,
        meta: {requiresAuth: true},
        children: [
            {
                path: '',
                name: 'Dashboard',
                component: Dashboard,
                meta: {title: 'Dashboard', menuItem: 'dashboard'}
            },
            {
                path: '/encryption-tools',
                name: 'EncryptionTools',
                component: EncryptionTools,
                meta: {title: 'Encryption Tools', menuItem: 'encryption-tools'}
            },
        ]
    },

    // Catch all route - redirect to dashboard
    {
        path: '/:pathMatch(.*)*',
        redirect: '/'
    }
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes
})

// Set page title
router.afterEach((to) => {
    document.title = to.meta.title ? `${to.meta.title} - Ditto` : 'Ditto'
})


export default router