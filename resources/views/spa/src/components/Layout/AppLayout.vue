<template>
  <div class="app-layout">
    <!-- Mobile Header -->
    <MobileHeader
        v-if="isMobile"
        @toggle-sidebar="toggleMobileSidebar"
    />

    <!-- Sidebar -->
    <Sidebar
        :is-collapsed="isCollapsed && !isMobile"
        :is-mobile="isMobile"
        :is-mobile-open="isMobileMenuOpen"
        :active-item="activeMenuItem"
        @toggle-sidebar="toggleSidebar"
        @set-active-item="handleMenuItemClick"
        @close-mobile="closeMobileSidebar"
    />

    <!-- Main Content Area -->
    <div :class="['main-content', { 'content-expanded': isCollapsed && !isMobile }]">
      <!-- Top Header -->
      <TopHeader
          v-if="!isMobile"
          @toggle-sidebar="toggleSidebar"
      />

      <!-- Page Content -->
      <div class="page-content">
        <router-view />
      </div>
    </div>

    <!-- Mobile Overlay -->
    <div
        v-if="isMobileMenuOpen && isMobile"
        class="mobile-overlay"
        @click="closeMobileSidebar"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import MobileHeader from './MobileHeader.vue'
import Sidebar from './Sidebar.vue'
import TopHeader from './TopHeader.vue'
import { useResponsive } from '@/composables/userResponsive.ts'

const router = useRouter()
const route = useRoute()

// Reactive states
const isCollapsed = ref(false)
const { isMobile, isMobileMenuOpen } = useResponsive()

// Get active menu item from route meta
const activeMenuItem = computed(() => {
  return route.meta.menuItem as string || 'dashboard'
})

// Methods
const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}

const toggleMobileSidebar = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

const closeMobileSidebar = () => {
  isMobileMenuOpen.value = false
}

const handleMenuItemClick = (menuItem: string) => {
  // Route mapping for menu items
  const routeMap: Record<string, string> = {
    dashboard: '/',
    bookmarks: '/bookmarks',
    team: '/team',
    messages: '/messages',
    calendar: '/calendar',
    frontend: '/networks/frontend',
    backend: '/networks/backend',
    uiux: '/networks/uiux'
  }

  const targetRoute = routeMap[menuItem]
  if (targetRoute && route.path !== targetRoute) {
    router.push(targetRoute)
  }

  if (isMobile.value) {
    closeMobileSidebar()
  }
}

// Watch for route changes to sync menu state
watch(route, (newRoute) => {
  // You can add additional logic here if needed
  console.log('Route changed to:', newRoute.path)
}, { immediate: true })
</script>

<style scoped>
.app-layout {
  display: flex;
  height: 100vh;
  background-color: #f8f9fa;
  position: relative;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  transition: margin-left 0.3s ease;
  overflow: hidden;
}

.page-content {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
  background-color: #f8f9fa;
}

.mobile-overlay {
  position: fixed;
  top: 60px;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 998;
}

@media (max-width: 767px) {
  .main-content {
    margin-left: 0;
    padding-top: 60px;
  }

  .page-content {
    padding: 1rem;
  }
}
</style>