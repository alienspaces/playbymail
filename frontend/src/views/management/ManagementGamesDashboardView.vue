<!--
  ManagementGamesDashboardView.vue
  Main dashboard for game management showing games and their instances.
-->
<template>
  <div class="games-dashboard">
    <PageHeader 
      title="Games & Instances" 
      titleLevel="h2" 
      :showIcon="false"
      subtitle="Manage your game instances and monitor player activity" 
    />

    <!-- Loading state -->
    <p v-if="gamesStore.loading">Loading games...</p>

    <!-- Error state -->
    <div v-else-if="gamesStore.error" class="error-state">
      <p>Error loading games: {{ gamesStore.error }}</p>
      <button @click="loadGames">Retry</button>
    </div>

    <!-- Games table -->
    <ResourceTable 
      v-else
      :columns="columns" 
      :rows="games" 
      :loading="gamesStore.loading"
      :error="gamesStore.error"
    >
      <template #cell-name="{ row }">
        <a href="#" class="game-link" @click.prevent="viewGameInstances(row)">{{ row.name }}</a>
      </template>

      <template #cell-instances="{ row }">
        {{ getGameInstanceCount(row.id) }}
      </template>

      <template #cell-active="{ row }">
        {{ getActiveInstanceCount(row.id) }}
      </template>

      <template #actions="{ row }">
        <TableActions :actions="getGameActions(row)" />
      </template>
    </ResourceTable>

    <!-- Empty state -->
    <p v-if="!gamesStore.loading && !gamesStore.error && games.length === 0" class="empty-message">
      You don't have access to any games yet. Contact an administrator to get access to games.
    </p>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useGamesStore } from '../../stores/games';
import { useGameInstancesStore } from '../../stores/gameInstances';
import PageHeader from '../../components/PageHeader.vue';
import ResourceTable from '../../components/ResourceTable.vue';
import TableActions from '../../components/TableActions.vue';

const router = useRouter();
const gamesStore = useGamesStore();
const gameInstancesStore = useGameInstancesStore();

const games = computed(() => gamesStore.games);

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'game_type', label: 'Type' },
  { key: 'instances', label: 'Instances' },
  { key: 'active', label: 'Active' }
];

onMounted(async () => {
  await loadGames();
});

const loadGames = async () => {
  try {
    // Filter to only show games where the user has Manager subscription
    await gamesStore.fetchGames({ subscriptionType: 'Manager' });
    await gameInstancesStore.fetchAllGameInstances();
  } catch (error) {
    console.error('Failed to load games:', error);
  }
};

const getGameInstanceCount = (gameId) => {
  return gameInstancesStore.gameInstances.filter(instance => instance.game_id === gameId).length;
};

const getActiveInstanceCount = (gameId) => {
  return gameInstancesStore.gameInstances.filter(instance =>
    instance.game_id === gameId &&
    ['started'].includes(instance.status)
  ).length;
};

const viewGameInstances = (game) => {
  gamesStore.setSelectedGame(game);
  router.push(`/admin/games/${game.id}/instances`);
};

const getGameActions = (game) => {
  return [
    {
      key: 'manage',
      label: 'Manage',
      primary: true,
      handler: () => viewGameInstances(game)
    }
  ];
};
</script>

<style scoped>
.games-dashboard {
  width: 100%;
}

.error-state {
  padding: var(--space-lg);
}

.error-state button {
  margin-top: var(--space-md);
  padding: var(--space-sm) var(--space-md);
  background: var(--color-primary);
  color: var(--color-text-light);
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
}

.game-link {
  color: var(--color-primary);
  text-decoration: none;
}

.game-link:hover {
  text-decoration: underline;
}

.empty-message {
  color: var(--color-text-muted);
  font-style: italic;
}
</style>
