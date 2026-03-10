<template>
  <div class="turn-sheet-view">

    <!-- Loading (includes silent token verification) -->
    <div v-if="loading" class="ts-loading" data-testid="ts-loading">
      <p>Loading your turn sheets...</p>
    </div>

    <!-- Expired / invalid token — request new link -->
    <div v-else-if="tokenExpired" class="ts-expired card" data-testid="ts-token-expired">
      <h1 class="hand-drawn-title">Link Expired</h1>
      <p>This turn sheet link is no longer valid. Enter your email to receive a fresh link.</p>
      <form @submit.prevent="onRequestNewLink" class="request-link-form">
        <div class="form-group">
          <label for="email">Email address</label>
          <input v-model="requestEmail" id="email" type="email" required autofocus data-testid="expired-email-input" />
        </div>
        <div class="form-actions">
          <button type="submit" :disabled="requestingLink" data-testid="expired-request-btn">
            {{ requestingLink ? 'Sending...' : 'Send New Link' }}
          </button>
        </div>
      </form>
      <p v-if="requestLinkMessage" class="success-inline" data-testid="expired-success">{{ requestLinkMessage }}</p>
      <p v-if="requestLinkError" class="error-message" data-testid="expired-error">{{ requestLinkError }}</p>
    </div>

    <!-- Error loading -->
    <div v-else-if="loadError" class="ts-error card" data-testid="ts-load-error">
      <p class="error-message">{{ loadError }}</p>
      <a href="/games" class="catalog-link" data-testid="link-browse-games">Browse games</a>
    </div>

    <!-- Submitted success -->
    <div v-else-if="submitted" class="ts-success card" data-testid="ts-success">
      <h1>Turn sheets submitted!</h1>
      <p class="success-message">
        Your turn sheets have been submitted successfully. The game manager will process them and
        you'll hear from us for the next turn.
      </p>
      <a href="/games" class="catalog-link" data-testid="link-browse-more">Browse more games</a>
    </div>

    <!-- Turn sheet list is empty -->
    <div v-else-if="currentTurnSheets.length === 0" class="ts-list" data-testid="ts-list">
      <h1 class="ts-title" data-testid="ts-title">Your Turn Sheets</h1>
      <p class="ts-empty" data-testid="ts-empty">
        No turn sheets are available for this turn yet.
      </p>
    </div>

    <!-- Stepper + inline viewer -->
    <div v-else class="ts-viewer" data-testid="ts-viewer">
      <h1 class="ts-title" data-testid="ts-title">Your Turn Sheets</h1>

      <!-- Stepper navigation bar -->
      <div class="ts-stepper" data-testid="ts-stepper">
        <button
          v-for="(sheet, idx) in currentTurnSheets"
          :key="sheet.id"
          class="ts-step"
          :class="{
            'ts-step--active': idx === activeIndex,
            'ts-step--ready': readySet.has(sheet.id),
            'ts-step--saved': hasSavedData(sheet) && !readySet.has(sheet.id),
          }"
          :data-testid="`ts-step-${idx}`"
          @click="activeIndex = idx"
        >
          <span class="ts-step-number">{{ idx + 1 }}</span>
          <span class="ts-step-label">{{ formatSheetType(sheet.sheet_type) }}</span>
          <span
            class="ts-step-status"
            :data-testid="`ts-step-status-${idx}`"
          >
            <span v-if="readySet.has(sheet.id)" class="status-dot status-ready" title="Ready"></span>
            <span v-else-if="hasSavedData(sheet)" class="status-dot status-saved" title="Saved"></span>
            <span v-else class="status-dot status-empty" title="Not started"></span>
          </span>
        </button>
      </div>

      <!-- Inline turn sheet viewer -->
      <div class="ts-iframe-container" data-testid="ts-iframe-container">
        <div v-if="loadingHTML" class="ts-iframe-loading">Loading turn sheet...</div>
        <iframe
          v-else-if="activeHTML"
          :srcdoc="activeHTML"
          class="ts-iframe"
          data-testid="ts-viewer-iframe"
          sandbox="allow-forms allow-scripts allow-same-origin"
          ref="iframeRef"
        ></iframe>
        <div v-else class="ts-iframe-empty">Select a turn sheet to view.</div>
      </div>

      <!-- Per-sheet actions -->
      <div class="ts-sheet-actions">
        <div class="ts-sheet-actions-left">
          <button
            class="secondary-button"
            :disabled="saving"
            data-testid="btn-save-sheet"
            @click="onSaveSheet"
          >
            {{ saving ? 'Saving...' : 'Save' }}
          </button>
          <button
            class="secondary-button"
            :class="{ 'btn-ready': isActiveReady }"
            data-testid="btn-mark-ready"
            @click="onToggleReady"
          >
            {{ isActiveReady ? 'Unmark Ready' : 'Mark Ready' }}
          </button>
        </div>
        <div class="ts-sheet-actions-right">
          <button
            class="secondary-button"
            :disabled="activeIndex === 0"
            data-testid="btn-prev-sheet"
            @click="activeIndex--"
          >
            &larr; Prev
          </button>
          <button
            class="secondary-button"
            :disabled="activeIndex >= currentTurnSheets.length - 1"
            data-testid="btn-next-sheet"
            @click="activeIndex++"
          >
            Next &rarr;
          </button>
        </div>
      </div>

      <p v-if="saveMessage" class="success-inline" data-testid="save-message">{{ saveMessage }}</p>
      <p v-if="saveError" class="error-message" data-testid="save-error">{{ saveError }}</p>

      <!-- Final submission -->
      <div class="ts-submit-section">
        <p v-if="submitError" class="error-message" data-testid="submit-error">{{ submitError }}</p>
        <button
          class="primary-button"
          :disabled="!allReady || submitting"
          data-testid="btn-submit-all"
          @click="onSubmit"
        >
          {{ submitting ? 'Submitting...' : 'Submit All Turn Sheets' }}
        </button>
        <p v-if="!allReady" class="ts-submit-hint">Mark all turn sheets as ready to enable submission.</p>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import {
  verifyGameSubscriptionToken,
  requestNewTurnSheetToken,
  getGameSubscriptionInstanceTurnSheets,
  getGameSubscriptionInstanceTurnSheetHTML,
  saveGameSubscriptionInstanceTurnSheet,
  submitGameSubscriptionInstanceTurnSheets,
} from '../api/player'

const route = useRoute()
const authStore = useAuthStore()

const loading = ref(true)
const loadError = ref(null)
const tokenExpired = ref(false)
const turnSheets = ref([])
const submitting = ref(false)
const submitError = ref(null)
const submitted = ref(false)

const requestEmail = ref('')
const requestingLink = ref(false)
const requestLinkMessage = ref(null)
const requestLinkError = ref(null)

// Stepper state
const activeIndex = ref(0)
const readySet = ref(new Set())
const loadingHTML = ref(false)
const activeHTML = ref('')
const saving = ref(false)
const saveMessage = ref(null)
const saveError = ref(null)
const iframeRef = ref(null)

// Derive current turn sheets (latest turn only)
const currentTurnSheets = computed(() => {
  if (turnSheets.value.length === 0) return []
  const maxTurn = Math.max(...turnSheets.value.map(s => s.turn_number))
  return turnSheets.value.filter(s => s.turn_number === maxTurn && !s.is_completed)
})

const activeSheet = computed(() => currentTurnSheets.value[activeIndex.value] || null)
const isActiveReady = computed(() => activeSheet.value && readySet.value.has(activeSheet.value.id))
const allReady = computed(() =>
  currentTurnSheets.value.length > 0 &&
  currentTurnSheets.value.every(s => readySet.value.has(s.id))
)

function formatSheetType(sheetType) {
  const labels = {
    adventure_game_join_game: 'Join Game',
    adventure_game_location_choice: 'Location Choice',
    adventure_game_inventory_management: 'Inventory Management',
  }
  return labels[sheetType] ?? sheetType.replace(/_/g, ' ')
}

function hasSavedData(sheet) {
  return sheet.scanned_data && Object.keys(sheet.scanned_data).length > 0
}

async function authenticateWithToken() {
  const token = route.params.turn_sheet_token
  if (!token) return true
  if (authStore.sessionToken) return true

  const gameSubscriptionInstanceId = route.params.game_subscription_instance_id
  try {
    const sessionToken = await verifyGameSubscriptionToken(gameSubscriptionInstanceId, token)
    authStore.setSessionToken(sessionToken)
    return true
  } catch {
    return false
  }
}

async function loadTurnSheets() {
  loading.value = true
  loadError.value = null
  tokenExpired.value = false
  try {
    const authenticated = await authenticateWithToken()
    if (!authenticated) {
      tokenExpired.value = true
      return
    }
    const res = await getGameSubscriptionInstanceTurnSheets(route.params.game_subscription_instance_id)
    turnSheets.value = res.turn_sheets ?? []
  } catch (err) {
    loadError.value = err.message || 'Failed to load turn sheets. Please try again.'
  } finally {
    loading.value = false
  }
}

async function loadSheetHTML(sheet) {
  if (!sheet) {
    activeHTML.value = ''
    return
  }
  loadingHTML.value = true
  try {
    const gameSubscriptionInstanceId = route.params.game_subscription_instance_id
    activeHTML.value = await getGameSubscriptionInstanceTurnSheetHTML(gameSubscriptionInstanceId, sheet.id)
  } catch (err) {
    activeHTML.value = `<p style="color:red;">Failed to load turn sheet: ${err.message}</p>`
  } finally {
    loadingHTML.value = false
  }
}

function extractFormData() {
  const iframe = iframeRef.value
  if (!iframe || !iframe.contentDocument) return null

  const formData = {}
  const doc = iframe.contentDocument

  // Extract all input, select, textarea values
  doc.querySelectorAll('input, select, textarea').forEach(el => {
    const name = el.name || el.id
    if (!name) return

    if (el.type === 'checkbox') {
      if (!formData[name]) formData[name] = []
      if (el.checked) formData[name].push(el.value)
    } else if (el.type === 'radio') {
      if (el.checked) formData[name] = el.value
    } else if (el.tagName === 'SELECT' && el.multiple) {
      formData[name] = Array.from(el.selectedOptions).map(o => o.value)
    } else {
      formData[name] = el.value
    }
  })

  return formData
}

async function onSaveSheet() {
  if (!activeSheet.value) return
  saving.value = true
  saveMessage.value = null
  saveError.value = null

  try {
    const formData = extractFormData()
    const gameSubscriptionInstanceId = route.params.game_subscription_instance_id
    await saveGameSubscriptionInstanceTurnSheet(gameSubscriptionInstanceId, activeSheet.value.id, formData || {})

    // Update local sheet's scanned_data
    const idx = turnSheets.value.findIndex(s => s.id === activeSheet.value.id)
    if (idx !== -1) {
      turnSheets.value[idx] = { ...turnSheets.value[idx], scanned_data: formData }
    }

    saveMessage.value = 'Turn sheet saved.'
    setTimeout(() => { saveMessage.value = null }, 3000)
  } catch (err) {
    saveError.value = err.message || 'Failed to save turn sheet.'
  } finally {
    saving.value = false
  }
}

function onToggleReady() {
  if (!activeSheet.value) return
  const id = activeSheet.value.id
  const newSet = new Set(readySet.value)
  if (newSet.has(id)) {
    newSet.delete(id)
  } else {
    newSet.add(id)
  }
  readySet.value = newSet
}

async function onRequestNewLink() {
  requestLinkError.value = null
  requestLinkMessage.value = null
  requestingLink.value = true
  try {
    await requestNewTurnSheetToken(route.params.game_subscription_instance_id, requestEmail.value)
    requestLinkMessage.value = 'A new link has been sent to your email. Please check your inbox.'
  } catch (err) {
    requestLinkError.value = err.message || 'Failed to send a new link. Please try again.'
  } finally {
    requestingLink.value = false
  }
}

async function onSubmit() {
  submitError.value = null
  submitting.value = true
  try {
    await submitGameSubscriptionInstanceTurnSheets(route.params.game_subscription_instance_id)
    submitted.value = true
  } catch (err) {
    submitError.value = err.message || 'Failed to submit. Please try again.'
  } finally {
    submitting.value = false
  }
}

// Load HTML when active sheet changes
watch(activeSheet, (sheet) => {
  saveMessage.value = null
  saveError.value = null
  loadSheetHTML(sheet)
})

// Load first sheet's HTML once current turn sheets are ready
watch(currentTurnSheets, (sheets) => {
  if (sheets.length > 0 && activeIndex.value === 0) {
    loadSheetHTML(sheets[0])
  }
}, { immediate: false })

onMounted(loadTurnSheets)
</script>

<style scoped>
.turn-sheet-view {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.card {
  background: var(--color-surface, #fff);
  border: 1px solid var(--color-border, #e2e8f0);
  border-radius: 12px;
  padding: 2rem;
  margin-bottom: 1.5rem;
}

.ts-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 1.5rem;
  color: var(--color-text, #11181c);
}

.ts-empty {
  color: var(--color-text-muted, #6b7280);
  text-align: center;
  padding: 2rem 0;
}

/* ---- Stepper ---- */
.ts-stepper {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.ts-step {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.5rem 1rem;
  border: 1.5px solid var(--color-border, #e2e8f0);
  border-radius: 9999px;
  background: var(--color-surface, #fff);
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-muted, #6b7280);
  transition: all 0.2s;
}

.ts-step:hover {
  border-color: #006ecd;
  color: #006ecd;
}

.ts-step--active {
  border-color: #006ecd;
  background: #e8f3fd;
  color: #006ecd;
}

.ts-step--ready {
  border-color: #059669;
  color: #059669;
}

.ts-step--ready.ts-step--active {
  background: #d1fae5;
}

.ts-step--saved:not(.ts-step--ready) {
  border-color: #d97706;
  color: #d97706;
}

.ts-step--saved.ts-step--active:not(.ts-step--ready) {
  background: #fef3c7;
}

.ts-step-number {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.3rem;
  height: 1.3rem;
  border-radius: 50%;
  background: currentColor;
  color: #fff;
  font-size: 0.7rem;
  font-weight: 700;
}

.ts-step-label {
  white-space: nowrap;
}

.ts-step-status {
  display: flex;
  align-items: center;
}

.status-dot {
  display: inline-block;
  width: 0.55rem;
  height: 0.55rem;
  border-radius: 50%;
}

.status-ready {
  background: #059669;
}

.status-saved {
  background: #d97706;
}

.status-empty {
  background: #d1d5db;
}

/* ---- Iframe viewer ---- */
.ts-iframe-container {
  border: 1.5px solid var(--color-border, #e2e8f0);
  border-radius: 12px;
  overflow: hidden;
  min-height: 400px;
  background: #fafafa;
  margin-bottom: 1rem;
}

.ts-iframe {
  width: 100%;
  height: 500px;
  border: none;
}

.ts-iframe-loading,
.ts-iframe-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  color: var(--color-text-muted, #6b7280);
}

/* ---- Sheet actions ---- */
.ts-sheet-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
  flex-wrap: wrap;
}

.ts-sheet-actions-left,
.ts-sheet-actions-right {
  display: flex;
  gap: 0.5rem;
}

.secondary-button {
  background: transparent;
  color: #006ecd;
  border: 1.5px solid #006ecd;
  padding: 0.4rem 1rem;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.secondary-button:hover:not(:disabled) {
  background: #e8f3fd;
}

.secondary-button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn-ready {
  border-color: #059669;
  color: #059669;
}

.btn-ready:hover:not(:disabled) {
  background: #d1fae5;
}

/* ---- Submit ---- */
.ts-submit-section {
  margin-top: 1.5rem;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.5rem;
}

.ts-submit-hint {
  font-size: 0.8rem;
  color: var(--color-text-muted, #6b7280);
}

.primary-button {
  background: #006ecd;
  color: #fff;
  border: none;
  padding: 0.75rem 2rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.primary-button:hover:not(:disabled) {
  background: #0055a0;
}

.primary-button:disabled {
  background: #93c5fd;
  cursor: not-allowed;
}

.error-message {
  color: #b91c1c;
  background: #fee2e2;
  border: 1px solid #fca5a5;
  border-radius: 8px;
  padding: 0.75rem 1rem;
  font-size: 0.9rem;
}

.success-message {
  color: var(--color-text, #11181c);
  font-size: 1rem;
  line-height: 1.6;
  margin-bottom: 1.5rem;
}

.catalog-link {
  display: inline-block;
  color: #006ecd;
  text-decoration: none;
  font-weight: 600;
}

.catalog-link:hover {
  text-decoration: underline;
}

.ts-loading {
  text-align: center;
  padding: 3rem 0;
  color: var(--color-text-muted, #6b7280);
}

.ts-expired {
  text-align: center;
}

.ts-expired .hand-drawn-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 1rem;
}

.request-link-form {
  max-width: 360px;
  margin: 1.5rem auto 0;
  text-align: left;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-weight: 600;
  margin-bottom: 0.3rem;
}

.form-group input {
  width: 100%;
  padding: 0.6rem 0.75rem;
  border: 1.5px solid var(--color-border, #e2e8f0);
  border-radius: 8px;
  font-size: 1rem;
}

.form-actions {
  margin-top: 1rem;
}

.form-actions button {
  width: 100%;
  background: #006ecd;
  color: #fff;
  border: none;
  padding: 0.75rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
}

.form-actions button:disabled {
  background: #93c5fd;
  cursor: not-allowed;
}

.success-inline {
  color: #065f46;
  background: #d1fae5;
  border-radius: 8px;
  padding: 0.75rem 1rem;
  margin-top: 0.5rem;
  font-size: 0.9rem;
}
</style>
