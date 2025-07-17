<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup lang="ts">
// No additional logic needed here - just the router outlet
</script>

<style>
/* Global styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

#app {
  height: 100vh;
  width: 100vw;
}

/* PrimeVue theme customization */
:root {
  --primary-color: #1976d2;
  --primary-color-text: #ffffff;
  --surface-0: #ffffff;
  --surface-50: #f8f9fa;
  --surface-100: #e9ecef;
  --surface-200: #dee2e6;
  --surface-300: #ced4da;
  --surface-400: #adb5bd;
  --surface-500: #6c757d;
  --surface-600: #495057;
  --surface-700: #343a40;
  --surface-800: #212529;
  --surface-900: #121416;
}
</style>