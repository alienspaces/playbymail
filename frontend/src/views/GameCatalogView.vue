<template>
  <div class="game-catalog-view card">
    <div class="catalog-header">
      <h1>Game Catalog</h1>
      <p>Browse play-by-mail games with open enrollment. Click <strong>Join Game</strong> on an available instance to play.</p>
    </div>

    <div v-if="loading" class="catalog-loading" data-testid="catalog-loading">
      Loading available games...
    </div>

    <div v-else-if="error" class="catalog-error" data-testid="catalog-error">
      <p>{{ error }}</p>
      <button class="retry-button" @click="fetchCatalog">Try again</button>
    </div>

    <div v-else-if="games.length === 0" class="catalog-empty" data-testid="catalog-empty">
      <p>No games are currently available for enrollment. Check back soon.</p>
    </div>

    <div v-else class="catalog-games" data-testid="catalog-games">
      <div
        v-for="game in games"
        :key="game.id"
        class="catalog-game card"
        :data-testid="`game-card-${game.id}`"
      >
        <div class="game-info">
          <h2 class="game-name">{{ game.name }}</h2>
          <p v-if="game.description" class="game-description">{{ game.description }}</p>
          <div class="game-meta">
            <span class="game-type badge">{{ formatGameType(game.game_type) }}</span>
            <span class="turn-duration">Turn: {{ game.turn_duration_hours }}h</span>
          </div>
        </div>

        <div class="game-instances">
          <h3 class="instances-heading">Available Instances</h3>
          <div
            v-for="instance in game.available_instances"
            :key="instance.id"
            class="instance-row"
            :data-testid="`instance-${instance.id}`"
          >
            <div class="instance-delivery">
              <span v-if="instance.delivery_email" class="delivery-badge">Email</span>
              <span v-if="instance.delivery_physical_local" class="delivery-badge">Local</span>
              <span v-if="instance.delivery_physical_post" class="delivery-badge">Post</span>
            </div>
            <div v-if="instance.required_player_count > 0" class="instance-capacity">
              {{ instance.player_count }} / {{ instance.required_player_count }} players
            </div>
            <a
              :href="`/player/join-game/${instance.id}`"
              class="join-button"
              :data-testid="`join-button-${instance.id}`"
            >
              Join Game
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listCatalogGames } from '../api/catalog'

const games = ref([])
const loading = ref(false)
const error = ref(null)

function formatGameType(gameType) {
  const types = { adventure: 'Adventure' }
  return types[gameType] ?? gameType
}

async function fetchCatalog() {
  loading.value = true
  error.value = null
  try {
    const res = await listCatalogGames()
    games.value = res.data ?? []
  } catch (err) {
    error.value = err.message || 'Failed to load the game catalog. Please try again.'
  } finally {
    loading.value = false
  }
}

onMounted(fetchCatalog)
</script>

<style scoped>
.game-catalog-view {
  max-width: 900px;
  margin: var(--space-lg) auto;
  padding: var(--space-xl);
}

.catalog-header {
  margin-bottom: var(--space-xl);
}

.catalog-header h1 {
  font-size: var(--font-size-xl);
  margin-bottom: var(--space-md);
}

.catalog-loading,
.catalog-error,
.catalog-empty {
  padding: var(--space-xl);
  text-align: center;
  color: var(--color-text-muted, #666);
}

.catalog-error .retry-button {
  margin-top: var(--space-md);
  padding: var(--space-sm) var(--space-lg);
  background: transparent;
  color: var(--color-button);
  border: 2px solid var(--color-button);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-weight: var(--font-weight-bold);
}

.catalog-games {
  display: flex;
  flex-direction: column;
  gap: var(--space-lg);
}

.catalog-game {
  padding: var(--space-lg);
  background: var(--color-bg-alt);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
}

.game-info {
  margin-bottom: var(--space-md);
}

.game-name {
  font-size: var(--font-size-lg);
  margin-bottom: var(--space-sm);
}

.game-description {
  font-size: var(--font-size-md);
  margin-bottom: var(--space-sm);
  color: var(--color-text-muted, #444);
}

.game-meta {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  font-size: var(--font-size-sm);
}

.badge {
  padding: 2px var(--space-sm);
  border-radius: var(--radius-sm);
  background: var(--color-primary, #3b82f6);
  color: #fff;
  font-size: var(--font-size-xs, 0.75rem);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.turn-duration {
  color: var(--color-text-muted, #666);
}

.instances-heading {
  font-size: var(--font-size-md);
  margin-bottom: var(--space-sm);
  font-weight: var(--font-weight-bold);
}

.instance-row {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  padding: var(--space-sm) 0;
  border-top: 1px solid var(--color-border);
}

.instance-delivery {
  display: flex;
  gap: var(--space-xs);
  flex: 1;
}

.delivery-badge {
  padding: 2px var(--space-sm);
  background: var(--color-bg, #f5f5f5);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs, 0.75rem);
}

.instance-capacity {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted, #666);
  white-space: nowrap;
}

.join-button {
  display: inline-block;
  padding: var(--space-sm) var(--space-lg);
  background: transparent;
  color: var(--color-button);
  border: 2px solid var(--color-button);
  border-radius: var(--radius-sm);
  text-decoration: none;
  font-weight: var(--font-weight-bold);
  font-size: var(--font-size-sm);
  white-space: nowrap;
  transition: all 0.2s;
}

.join-button:hover {
  background: var(--color-button);
  color: var(--color-text-light);
}
</style>
