<!--
  ManagementInstanceDetailView.vue
  Detailed view for managing individual game instances.
-->
<template>
  <div class="instance-detail-view">
    <div class="view-header">
      <div class="header-content">
        <h2>{{ selectedGame?.name }} - Instance Details</h2>
        <p>Manage game instance and monitor player activity</p>
      </div>
      <div class="header-actions">
        <button @click="goBack" class="btn-secondary">
          Back to Instances
        </button>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="loading-state">
      <p>Loading instance details...</p>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="error-state">
      <p>Error loading instance: {{ error }}</p>
      <button @click="loadInstance">Retry</button>
    </div>

    <!-- Instance details -->
    <div v-else-if="instance" class="instance-details">
      <!-- Status and Progress Section -->
      <div class="detail-section">
        <h3>Status & Progress</h3>
        <div class="status-grid">
          <div class="status-item">
            <span class="label">Status</span>
            <span :class="['status-badge', `status-${instance.status}`]">
              {{ getStatusLabel(instance.status) }}
            </span>
          </div>
          <div class="status-item">
            <span class="label">Current Turn</span>
            <span class="value">{{ instance.current_turn }}</span>
          </div>
          <div class="status-item" v-if="instance.next_turn_due_at">
            <span class="label">Next Turn Due</span>
            <span class="value">{{ formatDeadline(instance.next_turn_due_at) }}</span>
          </div>
        </div>
      </div>

      <!-- Timeline Section -->
      <div class="detail-section">
        <h3>Timeline</h3>
        <div class="timeline-grid">
          <div class="timeline-item">
            <span class="label">Created</span>
            <span class="value">{{ formatDate(instance.created_at) }}</span>
          </div>
          <div class="timeline-item" v-if="instance.started_at">
            <span class="label">Started</span>
            <span class="value">{{ formatDate(instance.started_at) }}</span>
          </div>
          <div class="timeline-item" v-if="instance.last_turn_processed_at">
            <span class="label">Last Turn Processed</span>
            <span class="value">{{ formatDate(instance.last_turn_processed_at) }}</span>
          </div>
          <div class="timeline-item" v-if="instance.completed_at">
            <span class="label">Completed</span>
            <span class="value">{{ formatDate(instance.completed_at) }}</span>
          </div>
        </div>
      </div>

      <!-- Runtime Controls Section -->
      <div class="detail-section">
        <h3>Runtime Controls</h3>
        <div class="controls-grid">
          <button 
            v-if="instance.status === 'created'" 
            @click="startInstance" 
            class="btn-primary"
            :disabled="controlLoading"
          >
            Start Game
          </button>
          <button 
            v-if="instance.status === 'running'" 
            @click="pauseInstance" 
            class="btn-warning"
            :disabled="controlLoading"
          >
            Pause Game
          </button>
          <button 
            v-if="instance.status === 'paused'" 
            @click="resumeInstance" 
            class="btn-success"
            :disabled="controlLoading"
          >
            Resume Game
          </button>
          <button 
            v-if="['created', 'running', 'paused'].includes(instance.status)" 
            @click="cancelInstance" 
            class="btn-danger"
            :disabled="controlLoading"
          >
            Cancel Game
          </button>
        </div>
      </div>

      <!-- Game Parameters Section -->
      <div class="detail-section">
        <h3>Game Parameters</h3>
        <div v-if="instanceParametersLoading" class="loading-content">
          <p>Loading parameters...</p>
        </div>
        <div v-else-if="instanceParameters.length === 0" class="placeholder-content">
          <p>No parameters configured for this instance.</p>
        </div>
        <div v-else class="parameters-grid">
          <div v-for="param in instanceParameters" :key="param.id" class="parameter-item">
            <span class="label">{{ param.config_key }}</span>
            <span class="value">{{ formatParameterValue(param) }}</span>
          </div>
        </div>
      </div>

      <!-- Player Activity Section (placeholder for future) -->
      <div class="detail-section">
        <h3>Player Activity</h3>
        <div class="placeholder-content">
          <p>Player activity monitoring will be implemented in future updates.</p>
          <p>This will include:</p>
          <ul>
            <li>Player turn submissions</li>
            <li>Turn processing status</li>
            <li>Player queries and support requests</li>
            <li>Game state changes</li>
          </ul>
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
import { useGameInstanceParametersStore } from '../../stores/gameInstanceParameters';

const route = useRoute();
const router = useRouter();
const gamesStore = useGamesStore();
const gameInstancesStore = useGameInstancesStore();
const gameInstanceParametersStore = useGameInstanceParametersStore();

const gameId = computed(() => route.params.gameId);
const instanceId = computed(() => route.params.instanceId);
const selectedGame = computed(() => gamesStore.games.find(g => g.id === gameId.value));

const loading = ref(false);
const controlLoading = ref(false);
const instanceParametersLoading = ref(false);
const error = ref('');
const instance = ref(null);
const instanceParameters = ref([]);

onMounted(async () => {
  await loadInstance();
  await loadInstanceParameters();
});

const loadInstance = async () => {
  loading.value = true;
  error.value = '';

  try {
    if (!selectedGame.value) {
      await gamesStore.fetchGames();
    }
    
    const instanceData = await gameInstancesStore.getGameInstance(gameId.value, instanceId.value);
    instance.value = instanceData;
  } catch (err) {
    error.value = err.message || 'Failed to load instance details';
  } finally {
    loading.value = false;
  }
};

const loadInstanceParameters = async () => {
  if (!instanceId.value) return;
  
  instanceParametersLoading.value = true;
  try {
    await gameInstanceParametersStore.fetchGameInstanceParameters(instanceId.value);
    instanceParameters.value = gameInstanceParametersStore.getParametersByGameInstanceId(instanceId.value);
  } catch (err) {
    console.error('Failed to load instance parameters:', err);
  } finally {
    instanceParametersLoading.value = false;
  }
};

const formatParameterValue = (param) => {
  if (param.value_type === 'boolean') {
    return param.value === 'true' ? 'Yes' : 'No';
  }
  if (param.value_type === 'json') {
    try {
      return JSON.stringify(JSON.parse(param.value), null, 2);
    } catch {
      return param.value;
    }
  }
  return param.value;
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
  return new Date(dateString).toLocaleString();
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

const startInstance = async () => {
  await performAction('startGameInstance', 'Failed to start instance');
};

const pauseInstance = async () => {
  await performAction('pauseGameInstance', 'Failed to pause instance');
};

const resumeInstance = async () => {
  await performAction('resumeGameInstance', 'Failed to resume instance');
};

const cancelInstance = async () => {
  if (!confirm('Are you sure you want to cancel this game instance? This action cannot be undone.')) {
    return;
  }
  await performAction('cancelGameInstance', 'Failed to cancel instance');
};

const performAction = async (action, errorMessage) => {
  controlLoading.value = true;
  try {
    await gameInstancesStore[action](gameId.value, instanceId.value);
    await loadInstance(); // Reload instance data
  } catch (err) {
    error.value = err.message || errorMessage;
  } finally {
    controlLoading.value = false;
  }
};

const goBack = () => {
  router.push(`/admin/games/${gameId.value}/instances`);
};
</script>

<style scoped>
.instance-detail-view {
  max-width: 1000px;
  margin: 0 auto;
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

.instance-details {
  display: flex;
  flex-direction: column;
  gap: var(--space-xl);
}

.detail-section {
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-lg);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.detail-section h3 {
  margin: 0 0 var(--space-lg) 0;
  font-size: var(--font-size-lg);
  color: var(--color-text);
  padding-bottom: var(--space-sm);
  border-bottom: 1px solid var(--color-border);
}

.status-grid,
.timeline-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--space-md);
}

.parameters-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--space-md);
}

.parameter-item {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
  padding: var(--space-md);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.status-item,
.timeline-item {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.label {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
  font-weight: var(--font-weight-medium);
}

.value {
  font-size: var(--font-size-md);
  color: var(--color-text);
  font-weight: var(--font-weight-medium);
}

.status-badge {
  padding: var(--space-xs) var(--space-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-bold);
  text-transform: uppercase;
  align-self: flex-start;
}

.status-created {
  background: var(--color-bg-light);
  color: var(--color-text-muted);
}

.status-starting {
  background: var(--color-warning-light);
  color: var(--color-warning);
}

.status-running {
  background: var(--color-success-light);
  color: var(--color-success);
}

.status-paused {
  background: var(--color-warning-light);
  color: var(--color-warning);
}

.status-completed {
  background: var(--color-success-light);
  color: var(--color-success);
}

.status-cancelled {
  background: var(--color-danger-light);
  color: var(--color-danger);
}

.config-display {
  background: var(--color-bg-light);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: var(--space-md);
  overflow-x: auto;
}

.config-display pre {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--color-text);
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.controls-grid {
  display: flex;
  gap: var(--space-md);
  flex-wrap: wrap;
}

.btn-primary,
.btn-secondary,
.btn-warning,
.btn-success,
.btn-danger {
  padding: var(--space-sm) var(--space-md);
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  transition: background 0.2s;
}

.btn-primary {
  background: var(--color-primary);
  color: var(--color-text-light);
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-dark);
}

.btn-secondary {
  background: var(--color-bg-light);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}

.btn-secondary:hover {
  background: var(--color-border);
}

.btn-warning {
  background: var(--color-warning);
  color: var(--color-text-light);
}

.btn-warning:hover:not(:disabled) {
  background: var(--color-warning-dark);
}

.btn-success {
  background: var(--color-success);
  color: var(--color-text-light);
}

.btn-success:hover:not(:disabled) {
  background: var(--color-success-dark);
}

.btn-danger {
  background: var(--color-danger);
  color: var(--color-text-light);
}

.btn-danger:hover:not(:disabled) {
  background: var(--color-danger-dark);
}

.btn-primary:disabled,
.btn-warning:disabled,
.btn-success:disabled,
.btn-danger:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.placeholder-content {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.placeholder-content p {
  margin-bottom: var(--space-md);
}

.placeholder-content ul {
  margin-left: var(--space-lg);
  padding-left: var(--space-sm);
}

.placeholder-content li {
  margin-bottom: var(--space-xs);
}
</style> 