<!--
  ManagementGameInstancesView.vue
  View for managing game instances of a specific game.
-->
<template>
  <div class="game-instances-view">
    <div class="view-header">
      <div class="header-content">
        <h2>{{ selectedGame?.name }} - Game Instances</h2>
        <p>Manage active game sessions and monitor player activity</p>
      </div>
      <div class="header-actions">
        <Button @click="createInstance" variant="primary">
          <svg class="icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
          </svg>
          Create Instance
        </Button>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="gameInstancesStore.loading" class="loading-state">
      <p>Loading game instances...</p>
    </div>

    <!-- Error state -->
    <div v-else-if="gameInstancesStore.error" class="error-state">
      <p>Error loading game instances: {{ gameInstancesStore.error }}</p>
      <button @click="loadGameInstances">Retry</button>
    </div>

    <!-- Active Instances Section -->
    <div class="instances-section">
      <h3>
        Active Instances
        <span class="count">({{ activeInstances.length }})</span>
      </h3>
      
      <div v-if="activeInstances.length === 0" class="empty-state">
        <h4>No Active Instances</h4>
        <p>Create a new game instance to get started.</p>
      </div>
      
      <div v-else class="instances-grid">
        <div v-for="instance in activeInstances" :key="instance.id" class="instance-card">
          <div class="instance-info">
            <div class="instance-header">
              <h4>Instance #{{ instance.id.slice(0, 8) }}</h4>
              <span class="status-badge status-{{ instance.status }}">
                {{ getStatusLabel(instance.status) }}
              </span>
            </div>
            <div class="info-row">
              <span class="label">Turn:</span>
              <span class="value">{{ instance.current_turn || 0 }}</span>
            </div>
            <div class="info-row">
              <span class="label">Next Turn Due:</span>
              <span class="value">{{ formatDeadline(instance.next_turn_due_at) }}</span>
            </div>
            <div class="info-row">
              <span class="label">Started:</span>
              <span class="value">{{ formatDate(instance.started_at) }}</span>
            </div>
          </div>

          <div class="instance-actions">
            <Button @click="viewInstance(instance)" variant="secondary" size="small">
              View Details
            </Button>
            <Button 
              v-if="instance.status === 'created'" 
              @click="startInstance(instance)" 
              variant="primary" 
              size="small"
            >
              Start
            </Button>
            <Button 
              v-if="instance.status === 'running'" 
              @click="pauseInstance(instance)" 
              variant="warning" 
              size="small"
            >
              Pause
            </Button>
            <Button 
              v-if="instance.status === 'paused'" 
              @click="resumeInstance(instance)" 
              variant="success" 
              size="small"
            >
              Resume
            </Button>
            <Button 
              v-if="['created', 'running', 'paused'].includes(instance.status)" 
              @click="cancelInstance(instance)" 
              variant="danger" 
              size="small"
            >
              Cancel
            </Button>
          </div>
        </div>
      </div>
    </div>

    <!-- Completed Instances Section -->
    <div class="instances-section">
      <h3>
        Completed Instances
        <span class="count">({{ completedInstances.length }})</span>
      </h3>
      
      <div v-if="completedInstances.length === 0" class="empty-state">
        <h4>No Completed Instances</h4>
        <p>Completed games will appear here.</p>
      </div>
      
      <div v-else class="instances-grid">
        <div v-for="instance in completedInstances" :key="instance.id" class="instance-card">
          <div class="instance-info">
            <div class="instance-header">
              <h4>Instance #{{ instance.id.slice(0, 8) }}</h4>
              <span class="status-badge status-{{ instance.status }}">
                {{ getStatusLabel(instance.status) }}
              </span>
            </div>
            <div class="info-row">
              <span class="label">Final Turn:</span>
              <span class="value">{{ instance.current_turn || 0 }}</span>
            </div>
            <div class="info-row">
              <span class="label">Completed:</span>
              <span class="value">{{ formatDate(instance.completed_at) }}</span>
            </div>
          </div>

          <div class="instance-actions">
            <Button @click="viewInstance(instance)" variant="secondary" size="small">
              View Details
            </Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useGamesStore } from '../../stores/games';
import { useGameInstancesStore } from '../../stores/gameInstances';
import Button from '../../components/Button.vue';

const route = useRoute();
const router = useRouter();
const gamesStore = useGamesStore();
const gameInstancesStore = useGameInstancesStore();

const gameId = computed(() => route.params.gameId);
const selectedGame = computed(() => gamesStore.games.find(g => g.id === gameId.value));

const gameInstances = computed(() => gameInstancesStore.gameInstances);
const activeInstances = computed(() => 
  gameInstances.value.filter(instance => 
    instance.game_id === gameId.value && 
    ['created', 'starting', 'running', 'paused'].includes(instance.status)
  )
);
const completedInstances = computed(() => 
  gameInstances.value.filter(instance => 
    instance.game_id === gameId.value && 
    ['completed', 'cancelled'].includes(instance.status)
  )
);

onMounted(async () => {
  // Load games first so selectedGame can be populated
  if (!gamesStore.games.length) {
    await gamesStore.fetchGames();
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
    'starting': 'Starting',
    'running': 'Running',
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
</script>

<style scoped>
.game-instances-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--space-lg);
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--space-xl);
  padding-bottom: var(--space-lg);
  border-bottom: 1px solid var(--color-border);
}

.header-content h2 {
  margin: 0 0 var(--space-sm) 0;
  font-size: var(--font-size-xl);
  color: var(--color-text);
}

.header-content p {
  margin: 0;
  color: var(--color-text-muted);
  font-size: var(--font-size-md);
}

.header-actions {
  display: flex;
  gap: var(--space-md);
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-sm) var(--space-md);
  background: var(--color-primary);
  color: var(--color-text-light);
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  transition: background 0.2s;
}

.btn-primary:hover {
  background: var(--color-primary-dark);
}

.icon {
  width: 16px;
  height: 16px;
}

.loading-state,
.error-state {
  text-align: center;
  padding: var(--space-xl);
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
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

.instances-section {
  margin-bottom: var(--space-xl);
}

.instances-section h3 {
  margin: 0 0 var(--space-md) 0;
  font-size: var(--font-size-lg);
  color: var(--color-text);
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.instances-section h3 .count {
  background: var(--color-primary);
  color: var(--color-text-light);
  padding: var(--space-xs) var(--space-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-normal);
}

.instances-grid {
  display: grid;
  gap: var(--space-md);
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
}

.instance-card {
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--space-lg);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: box-shadow 0.2s ease;
}

.instance-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.instance-card.completed {
  opacity: 0.7;
}

.instance-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-md);
}

.instance-header h4 {
  margin: 0;
  font-size: var(--font-size-md);
  color: var(--color-text);
  font-weight: var(--font-weight-bold);
}

.status-badge {
  padding: var(--space-xs) var(--space-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-bold);
  text-transform: uppercase;
  color: var(--color-text-light);
}

.status-created {
  background: var(--color-info);
}

.status-starting {
  background: var(--color-warning);
}

.status-running {
  background: var(--color-success);
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

.instance-info {
  margin-bottom: var(--space-lg);
}

.info-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: var(--space-sm);
}

.info-row:last-child {
  margin-bottom: 0;
}

.label {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}

.value {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text);
}

.instance-actions {
  display: flex;
  gap: var(--space-sm);
  flex-wrap: wrap;
}

.empty-state {
  text-align: center;
  padding: var(--space-xl);
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.empty-state h4 {
  margin: 0 0 var(--space-sm) 0;
  color: var(--color-text);
}

.empty-state p {
  margin: 0;
  color: var(--color-text-muted);
}
</style> 