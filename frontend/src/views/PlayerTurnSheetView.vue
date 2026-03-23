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
    <ConfirmationCard
      v-else-if="submitted"
      title="Turn sheets submitted!"
      message="Your turn sheets have been submitted successfully. The game manager will process them and you'll hear from us for the next turn."
      data-testid="ts-success"
    />

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
        <button v-for="(sheet, idx) in currentTurnSheets" :key="sheet.id" class="ts-step" :class="{
          'ts-step--active': idx === activeIndex,
          'ts-step--has-data': hasCachedData(sheet),
        }" :data-testid="`ts-step-${idx}`" @click="activeIndex = idx">
          <span class="ts-step-number">{{ idx + 1 }}</span>
          <span class="ts-step-label">{{ formatSheetType(sheet.sheet_type) }}</span>
          <span class="ts-step-status" :data-testid="`ts-step-status-${idx}`">
            <span v-if="hasCachedData(sheet)" class="status-dot status-has-data" title="Filled"></span>
            <span v-else class="status-dot status-empty" title="Not started"></span>
          </span>
        </button>
      </div>

      <!-- Inline turn sheet viewer -->
      <div class="ts-iframe-container" data-testid="ts-iframe-container">
        <div v-if="loadingHTML" class="ts-iframe-loading">Loading turn sheet...</div>
        <iframe v-else-if="activeHTML" :srcdoc="activeHTML" class="ts-iframe" data-testid="ts-viewer-iframe"
          sandbox="allow-forms allow-scripts allow-same-origin" ref="iframeRef" @load="onIframeLoad"></iframe>
        <div v-else class="ts-iframe-empty">Select a turn sheet to view.</div>
      </div>

      <!-- Navigation actions -->
      <div class="ts-sheet-actions">
        <button class="secondary-button" :disabled="activeIndex === 0" data-testid="btn-prev-sheet"
          @click="activeIndex--">
          &larr; Prev
        </button>
        <button class="secondary-button" :disabled="activeIndex >= currentTurnSheets.length - 1"
          data-testid="btn-next-sheet" @click="activeIndex++">
          Next &rarr;
        </button>
      </div>

      <!-- Final submission -->
      <div class="ts-submit-section">
        <p v-if="submitError" class="error-message" data-testid="submit-error">{{ submitError }}</p>
        <button class="primary-button" :disabled="submitting" data-testid="btn-submit-all" @click="onSubmit">
          {{ submitting ? 'Submitting...' : 'Submit All Turn Sheets' }}
        </button>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import ConfirmationCard from '../components/ConfirmationCard.vue'
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
const loadingHTML = ref(false)
const activeHTML = ref('')
const iframeRef = ref(null)

// In-memory form data cache: Map<sheetId, formData>
const formDataCache = ref(new Map())

// Canonical presentation order for turn sheet types.
const SHEET_PRESENTATION_ORDER = [
  'adventure_game_join_game',
  'adventure_game_location_choice',
  'adventure_game_inventory_management',
  'adventure_game_monster',
  'mecha_join_game',
  'mecha_orders',
]

// Derive current turn sheets (latest turn only), sorted by canonical presentation order.
const currentTurnSheets = computed(() => {
  if (turnSheets.value.length === 0) return []
  const maxTurn = Math.max(...turnSheets.value.map(s => s.turn_number))
  const sheets = turnSheets.value.filter(s => s.turn_number === maxTurn && !s.is_completed)
  return [...sheets].sort((a, b) => {
    const ai = SHEET_PRESENTATION_ORDER.indexOf(a.sheet_type)
    const bi = SHEET_PRESENTATION_ORDER.indexOf(b.sheet_type)
    const aOrder = ai === -1 ? 999 : ai
    const bOrder = bi === -1 ? 999 : bi
    return aOrder - bOrder
  })
})

const activeSheet = computed(() => currentTurnSheets.value[activeIndex.value] || null)

function formatSheetType(sheetType) {
  const labels = {
    adventure_game_join_game: 'Join Game',
    adventure_game_location_choice: 'Location Choice',
    adventure_game_inventory_management: 'Inventory Management',
    adventure_game_monster: 'Creature Encounter',
    mecha_join_game: 'Join Game',
    mecha_orders: 'Mech Orders',
  }
  return labels[sheetType] ?? sheetType.replace(/_/g, ' ')
}

function hasCachedData(sheet) {
  return formDataCache.value.has(sheet.id)
}

async function authenticateWithToken() {
  const token = route.params.turn_sheet_token
  if (!token) return true
  // Always exchange the URL token — an existing session may belong to a different account
  // (e.g. a designer browsing their own game), so we must not skip this step.
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

  // Convert per-item radio groups (inv_*, loc_*) into action arrays expected by backend.
  // Radio buttons are named inv_<itemId> for inventory items and loc_<itemId> for location
  // items, with the action as the value (equip, drop, unequip, pick_up).
  const equip = []
  const drop = []
  const pick_up = []
  const unequip = []

  for (const key of Object.keys(formData)) {
    if (key.startsWith('inv_')) {
      const itemId = key.slice(4)
      const action = formData[key]
      if (action === 'equip') equip.push(itemId)
      else if (action === 'drop') drop.push(itemId)
      else if (action === 'unequip') unequip.push(itemId)
      delete formData[key]
    } else if (key.startsWith('loc_')) {
      const itemId = key.slice(4)
      const action = formData[key]
      if (action === 'equip') equip.push(itemId)
      else if (action === 'pick_up') pick_up.push(itemId)
      delete formData[key]
    }
  }

  // Merge checkbox-based drop/pick_up (non-equippable items) with radio-derived arrays.
  if (formData.drop) { drop.push(...formData.drop); delete formData.drop }
  if (formData.pick_up) { pick_up.push(...formData.pick_up); delete formData.pick_up }

  if (equip.length) formData.equip = equip
  if (drop.length) formData.drop = drop
  if (pick_up.length) formData.pick_up = pick_up
  if (unequip.length) formData.unequip = unequip

  // Convert action_N / action_N_target flat fields into the structured actions array
  // expected by the monster encounter backend processor.
  const actions = []
  let actionIndex = 0
  while (`action_${actionIndex}` in formData) {
    const actionType = formData[`action_${actionIndex}`]
    const target = formData[`action_${actionIndex}_target`] || ''
    const action = { action_type: actionType }
    if (actionType === 'attack' && target) {
      action.target_creature_instance_id = target
    }
    actions.push(action)
    delete formData[`action_${actionIndex}`]
    delete formData[`action_${actionIndex}_target`]
    actionIndex++
  }
  if (actions.length > 0) {
    // Trim trailing do_nothing entries (matches OCR output behaviour)
    while (actions.length > 0 && actions[actions.length - 1].action_type === 'do_nothing') {
      actions.pop()
    }
    if (actions.length > 0) {
      formData.actions = actions
    }
  }

  // Convert mecha orders: move_to_<mechId> / attack_target_<mechId> selects
  // and mech_instance_id_<mechId> hidden inputs → mech_orders array.
  const mechOrders = []
  const mechKeys = []
  for (const key of Object.keys(formData)) {
    if (key.startsWith('mech_instance_id_')) {
      mechKeys.push(key)
      const mechId = formData[key]
      if (!mechOrders.find(o => o.mech_instance_id === mechId)) {
        mechOrders.push({
          mech_instance_id: mechId,
          move_to_sector_instance_id: formData[`move_to_${mechId}`] || '',
          attack_target_mech_instance_id: formData[`attack_target_${mechId}`] || '',
        })
      }
    }
  }
  if (mechOrders.length > 0) {
    for (const key of mechKeys) {
      const mechId = formData[key]
      delete formData[key]
      delete formData[`move_to_${mechId}`]
      delete formData[`attack_target_${mechId}`]
    }
    formData.mech_orders = mechOrders
  }

  // Remove empty arrays — an unchecked checkbox group sends [] which can
  // fail oneOf schema validation (e.g. equip field in inventory management).
  for (const key of Object.keys(formData)) {
    if (Array.isArray(formData[key]) && formData[key].length === 0) {
      delete formData[key]
    }
  }

  return formData
}

// Restore saved form data into the iframe after it loads
function applyFormData(data) {
  const iframe = iframeRef.value
  if (!iframe || !iframe.contentDocument || !data) return

  const doc = iframe.contentDocument

  // Restore standard inputs (text, select, plain checkboxes).
  // Per-item radio groups (inv_*, loc_*) are handled separately below.
  doc.querySelectorAll('input, select, textarea').forEach(el => {
    const name = el.name || el.id
    if (!name || !(name in data)) return
    const value = data[name]

    if (el.type === 'checkbox') {
      el.checked = Array.isArray(value) ? value.includes(el.value) : el.value === value
    } else if (el.type === 'radio') {
      el.checked = el.value === value
    } else if (el.tagName === 'SELECT' && el.multiple) {
      const selected = Array.isArray(value) ? value : [value]
      Array.from(el.options).forEach(opt => {
        opt.selected = selected.includes(opt.value)
      })
    } else {
      el.value = value ?? ''
    }
  })

  // Restore per-item radio groups from action arrays.
  // Inventory item radios are named inv_<itemId>; location item radios are named loc_<itemId>.
  const radioMappings = [
    { dataKey: 'equip', namePrefix: 'inv_', radioValue: 'equip' },
    { dataKey: 'equip', namePrefix: 'loc_', radioValue: 'equip' },
    { dataKey: 'drop', namePrefix: 'inv_', radioValue: 'drop' },
    { dataKey: 'pick_up', namePrefix: 'loc_', radioValue: 'pick_up' },
    { dataKey: 'unequip', namePrefix: 'inv_', radioValue: 'unequip' },
  ]

  for (const { dataKey, namePrefix, radioValue } of radioMappings) {
    const ids = data[dataKey]
    if (!Array.isArray(ids)) continue
    for (const itemId of ids) {
      const radio = doc.querySelector(
        `input[type="radio"][name="${namePrefix}${itemId}"][value="${radioValue}"]`
      )
      if (radio) radio.checked = true
    }
  }

  // Restore mecha orders: mech_orders array → move_to_<mechId> / attack_target_<mechId> selects.
  if (Array.isArray(data.mech_orders)) {
    for (const order of data.mech_orders) {
      if (!order.mech_instance_id) continue
      const moveSelect = doc.querySelector(`select[name="move_to_${order.mech_instance_id}"]`)
      if (moveSelect && order.move_to_sector_instance_id) {
        moveSelect.value = order.move_to_sector_instance_id
      }
      const attackSelect = doc.querySelector(`select[name="attack_target_${order.mech_instance_id}"]`)
      if (attackSelect && order.attack_target_mech_instance_id) {
        attackSelect.value = order.attack_target_mech_instance_id
      }
    }
  }

  // Restore monster encounter action slots from the structured actions array.
  if (Array.isArray(data.actions)) {
    data.actions.forEach((action, i) => {
      const actionRadio = doc.querySelector(
        `input[type="radio"][name="action_${i}"][value="${action.action_type}"]`
      )
      if (actionRadio) actionRadio.checked = true

      if (action.action_type === 'attack' && action.target_creature_instance_id) {
        const targetSelect = doc.querySelector(`select[name="action_${i}_target"]`)
        if (targetSelect) targetSelect.value = action.target_creature_instance_id
      }
    })
  }
}

// Called when the iframe finishes loading — size to content, then restore cached form data
function onIframeLoad() {
  const iframe = iframeRef.value
  if (iframe && iframe.contentDocument) {
    const docEl = iframe.contentDocument.documentElement
    const body = iframe.contentDocument.body
    const height = (docEl && docEl.scrollHeight) || (body && body.scrollHeight)
    if (height) iframe.style.height = height + 'px'
  }
  if (!activeSheet.value) return
  const cached = formDataCache.value.get(activeSheet.value.id)
  if (cached) {
    applyFormData(cached)
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

async function onSubmit() {
  submitError.value = null
  submitting.value = true

  try {
    // Cache the currently active sheet's form data before saving
    if (activeSheet.value) {
      const formData = extractFormData()
      if (formData && Object.keys(formData).length > 0) {
        const newCache = new Map(formDataCache.value)
        newCache.set(activeSheet.value.id, formData)
        formDataCache.value = newCache
      }
    }

    // Save all sheets that have cached data to the backend
    const gameSubscriptionInstanceId = route.params.game_subscription_instance_id
    for (const sheet of currentTurnSheets.value) {
      const cached = formDataCache.value.get(sheet.id)
      if (cached) {
        await saveGameSubscriptionInstanceTurnSheet(gameSubscriptionInstanceId, sheet.id, cached)
      }
    }

    // Submit all turn sheets (backend marks all as completed, with or without scanned data)
    await submitGameSubscriptionInstanceTurnSheets(gameSubscriptionInstanceId)
    submitted.value = true
  } catch (err) {
    submitError.value = err.message || 'Failed to submit. Please try again.'
  } finally {
    submitting.value = false
  }
}

// When the active sheet changes: cache outgoing sheet's form data, then load new sheet
watch(activeSheet, (newSheet, oldSheet) => {
  if (oldSheet) {
    const formData = extractFormData()
    if (formData && Object.keys(formData).length > 0) {
      const newCache = new Map(formDataCache.value)
      newCache.set(oldSheet.id, formData)
      formDataCache.value = newCache
    }
  }
  loadSheetHTML(newSheet)
})

// When current turn sheets load, pre-populate cache from any backend-saved scanned_data
watch(currentTurnSheets, (sheets) => {
  if (sheets.length > 0) {
    const newCache = new Map(formDataCache.value)
    for (const sheet of sheets) {
      if (sheet.scanned_data && Object.keys(sheet.scanned_data).length > 0 && !newCache.has(sheet.id)) {
        newCache.set(sheet.id, sheet.scanned_data)
      }
    }
    formDataCache.value = newCache

    if (activeIndex.value === 0) {
      loadSheetHTML(sheets[0])
    }
  }
}, { immediate: false })

onMounted(loadTurnSheets)
</script>

<style scoped>
.turn-sheet-view {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
  box-sizing: border-box;
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

.ts-step--has-data {
  border-color: #059669;
  color: #059669;
}

.ts-step--has-data.ts-step--active {
  background: #d1fae5;
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

.status-has-data {
  background: #059669;
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
  height: auto;
  min-height: 400px;
  border: none;
  display: block;
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
  justify-content: center;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
}

/* ---- Submit ---- */
.ts-submit-section {
  margin-top: 1.5rem;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
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
