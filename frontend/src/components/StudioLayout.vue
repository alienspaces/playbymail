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
        </ul>
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
import { storeToRefs } from 'pinia';
import { useGamesStore } from '../stores/games';
import SidebarLayout from './SidebarLayout.vue';
import StudioEntryView from '../views/StudioEntryView.vue';

const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);
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
