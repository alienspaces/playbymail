<!--
  ManagementLayout.vue
  Layout component for the Game Management interface.
  Uses the shared SidebarLayout component.
-->
<template>
  <SidebarLayout title="Game Management" icon-type="clipboard" icon-color="blue">
    <!-- Entry view for unauthenticated users -->
    <template #entry>
      <ManagementEntryView />
    </template>

    <!-- Sidebar navigation -->
    <template #sidebar>
      <ul>
        <li>
          <router-link to="/admin" active-class="active" exact>
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
            <router-link :to="`/admin/games/${selectedGame.id}/instances`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M4 6H2v14c0 1.1.9 2 2 2h14v-2H4V6zm16-4H8c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm-1 9h-4v4h-2v-4H9V9h4V5h2v4h4v2z" />
              </svg>
              Instances
            </router-link>
          </li>
          <li>
            <router-link :to="`/admin/games/${selectedGame.id}/turn-sheets`" active-class="active">
              <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z" />
              </svg>
              Turn Sheets
            </router-link>
          </li>
        </ul>
      </div>
    </template>

    <!-- Help content -->
    <template #help>
      <h2>Game Management Help</h2>
      <h3>Managing Game Instances</h3>
      <ul>
        <li><strong>Games:</strong> View all games you have access to and select one to manage</li>
        <li><strong>Instances:</strong> Create and manage game sessions for players</li>
        <li><strong>Turn Sheets:</strong> Download join forms and upload scanned submissions</li>
      </ul>

      <h3>Game Context</h3>
      <p>Select a game from the Games list to see game-specific options in the sidebar.</p>

      <h3>Game Types</h3>
      <ul>
        <li><strong>Adventure Games:</strong> Exploration and story-driven experiences</li>
        <li><strong>Future Types:</strong> Economic, sports, and other game types coming soon</li>
      </ul>
    </template>

    <!-- Main content -->
    <router-view />
  </SidebarLayout>
</template>

<script setup>
import { storeToRefs } from 'pinia';
import { useGamesStore } from '../stores/games';
import SidebarLayout from './SidebarLayout.vue';
import ManagementEntryView from '../views/ManagementEntryView.vue';

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
