<!--
  ManagementGamesDashboardView.vue
  Main dashboard for game management showing games and their instances.
-->
<template>
  <div class="games-dashboard">
    <div class="dashboard-header">
      <h2>Games & Instances</h2>
      <p>Manage your game instances and monitor player activity</p>
    </div>

    <!-- Loading state -->
    <div v-if="gamesStore.loading" class="loading-state">
      <p>Loading games...</p>
    </div>

    <!-- Error state -->
    <div v-else-if="gamesStore.error" class="error-state">
      <p>Error loading games: {{ gamesStore.error }}</p>
      <button @click="loadGames">Retry</button>
    </div>

    <!-- Games list -->
    <div v-else class="games-grid">
      <div v-for="game in games" :key="game.id" class="game-card">
        <div class="game-header">
          <h3>{{ game.name }}</h3>
          <span class="game-type">{{ game.game_type }}</span>
        </div>
        
        <div class="game-info">
          <p class="game-description">
            {{ getGameDescription(game.game_type) }}
          </p>
          <div class="game-stats">
            <div class="stat">
              <span class="stat-label">Instances:</span>
              <span class="stat-value">{{ getGameInstanceCount(game.id) }}</span>
            </div>
            <div class="stat">
              <span class="stat-label">Active:</span>
              <span class="stat-value">{{ getActiveInstanceCount(game.id) }}</span>
            </div>
          </div>
        </div>

        <div class="game-actions">
          <Button @click="viewGameInstances(game)" variant="primary" size="small">
            Manage
          </Button>
          <Button @click="createGameInstance(game)" variant="secondary" size="small">
            Create
          </Button>
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-if="!gamesStore.loading && !gamesStore.error && games.length === 0" class="empty-state">
      <h3>No Games Available</h3>
      <p>You don't have access to any games yet. Contact an administrator to get access to games.</p>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useGamesStore } from '../../stores/games';
import { useGameInstancesStore } from '../../stores/gameInstances';
import Button from '../../components/Button.vue'; // Added import for Button

const router = useRouter();
const gamesStore = useGamesStore();
const gameInstancesStore = useGameInstancesStore();

const games = computed(() => gamesStore.games);

onMounted(async () => {
  await loadGames();
});

const loadGames = async () => {
  try {
    await gamesStore.fetchGames();
    await gameInstancesStore.fetchAllGameInstances();
  } catch (error) {
    console.error('Failed to load games:', error);
  }
};

const getGameDescription = (gameType) => {
  const descriptions = {
    'adventure': 'Exploration and story-driven experiences with locations, items, and creatures.',
    'economic': 'Resource management and economic competition games.',
    'sports': 'Sports team management and competition games.',
    'mystery': 'Mystery solving and detective games.',
    'fantasy': 'Fantasy kingdom management and warfare games.'
  };
  return descriptions[gameType] || 'Custom game type';
};

const getGameInstanceCount = (gameId) => {
  // This would need to be implemented to get instance count per game
  // For now, return a placeholder
  return gameInstancesStore.gameInstances.filter(instance => instance.game_id === gameId).length;
};

const getActiveInstanceCount = (gameId) => {
  // This would need to be implemented to get active instance count per game
  // For now, return a placeholder
  return gameInstancesStore.gameInstances.filter(instance => 
    instance.game_id === gameId && 
    ['started'].includes(instance.status)
  ).length;
};

const viewGameInstances = (game) => {
  router.push(`/admin/games/${game.id}/instances`);
};

const createGameInstance = (game) => {
  router.push(`/admin/games/${game.id}/instances/create`);
};
</script>

<style scoped>
.games-dashboard {
  max-width: 1200px;
  margin: 0 auto;
}

.dashboard-header {
  margin-bottom: var(--space-xl);
  text-align: center;
}

.dashboard-header h2 {
  margin: 0 0 var(--space-sm) 0;
  font-size: var(--font-size-xl);
  color: var(--color-text);
}

.dashboard-header p {
  margin: 0;
  color: var(--color-text-muted);
  font-size: var(--font-size-md);
}

.loading-state,
.error-state,
.empty-state {
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

.games-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: var(--space-lg);
}

.game-card {
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-lg);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transition: transform 0.2s, box-shadow 0.2s;
  display: flex;
  flex-direction: column;
  min-height: 280px;
}

.game-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0,0,0,0.15);
}

.game-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--space-md);
}

.game-header h3 {
  margin: 0;
  font-size: var(--font-size-lg);
  color: var(--color-text);
  flex: 1;
  margin-right: var(--space-sm);
}

.game-type {
  background: var(--color-primary-light);
  color: var(--color-primary);
  padding: var(--space-xs) var(--space-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-bold);
  text-transform: uppercase;
  white-space: nowrap;
  flex-shrink: 0;
}

.game-info {
  margin-bottom: var(--space-lg);
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.game-description {
  margin: 0 0 var(--space-md) 0;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  line-height: 1.4;
  flex: 1;
}

.game-stats {
  display: flex;
  gap: var(--space-lg);
  margin-top: auto;
}

.stat {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-label {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  text-transform: uppercase;
  margin-bottom: var(--space-xs);
}

.stat-value {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
  color: var(--color-text);
}

.game-actions {
  display: flex;
  gap: var(--space-sm);
  flex-direction: column;
  width: 100%;
}

.game-actions .btn {
  width: 100%;
  justify-content: center;
  padding: var(--space-sm) var(--space-md);
  font-size: var(--font-size-sm);
  min-height: 44px;
}
</style> 