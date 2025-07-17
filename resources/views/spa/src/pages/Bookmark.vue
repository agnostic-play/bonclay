
<!-- Bookmarks.vue -->
<template>
  <div class="bookmarks-page">
    <div class="page-header">
      <h1>Bookmarks</h1>
      <p>Your saved items and important links</p>
    </div>

    <div class="bookmarks-content">
      <div class="bookmark-item" v-for="bookmark in bookmarks" :key="bookmark.id">
        <div class="bookmark-icon">
          <i :class="bookmark.icon"></i>
        </div>
        <div class="bookmark-info">
          <h3>{{ bookmark.title }}</h3>
          <p>{{ bookmark.description }}</p>
          <span class="bookmark-url">{{ bookmark.url }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const bookmarks = ref([
  {
    id: 1,
    title: 'Vue.js Documentation',
    description: 'Official Vue.js documentation and guides',
    url: 'https://vuejs.org',
    icon: 'pi pi-book'
  },
  {
    id: 2,
    title: 'PrimeVue Components',
    description: 'UI component library for Vue.js',
    url: 'https://primevue.org',
    icon: 'pi pi-palette'
  }
])
</script>

<style scoped>
.bookmarks-page {
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  color: #2c3e50;
  margin-bottom: 0.5rem;
}

.page-header p {
  color: #6c757d;
  font-size: 1.1rem;
}

.bookmarks-content {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.bookmark-item {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  display: flex;
  align-items: center;
  gap: 1rem;
}

.bookmark-icon {
  width: 48px;
  height: 48px;
  background: #e3f2fd;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #1976d2;
  font-size: 1.5rem;
}

.bookmark-info h3 {
  margin: 0 0 0.5rem 0;
  color: #2c3e50;
}

.bookmark-info p {
  margin: 0 0 0.5rem 0;
  color: #6c757d;
}

.bookmark-url {
  font-size: 0.875rem;
  color: #1976d2;
}
</style>