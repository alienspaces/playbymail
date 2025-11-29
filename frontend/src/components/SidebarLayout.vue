<!--
  SidebarLayout.vue
  Shared layout component with header, sidebar navigation, and main content area.
  Used by Studio, Account, and Management sections.
-->
<template>
  <div class="sidebar-layout">
    <!-- Show entry view for unauthenticated users if provided -->
    <slot v-if="!isLoggedIn && hasEntrySlot" name="entry"></slot>

    <!-- Show full interface for authenticated users (or if no entry slot) -->
    <div v-else class="layout-shell">
      <!-- Header expands entire screen width -->
      <div class="layout-header-wrapper">
        <!-- Burger menu (mobile only) -->
        <button 
          class="layout-burger icon-btn" 
          @click="toggleMenu" 
          aria-label="Open menu" 
          v-if="isMobile"
        >
          <span></span>
          <span></span>
          <span></span>
        </button>
        <MainPageHeader :title="title" :icon-type="iconType" :icon-color="iconColor">
          <template #actions>
            <slot name="header-actions">
              <button class="help-btn" @click="showHelp = true">
                <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                  <path
                    d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 17h-2v-2h2v2zm2.07-7.75l-.9.92C13.45 12.9 13 13.5 13 15h-2v-.5c0-1.1.45-2.1 1.17-2.83l1.24-1.26c.37-.36.59-.86.59-1.41 0-1.1-.9-2-2-2s-2 .9-2 2H8c0-2.21 1.79-4 4-4s4 1.79 4 4c0 .88-.36 1.68-.93 2.25z" />
                </svg>
                Help
              </button>
            </slot>
          </template>
        </MainPageHeader>
      </div>

      <!-- Flex row: sidebar + main content -->
      <div class="layout-body-row">
        <!-- Desktop sidebar -->
        <aside class="sidebar" v-if="!isMobile">
          <nav>
            <slot name="sidebar"></slot>
          </nav>
        </aside>

        <!-- Main content area -->
        <div class="main-content" @click="closeMenuOnMobile">
          <section class="layout-body">
            <slot></slot>
          </section>
        </div>
      </div>

      <!-- Help Panel emerges from the right of the screen -->
      <div v-if="showHelp" class="help-panel-overlay" @click.self="showHelp = false">
        <div class="help-panel">
          <button class="close-help" @click="showHelp = false">&times;</button>
          <slot name="help">
            <h2>Help</h2>
            <p>Context-sensitive help for this section.</p>
          </slot>
        </div>
      </div>

      <!-- Mobile sidebar overlay -->
      <div v-if="isMobile && menuOpen" class="mobile-overlay" @click="closeMenu"></div>

      <!-- Mobile sidebar -->
      <aside v-if="isMobile && menuOpen" class="sidebar mobile">
        <nav @click="closeMenu">
          <slot name="sidebar"></slot>
        </nav>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, useSlots } from 'vue';
import { storeToRefs } from 'pinia';
import { useAuthStore } from '../stores/auth';
import MainPageHeader from './MainPageHeader.vue';

defineProps({
  title: {
    type: String,
    required: true
  },
  iconType: {
    type: String,
    default: 'default'
  },
  iconColor: {
    type: String,
    default: 'blue'
  }
});

const slots = useSlots();
const authStore = useAuthStore();
const { sessionToken } = storeToRefs(authStore);

const isLoggedIn = computed(() => !!sessionToken.value);
const hasEntrySlot = computed(() => !!slots.entry);

const showHelp = ref(false);

// Responsive: detect mobile
const isMobile = ref(window.innerWidth <= 900);
const menuOpen = ref(false);

function handleResize() {
  isMobile.value = window.innerWidth <= 900;
  if (!isMobile.value) menuOpen.value = false;
}

onMounted(() => {
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});

function toggleMenu() {
  menuOpen.value = !menuOpen.value;
}

function closeMenu() {
  menuOpen.value = false;
}

function closeMenuOnMobile() {
  if (isMobile.value && menuOpen.value) {
    menuOpen.value = false;
  }
}
</script>

<style scoped>
.sidebar-layout {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-height: 0;
  height: 100%;
}

.layout-shell {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-height: 0;
  height: 100%;
}

.layout-header-wrapper {
  width: 100%;
  flex-shrink: 0;
  position: relative;
  z-index: 2998;
}

.layout-burger {
  position: absolute;
  left: var(--space-lg);
  top: 50%;
  transform: translateY(-50%);
  z-index: 10;
}

@media (max-width: 900px) {
  .layout-header-wrapper :deep(.main-page-header) {
    padding-left: calc(var(--space-xl) + 50px);
  }

  .layout-burger {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    width: 40px;
    height: 40px;
    margin-right: var(--space-lg);
    z-index: 3002;
  }
}

.layout-burger span {
  display: block;
  width: 28px;
  height: 4px;
  margin: 3px 0;
  background: var(--color-primary);
  border-radius: 2px;
  transition: 0.3s;
}

@media (max-width: 900px) {
  .sidebar {
    display: none;
  }

  .mobile-overlay {
    position: fixed;
    top: 70px;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.3);
    z-index: 2999;
  }

  .sidebar.mobile {
    display: flex;
    position: fixed;
    top: 70px;
    left: 0;
    height: calc(100vh - 70px);
    width: 280px;
    background: var(--color-bg-alt);
    border-right: 1px solid var(--color-border);
    padding: var(--space-lg) var(--space-md) 0 var(--space-md);
    flex-direction: column;
    z-index: 3000;
    box-shadow: 2px 0 16px rgba(0, 0, 0, 0.12);
    animation: slideInLeft 0.2s;
  }

  .layout-body-row {
    position: relative;
  }

  .main-content {
    position: relative;
    z-index: 1;
  }

  .layout-body {
    padding: var(--space-md) var(--space-sm);
  }

  .sidebar.mobile {
    width: 280px;
    padding: var(--space-lg) var(--space-md) 0 var(--space-md);
  }

  .sidebar :deep(ul) {
    font-size: 1rem;
  }
}

@media (max-width: 600px) {
  .layout-burger {
    width: 32px;
    height: 32px;
  }

  .layout-burger span {
    width: 22px;
    height: 3px;
    margin: 2px 0;
  }

  .layout-body {
    padding: var(--space-md) var(--space-sm);
    font-size: 1rem;
  }

  .sidebar,
  .sidebar.mobile {
    width: 180px;
    padding: var(--space-md) var(--space-sm) 0 var(--space-sm);
  }

  .sidebar :deep(ul) {
    font-size: 1rem;
  }

  .help-panel {
    min-width: 220px;
    padding: var(--space-md) var(--space-md);
    font-size: 1rem;
  }
}

@keyframes slideInLeft {
  from {
    transform: translateX(-100%);
  }

  to {
    transform: translateX(0);
  }
}

.layout-body-row {
  display: flex;
  flex: 1;
  min-height: 0;
  height: 100%;
  align-items: stretch;
  min-height: calc(100vh - 140px);
}

.sidebar {
  width: 200px;
  background: var(--color-bg-alt);
  border-right: 1px solid var(--color-border);
  padding: var(--space-lg) var(--space-md) 0 var(--space-md);
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-self: stretch;
  position: relative;
  min-height: calc(100vh - 140px);
}

.sidebar::after {
  content: "";
  position: absolute;
  left: 0;
  right: 0;
  top: 100%;
  background: var(--color-bg-alt);
  border-right: 1px solid var(--color-border);
  height: calc(var(--space-xl) + 70px);
  pointer-events: none;
  z-index: 0;
}

.sidebar :deep(ul) {
  list-style: none;
}

.sidebar :deep(li) {
  margin-bottom: var(--space-md);
}

.sidebar :deep(a) {
  color: var(--color-text);
  text-decoration: none;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.sidebar :deep(.nav-icon) {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.sidebar :deep(a.active) {
  color: var(--color-primary);
}

.sidebar :deep(a.disabled) {
  pointer-events: none;
  color: #aaa;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  background: transparent;
}

.layout-body {
  flex: 1;
  padding: var(--layout-content-padding);
  min-width: 0;
}

.help-panel-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.25);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.help-panel {
  background: var(--color-bg);
  padding: var(--space-lg) var(--space-xl);
  border-radius: var(--radius-lg);
  min-width: 320px;
  max-width: 90vw;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.18);
  position: relative;
}

.close-help {
  position: absolute;
  top: var(--space-md);
  right: var(--space-md);
  background: none;
  border: none;
  font-size: 2rem;
  color: #888;
  cursor: pointer;
}
</style>

<style>
/* Hide layout burger when global mobile menu is open */
.mobile-menu+* .layout-burger,
.mobile-menu~.layout-burger,
.mobile-menu .layout-burger,
.mobile-menu~* .layout-burger {
  display: none !important;
}
</style>

