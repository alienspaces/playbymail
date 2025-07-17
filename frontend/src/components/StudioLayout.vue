<template>
  <div class="studio-layout">
    <!-- Studio Header expands entire screen width-->
    <div class="studio-header">
      <!-- Studio burger menu (mobile only) -->
      <button
        class="studio-burger"
        @click="toggleStudioMenu"
        aria-label="Open studio menu"
        v-if="isMobile"
      >
        <span></span>
        <span></span>
        <span></span>
      </button>
      <h1>Game Designer Studio</h1>
    </div>
    <!-- Flex row: sidebar + main content -->
    <div class="studio-body-row">
      <!-- Sidebar left of screen-->
      <aside class="sidebar" v-if="!isMobile || studioMenuOpen">
        <nav>
          <ul>
            <li><router-link to="/studio" active-class="active">Games</router-link></li>
            <template v-if="showAdventureMenu">
              <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/locations` : '#'" :class="{disabled: !selectedGame}" active-class="active">Locations</router-link></li>
              <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/items` : '#'" :class="{disabled: !selectedGame}" active-class="active">Items</router-link></li>
              <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/creatures` : '#'" :class="{disabled: !selectedGame}" active-class="active">Creatures</router-link></li>
              <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/placement` : '#'" :class="{disabled: !selectedGame}" active-class="active">Placement</router-link></li>
            </template>
          </ul>
        </nav>
        <button class="help-btn" @click="showHelp = true">Help</button>
      </aside>
      <!-- Main content area right of sidebar-->
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
    <!-- Studio menu overlay for mobile -->
    <div v-if="isMobile && studioMenuOpen" class="studio-menu-overlay" @click.self="closeStudioMenu">
      <aside class="sidebar mobile">
        <nav>
          <ul>
            <li><router-link to="/studio" active-class="active" @click="closeStudioMenu">Games</router-link></li>
            <template v-if="showAdventureMenu">
              <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/locations` : '#'" :class="{disabled: !selectedGame}" active-class="active" @click="closeStudioMenu">Locations</router-link></li>
              <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/items` : '#'" :class="{disabled: !selectedGame}" active-class="active" @click="closeStudioMenu">Items</router-link></li>
              <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/creatures` : '#'" :class="{disabled: !selectedGame}" active-class="active" @click="closeStudioMenu">Creatures</router-link></li>
              <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/placement` : '#'" :class="{disabled: !selectedGame}" active-class="active" @click="closeStudioMenu">Placement</router-link></li>
            </template>
          </ul>
        </nav>
        <button class="help-btn" @click="showHelp = true; closeStudioMenu()">Help</button>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useGamesStore } from '../stores/games';
import { useRoute } from 'vue-router';

const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);
const showHelp = ref(false);
const route = useRoute();

// Show adventure menu only on /studio/:gameId/*
const showAdventureMenu = computed(() => {
  return (
    selectedGame.value &&
    selectedGame.value.game_type === 'adventure' &&
    /^\/studio\/[^/]+\//.test(route.path)
  );
});

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
  min-height: 100vh;
}
.studio-header {
  width: 100%;
  padding: 1.5rem 2rem 1rem 2rem;
  border-bottom: 1px solid #eee;
  background: #fff;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  position: relative;
}
.studio-header h1 {
  margin: 0 auto 0 0;
  font-size: 2rem;
}
.studio-burger {
  display: none;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 40px;
  height: 40px;
  background: none;
  border: none;
  cursor: pointer;
  margin-right: 1.5rem;
  z-index: 2100;
}
.studio-burger span {
  display: block;
  width: 28px;
  height: 4px;
  margin: 3px 0;
  background: #1976d2;
  border-radius: 2px;
  transition: 0.3s;
}
@media (max-width: 900px) {
  .studio-burger {
    display: flex;
  }
  .sidebar {
    display: none;
  }
  .sidebar.mobile {
    display: flex;
    position: fixed;
    top: 0;
    left: 0;
    height: 100vh;
    width: 240px;
    background: #f7f7f7;
    border-right: 1px solid #ddd;
    padding: 2rem 1rem 0 1rem;
    flex-direction: column;
    z-index: 2200;
    box-shadow: 2px 0 16px rgba(0,0,0,0.12);
    animation: slideInLeft 0.2s;
  }
  .studio-menu-overlay {
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    background: rgba(0,0,0,0.18);
    z-index: 2199;
    display: flex;
  }
}
@media (max-width: 600px) {
  .studio-header {
    padding: 1rem 1rem 0.5rem 1rem;
  }
  .studio-header h1 {
    font-size: 1.2rem;
    line-height: 1.2;
    margin: 0;
  }
  .studio-burger {
    width: 32px;
    height: 32px;
    margin-right: 1rem;
  }
  .studio-burger span {
    width: 22px;
    height: 3px;
    margin: 2px 0;
  }
  .studio-body {
    padding: 1rem 0.5rem;
    font-size: 1rem;
  }
  .sidebar, .sidebar.mobile {
    width: 180px;
    padding: 1rem 0.5rem 0 0.5rem;
  }
  .sidebar ul {
    font-size: 1rem;
  }
  .help-btn {
    font-size: 0.95rem;
    padding: 0.6rem 0;
  }
  .help-panel {
    min-width: 220px;
    padding: 1rem 1.2rem;
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
}
.sidebar {
  width: 200px;
  background: #f7f7f7;
  border-right: 1px solid #ddd;
  padding: 2rem 1rem 0 1rem;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
}
.sidebar ul {
  list-style: none;
}
.sidebar li {
  margin-bottom: 1rem;
}
.sidebar a {
  color: #222;
  text-decoration: none;
  font-weight: 500;
}
.sidebar a.active {
  color: #1976d2;
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
}
.studio-body {
  flex: 1;
  padding: 2rem;
  background: #fafbfc;
  min-width: 0;
}
.help-btn {
  margin-top: 2rem;
  width: 100%;
  padding: 0.75rem 0;
  background: #1976d2;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-weight: 600;
  cursor: pointer;
  font-size: 1rem;
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
  background: #fff;
  padding: 2rem 2.5rem;
  border-radius: 10px;
  min-width: 320px;
  max-width: 90vw;
  box-shadow: 0 2px 16px rgba(0,0,0,0.18);
  position: relative;
}
.close-help {
  position: absolute;
  top: 1rem;
  right: 1rem;
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