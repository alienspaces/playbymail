<template>
  <div class="join-game-view">

    <!-- Loading turn sheet -->
    <div v-if="loading" class="join-loading" data-testid="join-loading">
      <p>Loading join game turn sheet...</p>
    </div>

    <!-- Failed to load -->
    <div v-else-if="loadError" class="join-error card" data-testid="join-load-error">
      <p class="error-message">{{ loadError }}</p>
      <a href="/games" class="catalog-link">Browse other games</a>
    </div>

    <!-- Success confirmation -->
    <div v-else-if="step === 'success'" class="join-success card" data-testid="step-success">
      <h1>You're in!</h1>
      <p class="success-message">
        You have successfully joined the game.
        You will receive further instructions soon.
      </p>
      <a href="/games" class="catalog-link" data-testid="link-browse-more">Browse more games</a>
    </div>

    <!-- Turn sheet display -->
    <div v-else class="join-sheet-wrapper" data-testid="join-sheet">
      <iframe
        ref="sheetFrame"
        :srcdoc="turnSheetHtml"
        class="turn-sheet-frame"
        sandbox="allow-same-origin"
        data-testid="join-sheet-iframe"
        @load="onIframeLoad"
      ></iframe>

      <div class="join-actions">
        <p v-if="submitError" class="error-message" data-testid="submit-error">{{ submitError }}</p>
        <div class="action-buttons">
          <a href="/games" class="secondary-button" data-testid="btn-back">Back to Catalog</a>
          <button
            class="primary-button"
            :disabled="submitting"
            data-testid="btn-submit"
            @click="onSubmit"
          >
            {{ submitting ? 'Joining...' : 'Submit & Join Game' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getJoinSheet, submitJoinGame } from '../api/joinGame'

const route = useRoute()

console.log('[PlayerJoinGame] component setup, route params:', JSON.stringify(route.params))

const loading = ref(true)
const loadError = ref(null)
const step = ref('form')
const turnSheetHtml = ref('')
const submitting = ref(false)
const submitError = ref(null)
const sheetFrame = ref(null)

function extractFormData() {
  const doc = sheetFrame.value?.contentDocument
  if (!doc) return null

  const val = (name) => {
    const el = doc.querySelector(`[name="${name}"]`)
    if (!el) return ''
    if (el.type === 'radio') {
      const checked = doc.querySelector(`[name="${name}"]:checked`)
      return checked ? checked.value : ''
    }
    return el.value || ''
  }

  const deliveryMethod = val('delivery_method')

  return {
    email: val('email'),
    name: val('name'),
    postal_address_line1: val('postal_address_line1'),
    postal_address_line2: val('postal_address_line2') || undefined,
    state_province: val('state_province'),
    country: val('country'),
    postal_code: val('postal_code'),
    delivery_email: deliveryMethod === 'email',
    delivery_physical_local: deliveryMethod === 'local',
    delivery_physical_post: deliveryMethod === 'post',
  }
}

function onIframeLoad() {
  const iframe = sheetFrame.value
  if (!iframe) return
  try {
    const doc = iframe.contentDocument
    if (doc && doc.body) {
      iframe.style.height = doc.documentElement.scrollHeight + 'px'
    }
  } catch {
    // cross-origin guard
  }
}

async function loadSheet() {
  console.log('[PlayerJoinGame] loadSheet called, game_subscription_id:', route.params.game_subscription_id)
  loading.value = true
  loadError.value = null
  try {
    turnSheetHtml.value = await getJoinSheet(route.params.game_subscription_id)
    console.log('[PlayerJoinGame] loadSheet success, html length:', turnSheetHtml.value?.length)
  } catch (err) {
    console.error('[PlayerJoinGame] loadSheet error:', err)
    loadError.value = err.message || 'Failed to load the join game turn sheet. Please try again.'
  } finally {
    loading.value = false
    console.log('[PlayerJoinGame] loadSheet done, loading:', loading.value, 'loadError:', loadError.value, 'step:', step.value)
  }
}

async function onSubmit() {
  submitError.value = null

  const data = extractFormData()
  if (!data) {
    submitError.value = 'Unable to read form data. Please try again.'
    return
  }

  if (!data.email || !data.name || !data.postal_address_line1 || !data.state_province || !data.country || !data.postal_code) {
    submitError.value = 'Please fill in all required fields on the turn sheet.'
    return
  }

  submitting.value = true
  try {
    await submitJoinGame(route.params.game_subscription_id, data)
    step.value = 'success'
  } catch (err) {
    submitError.value = err.message || 'Failed to submit. Please try again.'
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  console.log('[PlayerJoinGame] onMounted, calling loadSheet')
  loadSheet()
})
</script>

<style scoped>
.join-game-view {
  width: 100%;
  max-width: 1200px;
  margin: var(--space-lg, 1.5rem) auto;
  padding: var(--space-md, 1rem);
  box-sizing: border-box;
}

.join-loading {
  padding: var(--space-xl, 2rem);
  text-align: center;
  color: var(--color-text-muted, #666);
}

.join-error {
  padding: var(--space-xl, 2rem);
  text-align: center;
}

.join-sheet-wrapper {
  display: flex;
  flex-direction: column;
  gap: var(--space-lg, 1.5rem);
}

.turn-sheet-frame {
  width: 100%;
  min-height: 800px;
  border: 1px solid var(--color-border, #e2e8f0);
  border-radius: var(--radius-md, 8px);
  background: #fff;
}

.join-actions {
  display: flex;
  flex-direction: column;
  gap: var(--space-md, 1rem);
}

.action-buttons {
  display: flex;
  gap: var(--space-md, 1rem);
  justify-content: flex-end;
  align-items: center;
}

.primary-button {
  padding: var(--space-sm, 0.5rem) var(--space-xl, 2rem);
  background: transparent;
  color: var(--color-button, #006ecd);
  border: 2px solid var(--color-button, #006ecd);
  border-radius: var(--radius-sm, 4px);
  font-size: var(--font-size-base, 1rem);
  font-weight: var(--font-weight-bold, 600);
  cursor: pointer;
  transition: all 0.2s;
}

.primary-button:hover {
  background: var(--color-button, #006ecd);
  color: var(--color-text-light, #fff);
}

.primary-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.secondary-button {
  padding: var(--space-sm, 0.5rem) var(--space-lg, 1.5rem);
  background: transparent;
  color: var(--color-text-muted, #666);
  border: 1px solid var(--color-border, #e2e8f0);
  border-radius: var(--radius-sm, 4px);
  font-size: var(--font-size-base, 1rem);
  text-decoration: none;
  cursor: pointer;
  transition: all 0.2s;
}

.secondary-button:hover {
  border-color: var(--color-text-muted, #666);
}

.error-message {
  color: var(--color-warning-dark, #c00);
  background: var(--color-warning-light, #fff3f3);
  padding: var(--space-sm, 0.5rem) var(--space-md, 1rem);
  border: 1px solid var(--color-warning, #f99);
  border-radius: var(--radius-sm, 4px);
  font-size: var(--font-size-sm, 0.875rem);
}

.join-success {
  padding: var(--space-xl, 2rem);
}

.join-success h1 {
  font-size: var(--font-size-xl, 1.5rem);
  margin-bottom: var(--space-md, 1rem);
}

.success-message {
  margin-bottom: var(--space-xl, 2rem);
  color: var(--color-text-muted, #444);
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
