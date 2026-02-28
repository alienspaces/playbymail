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

    <!-- Available Games for Subscription Section -->
    <div class="available-games-section">
      <h3>Available Games for Subscription</h3>
      <p class="section-description">Published games available for manager subscription.</p>

      <div v-if="availableGamesLoading" class="loading-state">
        <p>Loading available games...</p>
      </div>

      <div v-else-if="availableGamesError" class="error-state">
        <p>{{ availableGamesError }}</p>
        <button @click="loadAvailableGames">Retry</button>
      </div>

      <div v-else-if="availableGames.length > 0" class="available-games-grid">
        <DataCard v-for="game in availableGames" :key="game.id" :title="game.name" class="game-card">
          <div class="game-info">
            <DataItem label="Type" :value="game.game_type" />
            <DataItem label="Status" :value="game.status" />
            <DataItem label="Turn Duration" :value="formatTurnDuration(game.turn_duration_hours)" />
            <DataItem label="Description" :value="game.description" />
          </div>

          <template #primary>
            <AppButton @click="subscribeToGame(game)" variant="primary" size="small" :disabled="subscribing">
              {{ subscribing ? 'Subscribing...' : 'Subscribe as Manager' }}
            </AppButton>
          </template>
        </DataCard>
      </div>

      <div v-else class="empty-state">
        <p>No published games available for subscription.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useGamesStore } from '../../stores/games';
import { useGameInstancesStore } from '../../stores/gameInstances';
import { listGames } from '../../api/games';
import { createGameSubscription, getMyGameSubscriptions } from '../../api/gameSubscriptions';
import PageHeader from '../../components/PageHeader.vue';
import ResourceTable from '../../components/ResourceTable.vue';
import TableActions from '../../components/TableActions.vue';
import DataCard from '../../components/DataCard.vue';
import DataItem from '../../components/DataItem.vue';
import AppButton from '../../components/Button.vue';

const router = useRouter();
const gamesStore = useGamesStore();
const gameInstancesStore = useGameInstancesStore();

const games = computed(() => gamesStore.games);
const availableGames = ref([]);
const myGameSubscriptions = ref([]);
const availableGamesLoading = ref(false);
const availableGamesError = ref(null);
const subscribing = ref(false);

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'game_type', label: 'Type' },
  { key: 'instances', label: 'Instances' },
  { key: 'active', label: 'Active' }
];

onMounted(async () => {
  await loadGames();
  await loadAvailableGames();
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

const loadAvailableGames = async () => {
  try {
    availableGamesLoading.value = true;
    availableGamesError.value = null;

    // Get published games
    const publishedGamesResponse = await listGames({ status: 'published' });
    const publishedGames = publishedGamesResponse.data || [];

    // Get my game subscriptions to filter out games I'm already subscribed to
    const subscriptionsResponse = await getMyGameSubscriptions();
    myGameSubscriptions.value = subscriptionsResponse.data || [];
    const subscribedGameIds = new Set(myGameSubscriptions.value
      .filter(sub => sub.subscription_type === 'Manager')
      .map(sub => sub.game_id));

    // Filter to only show games I don't have manager subscription for
    availableGames.value = publishedGames.filter(game => !subscribedGameIds.has(game.id));
  } catch (error) {
    availableGamesError.value = error.message || 'Failed to load available games';
    console.error('Failed to load available games:', error);
  } finally {
    availableGamesLoading.value = false;
  }
};

const subscribeToGame = async (game) => {
  try {
    subscribing.value = true;
    await createGameSubscription(game.id, 'Manager');
    await loadAvailableGames();
  } catch (error) {
    availableGamesError.value = error.message || 'Failed to subscribe to game';
    console.error('Failed to subscribe to game:', error);
  } finally {
    subscribing.value = false;
  }
};

const formatTurnDuration = (hours) => {
  if (!hours) return 'Not set';
  if (hours % (24 * 7) === 0) {
    const weeks = hours / (24 * 7);
    return `${weeks} week${weeks === 1 ? '' : 's'}`;
  }
  if (hours % 24 === 0) {
    const days = hours / 24;
    return `${days} day${days === 1 ? '' : 's'}`;
  }
  return `${hours} hour${hours === 1 ? '' : 's'}`;
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

.available-games-section {
  margin-top: var(--space-xl);
  padding-top: var(--space-xl);
  border-top: 1px solid var(--color-border);
}

.available-games-section h3 {
  margin: 0 0 var(--space-sm) 0;
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text);
}

.section-description {
  margin: 0 0 var(--space-md) 0;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.available-games-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(600px, 1fr));
  gap: var(--space-lg);
  margin-top: var(--space-md);
}

.game-card {
  min-height: 200px;
}

.game-info {
  margin-bottom: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}
</style>
