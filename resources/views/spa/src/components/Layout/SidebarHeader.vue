<template>
  <div class="sidebar-header">
    <div class="logo">
      <i class="pi pi-slack logo-icon"></i>
      <span v-if="!isCollapsed || isMobile" class="logo-text">Ditto</span>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  isCollapsed: boolean
  isMobile: boolean
}

defineProps<Props>()
</script>

<style scoped>
.sidebar-header {
  padding: 1.5rem;
  border-bottom: 1px solid #e9ecef;
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 1.25rem;
  font-weight: 600;
  color: #2c3e50;
}

.logo-icon {
  font-size: 1.5rem;
  color: #2c3e50;
}
</style>