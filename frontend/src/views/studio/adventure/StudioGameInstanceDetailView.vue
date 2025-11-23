<!--
  StudioGameInstanceDetailView.vue
  This component manages a specific game instance with runtime controls.
-->
<template>
  <div>
    <div v-if="!selectedGame || !gameInstance">
      <p>Loading game instance...</p>
    </div>
    <div v-else class="game-instance-detail">
      <div class="header">
        <GameContext :gameName="selectedGame.name" />
        <div class="instance-header">
          <h1>{{ gameInstance.name }}</h1>
          <p class="description">{{ gameInstance.description || 'No description' }}</p>
        </div>
      </div>

      <!-- Status and Progress Section -->
      <div class="status-section">
        <h2>Game Status</h2>
        <div class="status-grid">
          <div class="status-item">
            <label>Status:</label>
            <span class="status-badge" :class="getStatusClass(gameInstance.status)">
              {{ gameInstance.status }}
            </span>
          </div>
          <div class="status-item">
            <label>Current Turn:</label>
            <span>{{ gameInstance.current_turn }}</span>
          </div>
          
          <div class="status-item" v-if="gameInstance.started_at">
            <label>Started:</label>
            <span>{{ formatDate(gameInstance.started_at) }}</span>
          </div>
          <div class="status-item" v-if="gameInstance.next_turn_due_at">
            <label>Next Turn Due:</label>
            <span>{{ formatDate(gameInstance.next_turn_due_at) }}</span>
          </div>
        </div>
      </div>

      <!-- Runtime Controls -->
      <div class="controls-section">
        <h2>Game Controls</h2>
        <div class="controls-grid">
          <button 
            v-if="gameInstance.status === 'created'"
            @click="startGame"
            :disabled="gameInstancesStore.loading"
            class="control-btn start-btn"
          >
            Start Game
          </button>
          <button 
            v-if="gameInstance.status === 'started'"
            @click="pauseGame"
            :disabled="gameInstancesStore.loading"
            class="control-btn pause-btn"
          >
            Pause Game
          </button>
          <button 
            v-if="gameInstance.status === 'paused'"
            @click="resumeGame"
            :disabled="gameInstancesStore.loading"
            class="control-btn resume-btn"
          >
            Resume Game
          </button>
          <button 
            v-if="['created', 'started', 'paused'].includes(gameInstance.status)"
            @click="cancelGame"
            :disabled="gameInstancesStore.loading"
            class="control-btn cancel-btn"
          >
            Cancel Game
          </button>
        </div>
      </div>

      <!-- Turn Management -->
      <div class="turn-section" v-if="gameInstance.status !== 'created'">
        <h2>Turn Management</h2>
        <div class="turn-info">
          <p><strong>Current Turn:</strong> {{ gameInstance.current_turn }}</p>
          <p v-if="gameInstance.last_turn_processed_at">
            <strong>Last Turn Processed:</strong> {{ formatDate(gameInstance.last_turn_processed_at) }}
          </p>
          <p v-if="gameInstance.next_turn_due_at">
            <strong>Next Turn Due:</strong> {{ formatDate(gameInstance.next_turn_due_at) }}
          </p>
        </div>
      </div>

      <!-- Error Display -->
      <div v-if="gameInstancesStore.error" class="error-section">
        <h3>Error</h3>
        <p class="error">{{ gameInstancesStore.error }}</p>
      </div>

      <!-- Loading Indicator -->
      <div v-if="gameInstancesStore.loading" class="loading-section">
        <p>Loading...</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useGameInstancesStore } from '../../../stores/gameInstances';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import GameContext from '../../../components/GameContext.vue';

const route = useRoute();
const router = useRouter();
const gameInstancesStore = useGameInstancesStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const gameInstance = ref(null);

// Computed properties
const gameId = computed(() => route.params.gameId);
const instanceId = computed(() => route.params.instanceId);

// Watch for route changes
watch(
  () => [gameId.value, instanceId.value],
  async ([newGameId, newInstanceId]) => {
    if (newGameId && newInstanceId) {
      await loadGameInstance();
    }
  },
  { immediate: true }
);

// Watch for game selection changes
watch(
  () => selectedGame.value,
  (newGame) => {
    if (newGame && newGame.id !== gameId.value) {
      // Redirect if game doesn't match
      router.push('/studio');
    }
  }
);

async function loadGameInstance() {
  try {
    gameInstance.value = await gameInstancesStore.getGameInstance(gameId.value, instanceId.value);
  } catch (err) {
    console.error('Failed to load game instance:', err);
  }
}

async function startGame() {
  try {
    await gameInstancesStore.startGameInstance(gameId.value, instanceId.value);
    await loadGameInstance(); // Refresh the instance data
  } catch (err) {
    console.error('Failed to start game:', err);
  }
}

async function pauseGame() {
  try {
    await gameInstancesStore.pauseGameInstance(gameId.value, instanceId.value);
    await loadGameInstance(); // Refresh the instance data
  } catch (err) {
    console.error('Failed to pause game:', err);
  }
}

async function resumeGame() {
  try {
    await gameInstancesStore.resumeGameInstance(gameId.value, instanceId.value);
    await loadGameInstance(); // Refresh the instance data
  } catch (err) {
    console.error('Failed to resume game:', err);
  }
}

async function cancelGame() {
  if (!confirm('Are you sure you want to cancel this game? This action cannot be undone.')) {
    return;
  }
  
  try {
    await gameInstancesStore.cancelGameInstance(gameId.value, instanceId.value);
    await loadGameInstance(); // Refresh the instance data
  } catch (err) {
    console.error('Failed to cancel game:', err);
  }
}

function getStatusClass(status) {
  const statusClasses = {
    'created': 'status-created',
    'started': 'status-started',
    'paused': 'status-paused',
    'completed': 'status-completed',
    'cancelled': 'status-cancelled'
  };
  return statusClasses[status] || 'status-unknown';
}

function formatDate(dateString) {
  if (!dateString) return 'N/A';
  return new Date(dateString).toLocaleString();
}

onMounted(async () => {
  if (gameId.value && instanceId.value) {
    await loadGameInstance();
  }
});
</script>

<style scoped>
.game-instance-detail {
  max-width: 800px;
  margin: 0 auto;
  padding: 1rem;
}

.header {
  margin-bottom: 2rem;
}

.instance-header h1 {
  margin: 1rem 0 0.5rem 0;
  color: #333;
}

.description {
  color: #666;
  font-style: italic;
  margin: 0;
}

.status-section,
.controls-section,
.turn-section {
  margin-bottom: 2rem;
  padding: 1.5rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: #f9f9f9;
}

.status-section h2,
.controls-section h2,
.turn-section h2 {
  margin-top: 0;
  margin-bottom: 1rem;
  color: #333;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.status-item {
  display: flex;
  flex-direction: column;
}

.status-item label {
  font-weight: 600;
  color: #555;
  margin-bottom: 0.25rem;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-created {
  background: var(--color-info-light);
  color: var(--color-primary);
}

.status-started {
  background: var(--color-success-light);
  color: var(--color-success);
}

.status-paused {
  background: #fff8e1;
  color: #fbc02d;
}

.status-completed {
  background: #f3e5f5;
  color: #7b1fa2;
}

.status-cancelled {
  background: #ffebee;
  color: #d32f2f;
}

.status-unknown {
  background: #f5f5f5;
  color: #757575;
}

.controls-grid {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.control-btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 600;
  transition: all 0.2s;
}

.control-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.start-btn {
  background: #28a745;
  color: white;
}

.start-btn:hover:not(:disabled) {
  background: #218838;
}

.pause-btn {
  background: #ffc107;
  color: #212529;
}

.pause-btn:hover:not(:disabled) {
  background: #e0a800;
}

.resume-btn {
  background: #17a2b8;
  color: white;
}

.resume-btn:hover:not(:disabled) {
  background: #138496;
}

.cancel-btn {
  background: #dc3545;
  color: white;
}

.cancel-btn:hover:not(:disabled) {
  background: #c82333;
}

.turn-info p {
  margin: 0.5rem 0;
  color: #333;
}

.error-section {
  margin-top: 2rem;
  padding: 1rem;
  background: #f8d7da;
  border: 1px solid #f5c6cb;
  border-radius: 4px;
}

.error-section h3 {
  margin-top: 0;
  color: #721c24;
}

.error {
  color: #721c24;
  margin: 0;
}

.loading-section {
  text-align: center;
  padding: 2rem;
  color: #666;
}
</style> 