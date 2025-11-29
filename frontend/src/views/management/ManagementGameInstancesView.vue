<!--
  ManagementGameInstancesView.vue
  View for managing game instances of a specific game using ResourceTable.
-->
<template>
  <div class="game-instances-view">
    <div class="view-header">
      <div class="header-content">
        <h2>{{ selectedGame?.name }} - Game Instances</h2>
        <p>Manage active game sessions and monitor player activity</p>
        <Button @click="goBack" variant="secondary" size="small" class="back-button">
          <svg class="icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M20 11H7.83l5.59-5.59L12 4l-8 8 8 8 1.41-1.41L7.83 13H20v-2z"/>
          </svg>
          Back to Games
        </Button>
      </div>
    </div>

    <!-- Active Instances Section -->
    <DataCard title="Active Instances" class="instances-card">
      <template #header-extra>
        <Button @click="createInstance" variant="primary" size="small">
          <svg class="icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
          </svg>
          Create Instance
        </Button>
      </template>

      <ResourceTable 
        :columns="columns" 
        :rows="activeInstances" 
        :loading="gameInstancesStore.loading"
        :error="gameInstancesStore.error"
      >
        <template #cell-id="{ row }">
          <router-link :to="`/admin/games/${gameId}/instances/${row.id}`" class="instance-id-link">
            {{ row.id.slice(0, 8) }}...
          </router-link>
        </template>

        <template #cell-status="{ row }">
          <span :class="['status-badge', `status-${row.status}`]">
            {{ getStatusLabel(row.status) }}
          </span>
        </template>

        <template #cell-current_turn="{ row }">
          {{ row.current_turn || 0 }}
        </template>

        <template #cell-next_turn_due_at="{ row }">
          {{ formatDeadline(row.next_turn_due_at) }}
        </template>

        <template #cell-started_at="{ row }">
          {{ formatDate(row.started_at) }}
        </template>

        <template #actions="{ row }">
          <TableActionsMenu :actions="getActiveInstanceActions(row)" />
        </template>
      </ResourceTable>

      <div v-if="!gameInstancesStore.loading && activeInstances.length === 0" class="empty-state">
        <p>No active instances. Create a new game instance to get started.</p>
      </div>
    </DataCard>

    <!-- Completed Instances Section -->
    <DataCard title="Completed Instances" class="instances-card">
      <ResourceTable 
        :columns="completedColumns" 
        :rows="completedInstances" 
        :loading="gameInstancesStore.loading"
        :error="null"
      >
        <template #cell-id="{ row }">
          <router-link :to="`/admin/games/${gameId}/instances/${row.id}`" class="instance-id-link">
            {{ row.id.slice(0, 8) }}...
          </router-link>
        </template>

        <template #cell-status="{ row }">
          <span :class="['status-badge', `status-${row.status}`]">
            {{ getStatusLabel(row.status) }}
          </span>
        </template>

        <template #cell-current_turn="{ row }">
          {{ row.current_turn || 0 }}
        </template>

        <template #cell-completed_at="{ row }">
          {{ formatDate(row.completed_at) }}
        </template>

        <template #actions="{ row }">
          <TableActionsMenu :actions="getCompletedInstanceActions(row)" />
        </template>
      </ResourceTable>

      <div v-if="!gameInstancesStore.loading && completedInstances.length === 0" class="empty-state">
        <p>No completed instances yet.</p>
      </div>
    </DataCard>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useGamesStore } from '../../stores/games';
import { useGameInstancesStore } from '../../stores/gameInstances';
import Button from '../../components/Button.vue';
import DataCard from '../../components/DataCard.vue';
import ResourceTable from '../../components/ResourceTable.vue';
import TableActionsMenu from '../../components/TableActionsMenu.vue';

const route = useRoute();
const router = useRouter();
const gamesStore = useGamesStore();
const gameInstancesStore = useGameInstancesStore();

const gameId = computed(() => route.params.gameId);
const selectedGame = computed(() => gamesStore.games.find(g => g.id === gameId.value));

const columns = [
  { key: 'id', label: 'Instance ID' },
  { key: 'status', label: 'Status' },
  { key: 'current_turn', label: 'Turn' },
  { key: 'next_turn_due_at', label: 'Next Turn Due' },
  { key: 'started_at', label: 'Started' }
];

const completedColumns = [
  { key: 'id', label: 'Instance ID' },
  { key: 'status', label: 'Status' },
  { key: 'current_turn', label: 'Final Turn' },
  { key: 'completed_at', label: 'Completed' }
];

const gameInstances = computed(() => gameInstancesStore.gameInstances);

const activeInstances = computed(() => 
  gameInstances.value.filter(instance => 
    instance.game_id === gameId.value && 
    ['created', 'started', 'paused'].includes(instance.status)
  )
);

const completedInstances = computed(() => 
  gameInstances.value.filter(instance => 
    instance.game_id === gameId.value && 
    ['completed', 'cancelled'].includes(instance.status)
  )
);

onMounted(async () => {
  if (!gamesStore.games.length) {
    await gamesStore.fetchGames();
  }
  if (selectedGame.value) {
    gamesStore.setSelectedGame(selectedGame.value);
  }
  await loadGameInstances();
});

const loadGameInstances = async () => {
  try {
    await gameInstancesStore.fetchGameInstances(gameId.value);
  } catch (error) {
    console.error('Failed to load game instances:', error);
  }
};

const getStatusLabel = (status) => {
  const labels = {
    'created': 'Created',
    'started': 'Running',
    'paused': 'Paused',
    'completed': 'Completed',
    'cancelled': 'Cancelled'
  };
  return labels[status] || status;
};

const formatDate = (dateString) => {
  if (!dateString) return 'N/A';
  return new Date(dateString).toLocaleDateString();
};

const formatDeadline = (deadlineString) => {
  if (!deadlineString) return 'N/A';
  const deadline = new Date(deadlineString);
  const now = new Date();
  const diff = deadline - now;
  
  if (diff < 0) return 'Overdue';
  if (diff < 24 * 60 * 60 * 1000) return 'Today';
  if (diff < 48 * 60 * 60 * 1000) return 'Tomorrow';
  
  return deadline.toLocaleDateString();
};

const goBack = () => {
  router.push('/admin');
};

const createInstance = () => {
  router.push(`/admin/games/${gameId.value}/instances/create`);
};

const viewInstance = (instance) => {
  router.push(`/admin/games/${gameId.value}/instances/${instance.id}`);
};

const startInstance = async (instance) => {
  try {
    await gameInstancesStore.startGameInstance(gameId.value, instance.id);
    await loadGameInstances();
  } catch (error) {
    console.error('Failed to start instance:', error);
  }
};

const pauseInstance = async (instance) => {
  try {
    await gameInstancesStore.pauseGameInstance(gameId.value, instance.id);
    await loadGameInstances();
  } catch (error) {
    console.error('Failed to pause instance:', error);
  }
};

const resumeInstance = async (instance) => {
  try {
    await gameInstancesStore.resumeGameInstance(gameId.value, instance.id);
    await loadGameInstances();
  } catch (error) {
    console.error('Failed to resume instance:', error);
  }
};

const cancelInstance = async (instance) => {
  if (!confirm(`Are you sure you want to cancel this game instance?`)) return;
  
  try {
    await gameInstancesStore.cancelGameInstance(gameId.value, instance.id);
    await loadGameInstances();
  } catch (error) {
    console.error('Failed to cancel instance:', error);
  }
};

const getActiveInstanceActions = (instance) => {
  const actions = [
    { key: 'view', label: 'View Details', handler: () => viewInstance(instance) }
  ];

  if (instance.status === 'created') {
    actions.push({ key: 'start', label: 'Start', handler: () => startInstance(instance) });
  } else if (instance.status === 'started') {
    actions.push({ key: 'pause', label: 'Pause', handler: () => pauseInstance(instance) });
  } else if (instance.status === 'paused') {
    actions.push({ key: 'resume', label: 'Resume', handler: () => resumeInstance(instance) });
  }

  if (['created', 'started', 'paused'].includes(instance.status)) {
    actions.push({ key: 'cancel', label: 'Cancel', danger: true, handler: () => cancelInstance(instance) });
  }

  return actions;
};

const getCompletedInstanceActions = (instance) => {
  return [
    { key: 'view', label: 'View Details', handler: () => viewInstance(instance) }
  ];
};
</script>

<style scoped>
.game-instances-view {
  max-width: 1200px;
  margin: 0 auto;
}

.view-header {
  margin-bottom: var(--space-lg);
  padding-bottom: var(--space-lg);
  border-bottom: 1px solid var(--color-border);
}

.header-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.header-content h2 {
  margin: 0;
  font-size: var(--font-size-xl);
  color: var(--color-text);
}

.header-content p {
  margin: 0;
  color: var(--color-text-muted);
  font-size: var(--font-size-md);
}

.back-button {
  margin-top: var(--space-sm);
  align-self: flex-start;
}

.icon {
  width: 16px;
  height: 16px;
}

.instances-card {
  margin-bottom: var(--space-xl);
}

.instance-id-link {
  font-family: monospace;
  font-size: var(--font-size-sm);
  color: var(--color-primary);
  text-decoration: none;
}

.instance-id-link:hover {
  text-decoration: underline;
}

.status-badge {
  padding: var(--space-xs) var(--space-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-bold);
  text-transform: uppercase;
  color: var(--color-text-light);
  white-space: nowrap;
  display: inline-block;
}

.status-created {
  background: var(--color-info);
}

.status-started {
  background: var(--color-success-light);
  color: var(--color-success);
}

.status-paused {
  background: var(--color-warning);
}

.status-completed {
  background: var(--color-success);
}

.status-cancelled {
  background: var(--color-danger);
}

.empty-state {
  text-align: center;
  padding: var(--space-lg);
  color: var(--color-text-muted);
}

.empty-state p {
  margin: 0;
}
</style>
