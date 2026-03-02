<template>
  <div class="join-game-view">

    <!-- Loading game info -->
    <div v-if="loading" class="join-loading" data-testid="join-loading">
      <p>Loading game information...</p>
    </div>

    <!-- Failed to load game info -->
    <div v-else-if="loadError" class="join-error card" data-testid="join-load-error">
      <p class="error-message">{{ loadError }}</p>
      <a href="/games" class="catalog-link">Browse other games</a>
    </div>

    <!-- Step 1: Game info + delivery method selection -->
    <div v-else-if="step === 'info'" class="join-step card" data-testid="step-info">
      <h1 class="game-title">{{ gameInfo.game_name }}</h1>
      <p v-if="gameInfo.game_description" class="game-description">{{ gameInfo.game_description }}</p>

      <div class="game-meta">
        <span class="badge">{{ formatGameType(gameInfo.game_type) }}</span>
        <span class="turn-duration">Turn: {{ gameInfo.turn_duration_hours }}h</span>
      </div>

      <div class="instance-info">
        <div v-if="gameInfo.total_capacity > 0" class="capacity-info">
          <span class="label">Players:</span>
          <span>{{ gameInfo.total_players }} / {{ gameInfo.total_capacity }}</span>
        </div>

        <div class="delivery-methods">
          <span class="label">Play by:</span>
          <span v-if="gameInfo.delivery_email" class="delivery-badge">Email</span>
          <span v-if="gameInfo.delivery_physical_local" class="delivery-badge">Local</span>
          <span v-if="gameInfo.delivery_physical_post" class="delivery-badge">Post</span>
        </div>
      </div>

      <button class="primary-button" @click="step = 'contact'" data-testid="btn-join">
        Join this Game
      </button>
    </div>

    <!-- Step 2: Contact details form -->
    <div v-else-if="step === 'contact'" class="join-step card" data-testid="step-contact">
      <h1>Your Details</h1>
      <p class="step-description">Enter your contact information to join <strong>{{ gameInfo.game_name }}</strong>.</p>

      <form @submit.prevent="onSubmit" class="join-form" novalidate>

        <div class="form-group">
          <label for="email">Email address <span class="required">*</span></label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            required
            autocomplete="email"
            data-testid="input-email"
          />
        </div>

        <div class="form-group">
          <label for="name">Full name <span class="required">*</span></label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            required
            autocomplete="name"
            data-testid="input-name"
          />
        </div>

        <div class="form-group">
          <label for="address-line1">Address line 1 <span class="required">*</span></label>
          <input
            id="address-line1"
            v-model="form.postal_address_line1"
            type="text"
            required
            autocomplete="address-line1"
            data-testid="input-address-line1"
          />
        </div>

        <div class="form-group">
          <label for="address-line2">Address line 2</label>
          <input
            id="address-line2"
            v-model="form.postal_address_line2"
            type="text"
            autocomplete="address-line2"
            data-testid="input-address-line2"
          />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="state">State / Province <span class="required">*</span></label>
            <input
              id="state"
              v-model="form.state_province"
              type="text"
              required
              autocomplete="address-level1"
              data-testid="input-state"
            />
          </div>
          <div class="form-group">
            <label for="postal-code">Postal code <span class="required">*</span></label>
            <input
              id="postal-code"
              v-model="form.postal_code"
              type="text"
              required
              autocomplete="postal-code"
              data-testid="input-postal-code"
            />
          </div>
        </div>

        <div class="form-group">
          <label for="country">Country <span class="required">*</span></label>
          <input
            id="country"
            v-model="form.country"
            type="text"
            required
            autocomplete="country-name"
            data-testid="input-country"
          />
        </div>

        <div v-if="availableDeliveryMethods.length > 1" class="form-group" data-testid="delivery-selection">
          <label>How would you like to play? <span class="required">*</span></label>
          <div class="delivery-options">
            <label v-if="gameInfo.delivery_email" class="delivery-option">
              <input
                type="radio"
                v-model="selectedDelivery"
                value="email"
                name="delivery"
                data-testid="delivery-email"
              />
              By Email
            </label>
            <label v-if="gameInfo.delivery_physical_local" class="delivery-option">
              <input
                type="radio"
                v-model="selectedDelivery"
                value="local"
                name="delivery"
                data-testid="delivery-local"
              />
              Local pickup
            </label>
            <label v-if="gameInfo.delivery_physical_post" class="delivery-option">
              <input
                type="radio"
                v-model="selectedDelivery"
                value="post"
                name="delivery"
                data-testid="delivery-post"
              />
              By Post
            </label>
          </div>
        </div>

        <p v-if="submitError" class="error-message" data-testid="submit-error">{{ submitError }}</p>

        <div class="form-actions">
          <button type="button" class="secondary-button" @click="step = 'info'" data-testid="btn-back">
            Back
          </button>
          <button type="submit" class="primary-button" :disabled="submitting" data-testid="btn-submit">
            {{ submitting ? 'Joining...' : 'Join Game' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Step 3: Success confirmation -->
    <div v-else-if="step === 'success'" class="join-step join-success card" data-testid="step-success">
      <h1>You're in!</h1>
      <p class="success-message">
        You have successfully joined <strong>{{ gameInfo.game_name }}</strong>.
        You will receive further instructions by {{ deliveryLabel }}.
      </p>
      <a href="/games" class="catalog-link" data-testid="link-browse-more">Browse more games</a>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getJoinGameInfo, submitJoinGame } from '../api/joinGame'

const route = useRoute()

const loading = ref(true)
const loadError = ref(null)
const step = ref('info')
const gameInfo = ref({})
const submitting = ref(false)
const submitError = ref(null)

const form = ref({
  email: '',
  name: '',
  postal_address_line1: '',
  postal_address_line2: '',
  state_province: '',
  postal_code: '',
  country: '',
})

const selectedDelivery = ref('')

const availableDeliveryMethods = computed(() => {
  const methods = []
  if (gameInfo.value.delivery_email) methods.push('email')
  if (gameInfo.value.delivery_physical_local) methods.push('local')
  if (gameInfo.value.delivery_physical_post) methods.push('post')
  return methods
})

const deliveryLabel = computed(() => {
  const labels = { email: 'email', local: 'local pickup', post: 'post' }
  return labels[selectedDelivery.value] || 'the selected method'
})

function formatGameType(gameType) {
  const types = { adventure: 'Adventure' }
  return types[gameType] ?? gameType
}

async function loadGameInfo() {
  loading.value = true
  loadError.value = null
  try {
    const res = await getJoinGameInfo(route.params.game_subscription_id)
    gameInfo.value = res.data ?? {}
    // Pre-select the only delivery method if there is just one
    if (availableDeliveryMethods.value.length === 1) {
      selectedDelivery.value = availableDeliveryMethods.value[0]
    }
  } catch (err) {
    loadError.value = err.message || 'Failed to load game information. Please try again.'
  } finally {
    loading.value = false
  }
}

async function onSubmit() {
  submitError.value = null
  submitting.value = true
  try {
    await submitJoinGame(route.params.game_subscription_id, {
      email: form.value.email,
      name: form.value.name,
      postal_address_line1: form.value.postal_address_line1,
      postal_address_line2: form.value.postal_address_line2 || undefined,
      state_province: form.value.state_province,
      postal_code: form.value.postal_code,
      country: form.value.country,
      delivery_email: selectedDelivery.value === 'email',
      delivery_physical_local: selectedDelivery.value === 'local',
      delivery_physical_post: selectedDelivery.value === 'post',
    })
    step.value = 'success'
  } catch (err) {
    submitError.value = err.message || 'Failed to submit. Please try again.'
  } finally {
    submitting.value = false
  }
}

onMounted(loadGameInfo)
</script>

<style scoped>
.join-game-view {
  max-width: 600px;
  margin: var(--space-lg) auto;
  padding: var(--space-md);
}

.join-loading {
  padding: var(--space-xl);
  text-align: center;
  color: var(--color-text-muted, #666);
}

.join-step {
  padding: var(--space-xl);
}

.join-error {
  padding: var(--space-xl);
  text-align: center;
}

.game-title {
  font-size: var(--font-size-xl);
  margin-bottom: var(--space-md);
}

.game-description {
  color: var(--color-text-muted, #444);
  margin-bottom: var(--space-md);
}

.game-meta {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  margin-bottom: var(--space-lg);
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

.instance-info {
  margin-bottom: var(--space-xl);
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  font-size: var(--font-size-sm);
}

.capacity-info,
.delivery-methods {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.label {
  font-weight: var(--font-weight-bold);
  min-width: 80px;
}

.delivery-badge {
  padding: 2px var(--space-sm);
  background: var(--color-bg, #f5f5f5);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs, 0.75rem);
}

.step-description {
  color: var(--color-text-muted, #444);
  margin-bottom: var(--space-xl);
}

.join-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-md);
}

label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-bold);
}

.required {
  color: var(--color-warning-dark, #c00);
}

input[type='text'],
input[type='email'] {
  padding: var(--space-sm) var(--space-md);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-base);
  background: var(--color-bg);
  color: var(--color-text);
}

input:focus {
  outline: none;
  border-color: var(--color-primary, #3b82f6);
}

.delivery-options {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.delivery-option {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  font-weight: normal;
  cursor: pointer;
}

.form-actions {
  display: flex;
  gap: var(--space-md);
  justify-content: flex-end;
  margin-top: var(--space-md);
}

.primary-button {
  padding: var(--space-sm) var(--space-xl);
  background: transparent;
  color: var(--color-button);
  border: 2px solid var(--color-button);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-bold);
  cursor: pointer;
  transition: all 0.2s;
}

.primary-button:hover {
  background: var(--color-button);
  color: var(--color-text-light);
}

.primary-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.secondary-button {
  padding: var(--space-sm) var(--space-lg);
  background: transparent;
  color: var(--color-text-muted, #666);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-base);
  cursor: pointer;
  transition: all 0.2s;
}

.secondary-button:hover {
  border-color: var(--color-text-muted, #666);
}

.error-message {
  color: var(--color-warning-dark, #c00);
  background: var(--color-warning-light, #fff3f3);
  padding: var(--space-sm) var(--space-md);
  border: 1px solid var(--color-warning, #f99);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
}

.join-success h1 {
  font-size: var(--font-size-xl);
  margin-bottom: var(--space-md);
}

.success-message {
  margin-bottom: var(--space-xl);
  color: var(--color-text-muted, #444);
}

.catalog-link {
  display: inline-block;
  padding: var(--space-sm) var(--space-xl);
  background: transparent;
  color: var(--color-button);
  border: 2px solid var(--color-button);
  border-radius: var(--radius-sm);
  text-decoration: none;
  font-weight: var(--font-weight-bold);
  font-size: var(--font-size-sm);
  transition: all 0.2s;
}

.catalog-link:hover {
  background: var(--color-button);
  color: var(--color-text-light);
}
</style>
