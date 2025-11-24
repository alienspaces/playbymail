<template>
  <div class="studio-layout">
    <!-- Show StudioEntryView for unauthenticated users -->
    <StudioEntryView v-if="!isLoggedIn" />
    
    <!-- Show full studio interface for authenticated users -->
    <div v-else class="studio-shell">
      <!-- Studio Header expands entire screen width-->
      <div class="studio-header-wrapper">
        <!-- Studio burger menu (mobile only) -->
        <button
          class="studio-burger icon-btn"
          @click="toggleStudioMenu"
          aria-label="Open studio menu"
          v-if="isMobile"
        >
          <span></span>
          <span></span>
          <span></span>
        </button>
        <MainPageHeader 
          title="Game Designer Studio" 
          icon-type="pencil" 
          icon-color="blue"
        >
          <template #actions>
            <button class="help-btn" @click="showHelp = true">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 17h-2v-2h2v2zm2.07-7.75l-.9.92C13.45 12.9 13 13.5 13 15h-2v-.5c0-1.1.45-2.1 1.17-2.83l1.24-1.26c.37-.36.59-.86.59-1.41 0-1.1-.9-2-2-2s-2 .9-2 2H8c0-2.21 1.79-4 4-4s4 1.79 4 4c0 .88-.36 1.68-.93 2.25z"/>
              </svg>
              Help
            </button>
          </template>
        </MainPageHeader>
      </div>
      <!-- Flex row: sidebar + main content -->
      <div class="studio-body-row">
        <!-- Desktop sidebar -->
        <aside class="sidebar" v-if="!isMobile">
          <nav>
            <ul>
              <li><router-link to="/studio" active-class="active">
                <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
                </svg>
                Games
              </router-link></li>
              <template v-if="selectedGame">
                <li><router-link :to="`/studio/${selectedGame.id}/locations`" active-class="active">
                  <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/>
                  </svg>
                  Locations
                </router-link></li>
                <li><router-link :to="`/studio/${selectedGame.id}/location-links`" active-class="active">
                  <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
                    <path d="M7 10l5 3 5-3"/>
                  </svg>
                  Location Links
                </router-link></li>
                <li><router-link :to="`/studio/${selectedGame.id}/items`" active-class="active">
                  <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                  </svg>
                  Items
                </router-link></li>
                <li><router-link :to="`/studio/${selectedGame.id}/item-placements`" active-class="active">
                  <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                    <circle cx="12" cy="12" r="3"/>
                  </svg>
                  Item Placements
                </router-link></li>
                <li><router-link :to="`/studio/${selectedGame.id}/creatures`" active-class="active">
                  <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                  </svg>
                  Creatures
                </router-link></li>
                <li><router-link :to="`/studio/${selectedGame.id}/creature-placements`" active-class="active">
                  <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                    <circle cx="12" cy="12" r="3"/>
                  </svg>
                  Creature Placements
                </router-link></li>

              </template>
            </ul>
          </nav>
        </aside>
        <!-- Main content area -->
        <div class="main-content" @click="closeStudioMenuOnMobile">
          <section class="studio-body">
            <router-view />
          </section>
        </div>
      </div>
    <!-- Help Panel emerges from the right of the screen -->
    <div v-if="showHelp" class="help-panel-overlay" @click.self="showHelp = false">
      <div class="help-panel">
        <button class="close-help" @click="showHelp = false">&times;</button>
        <h2>Studio Help</h2>
        <p>This is context-sensitive help for the current section. (Stub)</p>
      </div>
    </div>
    <!-- Mobile sidebar overlay -->
    <div v-if="isMobile && studioMenuOpen" class="mobile-overlay" @click="closeStudioMenu"></div>
    <!-- Mobile sidebar -->
    <aside v-if="isMobile && studioMenuOpen" class="sidebar mobile">
      <nav>
        <ul>
          <li><router-link to="/studio" active-class="active" @click="closeStudioMenu">
            <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
            </svg>
            Games
          </router-link></li>
          <template v-if="selectedGame">
            <li><router-link :to="`/studio/${selectedGame.id}/locations`" active-class="active" @click="closeStudioMenu">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/>
              </svg>
              Locations
            </router-link></li>
            <li><router-link :to="`/studio/${selectedGame.id}/location-links`" active-class="active" @click="closeStudioMenu">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
                <path d="M7 10l5 3 5-3"/>
              </svg>
              Location Links
            </router-link></li>
            <li><router-link :to="`/studio/${selectedGame.id}/items`" active-class="active" @click="closeStudioMenu">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
              </svg>
              Items
            </router-link></li>
            <li><router-link :to="`/studio/${selectedGame.id}/item-placements`" active-class="active" @click="closeStudioMenu">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
              Item Placements
            </router-link></li>
            <li><router-link :to="`/studio/${selectedGame.id}/creatures`" active-class="active" @click="closeStudioMenu">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
              </svg>
              Creatures
            </router-link></li>
            <li><router-link :to="`/studio/${selectedGame.id}/creature-placements`" active-class="active" @click="closeStudioMenu">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
              Creature Placements
            </router-link></li>

          </template>
        </ul>
      </nav>
    </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { storeToRefs } from 'pinia';
import { useGamesStore } from '../stores/games';
import { useAuthStore } from '../stores/auth';
import StudioEntryView from '../views/StudioEntryView.vue';
import MainPageHeader from './MainPageHeader.vue';

const gamesStore = useGamesStore();
const authStore = useAuthStore();
const { selectedGame } = storeToRefs(gamesStore);
const { sessionToken } = storeToRefs(authStore);
const showHelp = ref(false);

const isLoggedIn = computed(() => !!sessionToken.value);

// Responsive: detect mobile
const isMobile = ref(window.innerWidth <= 900);
const studioMenuOpen = ref(false);

function handleResize() {
  isMobile.value = window.innerWidth <= 900;
  if (!isMobile.value) studioMenuOpen.value = false;
}
onMounted(() => {
  window.addEventListener('resize', handleResize);
});
onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
function toggleStudioMenu() {
  studioMenuOpen.value = !studioMenuOpen.value;
}
function closeStudioMenu() {
  studioMenuOpen.value = false;
}
function closeStudioMenuOnMobile() {
  if (isMobile.value && studioMenuOpen.value) {
    studioMenuOpen.value = false;
  }
}
</script>

<style scoped>
.studio-layout {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-height: 0;
  height: 100%;
}
.studio-shell {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-height: 0;
  height: 100%;
}
.studio-header-wrapper {
  width: 100%;
  flex-shrink: 0;
  position: relative;
  z-index: 2998;
}

.studio-burger {
  position: absolute;
  left: var(--space-lg);
  top: 50%;
  transform: translateY(-50%);
  z-index: 10;
}

@media (max-width: 900px) {
  .studio-header-wrapper :deep(.main-page-header) {
    padding-left: calc(var(--space-xl) + 50px);
  }
  .studio-burger {
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
.studio-burger span {
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
    box-shadow: 2px 0 16px rgba(0,0,0,0.12);
    animation: slideInLeft 0.2s;
  }
  /* Prevent main content from being affected by mobile sidebar */
  .studio-body-row {
    position: relative;
  }
  .main-content {
    position: relative;
    z-index: 1;
  }
  .studio-body {
    padding: var(--space-md) var(--space-sm);
  }
  .sidebar.mobile {
    width: 280px;
    padding: var(--space-lg) var(--space-md) 0 var(--space-md);
  }
  .sidebar ul {
    font-size: 1rem;
  }
}

@media (max-width: 600px) {
  .studio-burger {
    width: 32px;
    height: 32px;
  }
  .studio-burger span {
    width: 22px;
    height: 3px;
    margin: 2px 0;
  }
  .studio-body {
    padding: var(--space-md) var(--space-sm);
    font-size: 1rem;
  }
  .sidebar, .sidebar.mobile {
    width: 180px;
    padding: var(--space-md) var(--space-sm) 0 var(--space-sm);
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
  from { transform: translateX(-100%); }
  to { transform: translateX(0); }
}
.studio-body-row {
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
.studio-body {
  flex: 1;
  padding: var(--layout-content-padding);
  min-width: 0;
}
.help-panel-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.25);
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
  box-shadow: 0 2px 16px rgba(0,0,0,0.18);
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
/* Hide studio burger when global mobile menu is open */
.mobile-menu + * .studio-burger,
.mobile-menu ~ .studio-burger,
.mobile-menu .studio-burger,
.mobile-menu ~ * .studio-burger {
  display: none !important;
}
</style> 