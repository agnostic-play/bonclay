<template>
  <div class="mobile-header">
    <Button
        icon="pi pi-bars"
        @click="$emit('toggle-sidebar')"
        class="p-button-text p-button-plain mobile-menu-button"
    />
    <div class="mobile-logo">
      <i class="pi pi-slack logo-icon"></i>
      <span>Ditto</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'

defineEmits<{
  'toggle-sidebar': []
}>()
</script>

<style scoped>
.mobile-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: white;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  align-items: center;
  padding: 0 1rem;
  z-index: 999;
}

.mobile-logo {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
  font-size: 1.2rem;
  margin-left: auto;
  margin-right: auto;
}

.mobile-logo .logo-icon {
  font-size: 1.5rem;
  color: #2c3e50;
}
</style>