<!--
  ManagementCreateInstanceView.vue
  View for creating new game instances.
-->
<template>
  <div class="create-instance-view">
    <div class="view-header">
      <div class="header-content">
        <h2>Create New Game Instance</h2>
        <p>Configure settings for a new game session</p>
      </div>
      <div class="header-actions">
        <button @click="goBack" class="btn-secondary">
          Cancel
        </button>
      </div>
    </div>

    <div class="form-container">
      <form @submit.prevent="createInstance" class="instance-form">
        <div class="form-section">
          <h3>Game Information</h3>
          <div class="game-info">
            <div class="info-item">
              <span class="label">Game:</span>
              <span class="value">{{ selectedGame?.name }}</span>
            </div>
            <div class="info-item">
              <span class="label">Type:</span>
              <span class="value">{{ selectedGame?.game_type }}</span>
            </div>
          </div>
        </div>

        <div class="form-section">
          <h3>Instance Configuration</h3>
          <!-- TODO: Add per-instance configuration fields here -->          
        </div>

        <!-- Game Type Configuration Section -->
        <div v-if="gameConfigurations.length > 0" class="form-section">
          <h3>Game Configuration</h3>
          <p class="section-description">
            Configure game-specific parameters for this instance
          </p>
          
          <div class="configuration-fields">
            <ConfigurationField
              v-for="config in gameConfigurations"
              :key="config.id"
              :config="config"
              :field-id="`config-${config.config_key}`"
              v-model="form.configurations[config.config_key]"
              @validation-error="handleConfigValidationError"
            />
          </div>
        </div>

        <div class="form-actions">
          <button type="submit" :disabled="loading" class="btn-primary">
            <span v-if="loading">Creating...</span>
            <span v-else>Create Instance</span>
          </button>
          <button type="button" @click="goBack" class="btn-secondary">
            Cancel
          </button>
        </div>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useGamesStore } from '../../stores/games';
import { useGameInstancesStore } from '../../stores/gameInstances';
import { useGameConfigurationsStore } from '../../stores/gameConfigurations';
import ConfigurationField from '../../components/ConfigurationField.vue';

const route = useRoute();
const router = useRouter();
const gamesStore = useGamesStore();
const gameInstancesStore = useGameInstancesStore();
const gameConfigurationsStore = useGameConfigurationsStore();

const gameId = computed(() => route.params.gameId);
const selectedGame = computed(() => gamesStore.games.find(g => g.id === gameId.value));
const gameConfigurations = computed(() => {
  if (!selectedGame.value) return [];
  return gameConfigurationsStore.getConfigurationsByGameType(selectedGame.value.game_type);
});

const loading = ref(false);
const error = ref('');

const form = ref({
  configurations: {}
});

const configValidationErrors = ref({});

onMounted(async () => {
  if (!selectedGame.value) {
    await gamesStore.fetchGames();
  }
  
  // Fetch game configurations for the selected game type
  if (selectedGame.value) {
    await gameConfigurationsStore.fetchGameConfigurationsByGameType(selectedGame.value.game_type);
  }
});

const createInstance = async () => {
  if (!selectedGame.value) {
    error.value = 'Game not found';
    return;
  }

  loading.value = true;
  error.value = '';

  try {
    const instanceData = {
      game_id: gameId.value,
      configurations: form.value.configurations
    };

    await gameInstancesStore.createGameInstance(gameId.value, instanceData);
    router.push(`/management/games/${gameId.value}/instances`);
  } catch (err) {
    error.value = err.message || 'Failed to create game instance';
  } finally {
    loading.value = false;
  }
};

const handleConfigValidationError = (configKey, error) => {
  if (error) {
    configValidationErrors.value[configKey] = error;
  } else {
    delete configValidationErrors.value[configKey];
  }
};

const goBack = () => {
  router.push(`/admin/games/${gameId.value}/instances`);
};
</script>

<style scoped>
.create-instance-view {
  max-width: 800px;
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

.form-container {
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-xl);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.instance-form {
  max-width: 600px;
}

.form-section {
  margin-bottom: var(--space-xl);
}

.form-section h3 {
  margin: 0 0 var(--space-lg) 0;
  font-size: var(--font-size-lg);
  color: var(--color-text);
  padding-bottom: var(--space-sm);
  border-bottom: 1px solid var(--color-border);
}

.game-info {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-md);
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.info-item .label {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
  font-weight: var(--font-weight-medium);
}

.info-item .value {
  font-size: var(--font-size-md);
  color: var(--color-text);
  font-weight: var(--font-weight-medium);
}

.form-group {
  margin-bottom: var(--space-lg);
}

.form-group label {
  display: block;
  margin-bottom: var(--space-sm);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text);
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: var(--space-sm);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
  font-family: inherit;
  background: var(--color-bg);
  color: var(--color-text);
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px var(--color-primary-light);
}

.form-group textarea {
  resize: vertical;
  min-height: 100px;
}

.help-text {
  margin: var(--space-xs) 0 0 0;
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  line-height: 1.4;
}

.section-description {
  margin-bottom: var(--space-lg);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.configuration-fields {
  display: grid;
  gap: var(--space-md);
}

.form-actions {
  display: flex;
  gap: var(--space-md);
  margin-top: var(--space-xl);
  padding-top: var(--space-lg);
  border-top: 1px solid var(--color-border);
}

.btn-primary,
.btn-secondary {
  padding: var(--space-sm) var(--space-lg);
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

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: var(--color-bg-light);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}

.btn-secondary:hover {
  background: var(--color-border);
}

.error-message {
  margin-top: var(--space-lg);
  padding: var(--space-md);
  background: var(--color-danger-light);
  color: var(--color-danger);
  border: 1px solid var(--color-danger);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
}
</style> 