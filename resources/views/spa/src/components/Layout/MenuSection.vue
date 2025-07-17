<template>
  <div class="menu-section">
    <div class="section-header" @click="toggleSection">
      <span class="section-title">
        <span v-if="!isCollapsed || isMobile">{{ section.title }}</span>
      </span>
      <i
          :class="['pi', section.isExpanded ? 'pi-chevron-down' : 'pi-chevron-right']"
          v-if="!isCollapsed || isMobile"
      />
    </div>

    <div v-if="section.isExpanded" class="menu-items">
      <MenuItem
          v-for="item in section.items"
          :key="item.id"
          :item="item"
          :is-collapsed="isCollapsed"
          :is-mobile="isMobile"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import MenuItem from './MenuItem.vue'
import type { MenuSectionType } from '../types/menu'

interface Props {
  section: MenuSectionType
  isCollapsed: boolean
  isMobile: boolean
}

const props = defineProps<Props>()

const toggleSection = () => {
  props.section.isExpanded = !props.section.isExpanded
}
</script>

<style scoped>
.menu-section {
  margin-bottom: 1.5rem;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 1.5rem;
  cursor: pointer;
  color: #6c757d;
  font-size: 0.875rem;
  font-weight: 500;
  transition: color 0.2s;
}

.section-header:hover {
  color: #495057;
}

.menu-items {
  margin-top: 0.5rem;
}
</style>