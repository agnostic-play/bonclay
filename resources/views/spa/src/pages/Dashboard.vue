<!-- Dashboard.vue -->
<template>
  <div class="dashboard-page">
    <div class="page-header">
      <h1>Dashboard</h1>
      <p>Welcome to your Ditto dashboard</p>
    </div>

    <div class="dashboard-grid">
      <div class="card">
        <h3>Total Projects</h3>
        <div class="stat-number">24</div>
      </div>

      <div class="card">
        <h3>Active Tasks</h3>
        <div class="stat-number">12</div>
      </div>

      <div class="card">
        <h3>Team Members</h3>
        <div class="stat-number">8</div>
      </div>

      <div class="card">
        <h3>Completed</h3>
        <div class="stat-number">156</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Page logic here
</script>

<style scoped>
.dashboard-page {
  max-width: 1200px;
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

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  text-align: center;
}

.card h3 {
  color: #6c757d;
  font-size: 0.9rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 1rem;
}

.stat-number {
  font-size: 2.5rem;
  font-weight: 700;
  color: #1976d2;
}
</style>
