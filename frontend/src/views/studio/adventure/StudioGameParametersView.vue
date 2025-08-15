<!--
  StudioGameParametersView.vue
  View for managing game parameters for adventure games in the studio.
-->
<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage game parameters.</p>
    </div>
    <div v-else class="game-table-section">
      <p class="game-context-name">Game: {{ selectedGame.name }}</p>
      <div class="section-header">
        <h2>Game Parameters</h2>
        <p class="section-description">
          Configure the available parameters for this adventure game. These parameters will be available 
          when game managers create instances of your game.
        </p>
      </div>

      <!-- Game Parameters Table -->
      <div class="parameters-section">
        <div class="section-actions">
          <button class="btn btn-primary" @click="openCreateModal" :disabled="parametersLoading">
            Create Parameter
          </button>
        </div>
        
        <div v-if="parametersLoading" class="loading-section">
          <p>Loading game parameters...</p>
        </div>
        
        <div v-else-if="mergedParameters.length > 0" class="parameters-table">
          <table>
            <thead>
              <tr>
                <th>Parameter</th>
                <th>Type</th>
                <th>Current Value</th>
                <th>Default Value</th>
                <th>Description</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="param in mergedParameters" :key="param.config_key">
                <td class="param-name">{{ formatParameterName(param.config_key) }}</td>
                <td class="param-type">{{ param.value_type }}</td>
                <td class="param-current-value">
                  <span v-if="param.current_value !== undefined" class="value">
                    {{ formatParameterValue(param.current_value, param.value_type) }}
                  </span>
                  <span v-else class="no-value">Not set</span>
                </td>
                <td class="param-default-value">
                  <span v-if="param.default_value" class="default">{{ param.default_value }}</span>
                  <span v-else class="no-default">None</span>
                </td>
                <td class="param-description">{{ param.description || 'No description' }}</td>
                <td class="param-actions">
                  <button 
                    v-if="param.current_value !== undefined" 
                    class="btn btn-sm btn-secondary" 
                    @click="openEditModal(param)"
                  >
                    Edit
                  </button>
                  <button 
                    v-else 
                    class="btn btn-sm btn-primary" 
                    @click="openCreateModalForParam(param)"
                  >
                    Set Value
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <div v-else class="no-parameters">
          <p>No parameters available for this game type.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { useGamesStore } from '../../../stores/games';

const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const parametersLoading = ref(false);
const availableParameterConfigurations = ref([]);
const gameParameters = ref([]);

// Fetch available parameter configurations from API
const fetchParameterConfigurations = async (gameType) => {
  try {
    const response = await fetch(`/api/v1/game-parameter-configurations?game_type=${gameType}`, {
      headers: {
        'Authorization': `Bearer ${gamesStore.sessionToken}`,
        'Content-Type': 'application/json'
      }
    });
    
    if (response.ok) {
      const data = await response.json();
      availableParameterConfigurations.value = data.data || [];
    } else {
      console.error('Failed to fetch parameter configurations');
      availableParameterConfigurations.value = [];
    }
  } catch (error) {
    console.error('Error fetching parameter configurations:', error);
    availableParameterConfigurations.value = [];
  }
};

// Fetch game-specific parameters from API
const fetchGameParameters = async (gameId) => {
  try {
    const response = await fetch(`/api/v1/games/${gameId}/parameters`, {
      headers: {
        'Authorization': `Bearer ${gamesStore.sessionToken}`,
        'Content-Type': 'application/json'
      }
    });
    
    if (response.ok) {
      const data = await response.json();
      gameParameters.value = data.data || [];
    } else {
      console.error('Failed to fetch game parameters');
      gameParameters.value = [];
    }
  } catch (error) {
    console.error('Error fetching game parameters:', error);
    gameParameters.value = [];
  }
};

// Merge available configurations with current game parameter values
const mergedParameters = computed(() => {
  if (!availableParameterConfigurations.value.length) return [];
  
  return availableParameterConfigurations.value.map(config => {
    // Find if this parameter has a value set for the current game
    const gameParam = gameParameters.value.find(p => p.config_key === config.config_key);
    
    return {
      ...config,
      current_value: gameParam?.value,
      parameter_id: gameParam?.id, // For editing existing parameters
    };
  });
});

// Helper function to format parameter names
const formatParameterName = (configKey) => {
  return configKey
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
};

// Helper function to format parameter values based on type
const formatParameterValue = (value, valueType) => {
  if (value === null || value === undefined) return 'Not set';
  
  switch (valueType) {
    case 'boolean':
      return value ? 'Yes' : 'No';
    case 'integer':
      return value.toString();
    case 'string':
      return value;
    case 'json':
      return JSON.stringify(value);
    default:
      return value.toString();
  }
};

// Modal and parameter management functions
const openCreateModal = () => {
  // TODO: Implement create parameter modal
  console.log('Open create parameter modal');
};

const openEditModal = (param) => {
  // TODO: Implement edit parameter modal
  console.log('Open edit modal for parameter:', param);
};

const openCreateModalForParam = (param) => {
  // TODO: Implement create value modal for specific parameter
  console.log('Open create value modal for parameter:', param);
};

// Watch for game selection changes
watch(
  () => selectedGame.value,
  async (newGame) => {
    if (newGame) {
      parametersLoading.value = true;
      try {
        // Fetch both parameter configurations and current game parameter values
        await Promise.all([
          fetchParameterConfigurations(newGame.game_type),
          fetchGameParameters(newGame.id)
        ]);
      } catch (error) {
        console.error('Failed to fetch game parameters:', error);
      } finally {
        parametersLoading.value = false;
      }
    }
  },
  { immediate: true }
);
</script>

<style scoped>
.game-table-section {
  margin-top: 20px;
}

.game-context-name {
  font-weight: bold;
  margin-bottom: 10px;
}

.section-header {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  margin-bottom: 20px;
}

.section-header h2 {
  margin: 0 0 10px 0;
}

.section-description, .subsection-description {
  color: var(--color-text-muted);
  font-size: 14px;
  margin-bottom: 15px;
  line-height: 1.5;
}

.parameters-section {
  margin-bottom: 30px;
}

.section-actions {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 20px;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 14px;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: background-color 0.2s ease;
}

.btn-primary {
  background: var(--color-primary);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.btn-secondary {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
  border: 1px solid var(--color-border);
}

.btn-secondary:hover:not(:disabled) {
  background: var(--color-bg-secondary);
}

.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.parameters-table {
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.parameters-table table {
  width: 100%;
  border-collapse: collapse;
}

.parameters-table th {
  background: var(--color-bg-secondary);
  padding: 12px;
  text-align: left;
  font-weight: 600;
  color: var(--color-text-primary);
  border-bottom: 1px solid var(--color-border);
}

.parameters-table td {
  padding: 12px;
  border-bottom: 1px solid var(--color-border);
  vertical-align: top;
}

.parameters-table tbody tr:last-child td {
  border-bottom: none;
}

.parameters-table tbody tr:hover {
  background: var(--color-bg-secondary);
}

.param-name {
  font-weight: 500;
  color: var(--color-text-primary);
}

.param-type {
  font-family: monospace;
  font-size: 12px;
  background: var(--color-bg-tertiary);
  padding: 2px 6px;
  border-radius: var(--radius-xs);
  color: var(--color-text-secondary);
}

.param-current-value .value {
  color: var(--color-text-primary);
  font-weight: 500;
}

.param-current-value .no-value {
  color: var(--color-text-muted);
  font-style: italic;
}

.param-default-value .default {
  color: var(--color-text-secondary);
  font-family: monospace;
  font-size: 12px;
}

.param-default-value .no-default {
  color: var(--color-text-muted);
  font-style: italic;
}

.param-description {
  color: var(--color-text-secondary);
  font-size: 14px;
  line-height: 1.4;
}

.param-actions {
  text-align: right;
  white-space: nowrap;
}

.loading-section,
.no-parameters {
  text-align: center;
  padding: 40px 20px;
  color: var(--color-text-muted);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}
</style>
