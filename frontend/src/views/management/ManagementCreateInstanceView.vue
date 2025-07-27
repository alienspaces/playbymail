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
          
          <div class="form-group">
            <label for="turn-deadline">Turn Deadline (hours)</label>
            <input
              id="turn-deadline"
              v-model.number="form.turnDeadlineHours"
              type="number"
              min="1"
              max="8760"
              required
            />
            <p class="help-text">How many hours players have to submit their turn (1-8760 hours)</p>
          </div>

          <div class="form-group">
            <label for="max-turns">Maximum Turns (optional)</label>
            <input
              id="max-turns"
              v-model.number="form.maxTurns"
              type="number"
              min="1"
              placeholder="Leave empty for unlimited"
            />
            <p class="help-text">Maximum number of turns before the game ends automatically</p>
          </div>

          <div class="form-group">
            <label for="game-config">Game Configuration (JSON)</label>
            <textarea
              id="game-config"
              v-model="form.gameConfig"
              rows="6"
              placeholder='{"player_limit": 10, "starting_location": "town_square"}'
            ></textarea>
            <p class="help-text">Optional JSON configuration specific to this game type</p>
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

const route = useRoute();
const router = useRouter();
const gamesStore = useGamesStore();
const gameInstancesStore = useGameInstancesStore();

const gameId = computed(() => route.params.gameId);
const selectedGame = computed(() => gamesStore.games.find(g => g.id === gameId.value));

const loading = ref(false);
const error = ref('');

const form = ref({
  turnDeadlineHours: 168, // 7 days default
  maxTurns: null,
  gameConfig: ''
});

onMounted(async () => {
  if (!selectedGame.value) {
    await gamesStore.fetchGames();
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
    // Parse game config if provided
    let gameConfig = null;
    if (form.value.gameConfig.trim()) {
      try {
        gameConfig = JSON.parse(form.value.gameConfig);
      } catch {
        error.value = 'Invalid JSON in game configuration';
        loading.value = false;
        return;
      }
    }

    const instanceData = {
      game_id: gameId.value,
      turn_deadline_hours: form.value.turnDeadlineHours,
      max_turns: form.value.maxTurns || null,
      game_config: gameConfig
    };

    await gameInstancesStore.createGameInstance(gameId.value, instanceData);
    
    // Navigate to the instances list
    router.push(`/admin/games/${gameId.value}/instances`);
  } catch (err) {
    error.value = err.message || 'Failed to create game instance';
  } finally {
    loading.value = false;
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