<template>
  <div class="turn-sheet-view">

    <!-- Loading -->
    <div v-if="loading" class="ts-loading" data-testid="ts-loading">
      <p>Loading your turn sheets...</p>
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
import { getGSITurnSheets, submitGSITurnSheets } from '../api/player'

const route = useRoute()

const loading = ref(true)
const loadError = ref(null)
const turnSheets = ref([])
const submitting = ref(false)
const submitError = ref(null)
const submitted = ref(false)

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

async function loadTurnSheets() {
  loading.value = true
  loadError.value = null
  try {
    const res = await getGSITurnSheets(route.params.game_subscription_instance_id)
    turnSheets.value = res.turn_sheets ?? []
  } catch (err) {
    loadError.value = err.message || 'Failed to load turn sheets. Please try again.'
  } finally {
    loading.value = false
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
</style>
