<template>
  <div
      :class="['sidebar', {
      'sidebar-collapsed': isCollapsed,
      'sidebar-mobile-open': isMobileOpen && isMobile,
      'sidebar-mobile': isMobile
    }]"
  >
    <!-- Logo Section -->
    <SidebarHeader
        :is-collapsed="isCollapsed"
        :is-mobile="isMobile"
    />

    <!-- Navigation Menu -->
    <div class="sidebar-content">
      <MenuSection
          v-for="section in menuSections"
          :key="section.id"
          :section="section"
          :is-collapsed="isCollapsed"
          :is-mobile="isMobile"
      />
    </div>

    <!-- Collapse Toggle Button (Desktop Only) -->
    <div class="sidebar-footer" v-if="!isMobile">
      <Button
          :icon="isCollapsed ? 'pi pi-angle-double-right' : 'pi pi-angle-double-left'"
          @click="$emit('toggle-sidebar')"
          class="p-button-text p-button-plain collapse-button"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import Button from 'primevue/button'
import SidebarHeader from './SidebarHeader.vue'
import MenuSection from './MenuSection.vue'
import type {MenuSectionType} from "@/types/menu.ts";

interface Props {
  isCollapsed: boolean
  isMobile: boolean
  isMobileOpen: boolean
}

defineProps<Props>()

defineEmits<{
  'toggle-sidebar': []
}>()

const menuSections = ref<MenuSectionType[]>([
  {
    id: 'APIMock',
    title: 'API Mock',
    isExpanded: false,
    items: [
      {
        id: 'squads',
        label: 'Squad',
        icon: 'pi pi-users',
        route: '/'
      },
      {
        id: 'histories',
        label: 'History',
        icon: 'pi pi-history',
        route: '/asd'
      },
    ]
  },
  {
    id: 'tools',
    title: 'Tools',
    isExpanded: false,
    items: [
      {
        id: 'encryption-tools',
        label: 'Encryption Tools',
        icon: 'pi pi-shield',
        route: '/encryption-tools'
      }
    ]
  },
])
</script>

<style scoped>
.sidebar {
  width: 280px;
  background-color: #ffffff;
  border-right: 1px solid #e9ecef;
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  position: relative;
  z-index: 1000;
}

.sidebar-collapsed {
  width: 70px;
}

.sidebar-content {
  flex: 1;
  overflow-y: auto;
  padding: 1rem 0;
}

.sidebar-footer {
  padding: 1rem;
  border-top: 1px solid #e9ecef;
}

.collapse-button {
  width: 100%;
  justify-content: center;
}

/* Mobile Styles */
@media (max-width: 767px) {
  .sidebar {
    position: fixed;
    top: 60px;
    left: -280px;
    bottom: 0;
    width: 280px;
    transition: left 0.3s ease;
  }

  .sidebar-mobile-open {
    left: 0;
  }
}

/* Scrollbar Styles */
.sidebar-content::-webkit-scrollbar {
  width: 6px;
}

.sidebar-content::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.sidebar-content::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 3px;
}

.sidebar-content::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>