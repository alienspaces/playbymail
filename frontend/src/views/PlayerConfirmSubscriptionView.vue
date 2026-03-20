<template>
  <div class="confirm-subscription-view">

    <div v-if="loading" class="confirm-loading" data-testid="confirm-loading">
      <p>Confirming your subscription...</p>
    </div>

    <div v-else-if="error" class="confirm-error card" data-testid="confirm-error">
      <h1>Something went wrong</h1>
      <p class="error-message" data-testid="error-message">{{ error }}</p>
      <a href="/games" class="catalog-link" data-testid="link-browse-games">Browse games</a>
    </div>

    <ConfirmationCard
      v-else
      title="You're confirmed!"
      message="Your subscription has been confirmed. You will receive further instructions soon."
      data-testid="confirm-success"
    />

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { approveSubscription } from '../api/approveSubscription'
import ConfirmationCard from '../components/ConfirmationCard.vue'

const route = useRoute()

const loading = ref(true)
const error = ref(null)

async function confirm() {
  const gameSubscriptionId = route.params.game_subscription_id
  const email = route.query.email

  if (!email) {
    error.value = 'Confirmation link is invalid. Please check your email and try again.'
    loading.value = false
    return
  }

  try {
    await approveSubscription(gameSubscriptionId, email)
  } catch (err) {
    error.value = err.message || 'Failed to confirm your subscription. Please try again.'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  confirm()
})
</script>

<style scoped>
.confirm-subscription-view {
  width: 100%;
  max-width: 600px;
  margin: var(--space-lg, 1.5rem) auto;
  padding: var(--space-md, 1rem);
  box-sizing: border-box;
}

.confirm-loading {
  padding: var(--space-xl, 2rem);
  text-align: center;
  color: var(--color-text-muted, #666);
}

.confirm-error {
  max-width: 600px;
  margin: 0 auto;
  padding: var(--space-xl, 2rem);
  text-align: center;
}

.confirm-error h1 {
  font-size: var(--font-size-xl, 1.5rem);
  margin-bottom: var(--space-md, 1rem);
}

.error-message {
  color: var(--color-warning-dark, #c00);
  background: var(--color-warning-light, #fff3f3);
  padding: var(--space-sm, 0.5rem) var(--space-md, 1rem);
  border: 1px solid var(--color-warning, #f99);
  border-radius: var(--radius-sm, 4px);
  font-size: var(--font-size-sm, 0.875rem);
  margin-bottom: var(--space-xl, 2rem);
}

.catalog-link {
  display: inline-block;
  padding: var(--space-sm, 0.5rem) var(--space-xl, 2rem);
  background: transparent;
  color: var(--color-button, #006ecd);
  border: 2px solid var(--color-button, #006ecd);
  border-radius: var(--radius-sm, 4px);
  text-decoration: none;
  font-weight: var(--font-weight-bold, 600);
  font-size: var(--font-size-sm, 0.875rem);
  transition: all 0.2s;
}

.catalog-link:hover {
  background: var(--color-button, #006ecd);
  color: var(--color-text-light, #fff);
}
</style>
