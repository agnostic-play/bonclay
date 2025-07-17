<template>
  <router-link
      :to="item.route"
      class="menu-item-link"
      :class="{ 'network-item': item.iconColor }"
  >
    <div
        v-if="item.iconColor"
        :class="['network-icon', item.iconColor]"
    >
      <i :class="item.icon"></i>
    </div>
    <i
        v-else
        :class="[item.icon, 'menu-icon']"
    />

    <span
        v-if="!isCollapsed || isMobile"
        class="menu-text"
    >
      {{ item.label }}
    </span>

    <Badge
        v-if="item.badge && (!isCollapsed || isMobile)"
        :value="item.badge.value"
        :severity="item.badge.severity"
        class="menu-badge"
    />
  </router-link>
</template>

<script setup lang="ts">
import Badge from 'primevue/badge'
import type { MenuItemType } from '../types/menu'

interface Props {
  item: MenuItemType
  isCollapsed: boolean
  isMobile: boolean
}

defineProps<Props>()
</script>

<style scoped>
.menu-item-link {
  display: flex;
  align-items: center;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
  color: #495057;
  text-decoration: none;
}

.menu-item-link:hover {
  background-color: #f8f9fa;
  color: #2c3e50;
  text-decoration: none;
}

.menu-item-link.router-link-active {
  background-color: #e3f2fd;
  color: #1976d2;
}

.menu-item-link.router-link-active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background-color: #1976d2;
}

.menu-icon {
  font-size: 1.1rem;
  width: 24px;
}

.menu-text {
  margin-left: 0.75rem;
  flex: 1;
}

.menu-badge {
  margin-left: auto;
}

.network-icon {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 0.875rem;
}

.network-icon.purple {
  background-color: #9c27b0;
}

.network-icon.orange {
  background-color: #ff9800;
}

.network-icon.green {
  background-color: #4caf50;
}
</style>