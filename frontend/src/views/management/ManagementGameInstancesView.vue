<!--
  ManagementGameInstancesView.vue
  View for managing game instances of a specific game using ResourceTable.
-->
<template>
  <div class="game-instances-view">
    <PageHeader 
      :title="`${selectedGame?.name || ''} - Game Instances`"
      subtitle="Manage active game sessions and monitor player activity"
      :showIcon="false"
      titleLevel="h2"
    />
    
    <Button @click="goBack" variant="secondary" size="small" class="back-button">
      <svg class="icon" viewBox="0 0 24 24" fill="currentColor">
        <path d="M20 11H7.83l5.59-5.59L12 4l-8 8 8 8 1.41-1.41L7.83 13H20v-2z"/>
      </svg>
      Back to Games
    </Button>

    <!-- Active Instances Section -->
    <div class="section">
      <PageHeader 
        title="Active Instances"
        actionText="Create Instance"
        :showIcon="false"
        titleLevel="h3"
        @action="createInstance"
      />

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

        <template #cell-player_count="{ row }">
          {{ row.player_count || 0 }}{{ row.required_player_count > 0 ? ` / ${row.required_player_count}` : '' }}
        </template>

        <template #cell-delivery_methods="{ row }">
          <span class="delivery-methods">
            <span v-if="row.delivery_physical_post" class="delivery-badge">Post</span>
            <span v-if="row.delivery_physical_local" class="delivery-badge">Local</span>
            <span v-if="row.delivery_email" class="delivery-badge">Email</span>
          </span>
        </template>

        <template #cell-next_turn_due_at="{ row }">
          {{ formatDeadline(row.next_turn_due_at) }}
        </template>

        <template #cell-started_at="{ row }">
          {{ formatDate(row.started_at) }}
        </template>

        <template #actions="{ row }">
          <TableActions :actions="getActiveInstanceActions(row)" />
        </template>
      </ResourceTable>

      <p v-if="!gameInstancesStore.loading && activeInstances.length === 0" class="empty-message">
        No active instances. Create a new game instance to get started.
      </p>
    </div>

    <!-- Completed Instances Section -->
    <div class="section">
      <PageHeader 
        title="Completed Instances"
        :showIcon="false"
        titleLevel="h3"
      />

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
          <TableActions :actions="getCompletedInstanceActions(row)" />
        </template>
      </ResourceTable>

      <p v-if="!gameInstancesStore.loading && completedInstances.length === 0" class="empty-message">
        No completed instances yet.
      </p>
    </div>

    <!-- Create Instance Modal -->
    <ResourceModalForm
      :visible="showCreateModal"
      mode="create"
      title="Game Instance"
      :fields="instanceFields"
      :model-value="instanceForm"
      :error="createModalError"
      @submit="handleCreateInstance"
      @cancel="closeCreateModal"
    />
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useGamesStore } from '../../stores/games';
import { useGameInstancesStore } from '../../stores/gameInstances';
import { ref } from 'vue';
import Button from '../../components/Button.vue';
import PageHeader from '../../components/PageHeader.vue';
import ResourceTable from '../../components/ResourceTable.vue';
import TableActions from '../../components/TableActions.vue';
import ResourceModalForm from '../../components/ResourceModalForm.vue';

const route = useRoute();
const router = useRouter();
const gamesStore = useGamesStore();
const gameInstancesStore = useGameInstancesStore();

const gameId = computed(() => route.params.gameId);
const selectedGame = computed(() => gamesStore.games.find(g => g.id === gameId.value));

// Create instance modal state
const showCreateModal = ref(false);
const createModalError = ref('');
const instanceForm = ref({
  delivery_email: true,
  delivery_physical_post: false,
  delivery_physical_local: false,
  required_player_count: 1,
  is_closed_testing: false
});

const isDraftGame = computed(() => selectedGame.value?.status === 'draft');

const instanceFields = computed(() => {
  const fields = [
    {
      key: 'delivery_email',
      label: 'Email Delivery',
      type: 'checkbox',
      checkboxLabel: 'Enable email delivery (web-based turn sheet viewer)'
    },
    {
      key: 'delivery_physical_post',
      label: 'Physical Post Delivery',
      type: 'checkbox',
      checkboxLabel: 'Enable physical post delivery (traditional mail-based)'
    },
    {
      key: 'delivery_physical_local',
      label: 'Physical Local Delivery',
      type: 'checkbox',
      checkboxLabel: 'Enable physical local delivery (convention/classroom - game master prints locally, players fill at table, manual scanning/submission)'
    },
    {
      key: 'required_player_count',
      label: 'Required Player Count',
      type: 'number',
      required: true,
      min: 1,
      placeholder: 'Minimum number of players required before game can start'
    },
  ];

  if (isDraftGame.value) {
    fields.push({
      key: 'closed_testing_notice',
      type: 'info',
      text: 'This game is unpublished. Instances are restricted to closed testing \u2014 players must be invited to join.'
    });
  } else {
    fields.push({
      key: 'is_closed_testing',
      label: 'Closed Testing',
      type: 'checkbox',
      checkboxLabel: 'Enable closed testing mode (requires join game key for players to join)'
    });
  }

  return fields;
});

const columns = [
  { key: 'id', label: 'Instance ID' },
  { key: 'status', label: 'Status' },
  { key: 'current_turn', label: 'Turn' },
  { key: 'player_count', label: 'Players' },
  { key: 'delivery_methods', label: 'Delivery Methods' },
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
  instanceForm.value = {
    delivery_email: true,
    delivery_physical_post: false,
    delivery_physical_local: false,
    required_player_count: 1,
    is_closed_testing: isDraftGame.value ? true : false
  };
  createModalError.value = '';
  showCreateModal.value = true;
};

const closeCreateModal = () => {
  showCreateModal.value = false;
  createModalError.value = '';
};

const handleCreateInstance = async (formData) => {
  createModalError.value = '';
  
  // Ensure boolean values are properly set (checkboxes can be undefined)
  const deliveryPhysicalPost = Boolean(formData.delivery_physical_post);
  const deliveryPhysicalLocal = Boolean(formData.delivery_physical_local);
  const deliveryEmail = Boolean(formData.delivery_email);
  
  // Validate at least one delivery method is selected
  if (!deliveryPhysicalPost && !deliveryPhysicalLocal && !deliveryEmail) {
    createModalError.value = 'At least one delivery method must be enabled';
    return;
  }

  try {
    const instanceData = {
      game_id: gameId.value,
      delivery_physical_post: deliveryPhysicalPost,
      delivery_physical_local: deliveryPhysicalLocal,
      delivery_email: deliveryEmail,
      required_player_count: formData.required_player_count || 1,
      is_closed_testing: Boolean(formData.is_closed_testing)
    };

    const createdInstance = await gameInstancesStore.createGameInstance(gameId.value, instanceData);
    
    // Close modal and refresh list
    closeCreateModal();
    await loadGameInstances();
    
    // Navigate to instance details
    if (createdInstance && createdInstance.id) {
      router.push(`/admin/games/${gameId.value}/instances/${createdInstance.id}`);
    }
  } catch (err) {
    console.error('Failed to create instance:', err);
    createModalError.value = err.message || 'Failed to create game instance';
  }
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
  width: 100%;
}

.back-button {
  margin-bottom: var(--space-lg);
}

.icon {
  width: 16px;
  height: 16px;
}

.section {
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

.empty-message {
  color: var(--color-text-muted);
  font-style: italic;
}

.delivery-methods {
  display: flex;
  gap: var(--space-xs);
  flex-wrap: wrap;
}

.delivery-badge {
  display: inline-block;
  padding: 2px 6px;
  border-radius: var(--radius-xs);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  background: var(--color-bg-light);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}
</style>
