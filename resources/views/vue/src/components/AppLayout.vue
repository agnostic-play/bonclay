<template>
  <div class="app-layout">
    <!-- Mobile Menu Button -->
    <div class="mobile-header">
      <Button 
        icon="pi pi-bars" 
        @click="toggleMobileSidebar"
        class="p-button-text p-button-plain mobile-menu-button"
        v-if="isMobile"
      />
      <div class="mobile-logo">
        <i class="pi pi-slack logo-icon"></i>
        <span>ZenTrail</span>
      </div>
    </div>

    <!-- Sidebar -->
    <div 
      :class="['sidebar', { 
        'sidebar-collapsed': isCollapsed && !isMobile,
        'sidebar-mobile-open': isMobileMenuOpen && isMobile,
        'sidebar-mobile': isMobile 
      }]"
    >
      <!-- Logo Section -->
      <div class="sidebar-header">
        <div class="logo">
          <i class="pi pi-slack logo-icon"></i>
          <span v-if="!isCollapsed || isMobile" class="logo-text">ZenTrail</span>
        </div>
      </div>

      <!-- Navigation Menu -->
      <div class="sidebar-content">
        <!-- Home Section -->
        <div class="menu-section">
          <div class="section-header" @click="toggleSection('home')">
            <span class="section-title">
              <span v-if="!isCollapsed || isMobile">Home</span>
            </span>
            <i :class="['pi', sections.home ? 'pi-chevron-down' : 'pi-chevron-right']" 
               v-if="!isCollapsed || isMobile"></i>
          </div>
          
          <div v-if="sections.home" class="menu-items">
            <div 
              class="menu-item" 
              :class="{ 'active': activeItem === 'dashboard' }"
              @click="setActiveItem('dashboard')"
            >
              <i class="pi pi-home menu-icon"></i>
              <span v-if="!isCollapsed || isMobile" class="menu-text">Dashboard</span>
            </div>
            
            <div 
              class="menu-item"
              :class="{ 'active': activeItem === 'bookmarks' }"
              @click="setActiveItem('bookmarks')"
            >
              <i class="pi pi-bookmark menu-icon"></i>
              <span v-if="!isCollapsed || isMobile" class="menu-text">Bookmarks</span>
              <Badge v-if="!isCollapsed || isMobile" value="8" class="menu-badge" />
            </div>
            
            <div 
              class="menu-item"
              :class="{ 'active': activeItem === 'team' }"
              @click="setActiveItem('team')"
            >
              <i class="pi pi-users menu-icon"></i>
              <span v-if="!isCollapsed || isMobile" class="menu-text">Team</span>
            </div>
            
            <div 
              class="menu-item"
              :class="{ 'active': activeItem === 'messages' }"
              @click="setActiveItem('messages')"
            >
              <i class="pi pi-comments menu-icon"></i>
              <span v-if="!isCollapsed || isMobile" class="menu-text">Messages</span>
              <Badge v-if="!isCollapsed || isMobile" value="2" severity="danger" class="menu-badge" />
            </div>
            
            <div 
              class="menu-item"
              :class="{ 'active': activeItem === 'calendar' }"
              @click="setActiveItem('calendar')"
            >
              <i class="pi pi-calendar menu-icon"></i>
              <span v-if="!isCollapsed || isMobile" class="menu-text">Calendar</span>
            </div>
          </div>
        </div>

        <!-- Networks Section -->
        <div class="menu-section">
          <div class="section-header" @click="toggleSection('networks')">
            <span class="section-title">
              <span v-if="!isCollapsed || isMobile">Your Networks</span>
            </span>
            <i :class="['pi', sections.networks ? 'pi-chevron-down' : 'pi-chevron-right']"
               v-if="!isCollapsed || isMobile"></i>
          </div>
          
          <div v-if="sections.networks" class="menu-items">
            <div 
              class="menu-item network-item"
              :class="{ 'active': activeItem === 'frontend' }"
              @click="setActiveItem('frontend')"
            >
              <div class="network-icon purple">
                <i class="pi pi-code"></i>
              </div>
              <span v-if="!isCollapsed || isMobile" class="menu-text">Front-End Developers</span>
            </div>
            
            <div 
              class="menu-item network-item"
              :class="{ 'active': activeItem === 'backend' }"
              @click="setActiveItem('backend')"
            >
              <div class="network-icon orange">
                <i class="pi pi-server"></i>
              </div>
              <span v-if="!isCollapsed || isMobile" class="menu-text">Back-End Developers</span>
            </div>
            
            <div 
              class="menu-item network-item"
              :class="{ 'active': activeItem === 'uiux' }"
              @click="setActiveItem('uiux')"
            >
              <div class="network-icon green">
                <i class="pi pi-palette"></i>
              </div>
              <span v-if="!isCollapsed || isMobile" class="menu-text">UI/UX Designers</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Collapse Toggle Button (Desktop Only) -->
      <div class="sidebar-footer" v-if="!isMobile">
        <Button 
          :icon="isCollapsed ? 'pi pi-angle-double-right' : 'pi pi-angle-double-left'"
          @click="toggleSidebar"
          class="p-button-text p-button-plain collapse-button"
        />
      </div>
    </div>

    <!-- Main Content Area -->
    <div :class="['main-content', { 'content-expanded': isCollapsed && !isMobile }]">
      <!-- Top Header -->
      <div class="content-header">
        <div class="header-left">
          <Button 
            icon="pi pi-bars" 
            @click="toggleSidebar"
            class="p-button-text p-button-plain"
            v-if="!isMobile"
          />
        </div>
        <div class="header-right">
          <Button icon="pi pi-bell" class="p-button-text p-button-plain" badge="3" />
          <Avatar 
            image="https://primefaces.org/cdn/primevue/images/avatar/amyelsner.png" 
            shape="circle"
            class="user-avatar"
          />
        </div>
      </div>

      <!-- Page Content -->
      <div class="page-content">
        <slot></slot>
      </div>
    </div>

    <!-- Mobile Overlay -->
    <div 
      v-if="isMobileMenuOpen && isMobile" 
      class="mobile-overlay"
      @click="closeMobileSidebar"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import Button from 'primevue/button';
import Badge from 'primevue/badge';
import Avatar from 'primevue/avatar';

// Reactive states
const isCollapsed = ref(false);
const isMobile = ref(false);
const isMobileMenuOpen = ref(false);
const activeItem = ref('dashboard');
const sections = ref({
  home: true,
  networks: true
});

// Methods
const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value;
};

const toggleMobileSidebar = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value;
};

const closeMobileSidebar = () => {
  isMobileMenuOpen.value = false;
};

const setActiveItem = (item: string) => {
  activeItem.value = item;
  if (isMobile.value) {
    closeMobileSidebar();
  }
};

const toggleSection = (section: 'home' | 'networks') => {
  sections.value[section] = !sections.value[section];
};

// Handle responsive behavior
const checkMobile = () => {
  isMobile.value = window.innerWidth < 768;
  if (!isMobile.value) {
    isMobileMenuOpen.value = false;
  }
};

// Lifecycle hooks
onMounted(() => {
  checkMobile();
  window.addEventListener('resize', checkMobile);
});

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile);
});
</script>

<style scoped>
.app-layout {
  display: flex;
  height: 100vh;
  background-color: #f8f9fa;
  position: relative;
}

/* Mobile Header */
.mobile-header {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: white;
  border-bottom: 1px solid #e9ecef;
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

/* Sidebar Styles */
.sidebar {
  width: 280px;
  background-color: #ffffff;
  border-right: 1px solid #e9ecef;
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  position: relative;
  z-index: 1000;
}

.sidebar-collapsed {
  width: 70px;
}

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

.sidebar-collapsed .logo {
  justify-content: center;
}

.sidebar-content {
  flex: 1;
  overflow-y: auto;
  padding: 1rem 0;
}

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

.sidebar-collapsed .section-header {
  justify-content: center;
}

.menu-items {
  margin-top: 0.5rem;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
  color: #495057;
}

.menu-item:hover {
  background-color: #f8f9fa;
  color: #2c3e50;
}

.menu-item.active {
  background-color: #e3f2fd;
  color: #1976d2;
}

.menu-item.active::before {
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

.sidebar-collapsed .menu-item {
  justify-content: center;
  padding: 0.75rem;
}

.network-item {
  display: flex;
  align-items: center;
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

/* Sidebar Footer */
.sidebar-footer {
  padding: 1rem;
  border-top: 1px solid #e9ecef;
}

.collapse-button {
  width: 100%;
  justify-content: center;
}

/* Main Content */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  transition: margin-left 0.3s ease;
  overflow: hidden;
}

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

.page-content {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
  background-color: #f8f9fa;
}

/* Mobile Styles */
@media (max-width: 767px) {
  .mobile-header {
    display: flex;
  }

  .sidebar {
    position: fixed;
    top: 60px;
    left: -280px;
    bottom: 0;
    width: 280px;
    transition: left 0.3s ease;
  }

  .sidebar-mobile-open {
    left: 0;
  }

  .main-content {
    margin-left: 0;
    padding-top: 60px;
  }

  .content-header {
    display: none;
  }

  .mobile-overlay {
    position: fixed;
    top: 60px;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 998;
  }

  .page-content {
    padding: 1rem;
  }
}

/* Scrollbar Styles */
.sidebar-content::-webkit-scrollbar {
  width: 6px;
}

.sidebar-content::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.sidebar-content::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 3px;
}

.sidebar-content::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>