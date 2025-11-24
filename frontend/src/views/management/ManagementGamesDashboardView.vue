<!--
  ManagementGamesDashboardView.vue
  Main dashboard for game management showing games and their instances.
-->
<template>
  <div class="games-dashboard">
    <PageHeader
      title="Games & Instances"
      titleLevel="h2"
      :showIcon="false"
      subtitle="Manage your game instances and monitor player activity"
    />

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
      <DataCard
        v-for="game in games"
        :key="game.id"
        :title="game.name"
        class="game-card"
      >
        <div class="game-info">
          <DataItem label="Game Type" :value="game.game_type" />
          <DataItem label="Description" :value="getGameDescription(game.game_type)" />
          <div class="game-stats">
            <DataItem label="Instances" :value="getGameInstanceCount(game.id)" />
            <DataItem label="Active" :value="getActiveInstanceCount(game.id)" />
          </div>
        </div>
        
        <template #primary>
          <Button @click="viewGameInstances(game)" variant="primary" size="small">
            Manage
          </Button>
        </template>
      </DataCard>
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
import Button from '../../components/Button.vue';
import DataCard from '../../components/DataCard.vue';
import DataItem from '../../components/DataItem.vue';
import PageHeader from '../../components/PageHeader.vue';

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
</script>

<style scoped>
.games-dashboard {
  max-width: 1200px;
  margin: 0 auto;
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
  min-height: 280px;
}

.game-info {
  margin-bottom: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.game-stats {
  display: flex;
  gap: var(--space-lg);
  margin-top: var(--space-md);
  padding-top: var(--space-md);
  border-top: 1px solid var(--color-border);
}
</style> 