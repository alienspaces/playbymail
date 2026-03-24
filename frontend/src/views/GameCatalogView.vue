<template>
  <div class="game-catalog-view card">
    <div class="catalog-header">
      <h1>Game Catalog</h1>
      <p>Browse play-by-mail games with open enrollment. Click <strong>Join Game</strong> to play.</p>
    </div>

    <div v-if="loading" class="catalog-loading" data-testid="catalog-loading">
      Loading available games...
    </div>

    <div v-else-if="error" class="catalog-error" data-testid="catalog-error">
      <p>{{ error }}</p>
      <button class="retry-button" @click="fetchCatalog">Try again</button>
    </div>

    <div v-else-if="instances.length === 0" class="catalog-empty" data-testid="catalog-empty">
      <p>No games are currently available for enrollment. Check back soon.</p>
    </div>

    <div v-else class="catalog-games" data-testid="catalog-games">
      <div
        v-for="entry in instances"
        :key="entry.game_instance_id"
        class="catalog-game card"
        :data-testid="`instance-card-${entry.game_instance_id}`"
      >
        <div class="game-info">
          <h2 class="game-name">{{ entry.game_name }}</h2>
          <p v-if="entry.game_description" class="game-description">{{ entry.game_description }}</p>
          <p v-if="entry.account_name" class="game-host">Hosted by {{ entry.account_name }}</p>
          <div class="game-meta">
            <span class="game-type badge">{{ formatGameType(entry.game_type) }}</span>
            <span class="turn-duration">Turn: {{ entry.turn_duration_hours }}h</span>
            <span v-if="entry.required_player_count > 0" class="capacity">{{ entry.remaining_capacity }} {{ entry.remaining_capacity === 1 ? 'player' : 'players' }} needed</span>
          </div>
        </div>

        <div class="subscription-details">
          <div class="delivery-methods">
            <span v-if="entry.delivery_email" class="delivery-badge">Email</span>
            <span v-if="entry.delivery_physical_local" class="delivery-badge">Local</span>
            <span v-if="entry.delivery_physical_post" class="delivery-badge">Post</span>
          </div>
          <a
            :href="`/player/join-game/${entry.game_subscription_id}`"
            class="join-button"
            :data-testid="`join-button-${entry.game_subscription_id}`"
          >
            Join Game
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listCatalogGameInstances } from '../api/catalog'

const instances = ref([])
const loading = ref(false)
const error = ref(null)

function formatGameType(gameType) {
  const types = { adventure: 'Adventure', mecha: 'Mecha' }
  return types[gameType] ?? gameType
}

async function fetchCatalog() {
  loading.value = true
  error.value = null
  try {
    const res = await listCatalogGameInstances()
    instances.value = res.data ?? []
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
  width: 100%;
  margin: var(--space-lg) auto;
  padding: var(--space-xl);
}

@media (max-width: 600px) {
  .game-catalog-view {
    padding: var(--space-md);
    margin: var(--space-sm) auto;
  }

  .catalog-game {
    padding: var(--space-md);
  }

  .game-meta {
    flex-wrap: wrap;
  }

  .subscription-details {
    flex-wrap: wrap;
  }
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

.game-host {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted, #666);
  margin-bottom: var(--space-xs);
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

.subscription-details {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  padding-top: var(--space-sm);
  border-top: 1px solid var(--color-border);
}

.delivery-methods {
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

.capacity {
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
