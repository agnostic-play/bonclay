<template>
  <div class="content-header">
    <div class="header-left">
      <Button
          icon="pi pi-bars"
          @click="$emit('toggle-sidebar')"
          class="p-button-text p-button-plain"
      />
    </div>
    <div class="header-right">

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
.content-header {
  height: 60px;
  background-color: #ffffff;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1.5rem;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-avatar {
  cursor: pointer;
}
</style>