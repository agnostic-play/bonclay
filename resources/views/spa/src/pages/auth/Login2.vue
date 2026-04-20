<template>
  <div class="layout-wrapper">
    <!-- Mobile Overlay -->
    <div 
      v-if="sidebarOpen && isMobile" 
      class="mobile-overlay" 
      @click="closeMobileSidebar"
    ></div>

    <!-- Mobile Menu Button -->
    <!-- <Button 
      v-if="isMobile"
      icon="pi pi-bars" 
      text 
      rounded 
      size="small"
      @click="toggleMobileSidebar"
      class="mobile-menu-btn"
    /> -->

    <!-- Sidebar -->
    <div 
      class="layout-sidebar" 
      :class="{ 
        'layout-sidebar-collapsed': sidebarCollapsed && !isMobile,
        'layout-sidebar-mobile': isMobile,
        'layout-sidebar-mobile-open': sidebarOpen && isMobile,
        'layout-sidebar-dark': isDarkMode
      }"
    >
      <!-- Sidebar Header -->
      <div class="sidebar-header">
        <div class="flex align-items-center gap-2">
          <i class="pi pi-home text-xl"></i>
          <span v-if="!sidebarCollapsed" class="font-medium">AlloDitoo</span>
        </div>
        <Button 
          v-if="!sidebarCollapsed && !isMobile"
          icon="pi pi-angle-left" 
          text 
          rounded 
          size="small"
          @click="toggleSidebar"
          class="toggle-btn"
        />
      </div>

      <!-- Navigation Menu -->
      <div class="sidebar-content">
        <!-- Main Navigation -->
        <div class="nav-section">
          <div class="nav-item" :class="{ active: activeItem === 'search' }" @click="setActive('search')">
            <i class="pi pi-search" :class="{ 'collapsed-icon': sidebarCollapsed && !isMobile }"></i>
            <span v-if="!sidebarCollapsed || isMobile">Search</span>
          </div>
          <div class="nav-item" :class="{ active: activeItem === 'home' }" @click="setActive('home')">
            <i class="pi pi-home" :class="{ 'collapsed-icon': sidebarCollapsed && !isMobile }"></i>
            <span v-if="!sidebarCollapsed || isMobile">Home</span>
          </div>
          <div class="nav-item" :class="{ active: activeItem === 'inbox' }" @click="setActive('inbox')">
            <i class="pi pi-inbox" :class="{ 'collapsed-icon': sidebarCollapsed && !isMobile }"></i>
            <span v-if="!sidebarCollapsed || isMobile">Inbox</span>
          </div>
        </div>

        <!-- Private Section -->
        <div class="nav-section" v-if="!sidebarCollapsed || isMobile">
          <div class="section-header">
            <span class="section-title">Private</span>
          </div>
          <div class="nav-item" :class="{ active: activeItem === 'notes' }" @click="setActive('notes')">
            <i class="pi pi-file"></i>
            <span>1:1 notes</span>
          </div>
          <div class="nav-item" :class="{ active: activeItem === 'getting-started' }" @click="setActive('getting-started')">
            <i class="pi pi-star"></i>
            <span>Getting Started</span>
          </div>
          <div class="nav-item" :class="{ active: activeItem === 'scratchpad' }" @click="setActive('scratchpad')">
            <i class="pi pi-pencil"></i>
            <span>Scratchpad</span>
          </div>
        </div>

        <!-- Teamspaces Section -->
        <div class="nav-section" v-if="!sidebarCollapsed || isMobile">
          <div class="section-header">
            <span class="section-title">Teamspaces</span>
          </div>
          <div class="nav-item" :class="{ active: activeItem === 'research-workspace' }" @click="setActive('research-workspace')">
            <i class="pi pi-users"></i>
            <span>AlloDitoo</span>
          </div>
          <div class="nav-item" :class="{ active: activeItem === 'tasks' }" @click="setActive('tasks')">
            <i class="pi pi-check-square"></i>
            <span>Tasks</span>
          </div>
          <div class="nav-item" :class="{ active: activeItem === 'projects' }" @click="setActive('projects')">
            <i class="pi pi-folder"></i>
            <span>Projects</span>
          </div>
          <div class="nav-add-item" @click="addNew">
            <i class="pi pi-plus"></i>
            <span>Add new</span>
          </div>
        </div>

        <!-- Collapsed Navigation Icons -->
        <div v-if="sidebarCollapsed && !isMobile" class="collapsed-nav-section">
          <div class="nav-item" :class="{ active: activeItem === 'notes' }" @click="setActive('notes')" title="1:1 notes">
            <i class="pi pi-file collapsed-icon"></i>
          </div>
          <div class="nav-item" :class="{ active: activeItem === 'getting-started' }" @click="setActive('getting-started')" title="Getting Started">
            <i class="pi pi-star collapsed-icon"></i>
          </div>
          <div class="nav-item" :class="{ active: activeItem === 'scratchpad' }" @click="setActive('scratchpad')" title="Scratchpad">
            <i class="pi pi-pencil collapsed-icon"></i>
          </div>
          <div class="nav-item" :class="{ active: activeItem === 'research-workspace' }" @click="setActive('research-workspace')" title="AlloDitoo">
            <i class="pi pi-users collapsed-icon"></i>
          </div>
          <div class="nav-item" :class="{ active: activeItem === 'tasks' }" @click="setActive('tasks')" title="Tasks">
            <i class="pi pi-check-square collapsed-icon"></i>
          </div>
          <div class="nav-item" :class="{ active: activeItem === 'projects' }" @click="setActive('projects')" title="Projects">
            <i class="pi pi-folder collapsed-icon"></i>
          </div>
        </div>
      </div>

      <!-- Sidebar Footer -->
      <div class="sidebar-footer" v-if="!sidebarCollapsed || isMobile">
        <div class="nav-item" @click="setActive('settings')">
          <i class="pi pi-cog"></i>
          <span>Settings</span>
        </div>
        <div class="nav-item" @click="setActive('templates')">
          <i class="pi pi-clone"></i>
          <span>Templates</span>
        </div>
        <div class="nav-item" @click="setActive('trash')">
          <i class="pi pi-trash"></i>
          <span>Trash</span>
        </div>
        <div class="upgrade-section">
          <Button label="Upgrade plan" size="small" outlined class="w-full" />
        </div>
      </div>

      <!-- Collapsed Footer Icons -->
      <div v-if="sidebarCollapsed && !isMobile" class="collapsed-footer">
        <div class="nav-item" @click="setActive('settings')" title="Settings">
          <i class="pi pi-cog collapsed-icon"></i>
        </div>
        <div class="nav-item" @click="setActive('templates')" title="Templates">
          <i class="pi pi-clone collapsed-icon"></i>
        </div>
        <div class="nav-item" @click="setActive('trash')" title="Trash">
          <i class="pi pi-trash collapsed-icon"></i>
        </div>
      </div>

      <!-- Collapsed sidebar toggle -->
      <Button 
        v-if="sidebarCollapsed && !isMobile"
        icon="pi pi-angle-right" 
        text 
        rounded 
        size="small"
        @click="toggleSidebar"
        class="collapsed-toggle"
      />
    </div>

    <!-- Main Content Area -->
    <div class="layout-main">
      <!-- Top Bar -->
      <div class="layout-topbar">
        <div class="flex align-items-center gap-3">
          <!-- <Button 
            v-if="!isMobile && sidebarCollapsed"
            icon="pi pi-bars" 
            text 
            rounded 
            size="small"
            @click="toggleSidebar"
            class="desktop-menu-btn"
          /> -->
          <!-- <Breadcrumb :model="breadcrumbItems" v-if="!isMobile" /> -->
        </div>
        <div class="flex align-items-center gap-3">
          <!-- Mobile Navigation Pills -->
          <div v-if="isMobile" class="mobile-nav-pills">
            <Button 
              icon="pi pi-bars" 
              text 
              rounded 
              size="small"
              @click="toggleMobileSidebar"
              class="mobile-nav-btn"
              :class="{ 'active': sidebarOpen }"
            />
          </div>
          <div class="theme-toggle">
            <Button 
              :icon="isDarkMode ? 'pi pi-sun' : 'pi pi-moon'" 
              text 
              rounded 
              size="small"
              @click="toggleDarkMode"
              :title="isDarkMode ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
              class="theme-toggle-btn"
            />
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <div class="layout-content">
        <!-- Mobile Active Section Pills -->
        <div v-if="isMobile" class="mobile-active-section">
          <div class="active-pill">
            <i class="pi pi-home text-sm"></i>
            <span class="font-medium ml-1">AlloDitoo</span>
            <i class="pi pi-angle-right text-xs mx-1"></i>
            <span class="font-medium">{{ activePageTitle }}</span>
          </div>
        </div>
        
        <div class="content-container">
          <div class="content-header">
            <h1 class="content-title">{{ activePageTitle }}</h1>
            <div class="content-subtitle" v-if="activeItem === 'getting-started'">
              <i class="pi pi-star text-yellow-500"></i>
              <span>Welcome to your team's workspace!</span>
            </div>
          </div>

          <div class="content-body" v-if="activeItem === 'getting-started'">
            <p class="mb-4">Here are the basics:</p>
            
            <div class="checklist">
              <div class="checklist-item">
                <Checkbox v-model="checkItems.typing" binary />
                <span>Click anywhere and just start typing</span>
              </div>
              <div class="checklist-item">
                <Checkbox v-model="checkItems.content" binary />
                <span>Hit / to see all the types of content you can create - headers, videos, sub pages, etc.</span>
              </div>
              <div class="checklist-item">
                <Checkbox v-model="checkItems.highlight" binary />
                <span>Highlight any text, and use the menu that pops up to <strong>style</strong> <em>your</em> <span class="text-red-500">writing</span> however you like</span>
              </div>
              <div class="checklist-item">
                <Checkbox v-model="checkItems.drag" binary />
                <span>See the ⠿ to the left of this checkbox on hover? Click and drag to move this line</span>
              </div>
              <div class="checklist-item">
                <Checkbox v-model="checkItems.sidebar" binary />
                <span>Explore the sidebar to your left</span>
              </div>
            </div>

            <div class="mt-4">
              <div class="bullet-point">
                <span>Click <strong>All teamspaces</strong> at the top of your sidebar to browse and join teamspaces</span>
              </div>
              <div class="bullet-point">
                <span>Pages in <Badge value="Teamspaces" severity="info" /> are shared with different areas of the company</span>
              </div>
              <div class="bullet-point">
                <span>Pages in <Badge value="Private" severity="secondary" /> are not shared with any other workspace members</span>
              </div>
            </div>

            <div class="mt-4">
              <div class="checklist-item">
                <Checkbox v-model="checkItems.search" binary />
                <span>Click <strong>Search</strong> in your sidebar to find content</span>
              </div>
              <div class="checklist-item">
                <Checkbox v-model="checkItems.calendar" binary />
                <span>Click <strong>Calendar</strong> in your sidebar to manage your time and work together. Notion's calendar app is integrated and synced with all your Google Calendar events!</span>
              </div>
              <div class="checklist-item">
                <Checkbox v-model="checkItems.alanKay" binary />
                <span>Highlight <strong>Alan Kay</strong>'s name and ask AI to explain who Notion's favorite scientist is</span>
              </div>
            </div>

            <div class="toggle-section mt-4">
              <div class="flex align-items-center gap-2 cursor-pointer" @click="toggleExpanded = !toggleExpanded">
                <i :class="toggleExpanded ? 'pi pi-chevron-down' : 'pi pi-chevron-right'"></i>
                <span>This is a toggle block. Click the little triangle to see more useful links!</span>
              </div>
              <div v-if="toggleExpanded" class="mt-2 ml-4">
                <p>Here are some useful resources to get you started!</p>
              </div>
            </div>
          </div>

          <!-- Other page content can be added here based on activeItem -->
          <div v-else class="content-body">
            <Card>
              <template #content>
                <p>Content for {{ activePageTitle }} goes here.</p>
                <p>This is where you would implement the specific functionality for each section.</p>
              </template>
            </Card>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useTheme } from '@primeuix/themes'
import Button from 'primevue/button'
import Breadcrumb from 'primevue/breadcrumb'
import Checkbox from 'primevue/checkbox'
import Card from 'primevue/card'
import Badge from 'primevue/badge'

// Reactive state
const sidebarCollapsed = ref(false)
const sidebarOpen = ref(false)
const activeItem = ref('getting-started')
const toggleExpanded = ref(false)
const isMobile = ref(false)
const isDarkMode = ref(false)

// PrimeVue theme composable
const { switchTheme } = useTheme()

// Checklist items
const checkItems = ref({
  typing: false,
  content: false,
  highlight: false,
  drag: false,
  sidebar: false,
  search: false,
  calendar: false,
  alanKay: false
})

// Mobile detection
const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
  if (!isMobile.value) {
    sidebarOpen.value = false
  }
}

// Dark mode functionality
const initializeDarkMode = () => {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme) {
    isDarkMode.value = savedTheme === 'dark'
  } else {
    isDarkMode.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  applyTheme()
}

const applyTheme = () => {
  // Toggle dark class on document for CSS targeting
  if (isDarkMode.value) {
    document.documentElement.classList.add('dark-theme')
    document.documentElement.setAttribute('data-theme', 'dark')
  } else {
    document.documentElement.classList.remove('dark-theme')
    document.documentElement.setAttribute('data-theme', 'light')
  }
  
  // Use PrimeVue's switchTheme if available
  if (switchTheme) {
    switchTheme(isDarkMode.value ? 'dark' : 'light')
  }
}

const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value
  localStorage.setItem('theme', isDarkMode.value ? 'dark' : 'light')
  applyTheme()
}

onMounted(() => {
  checkMobile()
  initializeDarkMode()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

// Computed properties
const activePageTitle = computed(() => {
  const titles: Record<string, string> = {
    'search': 'Search',
    'home': 'Home',
    'inbox': 'Inbox',
    'notes': '1:1 Notes',
    'getting-started': 'Getting Started',
    'scratchpad': 'Scratchpad',
    'research-workspace': 'AlloDitoo',
    'tasks': 'Tasks',
    'projects': 'Projects',
    'settings': 'Settings',
    'templates': 'Templates',
    'trash': 'Trash'
  }
  return titles[activeItem.value] || 'Dashboard'
})

const breadcrumbItems = computed(() => [
  { label: 'AlloDitoo', icon: 'pi pi-home' },
  { label: activePageTitle.value }
])

// Methods
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

const toggleMobileSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

const closeMobileSidebar = () => {
  sidebarOpen.value = false
}

const setActive = (item: string) => {
  activeItem.value = item
  if (isMobile.value) {
    closeMobileSidebar()
  }
}

const addNew = () => {
  console.log('Add new item')
}

// Get icon for active item
const getActiveIcon = (item: string) => {
  const icons: Record<string, string> = {
    'search': 'pi pi-search',
    'home': 'pi pi-home',
    'inbox': 'pi pi-inbox',
    'notes': 'pi pi-file',
    'getting-started': 'pi pi-star',
    'scratchpad': 'pi pi-pencil',
    'research-workspace': 'pi pi-users',
    'tasks': 'pi pi-check-square',
    'projects': 'pi pi-folder',
    'settings': 'pi pi-cog',
    'templates': 'pi pi-clone',
    'trash': 'pi pi-trash'
  }
  return icons[item] || 'pi pi-circle'
}
</script>

<style scoped>
.layout-wrapper {
  display: flex;
  height: 100vh;
  background: var(--p-surface-ground);
  position: relative;
  color: var(--p-text-color);
  transition: background-color 0.3s ease, color 0.3s ease;
}

/* Mobile Overlay */
.mobile-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 998;
}

/* Mobile Menu Button */
.mobile-menu-btn {
  position: fixed;
  top: 1rem;
  left: 1rem;
  z-index: 1000;
  background: var(--p-surface-card);
  border: 1px solid var(--p-surface-border);
}

.layout-sidebar {
  width: 255px;
  background: #F8F8F7; /* Light mode: warm off-white */
  border-right: 1px solid var(--p-surface-border);
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease, background-color 0.3s ease;
  position: relative;
  color: #212529; /* Light mode text color */
}

/* Dark mode sidebar */
.layout-sidebar-dark {
  background: #050505 !important; /* Dark mode: very dark background */
  color: #ffffff !important; /* Dark mode text color */
}

.layout-sidebar-collapsed {
  width: 70px;
}

/* Mobile Sidebar Styles */
.layout-sidebar-mobile {
  position: fixed;
  top: 0;
  left: -255px;
  height: 100vh;
  z-index: 999;
  transition: left 0.3s ease;
}

.layout-sidebar-mobile-open {
  left: 0;
}

.sidebar-header {
  padding: 1rem;
  border-bottom: 1px solid var(--p-surface-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.toggle-btn {
  opacity: 0.7;
}

.toggle-btn:hover {
  opacity: 1;
}

.collapsed-toggle {
  position: absolute;
  bottom: 1rem;
  left: 50%;
  transform: translateX(-50%);
  background: var(--p-surface-card);
  border: 1px solid var(--p-surface-border);
}

.desktop-menu-btn {
  background: var(--p-surface-card);
  border: 1px solid var(--p-surface-border);
}

.sidebar-content {
  flex: 1;
  padding: 0.5rem;
  overflow-y: auto;
}

.nav-section {
  margin-bottom: 1.5rem;
}

.collapsed-nav-section {
  margin-bottom: 1.5rem;
}

.section-header {
  padding: 0.5rem 0.75rem;
  margin-bottom: 0.25rem;
}

.section-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--p-text-muted-color);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

/* Dark mode section title */
.layout-sidebar-dark .section-title {
  color: #b3b3b3 !important;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1rem; /* Increased padding for more space */
  margin: 0.125rem 0; /* Added margin between items */
  border-radius: 0.5rem; /* Slightly more rounded */
  cursor: pointer;
  transition: background-color 0.2s;
  font-size: 0.875rem;
  justify-content: flex-start;
}

.layout-sidebar-collapsed .nav-item {
  justify-content: center;
  padding: 1rem 0.5rem;
}

.nav-item:hover {
  background: var(--surface-hover);
}

.nav-item.active {
  background: var(--primary-color);
  color: var(--primary-color-text);
}

.nav-item i {
  width: 16px;
  text-align: center;
  font-size: 0.875rem;
}

.collapsed-icon {
  font-size: 1.125rem !important;
  width: 20px !important;
}

.nav-add-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: background-color 0.2s;
  font-size: 0.875rem;
  color: var(--text-color-secondary);
}

.nav-add-item:hover {
  background: var(--surface-hover);
}

.sidebar-footer {
  padding: 0.5rem;
  border-top: 1px solid var(--p-surface-border);
}

.collapsed-footer {
  padding: 0.5rem;
  border-top: 1px solid var(--p-surface-border);
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.upgrade-section {
  margin-top: 1rem;
  padding: 0 0.25rem;
}

.layout-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.layout-topbar {
  padding: 0.75rem 1.5rem;
  border-bottom: 1px solid var(--p-surface-border);
  background: var(--p-surface-card);
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
  width: 100%;
  box-sizing: border-box;
}

.layout-content {
  flex: 1;
  overflow-y: auto;
  padding: 2rem;
  width: 100%;
  box-sizing: border-box;
}

.content-container {
  max-width: 1000px;
  margin: 0 auto;
}

.content-header {
  margin-bottom: 2rem;
}

.content-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin: 0 0 1rem 0;
  color: var(--p-text-color);
}

.content-subtitle {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1rem;
  color: var(--p-text-muted-color);
}

.content-body {
  line-height: 1.6;
}

.checklist {
  margin: 1.5rem 0;
}

.checklist-item {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
  padding: 0.25rem 0;
}

.checklist-item span {
  flex: 1;
}

.bullet-point {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
  padding-left: 1rem;
}

.bullet-point::before {
  content: "•";
  font-weight: bold;
  margin-top: 0.1rem;
}

.toggle-section {
  padding: 0.5rem 0;
}

/* Mobile Responsive Styles */
@media (max-width: 768px) {
  .layout-main {
    margin-left: 0;
    width: 100%;
  }
  
  .layout-topbar {
    padding: 0.75rem 1rem;
    position: relative;
    justify-content: space-between;
  }
  
  .layout-topbar .flex:first-child {
    flex: 1;
    justify-content: flex-start;
  }
  
  .layout-topbar .flex:last-child {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-shrink: 0;
  }
  
  /* Hide desktop breadcrumb on mobile */
  .layout-topbar .flex:first-child .p-breadcrumb {
    display: none;
  }
  
  /* Hide desktop menu button on mobile */
  .desktop-menu-btn {
    display: none;
  }
  
  .content-title {
    font-size: 2rem;
    text-align: left;
  }
  
  .layout-content {
    padding: 1rem;
    width: 100%;
    box-sizing: border-box;
  }
  
  .content-container {
    max-width: 100%;
    width: 100%;
    margin: 0;
    padding: 0;
  }
  
  .hidden-mobile {
    display: none;
  }
  
  .mobile-nav-pills {
    display: flex;
  }
}

@media (min-width: 769px) {
  .mobile-nav-pills {
    display: none;
  }
  
  .mobile-active-section {
    display: none;
  }
}

@media (max-width: 480px) {
  .content-title {
    font-size: 1.75rem;
  }
  
  .layout-topbar {
    padding: 0.5rem 0.75rem;
  }
  
  .layout-content {
    padding: 0.75rem;
  }
  
  .checklist-item {
    flex-direction: column;
    gap: 0.5rem;
    align-items: flex-start;
  }
  
  .content-header {
    margin-bottom: 1.5rem;
  }
  
  .content-body {
    font-size: 0.9rem;
    line-height: 1.5;
  }
  
  .active-pill {
    padding: 0;
    margin-bottom: 1rem;
    font-size: 0.8rem;
  }
  
  .mobile-nav-pills {
    gap: 0.25rem;
  }
}
</style>