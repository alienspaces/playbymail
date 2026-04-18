<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage sectors.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Sectors" actionText="Create New Sector" :showIcon="false" titleLevel="h2"
        @action="openCreate" />
      <ResourceTable :columns="columns" :rows="formattedSectors" :loading="store.loading" :error="store.error">
        <template #cell-name="{ row }">
          <a href="#" class="edit-link" @click.prevent="openEdit(row)">{{ row.name }}</a>
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>
      <TablePagination :pageNumber="store.pageNumber" :hasMore="store.hasMore"
        @page-change="(p) => store.fetchMechaGameSectors(selectedGame.id, p)" />
    </div>

    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal">
          <h2>{{ modalMode === 'create' ? 'Create Sector' : 'Edit Sector' }}</h2>
          <form @submit.prevent="handleSubmit(modalForm)" class="modal-form">
            <div class="form-group">
              <label>Name <span class="required">*</span></label>
              <input v-model="modalForm.name" required maxlength="100" autocomplete="off" />
            </div>
            <div class="form-group">
              <label>Description</label>
              <textarea v-model="modalForm.description" rows="3" maxlength="4096"></textarea>
            </div>
            <div class="form-row">
              <div class="form-group half">
                <label>Elevation</label>
                <select v-model.number="modalForm.elevation">
                  <option v-for="v in elevationOptions" :key="v" :value="v">{{ v }}</option>
                </select>
                <FieldHint>Relative height (-10 to 10). Higher elevation is preferred by defensive AI opponents and used as a tie-breaker when picking cover.</FieldHint>
              </div>
              <div class="form-group half">
                <label>Cover Modifier</label>
                <select v-model.number="modalForm.cover_modifier">
                  <option v-for="v in coverModifierOptions" :key="v" :value="v">{{ formatSigned(v) }}</option>
                </select>
                <FieldHint>Added directly to attacker hit chance (-50 to +50, step 5). Negative = harder to hit; 0 = no effect; positive = easier to hit.</FieldHint>
              </div>
            </div>
            <div class="form-group checkbox-group">
              <label class="checkbox-label">
                <input type="checkbox" v-model="modalForm.is_starting_sector" />
                This is the starting sector for new players
              </label>
            </div>
            <div class="modal-actions">
              <button type="submit">{{ modalMode === 'create' ? 'Create' : 'Save' }}</button>
              <button type="button" @click="closeModal">Cancel</button>
            </div>
          </form>
          <div v-if="modalError" class="error"><p>{{ modalError }}</p></div>
        </div>
      </div>
    </Teleport>

    <ConfirmationModal :visible="showDeleteModal" title="Delete Sector"
      :message="`Are you sure you want to delete '${toDelete?.name}'?`"
      @confirm="handleDelete" @cancel="showDeleteModal = false" />
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useMechaGameSectorsStore } from '../../../stores/mechaGameSectors'
import { useGamesStore } from '../../../stores/games'
import ResourceTable from '../../../components/ResourceTable.vue'
import ConfirmationModal from '../../../components/ConfirmationModal.vue'
import PageHeader from '../../../components/PageHeader.vue'
import GameContext from '../../../components/GameContext.vue'
import TableActions from '../../../components/TableActions.vue'
import TablePagination from '../../../components/TablePagination.vue'
import FieldHint from '../../../components/FieldHint.vue'

const store = useMechaGameSectorsStore()
const gamesStore = useGamesStore()
const { selectedGame } = storeToRefs(gamesStore)

const formattedSectors = computed(() =>
  store.sectors.map(s => ({ ...s, is_starting_sector: s.is_starting_sector ? 'Yes' : 'No' }))
)

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'description', label: 'Description' },
  { key: 'elevation', label: 'Elevation' },
  { key: 'cover_modifier', label: 'Cover Mod.' },
  { key: 'is_starting_sector', label: 'Starting' },
]

const showModal = ref(false)
const modalMode = ref('create')
const modalForm = ref({ name: '', description: '', elevation: 0, cover_modifier: 0, is_starting_sector: false })
const modalError = ref('')
const showDeleteModal = ref(false)
const toDelete = ref(null)

// Build a list of numeric options from `min` to `max` stepping by `step`. If
// `current` is a value within the range but not on the step grid (e.g. a
// legacy record that predates the dropdown), it is inserted so the <select>
// can display it in edit mode without silently changing the data.
function buildOptions(min, max, step, current) {
  const opts = []
  for (let v = min; v <= max; v += step) opts.push(v)
  if (typeof current === 'number' && current >= min && current <= max && !opts.includes(current)) {
    opts.push(current)
    opts.sort((a, b) => a - b)
  }
  return opts
}

// Render signed numeric labels so cover-modifier polarity is unambiguous
// (e.g. "+5" for easier-to-hit, "-5" for harder-to-hit).
function formatSigned(v) {
  return v > 0 ? `+${v}` : String(v)
}

const elevationOptions = computed(() => buildOptions(-10, 10, 1, modalForm.value.elevation))
const coverModifierOptions = computed(() => buildOptions(-50, 50, 5, modalForm.value.cover_modifier))

watch(() => selectedGame.value, (g) => { if (g) store.fetchMechaGameSectors(g.id) }, { immediate: true })

function openCreate() {
  modalMode.value = 'create'
  modalForm.value = { name: '', description: '', elevation: 0, cover_modifier: 0, is_starting_sector: false }
  modalError.value = ''
  showModal.value = true
}

function openEdit(row) {
  modalMode.value = 'edit'
  const original = store.sectors.find(s => s.id === row.id)
  modalForm.value = { ...original }
  modalError.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  modalError.value = ''
}

async function handleSubmit(formData) {
  modalError.value = ''
  const allowed = ['name', 'description', 'elevation', 'cover_modifier', 'is_starting_sector']
  const data = Object.fromEntries(allowed.map(k => [k, formData[k]]))
  try {
    if (modalMode.value === 'create') {
      await store.createMechaGameSector(data)
    } else {
      await store.updateMechaGameSector(modalForm.value.id, data)
    }
    closeModal()
  } catch (e) {
    modalError.value = e.message || 'Failed to save.'
  }
}

async function handleDelete() {
  try {
    await store.deleteMechaGameSector(toDelete.value.id)
    showDeleteModal.value = false
    toDelete.value = null
  } catch (e) {
    console.error('Failed to delete sector:', e)
  }
}

function getActions(row) {
  return [
    { key: 'edit', label: 'Edit', handler: () => openEdit(row) },
    { key: 'delete', label: 'Delete', danger: true, handler: () => { toDelete.value = row; showDeleteModal.value = true } },
  ]
}
</script>

<style scoped>
.edit-link { color: var(--color-primary); text-decoration: none; }
.edit-link:hover { text-decoration: underline; }
.modal-form { display: flex; flex-direction: column; gap: 0.25rem; }
.form-group { margin-bottom: var(--space-sm); }
.form-group label { display: block; margin-bottom: var(--space-xs); font-weight: var(--font-weight-semibold); }
.form-group input, .form-group textarea { width: 100%; padding: var(--space-sm); border: 1px solid var(--color-border); border-radius: var(--radius-sm); font-size: var(--font-size-base); }
.form-group textarea { resize: vertical; }
.form-row { display: flex; gap: var(--space-sm); }
.form-row .half { flex: 1; }
.checkbox-group { margin-top: var(--space-sm); }
.checkbox-label { display: flex; align-items: center; gap: var(--space-sm); cursor: pointer; font-weight: normal; }
.checkbox-label input[type="checkbox"] { width: auto; }
.required { color: var(--color-danger); }
.error { color: var(--color-warning-dark); background: var(--color-warning-light); padding: var(--space-sm) var(--space-md); border-radius: var(--radius-sm); border: 1px solid var(--color-warning); margin-top: var(--space-md); }
.error p { margin: 0; }
</style>
