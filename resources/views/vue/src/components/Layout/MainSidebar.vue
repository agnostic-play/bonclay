<template>
    <div class="main-layout">
        <!-- Main Sidebar Drawer -->
        <Sidebar 
            v-model:visible="visible" 
            position="left"
            :modal="false"
            :dismissableMask="false"
            :closeOnEscape="false"
            :showCloseIcon="false"
            class="main-sidebar"
            :style="{ width: sidebarWidth }"
        >
            <template>
                <div class="flex flex-column h-full">
                    <!-- Header with Logo -->
                    <div class="flex align-items-center justify-content-between px-4 pt-3 flex-shrink-0">
                        <span class="inline-flex align-items-center gap-2">
                            <svg width="35" height="40" viewBox="0 0 35 40" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <path
                                    d="M25.87 18.05L23.16 17.45L25.27 20.46V29.78L32.49 23.76V13.53L29.18 14.73L25.87 18.04V18.05ZM25.27 35.49L29.18 31.58V27.67L25.27 30.98V35.49ZM20.16 17.14H20.03H20.17H20.16ZM30.1 5.19L34.89 4.81L33.08 12.33L24.1 15.67L30.08 5.2L30.1 5.19ZM5.72 14.74L2.41 13.54V23.77L9.63 29.79V20.47L11.74 17.46L9.03 18.06L5.72 14.75V14.74ZM9.63 30.98L5.72 27.67V31.58L9.63 35.49V30.98ZM4.8 5.2L10.78 15.67L1.81 12.33L0 4.81L4.79 5.19L4.8 5.2ZM24.37 21.05V34.59L22.56 37.29L20.46 39.4H14.44L12.34 37.29L10.53 34.59V21.05L12.42 18.23L17.45 26.8L22.48 18.23L24.37 21.05ZM22.85 0L22.57 0.69L17.45 13.08L12.33 0.69L12.05 0H22.85Z"
                                    fill="var(--primary-color)"
                                />
                                <path
                                    d="M30.69 4.21L24.37 4.81L22.57 0.69L22.86 0H26.48L30.69 4.21ZM23.75 5.67L22.66 3.08L18.05 14.24V17.14H19.7H20.03H20.16H20.2L24.1 15.7L30.11 5.19L23.75 5.67ZM4.21002 4.21L10.53 4.81L12.33 0.69L12.05 0H8.43002L4.22002 4.21H4.21002ZM21.9 17.4L20.6 18.2H14.3L13 17.4L12.4 18.2L12.42 18.23L17.45 26.8L22.48 18.23L22.5 18.2L21.9 17.4ZM4.79002 5.19L10.8 15.7L14.7 17.14H14.74H15.2H16.85V14.24L12.24 3.09L11.15 5.68L4.79002 5.2V5.19Z"
                                    fill="var(--text-color)"
                                />
                            </svg>
                            <span class="font-semibold text-2xl text-primary">Your Logo</span>
                        </span>
                        <span v-if="collapsible">
                            <Button 
                                type="button" 
                                @click="toggleSidebar" 
                                :icon="isCollapsed ? 'pi pi-angle-right' : 'pi pi-angle-left'" 
                                rounded 
                                outlined 
                                class="h-2rem w-2rem"
                            />
                        </span>
                    </div>

                    <!-- Navigation Content -->
                    <div class="overflow-y-auto flex-1">
                        <!-- Favorites Section -->
                        <ul class="list-none p-3 m-0">
                            <li>
                                <div
                                    v-ripple
                                    v-styleclass="{
                                        selector: '@next',
                                        enterClass: 'hidden',
                                        enterActiveClass: 'slidedown',
                                        leaveToClass: 'hidden',
                                        leaveActiveClass: 'slideup'
                                    }"
                                    class="p-3 flex align-items-center justify-content-between text-600 cursor-pointer p-ripple"
                                    @click="toggleSection('favorites')"
                                >
                                    <span class="font-medium">FAVORITES</span>
                                    <i :class="sections.favorites ? 'pi pi-chevron-up' : 'pi pi-chevron-down'"></i>
                                </div>
                                <ul class="list-none p-0 m-0 overflow-hidden" :class="{ 'hidden': !sections.favorites }">
                                    <li v-for="item in favoriteItems" :key="item.id">
                                        <a 
                                            v-ripple 
                                            class="flex align-items-center cursor-pointer p-3 border-round text-700 hover:surface-100 transition-duration-150 transition-colors p-ripple"
                                            :class="{ 'surface-200': activeItem === item.id }"
                                            @click="selectMenuItem(item)"
                                        >
                                            <i :class="item.icon + ' mr-2'"></i>
                                            <span class="font-medium">{{ item.label }}</span>
                                            <span 
                                                v-if="item.badge" 
                                                class="inline-flex align-items-center justify-content-center ml-auto bg-primary border-circle text-white text-xs"
                                                style="min-width: 1.5rem; height: 1.5rem"
                                            >
                                                {{ item.badge }}
                                            </span>
                                        </a>

                                        <!-- Submenu for Reports -->
                                        <ul 
                                            v-if="item.children && expandedItems.includes(item.id)" 
                                            class="list-none py-0 pl-3 pr-0 m-0 overflow-y-hidden transition-all transition-duration-400 transition-ease-in-out"
                                        >
                                            <li v-for="child in item.children" :key="child.id">
                                                <a
                                                    v-ripple
                                                    class="flex align-items-center cursor-pointer p-3 border-round text-700 hover:surface-100 transition-duration-150 transition-colors p-ripple"
                                                    :class="{ 'surface-200': activeItem === child.id }"
                                                    @click="selectMenuItem(child)"
                                                >
                                                    <i :class="child.icon + ' mr-2'"></i>
                                                    <span class="font-medium">{{ child.label }}</span>
                                                    <i v-if="child.children" class="pi pi-chevron-down ml-auto"></i>
                                                </a>

                                                <!-- Third level submenu -->
                                                <ul 
                                                    v-if="child.children && expandedItems.includes(child.id)" 
                                                    class="list-none py-0 pl-3 pr-0 m-0 overflow-y-hidden transition-all transition-duration-400 transition-ease-in-out"
                                                >
                                                    <li v-for="grandchild in child.children" :key="grandchild.id">
                                                        <a 
                                                            v-ripple 
                                                            class="flex align-items-center cursor-pointer p-3 border-round text-700 hover:surface-100 transition-duration-150 transition-colors p-ripple"
                                                            :class="{ 'surface-200': activeItem === grandchild.id }"
                                                            @click="selectMenuItem(grandchild)"
                                                        >
                                                            <i :class="grandchild.icon + ' mr-2'"></i>
                                                            <span class="font-medium">{{ grandchild.label }}</span>
                                                        </a>
                                                    </li>
                                                </ul>
                                            </li>
                                        </ul>
                                    </li>
                                </ul>
                            </li>
                        </ul>

                        <!-- Application Section -->
                        <ul class="list-none p-3 m-0">
                            <li>
                                <div
                                    v-ripple
                                    v-styleclass="{
                                        selector: '@next',
                                        enterClass: 'hidden',
                                        enterActiveClass: 'slidedown',
                                        leaveToClass: 'hidden',
                                        leaveActiveClass: 'slideup'
                                    }"
                                    class="p-3 flex align-items-center justify-content-between text-600 cursor-pointer p-ripple"
                                    @click="toggleSection('application')"
                                >
                                    <span class="font-medium">APPLICATION</span>
                                    <i :class="sections.application ? 'pi pi-chevron-up' : 'pi pi-chevron-down'"></i>
                                </div>
                                <ul class="list-none p-0 m-0 overflow-hidden" :class="{ 'hidden': !sections.application }">
                                    <li v-for="item in applicationItems" :key="item.id">
                                        <a 
                                            v-ripple 
                                            class="flex align-items-center cursor-pointer p-3 border-round text-700 hover:surface-100 transition-duration-150 transition-colors p-ripple"
                                            :class="{ 'surface-200': activeItem === item.id }"
                                            @click="selectMenuItem(item)"
                                        >
                                            <i :class="item.icon + ' mr-2'"></i>
                                            <span class="font-medium">{{ item.label }}</span>
                                        </a>
                                    </li>
                                </ul>
                            </li>
                        </ul>
                    </div>

                    <!-- User Profile Footer -->
                    <div class="mt-auto">
                        <hr class="mb-3 mx-3 border-top-1 border-none surface-border" />
                        <a 
                            v-ripple 
                            class="m-3 flex align-items-center cursor-pointer p-3 gap-2 border-round text-700 hover:surface-100 transition-duration-150 transition-colors p-ripple"
                            @click="openUserMenu"
                        >
                            <Avatar :image="userProfile.avatar" shape="circle" />
                            <span class="font-bold">{{ userProfile.name }}</span>
                        </a>
                    </div>
                </div>
            </template>
        </Sidebar>

        <!-- Main Content Area -->
        <div class="main-content" :class="{ 'sidebar-open': visible, 'sidebar-collapsed': isCollapsed }">
            <!-- Toggle Button for Mobile/Collapsed State -->
            <div class="content-header p-3 border-bottom-1 surface-border">
                <Button 
                    v-if="!visible || isCollapsed" 
                    icon="pi pi-bars" 
                    @click="toggleSidebar" 
                    rounded
                    outlined
                    class="mr-3"
                />
                <h2 class="m-0">{{ currentPageTitle }}</h2>
            </div>

            <!-- Your main content goes here -->
            <div class="p-4">
                <div class="card">
                    <h3>Main Content Area</h3>
                    <p>This is where your main application content will be displayed.</p>
                    <p>Current active menu: <strong>{{ activeMenuItem?.label || 'None' }}</strong></p>
                    
                    <!-- Sample content -->
                    <div class="grid">
                        <div class="col-12 md:col-6 lg:col-3">
                            <div class="surface-0 shadow-2 p-3 border-1 border-50 border-round">
                                <div class="flex justify-content-between mb-3">
                                    <div>
                                        <span class="block text-500 font-medium mb-3">Orders</span>
                                        <div class="text-900 font-medium text-xl">152</div>
                                    </div>
                                    <div class="flex align-items-center justify-content-center bg-blue-100 border-round" style="width:2.5rem;height:2.5rem">
                                        <i class="pi pi-shopping-cart text-blue-500 text-xl"></i>
                                    </div>
                                </div>
                                <span class="text-green-500 font-medium">24 new </span>
                                <span class="text-500">since last visit</span>
                            </div>
                        </div>
                        <div class="col-12 md:col-6 lg:col-3">
                            <div class="surface-0 shadow-2 p-3 border-1 border-50 border-round">
                                <div class="flex justify-content-between mb-3">
                                    <div>
                                        <span class="block text-500 font-medium mb-3">Revenue</span>
                                        <div class="text-900 font-medium text-xl">$2,100</div>
                                    </div>
                                    <div class="flex align-items-center justify-content-center bg-orange-100 border-round" style="width:2.5rem;height:2.5rem">
                                        <i class="pi pi-map-marker text-orange-500 text-xl"></i>
                                    </div>
                                </div>
                                <span class="text-green-500 font-medium">%52+ </span>
                                <span class="text-500">since last week</span>
                            </div>
                        </div>
                        <div class="col-12 md:col-6 lg:col-3">
                            <div class="surface-0 shadow-2 p-3 border-1 border-50 border-round">
                                <div class="flex justify-content-between mb-3">
                                    <div>
                                        <span class="block text-500 font-medium mb-3">Customers</span>
                                        <div class="text-900 font-medium text-xl">28441</div>
                                    </div>
                                    <div class="flex align-items-center justify-content-center bg-cyan-100 border-round" style="width:2.5rem;height:2.5rem">
                                        <i class="pi pi-inbox text-cyan-500 text-xl"></i>
                                    </div>
                                </div>
                                <span class="text-green-500 font-medium">520  </span>
                                <span class="text-500">newly registered</span>
                            </div>
                        </div>
                        <div class="col-12 md:col-6 lg:col-3">
                            <div class="surface-0 shadow-2 p-3 border-1 border-50 border-round">
                                <div class="flex justify-content-between mb-3">
                                    <div>
                                        <span class="block text-500 font-medium mb-3">Comments</span>
                                        <div class="text-900 font-medium text-xl">152 Unread</div>
                                    </div>
                                    <div class="flex align-items-center justify-content-center bg-purple-100 border-round" style="width:2.5rem;height:2.5rem">
                                        <i class="pi pi-comment text-purple-500 text-xl"></i>
                                    </div>
                                </div>
                                <span class="text-green-500 font-medium">85 </span>
                                <span class="text-500">responded</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import Sidebar from 'primevue/sidebar'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'

// Types
interface MenuItem {
    id: string
    label: string
    icon: string
    route?: string
    badge?: number
    children?: MenuItem[]
}

interface UserProfile {
    name: string
    avatar: string
    email?: string
}

interface SidebarSections {
    favorites: boolean
    application: boolean
}

// Props
interface SidebarProps {
    defaultOpen?: boolean
    collapsible?: boolean
    width?: string
    collapsedWidth?: string
}

const props = withDefaults(defineProps<SidebarProps>(), {
    defaultOpen: true,
    collapsible: true,
    width: '20rem',
    collapsedWidth: '4rem'
})

// Emits
interface SidebarEmits {
    'menu-select': [item: MenuItem]
    'toggle': [isOpen: boolean]
    'user-menu': []
}

const emit = defineEmits<SidebarEmits>()

// Reactive state
const visible = ref<boolean>(props.defaultOpen)
const isCollapsed = ref<boolean>(false)
const activeItem = ref<string>('')
const expandedItems = ref<string[]>(['reports'])
const currentPageTitle = ref<string>('Dashboard')

const sections = reactive<SidebarSections>({
    favorites: true,
    application: true
})

const userProfile = reactive<UserProfile>({
    name: 'Amy Elsner',
    avatar: 'https://primefaces.org/cdn/primevue/images/avatar/amyelsner.png'
})

// Menu items
const favoriteItems = ref<MenuItem[]>([
    { id: 'dashboard', label: 'Dashboard', icon: 'pi pi-home', route: '/dashboard' },
    { id: 'bookmarks', label: 'Bookmarks', icon: 'pi pi-bookmark', route: '/bookmarks' },
    { 
        id: 'reports', 
        label: 'Reports', 
        icon: 'pi pi-chart-line', 
        route: '/reports',
        children: [
            { 
                id: 'revenue', 
                label: 'Revenue', 
                icon: 'pi pi-chart-line', 
                route: '/reports/revenue',
                children: [
                    { id: 'revenue-view', label: 'View', icon: 'pi pi-table', route: '/reports/revenue/view' },
                    { id: 'revenue-search', label: 'Search', icon: 'pi pi-search', route: '/reports/revenue/search' }
                ]
            },
            { id: 'expenses', label: 'Expenses', icon: 'pi pi-chart-line', route: '/reports/expenses' }
        ]
    },
    { id: 'team', label: 'Team', icon: 'pi pi-users', route: '/team' },
    { id: 'messages', label: 'Messages', icon: 'pi pi-comments', route: '/messages', badge: 3 },
    { id: 'calendar', label: 'Calendar', icon: 'pi pi-calendar', route: '/calendar' },
    { id: 'settings', label: 'Settings', icon: 'pi pi-cog', route: '/settings' }
])

const applicationItems = ref<MenuItem[]>([
    { id: 'projects', label: 'Projects', icon: 'pi pi-folder', route: '/projects' },
    { id: 'performance', label: 'Performance', icon: 'pi pi-chart-bar', route: '/performance' },
    { id: 'app-settings', label: 'Settings', icon: 'pi pi-cog', route: '/app-settings' }
])

// Computed
const sidebarWidth = computed(() => {
    if (isCollapsed.value) return props.collapsedWidth
    return props.width
})

const activeMenuItem = computed(() => {
    const findMenuItem = (items: MenuItem[]): MenuItem | null => {
        for (const item of items) {
            if (item.id === activeItem.value) return item
            if (item.children) {
                const found = findMenuItem(item.children)
                if (found) return found
            }
        }
        return null
    }
    
    return findMenuItem([...favoriteItems.value, ...applicationItems.value])
})

// Methods
const toggleSidebar = (): void => {
    if (props.collapsible && visible.value) {
        isCollapsed.value = !isCollapsed.value
    } else {
        visible.value = !visible.value
    }
    emit('toggle', visible.value && !isCollapsed.value)
}

const toggleSection = (section: keyof SidebarSections): void => {
    sections[section] = !sections[section]
}

const selectMenuItem = (item: MenuItem): void => {
    activeItem.value = item.id
    currentPageTitle.value = item.label
    
    // Handle submenu expansion
    if (item.children) {
        const index = expandedItems.value.indexOf(item.id)
        if (index > -1) {
            expandedItems.value.splice(index, 1)
        } else {
            expandedItems.value.push(item.id)
        }
    }
    
    emit('menu-select', item)
    
    // Auto-collapse on mobile after selection
    if (window.innerWidth < 768) {
        visible.value = false
    }
}

const openUserMenu = (): void => {
    emit('user-menu')
}

// Lifecycle
onMounted(() => {
    // Set initial active item
    activeItem.value = 'dashboard'
    
    // Handle responsive behavior
    const handleResize = () => {
        if (window.innerWidth < 768) {
            visible.value = false
            isCollapsed.value = false
        } else {
            visible.value = props.defaultOpen
        }
    }
    
    handleResize()
    window.addEventListener('resize', handleResize)
})
</script>

<style scoped>
.main-layout {
    display: flex;
    min-height: 100vh;
    width: 100%;
}

.main-sidebar {
    position: fixed;
    left: 0;
    top: 0;
    height: 100vh;
    z-index: 1000;
    transition: transform 0.3s ease;
}

.main-content {
    flex: 1;
    transition: margin-left 0.3s ease;
    margin-left: 0;
}

.main-content.sidebar-open {
    margin-left: v-bind('props.width');
}

.main-content.sidebar-collapsed {
    margin-left: v-bind('props.collapsedWidth');
}

.content-header {
    display: flex;
    align-items: center;
    background: var(--surface-0);
    position: sticky;
    top: 0;
    z-index: 100;
}

/* Mobile responsive */
@media (max-width: 768px) {
    .main-content.sidebar-open {
        margin-left: 0;
    }
    
    .main-content.sidebar-collapsed {
        margin-left: 0;
    }
    
    .main-sidebar {
        transform: translateX(-100%);
    }
    
    .main-sidebar.p-sidebar-visible {
        transform: translateX(0);
    }
}

/* Custom animations for slidedown/slideup */
@keyframes slidedown {
    from {
        opacity: 0;
        transform: translateY(-10px);
        max-height: 0;
    }
    to {
        opacity: 1;
        transform: translateY(0);
        max-height: 1000px;
    }
}

@keyframes slideup {
    from {
        opacity: 1;
        transform: translateY(0);
        max-height: 1000px;
    }
    to {
        opacity: 0;
        transform: translateY(-10px);
        max-height: 0;
    }
}

.slidedown {
    animation: slidedown 0.3s ease forwards;
}

.slideup {
    animation: slideup 0.3s ease forwards;
}

/* Active menu item styling */
.surface-200 {
    background-color: var(--surface-200) !important;
    border-left: 3px solid var(--primary-color);
}

/* Smooth transitions */
.p-sidebar {
    transition: width 0.3s ease;
}

/* Hide text when collapsed */
.sidebar-collapsed .font-medium {
    display: none;
}

.sidebar-collapsed .ml-auto {
    display: none;
}
</style>