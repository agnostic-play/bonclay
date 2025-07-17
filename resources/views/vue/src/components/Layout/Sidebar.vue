<template>
  <aside class="w-64 bg-white shadow h-full flex flex-col justify-between">
    <div>
      <div class="p-4 font-bold text-lg border-b">ZenTrail</div>
      <Menu :model="items" class="p-2" />
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Menu from 'primevue/menu'

const items = ref([
  {
    label: 'Home',
    items: [
      { label: 'Dashboard', icon: 'pi pi-home' },
      { label: 'Bookmarks', icon: 'pi pi-bookmark' },
      { label: 'Team', icon: 'pi pi-users' },
      { label: 'Messages', icon: 'pi pi-envelope' },
      { label: 'Calendar', icon: 'pi pi-calendar' }
    ]
  },
  {
    label: 'Your Networks',
    items: [
      { label: 'Front-End Developers', icon: 'pi pi-code' },
      { label: 'Back-End Developers', icon: 'pi pi-database' },
      { label: 'UI/UX Designers', icon: 'pi pi-palette' }
    ]
  }
])
</script>

<style scoped>
.p-menu {
  border: none;
}
</style>
