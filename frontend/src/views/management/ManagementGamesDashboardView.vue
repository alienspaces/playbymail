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
    <p v-if="managedGamesLoading">Loading games...</p>

    <!-- Error state -->
    <div v-else-if="managedGamesError" class="error-state">
      <p>Error loading games: {{ managedGamesError }}</p>
      <button @click="loadManagedGames">Retry</button>
    </div>

    <!-- Games table -->
    <ResourceTable 
      v-else
      :columns="columns" 
      :rows="managedGames" 
      :loading="managedGamesLoading"
      :error="managedGamesError"
    >
      <template #cell-name="{ row }">
        <a href="#" class="game-link" @click.prevent="viewGameInstances(row)">{{ row.game_name }}</a>
      </template>

      <template #cell-game_type="{ row }">
        {{ row.game_type }}
      </template>

      <template #cell-instances="{ row }">
        {{ row.instance_count }}
      </template>

      <template #cell-active="{ row }">
        {{ row.active_count }}
      </template>

      <template #actions="{ row }">
        <TableActions :actions="getGameActions(row)" />
      </template>
    </ResourceTable>

    <!-- Empty state -->
    <p v-if="!managedGamesLoading && !managedGamesError && managedGames.length === 0" class="empty-message">
      You don't have access to any games yet. Subscribe to a published game below.
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
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useGamesStore } from '../../stores/games';
import { listAllGameInstances } from '../../api/gameInstances';
import { listGames } from '../../api/games';
import { createGameSubscription } from '../../api/gameSubscriptions';
import PageHeader from '../../components/PageHeader.vue';
import ResourceTable from '../../components/ResourceTable.vue';
import TableActions from '../../components/TableActions.vue';
import DataCard from '../../components/DataCard.vue';
import DataItem from '../../components/DataItem.vue';
import AppButton from '../../components/Button.vue';

const router = useRouter();
const gamesStore = useGamesStore();

const managedGames = ref([]);
const managedGamesLoading = ref(false);
const managedGamesError = ref(null);
const availableGames = ref([]);
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
  await loadManagedGames();
  await loadAvailableGames();
});

const loadManagedGames = async () => {
  try {
    managedGamesLoading.value = true;
    managedGamesError.value = null;

    const response = await listAllGameInstances();
    const rows = response.data || [];

    const gameMap = new Map();
    for (const row of rows) {
      if (!gameMap.has(row.game_id)) {
        gameMap.set(row.game_id, {
          game_id: row.game_id,
          game_name: row.game_name,
          game_type: row.game_type,
          game_description: row.game_description,
          game_subscription_id: row.game_subscription_id,
          created_at: row.created_at,
          instance_count: 0,
          active_count: 0,
        });
      }
      const game = gameMap.get(row.game_id);
      if (row.game_instance_id) {
        game.instance_count++;
        if (row.instance_status === 'started') {
          game.active_count++;
        }
      }
    }

    managedGames.value = Array.from(gameMap.values());
  } catch (error) {
    managedGamesError.value = error.message || 'Failed to load managed games';
    console.error('Failed to load managed games:', error);
  } finally {
    managedGamesLoading.value = false;
  }
};

const loadAvailableGames = async () => {
  try {
    availableGamesLoading.value = true;
    availableGamesError.value = null;

    const response = await listGames({ canManage: true });
    availableGames.value = response.data || [];
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
    await createGameSubscription(game.id, 'manager');
    await loadManagedGames();
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

const viewGameInstances = (game) => {
  gamesStore.setSelectedGame({ id: game.game_id, name: game.game_name, game_type: game.game_type });
  router.push(`/admin/games/${game.game_id}/instances`);
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
