<template>
  <div class="account-layout">
    <!-- Account Header expands entire screen width-->
    <div class="account-header-wrapper">
      <!-- Account burger menu (mobile only) -->
      <button class="account-burger icon-btn" @click="toggleAccountMenu" aria-label="Open account menu" v-if="isMobile">
        <span></span>
        <span></span>
        <span></span>
      </button>
      <MainPageHeader title="Account" icon-type="person" icon-color="blue">
        <template #actions>
          <button class="help-btn" @click="showHelp = true">
            <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
              <path
                d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 17h-2v-2h2v2zm2.07-7.75l-.9.92C13.45 12.9 13 13.5 13 15h-2v-.5c0-1.1.45-2.1 1.17-2.83l1.24-1.26c.37-.36.59-.86.59-1.41 0-1.1-.9-2-2-2s-2 .9-2 2H8c0-2.21 1.79-4 4-4s4 1.79 4 4c0 .88-.36 1.68-.93 2.25z" />
            </svg>
            Help
          </button>
        </template>
      </MainPageHeader>
    </div>
    <!-- Flex row: sidebar + main content -->
    <div class="account-shell">
      <div class="account-body-row">
        <!-- Desktop sidebar -->
        <aside class="sidebar" v-if="!isMobile">
          <nav>
            <ul>
              <li>
                <router-link to="/account" active-class="active" exact>
                  <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                    <path
                      d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z" />
                  </svg>
                  Profile
                </router-link>
              </li>
              <li>
                <router-link to="/account/contacts" active-class="active">
                  <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                    <path
                      d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z" />
                  </svg>
                  Contacts
                </router-link>
              </li>
            </ul>
          </nav>
        </aside>
        <!-- Main content area -->
        <div class="main-content" @click="closeAccountMenuOnMobile">
          <section class="account-body">
            <router-view />
          </section>
        </div>
      </div>
    </div>
    <!-- Help Panel emerges from the right of the screen -->
    <div v-if="showHelp" class="help-panel-overlay" @click.self="showHelp = false">
      <div class="help-panel">
        <button class="close-help" @click="showHelp = false">&times;</button>
        <h2>Account Help</h2>
        <p>This is context-sensitive help for the account section. (Stub)</p>
      </div>
    </div>
    <!-- Mobile sidebar overlay -->
    <div v-if="isMobile && accountMenuOpen" class="mobile-overlay" @click="closeAccountMenu"></div>
    <!-- Mobile sidebar -->
    <aside v-if="isMobile && accountMenuOpen" class="sidebar mobile">
      <nav>
        <ul>
          <li>
            <router-link to="/account" active-class="active" exact @click="closeAccountMenu">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path
                  d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z" />
              </svg>
              Profile
            </router-link>
          </li>
          <li>
            <router-link to="/account/contacts" active-class="active" @click="closeAccountMenu">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path
                  d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z" />
              </svg>
              Contacts
            </router-link>
          </li>
        </ul>
      </nav>
    </aside>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import MainPageHeader from './MainPageHeader.vue';

// Responsive: detect mobile
const isMobile = ref(window.innerWidth <= 900);
const accountMenuOpen = ref(false);
const showHelp = ref(false);

function handleResize() {
  isMobile.value = window.innerWidth <= 900;
  if (!isMobile.value) accountMenuOpen.value = false;
}
onMounted(() => {
  window.addEventListener('resize', handleResize);
});
onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
function toggleAccountMenu() {
  accountMenuOpen.value = !accountMenuOpen.value;
}
function closeAccountMenu() {
  accountMenuOpen.value = false;
}
function closeAccountMenuOnMobile() {
  if (isMobile.value && accountMenuOpen.value) {
    accountMenuOpen.value = false;
  }
}
</script>

<style scoped>
.account-layout {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-height: 0;
  height: 100%;
}

.account-shell {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-height: 0;
  height: 100%;
}

.account-header-wrapper {
  width: 100%;
  flex-shrink: 0;
  position: relative;
  z-index: 2998;
}

.account-burger {
  position: absolute;
  left: var(--space-lg);
  top: 50%;
  transform: translateY(-50%);
  z-index: 10;
}

@media (max-width: 900px) {
  .account-header-wrapper :deep(.main-page-header) {
    padding-left: calc(var(--space-xl) + 50px);
  }

  .account-burger {
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

.account-burger span {
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

  /* Prevent main content from being affected by mobile sidebar */
  .account-body-row {
    position: relative;
  }

  .main-content {
    position: relative;
    z-index: 1;
  }

  .account-body {
    padding: var(--space-md) var(--space-sm);
  }

  .sidebar.mobile {
    width: 280px;
    padding: var(--space-lg) var(--space-md) var(--space-lg) var(--space-md);
  }

  .sidebar ul {
    font-size: 1rem;
  }
}

@media (max-width: 600px) {
  .account-burger {
    width: 32px;
    height: 32px;
  }

  .account-burger span {
    width: 22px;
    height: 3px;
    margin: 2px 0;
  }

  .account-body {
    padding: var(--space-md) var(--space-sm);
    font-size: 1rem;
  }

  .sidebar,
  .sidebar.mobile {
    width: 180px;
    padding: var(--space-md) var(--space-sm) var(--space-md) var(--space-sm);
  }

  .sidebar ul {
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

.account-body-row {
  display: flex;
  flex: 1;
  min-height: 0;
  height: 100%;
  align-items: stretch;
}

.sidebar {
  width: 200px;
  background: var(--color-bg-alt);
  border-right: 1px solid var(--color-border);
  padding: var(--space-lg) var(--space-md) var(--space-lg) var(--space-md);
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

.sidebar ul {
  list-style: none;
}

.sidebar li {
  margin-bottom: var(--space-md);
}

.sidebar a {
  color: var(--color-text);
  text-decoration: none;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.nav-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.sidebar a.active {
  color: var(--color-primary);
}

.sidebar a.disabled {
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

.account-body {
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
