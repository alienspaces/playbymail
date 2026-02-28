<!--
  GameIssuesPanel.vue
  Compact sidebar panel showing game validation issues.
  Auto-fetches validation when the selected game changes.
-->
<template>
  <div class="issues-panel">
    <div class="issues-header">
      <span class="issues-label">Game Readiness</span>
      <button class="refresh-btn" @click="refresh" :disabled="gamesStore.validationLoading" aria-label="Refresh validation">
        <svg class="refresh-icon" :class="{ spinning: gamesStore.validationLoading }" viewBox="0 0 24 24"
          fill="currentColor">
          <path
            d="M17.65 6.35A7.958 7.958 0 0012 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0112 18c-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z" />
        </svg>
      </button>
    </div>

    <div v-if="gamesStore.validationLoading && !issues.length" class="issues-loading">
      Checking...
    </div>

    <div v-else-if="gamesStore.validationError" class="issues-error">
      Failed to check
    </div>

    <div v-else-if="issues.length === 0" class="issues-valid">
      <svg class="status-icon valid" viewBox="0 0 24 24" fill="currentColor">
        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" />
      </svg>
      <span>Ready to create instance</span>
    </div>

    <ul v-else class="issues-list">
      <li v-for="(issue, idx) in issues" :key="idx" class="issue-item" :class="issue.severity">
        <svg v-if="issue.severity === 'error'" class="status-icon error" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z" />
        </svg>
        <svg v-else class="status-icon warning" viewBox="0 0 24 24" fill="currentColor">
          <path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z" />
        </svg>
        <router-link v-if="issueLink(issue)" :to="issueLink(issue)" class="issue-text">
          {{ issue.message }}
        </router-link>
        <span v-else class="issue-text">{{ issue.message }}</span>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { computed, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { useGamesStore } from '../stores/games';

const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const issues = computed(() => gamesStore.validationIssues);

const fieldToRoute = {
  locations: 'locations',
  starting_location: 'locations',
  location_links: 'location-links',
  items: 'items',
  item_placements: 'item-placements',
  creatures: 'creatures',
  creature_placements: 'creature-placements',
};

function issueLink(issue) {
  if (!selectedGame.value || !issue.field) return null;
  const route = fieldToRoute[issue.field];
  if (!route) return null;
  return `/studio/${selectedGame.value.id}/${route}`;
}

function refresh() {
  if (selectedGame.value) {
    gamesStore.fetchValidation(selectedGame.value.id);
  }
}

watch(
  () => selectedGame.value?.id,
  (newId) => {
    if (newId) {
      gamesStore.fetchValidation(newId);
    } else {
      gamesStore.clearValidation();
    }
  },
  { immediate: true }
);
</script>

<style scoped>
.issues-panel {
  margin-top: var(--space-md);
  padding-top: var(--space-md);
  border-top: 1px solid var(--color-border);
}

.issues-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-sm);
}

.issues-label {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.refresh-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 2px;
  color: var(--color-text-muted);
  display: flex;
  align-items: center;
}

.refresh-btn:hover {
  color: var(--color-primary);
}

.refresh-btn:disabled {
  cursor: default;
  opacity: 0.5;
}

.refresh-icon {
  width: 14px;
  height: 14px;
}

.spinning {
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.issues-loading,
.issues-error {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
  padding: var(--space-xs) 0;
}

.issues-error {
  color: var(--color-danger);
}

.issues-valid {
  display: flex;
  align-items: center;
  gap: var(--space-xs);
  font-size: var(--font-size-sm);
  color: var(--color-success);
  padding: var(--space-xs) 0;
}

.issues-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.issue-item {
  display: flex;
  align-items: flex-start;
  gap: var(--space-xs);
  padding: var(--space-xs) 0;
  font-size: var(--font-size-sm);
  line-height: 1.3;
}

.status-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  margin-top: 1px;
}

.status-icon.valid {
  color: var(--color-success);
}

.status-icon.error {
  color: var(--color-danger);
}

.status-icon.warning {
  color: var(--color-warning);
}

.issue-text {
  color: var(--color-text);
}

a.issue-text {
  color: var(--color-primary);
  text-decoration: none;
}

a.issue-text:hover {
  text-decoration: underline;
}
</style>
