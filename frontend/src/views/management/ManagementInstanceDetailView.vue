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
        <Button @click="goBack" variant="secondary">
          Back to Instances
        </Button>
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
          <Button 
            v-if="instance.status === 'created'" 
            @click="startInstance" 
            variant="primary"
            :disabled="controlLoading"
          >
            Start Game
          </Button>
          <Button 
            v-if="instance.status === 'running'" 
            @click="pauseInstance" 
            variant="warning"
            :disabled="controlLoading"
          >
            Pause Game
          </Button>
          <Button 
            v-if="instance.status === 'paused'" 
            @click="resumeInstance" 
            variant="success"
            :disabled="controlLoading"
          >
            Resume Game
          </Button>
          <Button 
            v-if="['created', 'running', 'paused'].includes(instance.status)" 
            @click="cancelInstance" 
            variant="danger"
            :disabled="controlLoading"
          >
            Cancel Game
          </Button>
        </div>
      </div>

      <!-- Game Instance Parameters Section -->
      <div class="detail-section">
        <h3>Game Instance Parameters</h3>
        
        <div v-if="instanceParametersLoading" class="loading-content">
          <p>Loading parameters...</p>
        </div>
        <div v-else class="parameters-table-container">
          <table class="parameters-table">
            <thead>
              <tr>
                <th>Parameter</th>
                <th>Type</th>
                <th>Current Value</th>
                <th>Default Value</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="param in allAvailableParameters" :key="param.config_key" class="parameter-row">
                <td class="parameter-name">
                  <div class="param-info">
                    <strong>{{ param.description }}</strong>
                    <small class="param-key">{{ param.config_key }}</small>
                  </div>
                </td>
                <td class="parameter-type">
                  <span class="type-badge" :class="'type-' + param.value_type">
                    {{ param.value_type }}
                  </span>
                </td>
                <td class="parameter-value">
                  <div v-if="getCurrentParameterValue(param.config_key)" class="current-value">
                    {{ getCurrentParameterValue(param.config_key) }}
                  </div>
                  <div v-else class="no-value">
                    <em>Not set</em>
                  </div>
                </td>
                <td class="parameter-default">
                  <span v-if="param.default_value" class="default-value">
                    {{ param.default_value }}
                  </span>
                  <span v-else class="no-default">
                    <em>None</em>
                  </span>
                </td>
                <td class="parameter-actions">
                  <button 
                    v-if="instance.status === 'created'"
                    @click="editParameterInline(param)"
                    class="btn-secondary btn-sm"
                  >
                    {{ getCurrentParameterValue(param.config_key) ? 'Edit' : 'Set' }}
                  </button>
                  <button 
                    v-if="getCurrentParameterValue(param.config_key) && instance.status === 'created'"
                    @click="removeParameterByKey(param.config_key)"
                    class="btn-danger btn-sm"
                  >
                    Remove
                  </button>
                  <span v-else-if="instance.status !== 'created'" class="locked-message">
                    <em>Locked</em>
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
          
          <div v-if="allAvailableParameters.length === 0" class="no-parameters">
            <p>No parameters available for this game type.</p>
          </div>
        </div>
      </div>

      <!-- Add/Edit Parameter Modal -->
      <div v-if="showEditParameterModal" class="modal-overlay" @click="closeParameterModal">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h3>Configure Parameter: {{ getEditingParameterDescription() }}</h3>
            <button @click="closeParameterModal" class="btn-close">&times;</button>
          </div>
          
          <form @submit.prevent="saveParameter" class="parameter-form">
            <div class="form-group">
              <label for="parameterKey">Parameter</label>
              <input 
                id="parameterKey" 
                :value="parameterForm.parameter_key" 
                type="text" 
                disabled
                class="disabled-input"
              />
              <small class="help-text">{{ getEditingParameterType() }}</small>
            </div>
            
            <div class="form-group">
              <label for="parameterValue">Value</label>
              <input 
                v-if="selectedParameterType === 'string'"
                id="parameterValue" 
                v-model="parameterForm.parameter_value" 
                type="text" 
                required
                :placeholder="getParameterPlaceholder()"
              />
              <input 
                v-else-if="selectedParameterType === 'integer'"
                id="parameterValue" 
                v-model="parameterForm.parameter_value" 
                type="number" 
                required
                :placeholder="getParameterPlaceholder()"
              />
              <select 
                v-else-if="selectedParameterType === 'boolean'"
                id="parameterValue" 
                v-model="parameterForm.parameter_value" 
                required
              >
                <option value="">Select value...</option>
                <option value="true">True</option>
                <option value="false">False</option>
              </select>
              <input 
                v-else
                id="parameterValue" 
                v-model="parameterForm.parameter_value" 
                type="text" 
                required
                :placeholder="getParameterPlaceholder()"
              />
              <small v-if="getEditingParameterDefault()" class="help-text">
                Default: {{ getEditingParameterDefault() }}
              </small>
            </div>
            
            <div class="form-actions">
              <Button type="submit" :disabled="savingParameter" variant="primary">
                {{ savingParameter ? 'Saving...' : 'Save Parameter' }}
              </Button>
              <Button type="button" @click="closeParameterModal" variant="secondary">
                Cancel
              </Button>
            </div>
            
            <div v-if="parameterError" class="error-message">
              {{ parameterError }}
            </div>
          </form>
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
import { useGameParametersStore } from '../../stores/gameParameters';
import Button from '../../components/Button.vue';

const route = useRoute();
const router = useRouter();
const gamesStore = useGamesStore();
const gameInstancesStore = useGameInstancesStore();
const gameInstanceParametersStore = useGameInstanceParametersStore();
const gameParametersStore = useGameParametersStore();

const gameId = computed(() => route.params.gameId);
const instanceId = computed(() => route.params.instanceId);
const selectedGame = computed(() => gamesStore.games.find(g => g.id === gameId.value));

const loading = ref(false);
const controlLoading = ref(false);
const instanceParametersLoading = ref(false);
const error = ref('');
const instance = ref(null);
const instanceParameters = ref([]);

// Parameter management state
const showEditParameterModal = ref(false);
const savingParameter = ref(false);
const parameterError = ref('');
const parameterForm = ref({
  parameter_key: '',
  parameter_value: ''
});
const editingParameterId = ref(null);

// Available parameters for the game type
const availableParameters = computed(() => {
  if (!selectedGame.value) return [];
  return gameParametersStore.getParametersByGameType(selectedGame.value.game_type);
});

// All available parameters, including those not yet configured for this instance
const allAvailableParameters = computed(() => {
  const allParams = [...availableParameters.value];
  // Add a 'configured' property to indicate if it's already in the instance's parameters
  allParams.forEach(param => {
    param.configured = instanceParameters.value.some(ip => ip.parameter_key === param.config_key);
  });
  return allParams;
});

// Selected parameter type for form validation
const selectedParameterType = computed(() => {
  if (!parameterForm.value.parameter_key) return '';
  const param = availableParameters.value.find(p => p.config_key === parameterForm.value.parameter_key);
  return param ? param.value_type : '';
});

onMounted(async () => {
  await loadInstance();
  await loadInstanceParameters();
  await loadGameParameters();
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
    await gameInstanceParametersStore.fetchGameInstanceParameters(gameId.value, instanceId.value);
    instanceParameters.value = gameInstanceParametersStore.getParametersByGameInstanceId(instanceId.value);
  } catch (err) {
    console.error('Failed to load instance parameters:', err);
  } finally {
    instanceParametersLoading.value = false;
  }
};

const loadGameParameters = async () => {
  if (!selectedGame.value) return;
  
  try {
    await gameParametersStore.fetchGameParameters();
  } catch (err) {
    console.error('Failed to load game parameters:', err);
  }
};

const formatParameterValue = (param) => {
  // For game instance parameters, we only have parameter_value as a string
  // The type validation is handled by the backend
  return param.parameter_value || 'No value set';
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

const getParameterPlaceholder = () => {
  const param = availableParameters.value.find(p => p.config_key === parameterForm.value.parameter_key);
  if (param) {
    if (param.value_type === 'string') return 'Enter a value (e.g., "Hello World")';
    if (param.value_type === 'integer') return 'Enter a number (e.g., 123)';
    if (param.value_type === 'boolean') return 'Select a value (true or false)';
    return 'Enter a value';
  }
  return 'Enter a value';
};

const isParameterAlreadyConfigured = (key) => {
  return instanceParameters.value.some(p => p.parameter_key === key);
};

const editParameterInline = (param) => {
  editingParameterId.value = null; // Clear any existing editing ID
  parameterForm.value.parameter_key = param.config_key;
  
  // Get current value if parameter is already configured, otherwise use default or empty
  const currentParam = instanceParameters.value.find(p => p.parameter_key === param.config_key);
  if (currentParam) {
    parameterForm.value.parameter_value = currentParam.parameter_value;
  } else {
    // Use default value if available, otherwise empty
    parameterForm.value.parameter_value = param.default_value || '';
  }
  
  showEditParameterModal.value = true;
};

const closeParameterModal = () => {
  showEditParameterModal.value = false;
  parameterForm.value = { parameter_key: '', parameter_value: '' };
  editingParameterId.value = null;
  parameterError.value = '';
};

const saveParameter = async () => {
  if (!parameterForm.value.parameter_key) {
    parameterError.value = 'Please select a parameter.';
    return;
  }

  if (!parameterForm.value.parameter_value) {
    parameterError.value = 'Parameter value cannot be empty.';
    return;
  }

  savingParameter.value = true;
  parameterError.value = '';

  try {
    // Check if parameter already exists for this instance
    const existingParam = instanceParameters.value.find(p => p.parameter_key === parameterForm.value.parameter_key);
    
    if (existingParam) {
      // Update existing parameter
      await gameInstanceParametersStore.updateGameInstanceParameter(gameId.value, instanceId.value, existingParam.id, {
        parameter_key: parameterForm.value.parameter_key,
        parameter_value: parameterForm.value.parameter_value
      });
      alert('Parameter updated successfully!');
    } else {
      // Create new parameter
      await gameInstanceParametersStore.createGameInstanceParameter(gameId.value, instanceId.value, {
        parameter_key: parameterForm.value.parameter_key,
        parameter_value: parameterForm.value.parameter_value
      });
      alert('Parameter added successfully!');
    }
    
    await loadInstanceParameters(); // Reload parameters after save
    closeParameterModal();
  } catch (err) {
    parameterError.value = err.message || 'Failed to save parameter.';
  } finally {
    savingParameter.value = false;
  }
};

const removeParameterByKey = async (key) => {
  if (!confirm('Are you sure you want to remove this parameter?')) {
    return;
  }
  try {
    const parameterToRemove = instanceParameters.value.find(p => p.parameter_key === key);
    if (parameterToRemove) {
      await gameInstanceParametersStore.deleteGameInstanceParameter(gameId.value, instanceId.value, parameterToRemove.id);
      alert('Parameter removed successfully!');
      await loadInstanceParameters();
    }
  } catch (err) {
    alert('Failed to remove parameter: ' + err.message);
  }
};

const getCurrentParameterValue = (key) => {
  const param = instanceParameters.value.find(p => p.parameter_key === key);
  return param ? param.parameter_value : null;
};

const getEditingParameterDescription = () => {
  const param = availableParameters.value.find(p => p.config_key === parameterForm.value.parameter_key);
  return param ? param.description : 'N/A';
};

const getEditingParameterType = () => {
  const param = availableParameters.value.find(p => p.config_key === parameterForm.value.parameter_key);
  return param ? param.value_type : 'N/A';
};

const getEditingParameterDefault = () => {
  const param = availableParameters.value.find(p => p.config_key === parameterForm.value.parameter_key);
  return param ? param.default_value : 'N/A';
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
  justify-content: space-between;
  align-items: center;
  padding: var(--space-md);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.parameter-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.parameter-actions {
  display: flex;
  gap: var(--space-xs);
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
  margin-top: var(--space-md);
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

.hint {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  margin-top: var(--space-xs);
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
  width: 90%;
  max-width: 500px;
  max-height: 90%;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-lg);
  border-bottom: 1px solid var(--color-border);
  background: var(--color-bg-light);
}

.modal-header h3 {
  margin: 0;
  font-size: var(--font-size-lg);
  color: var(--color-text);
}

.btn-close {
  background: none;
  border: none;
  font-size: var(--font-size-xl);
  color: var(--color-text-muted);
  cursor: pointer;
  padding: var(--space-xs);
}

.btn-close:hover {
  color: var(--color-text);
}

.parameter-form {
  padding: var(--space-lg);
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-group label {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
  margin-bottom: var(--space-xs);
  font-weight: var(--font-weight-medium);
}

.form-group select,
.form-group input {
  padding: var(--space-sm);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
  color: var(--color-text);
  background: var(--color-bg-secondary);
}

.form-group select:focus,
.form-group input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px var(--color-primary-light);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-md);
}

.error-message {
  color: var(--color-danger);
  font-size: var(--font-size-sm);
  margin-top: var(--space-sm);
}

.section-header-with-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-md);
}

.section-header-with-actions h3 {
  margin: 0;
}

.parameters-table-container {
  overflow-x: auto;
  margin-top: var(--space-md);
}

.parameters-table {
  width: 100%;
  border-collapse: collapse;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.parameters-table th,
.parameters-table td {
  padding: var(--space-md);
  text-align: left;
  border-bottom: 1px solid var(--color-border-light);
  vertical-align: top;
}

.parameters-table th {
  background: var(--color-bg-light);
  font-weight: var(--font-weight-bold);
  color: var(--color-text);
  position: sticky;
  top: 0;
  z-index: 1;
  border-bottom: 2px solid var(--color-border);
}

.parameter-row:last-child td {
  border-bottom: none;
}

.parameter-row:hover {
  background-color: var(--color-bg-light);
}

.parameter-name {
  min-width: 200px;
}

.param-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.param-key {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  font-weight: var(--font-weight-medium);
  font-family: monospace;
  background: var(--color-bg-light);
  padding: 2px 6px;
  border-radius: var(--radius-xs);
  display: inline-block;
}

.parameter-type {
  text-align: center;
  min-width: 80px;
}

.type-badge {
  padding: var(--space-xs) var(--space-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-bold);
  text-transform: uppercase;
  color: white;
  display: inline-block;
  min-width: 60px;
  text-align: center;
}

.type-string { background: var(--color-info); }
.type-integer { background: var(--color-success); }
.type-boolean { background: var(--color-warning); }
.type-object { background: var(--color-primary); }
.type-array { background: var(--color-info); }
.type-number { background: var(--color-success); }

.parameter-value,
.parameter-default {
  font-size: var(--font-size-sm);
  color: var(--color-text);
  font-weight: var(--font-weight-medium);
  min-width: 120px;
}

.current-value {
  color: var(--color-success);
  font-weight: var(--font-weight-bold);
}

.no-value {
  color: var(--color-text-muted);
  font-style: italic;
}

.default-value {
  color: var(--color-info);
  font-weight: var(--font-weight-medium);
}

.no-default {
  color: var(--color-text-muted);
  font-style: italic;
}

.parameter-actions {
  min-width: 120px;
  text-align: center;
}

.locked-message {
  color: var(--color-text-muted);
  font-style: italic;
  font-size: var(--font-size-sm);
}

.no-parameters {
  text-align: center;
  padding: var(--space-xl);
  color: var(--color-text-muted);
}

.disabled-input {
  background-color: var(--color-bg-light);
  color: var(--color-text-muted);
  cursor: not-allowed;
  border-color: var(--color-border-light);
}

.help-text {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  margin-top: var(--space-xs);
}
</style> 