<template>
  <AppLayout>
    <!-- Your page content goes here -->
    <div class="dashboard-content">
      <h1>Dashboard</h1>
      <p>Welcome to your dashboard! This is where your main content will be displayed.</p>
      
      <!-- Example Cards -->
      <div class="card-grid">
        <Card v-for="i in 4" :key="i">
          <template #title>
            Card Title {{ i }}
          </template>
          <template #content>
            <p>This is some example content for card {{ i }}.</p>
          </template>
        </Card>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import AppLayout from './components/AppLayout.vue';
import Card from 'primevue/card';
</script>

<style scoped>
.dashboard-content {
  max-width: 1200px;
  margin: 0 auto;
}

.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-top: 2rem;
}

h1 {
  color: #2c3e50;
  margin-bottom: 1rem;
}
</style>