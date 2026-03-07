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

    <!-- Turn sheet list -->
    <div v-else class="ts-list" data-testid="ts-list">
      <h1 class="ts-title" data-testid="ts-title">Your Turn Sheets</h1>

      <p v-if="turnSheets.length === 0" class="ts-empty" data-testid="ts-empty">
        No turn sheets are available for this turn yet.
      </p>

      <div v-else class="ts-cards">
        <div
          v-for="sheet in turnSheets"
          :key="sheet.id"
          class="ts-card card"
          :data-testid="`ts-card-${sheet.id}`"
        >
          <div class="ts-card-header">
            <span class="ts-type" data-testid="ts-sheet-type">{{ formatSheetType(sheet.sheet_type) }}</span>
            <span
              class="ts-status"
              :class="sheet.is_completed ? 'status-complete' : 'status-pending'"
              :data-testid="`ts-status-${sheet.id}`"
            >
              {{ sheet.is_completed ? 'Completed' : 'Pending' }}
            </span>
          </div>
          <div class="ts-card-meta">
            <span>Turn {{ sheet.turn_number }}</span>
          </div>
          <div class="ts-card-actions">
            <button
              class="secondary-button"
              :data-testid="`btn-download-${sheet.id}`"
              @click="onDownload(sheet.id)"
            >
              Download PDF
            </button>
            <label
              v-if="!sheet.is_completed"
              class="secondary-button upload-label"
              :data-testid="`btn-upload-${sheet.id}`"
            >
              Upload Scan
              <input
                type="file"
                accept="image/png,image/jpeg,image/webp"
                class="upload-input"
                :data-testid="`input-upload-${sheet.id}`"
                @change="onUploadScan(sheet.id, $event)"
              />
            </label>
          </div>
        </div>
      </div>

      <div v-if="turnSheets.length > 0" class="ts-actions">
        <p v-if="submitError" class="error-message" data-testid="submit-error">{{ submitError }}</p>
        <button
          class="primary-button"
          :disabled="submitting || allCompleted"
          @click="onSubmit"
          data-testid="btn-submit"
        >
          {{ submitting ? 'Submitting...' : allCompleted ? 'Already submitted' : 'Submit turn sheets' }}
        </button>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import {
  verifyGameSubscriptionToken,
  requestNewTurnSheetToken,
  getGSITurnSheets,
  submitGSITurnSheets,
  downloadGSITurnSheetPDF,
  uploadGSITurnSheetScan,
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

const allCompleted = computed(() =>
  turnSheets.value.length > 0 && turnSheets.value.every((s) => s.is_completed)
)

function formatSheetType(sheetType) {
  const labels = {
    join_game: 'Join Game',
    location_choice: 'Location Choice',
    inventory_management: 'Inventory Management',
  }
  return labels[sheetType] ?? sheetType
}

async function authenticateWithToken() {
  const token = route.params.turn_sheet_token
  if (!token) return true

  if (authStore.sessionToken) return true

  const gsiId = route.params.game_subscription_instance_id
  try {
    const sessionToken = await verifyGameSubscriptionToken(gsiId, token)
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
    const res = await getGSITurnSheets(route.params.game_subscription_instance_id)
    turnSheets.value = res.turn_sheets ?? []
  } catch (err) {
    loadError.value = err.message || 'Failed to load turn sheets. Please try again.'
  } finally {
    loading.value = false
  }
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

async function onUploadScan(turnSheetId, event) {
  const file = event.target.files[0]
  if (!file) return

  submitError.value = null
  try {
    await uploadGSITurnSheetScan(route.params.game_subscription_instance_id, turnSheetId, file)
    await loadTurnSheets()
  } catch (err) {
    submitError.value = err.message || 'Failed to upload scanned image. Please try again.'
  }
}

async function onDownload(turnSheetId) {
  try {
    const res = await downloadGSITurnSheetPDF(route.params.game_subscription_instance_id, turnSheetId)
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `turn-sheet-${turnSheetId}.pdf`
    a.click()
    URL.revokeObjectURL(url)
  } catch (err) {
    submitError.value = err.message || 'Failed to download PDF.'
  }
}

async function onSubmit() {
  submitError.value = null
  submitting.value = true
  try {
    await submitGSITurnSheets(route.params.game_subscription_instance_id)
    submitted.value = true
  } catch (err) {
    submitError.value = err.message || 'Failed to submit. Please try again.'
  } finally {
    submitting.value = false
  }
}

onMounted(loadTurnSheets)
</script>

<style scoped>
.turn-sheet-view {
  max-width: 700px;
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

.ts-cards {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.ts-card {
  padding: 1rem 1.25rem;
}

.ts-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.4rem;
}

.ts-type {
  font-weight: 600;
  font-size: 1rem;
  color: var(--color-text, #11181c);
}

.ts-status {
  font-size: 0.8rem;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: 9999px;
}

.status-complete {
  background: #d1fae5;
  color: #065f46;
}

.status-pending {
  background: #fef3c7;
  color: #92400e;
}

.ts-card-meta {
  font-size: 0.875rem;
  color: var(--color-text-muted, #6b7280);
}

.ts-card-actions {
  margin-top: 0.75rem;
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

.secondary-button:hover {
  background: #e8f3fd;
}

.upload-label {
  cursor: pointer;
  display: inline-flex;
  align-items: center;
}

.upload-input {
  display: none;
}

.ts-actions {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.75rem;
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
  margin-top: 1rem;
  font-size: 0.9rem;
}
</style>
