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
        <Button @click="goBack" variant="secondary">
          Cancel
        </Button>
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
            <div class="info-item">
              <span class="label">Default Turn Duration:</span>
              <span class="value">{{ formatTurnDuration(selectedGame?.turn_duration_hours) }}</span>
            </div>
          </div>
        </div>

        <div class="form-section">
          <h3>Instance Configuration</h3>
          <p class="section-description">
            Basic instance settings. Game-specific parameters can be configured after the instance is created.
          </p>
          <!-- TODO: Add basic instance configuration fields here (e.g., name, description) -->
        </div>

        <!-- Game Type Configuration Section - REMOVED -->
        <!-- Game parameters are now configured after instance creation, not during creation -->

        <div class="form-actions">
          <Button type="submit" :disabled="loading" variant="primary">
            <span v-if="loading">Creating...</span>
            <span v-else>Create Instance</span>
          </Button>
          <Button type="button" @click="goBack" variant="secondary">
            Cancel
          </Button>
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
import Button from '../../components/Button.vue';
// Removed: import { useGameParametersStore } from '../../stores/gameParameters';
// Removed: import ConfigurationField from '../../components/ConfigurationField.vue';

const route = useRoute();
const router = useRouter();
const gamesStore = useGamesStore();
const gameInstancesStore = useGameInstancesStore();
// Removed: const gameParametersStore = useGameParametersStore();

const gameId = computed(() => route.params.gameId);
const selectedGame = computed(() => gamesStore.games.find(g => g.id === gameId.value));
// Removed: const gameParameters = computed(() => {
//   if (!selectedGame.value) return [];
//   return gameParametersStore.getParametersByGameType(selectedGame.value.game_type);
// });

const loading = ref(false);
const error = ref('');

// Removed: const form = ref({
//   // Basic instance configuration
//   name: '',
//   description: '',
//   // Removed: parameters: {} // Game parameters are now configured after instance creation
// });

// Removed: const configValidationErrors = ref({});

// Helper function to format turn duration
const formatTurnDuration = (hours) => {
  if (!hours) return 'Not set'
  if (hours % (24 * 7) === 0) {
    const weeks = hours / (24 * 7)
    return `${weeks} week${weeks === 1 ? '' : 's'}`
  }
  if (hours % 24 === 0) {
    const days = hours / 24
    return `${days} day${days === 1 ? '' : 's'}`
  }
  return `${hours} hour${hours === 1 ? '' : 's'}`
};

onMounted(async () => {
  if (!selectedGame.value) {
    await gamesStore.fetchGames();
  }
  
  // Fetch game parameters for the selected game type
  if (selectedGame.value) {
    // Removed: await gameParametersStore.fetchGameParametersByGameType(selectedGame.value.game_type);
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
      // Game parameters are now configured after instance creation
    };

    const createdInstance = await gameInstancesStore.createGameInstance(gameId.value, instanceData);
    
    // Redirect to the instance details page instead of the instances list
    if (createdInstance && createdInstance.id) {
      console.log('Instance created successfully, redirecting to:', createdInstance.id);
      router.push(`/admin/games/${gameId.value}/instances/${createdInstance.id}`);
    } else {
      console.warn('Instance created but no ID returned, falling back to instances list');
      // Fallback to instances list if we don't get the instance ID
      router.push(`/admin/games/${gameId.value}/instances`);
    }
  } catch (err) {
    console.error('Failed to create instance:', err);
    error.value = err.message || 'Failed to create game instance';
  } finally {
    loading.value = false;
  }
};

// Removed: const handleConfigValidationError = (configKey, error) => {
//   if (error) {
//     configValidationErrors.value[configKey] = error;
//   } else {
//     delete configValidationErrors.value[configKey];
//   }
// };

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