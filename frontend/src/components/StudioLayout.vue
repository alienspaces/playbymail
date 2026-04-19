<!--
  StudioLayout.vue
  Layout component for the Game Designer Studio interface.
  Uses the shared SidebarLayout component.
-->
<template>
  <SidebarLayout title="Game Designer Studio" icon-type="pencil" icon-color="blue">
    <!-- Entry view for unauthenticated users -->
    <template #entry>
      <StudioEntryView />
    </template>

    <!-- Sidebar navigation -->
    <template #sidebar>
      <ul>
        <li>
          <router-link to="/studio" active-class="active">
            <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
            </svg>
            Games
          </router-link>
        </li>
      </ul>

      <!-- Game Context Section -->
      <div v-if="selectedGame" class="game-context">
        <div class="context-header">
          <span class="context-label">Selected Game</span>
          <span class="context-name">{{ selectedGame.name }}</span>
        </div>

        <!-- Common links shown for all game types -->
        <ul>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/turn-sheet-backgrounds`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path
                  d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5V5h14v14zm-5.04-6.71l-2.75 3.54-1.96-2.36L6.5 17h11l-3.54-4.71z" />
              </svg>
              Turn Sheets
            </router-link>
          </li>
        </ul>

        <!-- Adventure game specific links -->
        <ul v-if="isAdventureGame">
          <li>
            <router-link :to="`/studio/${selectedGame.id}/locations`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path
                  d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z" />
              </svg>
              Locations
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/location-links`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
                <path d="M7 10l5 3 5-3" />
              </svg>
              Location Links
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/location-link-requirements`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z" />
              </svg>
              Link Requirements
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/items`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path
                  d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
              </svg>
              Items
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/item-placements`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path
                  d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
                <circle cx="12" cy="12" r="3" />
              </svg>
              Item Placements
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/item-effects`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z" />
              </svg>
              Item Effects
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/creatures`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path
                  d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" />
              </svg>
              Creatures
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/creature-placements`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path
                  d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" />
                <circle cx="12" cy="12" r="3" />
              </svg>
              Creature Placements
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/location-objects`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 14l-5-5 1.41-1.41L12 14.17l7.59-7.59L21 8l-9 9z" />
              </svg>
              Location Objects
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/location-object-effects`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z" />
              </svg>
              Object Effects
            </router-link>
          </li>
        </ul>

        <!-- MechaGame specific links -->
        <ul v-if="isMechaGame">
          <li>
            <router-link :to="`/studio/${selectedGame.id}/chassis`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2L2 12h3v8h6v-5h2v5h6v-8h3L12 2z" />
              </svg>
              Chassis
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/weapons`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M22 9V7h-2V5c0-1.1-.9-2-2-2H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-2h2v-2h-2v-2h2v-2h-2V9h2zm-4 10H4V5h14v14z" />
              </svg>
              Weapons
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/equipment`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M22.7 19l-9.1-9.1c.9-2.3.4-5-1.5-6.9-2-2-5-2.4-7.4-1.3L9 6 6 9 1.6 4.7C.4 7.1.9 10.1 2.9 12.1c1.9 1.9 4.6 2.4 6.9 1.5l9.1 9.1c.4.4 1 .4 1.4 0l2.3-2.3c.5-.4.5-1.1.1-1.4z" />
              </svg>
              Equipment
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/sectors`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" />
              </svg>
              Sectors
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/sector-links`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
              </svg>
              Sector Links
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/squads`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5z" />
              </svg>
              Squads
            </router-link>
          </li>
          <li>
            <router-link :to="`/studio/${selectedGame.id}/computer-opponents`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M21 3H3c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h5v2H6v2h12v-2h-2v-2h5c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 14H3V5h18v12z" />
              </svg>
              Computer Opponents
            </router-link>
          </li>
        </ul>

        <GameIssuesPanel />
      </div>
    </template>

    <!-- Help content -->
    <template #help>
      <h2>Studio Help</h2>
      <p>This is context-sensitive help for the current section. (Stub)</p>
    </template>

    <!-- Main content -->
    <router-view />
  </SidebarLayout>
</template>

<script setup>
import { computed } from 'vue';
import { storeToRefs } from 'pinia';
import { useGamesStore } from '../stores/games';
import SidebarLayout from './SidebarLayout.vue';
import StudioEntryView from '../views/StudioEntryView.vue';
import GameIssuesPanel from './GameIssuesPanel.vue';

const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const isAdventureGame = computed(() => selectedGame.value?.game_type === 'adventure')
const isMechaGame = computed(() => selectedGame.value?.game_type === 'mecha')
</script>

<style scoped>
.game-context {
  margin-top: var(--space-md);
  padding-top: var(--space-md);
  border-top: 1px solid var(--color-border);
}

.context-header {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
  margin-bottom: var(--space-md);
  padding: var(--space-sm);
  background: var(--color-bg-light);
  border-radius: var(--radius-sm);
}

.context-label {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.context-name {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.game-context ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.game-context li {
  margin-bottom: var(--menu-item-spacing);
}

.game-context a {
  color: var(--color-text);
  text-decoration: none;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.game-context a.active {
  color: var(--color-primary);
}

.nav-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}
</style>
