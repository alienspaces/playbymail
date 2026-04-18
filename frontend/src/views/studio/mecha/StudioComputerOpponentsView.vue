<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage computer opponents.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Computer Opponents" actionText="Create New Opponent" :showIcon="false" titleLevel="h2"
        @action="openCreate" />
      <ResourceTable :columns="columns" :rows="store.opponents" :loading="store.loading" :error="store.error">
        <template #cell-name="{ row }">
          <a href="#" class="edit-link" @click.prevent="openEdit(row)">{{ row.name }}</a>
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>
      <TablePagination :pageNumber="store.pageNumber" :hasMore="store.hasMore"
        @page-change="(p) => store.fetchMechaGameComputerOpponents(selectedGame.id, p)" />
    </div>

    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal">
          <h2>{{ modalMode === 'create' ? 'Create Computer Opponent' : 'Edit Computer Opponent' }}</h2>
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
                <label>Aggression (1–10)</label>
                <input v-model.number="modalForm.aggression" type="number" min="1" max="10" />
                <p class="field-hint">1 = defensive, 10 = all-out assault</p>
              </div>
              <div class="form-group half">
                <label>IQ (1–10)</label>
                <input v-model.number="modalForm.iq" type="number" min="1" max="10" />
                <p class="field-hint">1 = predictable, 10 = expert tactics</p>
              </div>
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

    <ConfirmationModal :visible="showDeleteModal" title="Delete Computer Opponent"
      :message="`Are you sure you want to delete '${toDelete?.name}'?`"
      @confirm="handleDelete" @cancel="showDeleteModal = false" />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useMechaGameComputerOpponentsStore } from '../../../stores/mechaGameComputerOpponents'
import { useGamesStore } from '../../../stores/games'
import ResourceTable from '../../../components/ResourceTable.vue'
import ConfirmationModal from '../../../components/ConfirmationModal.vue'
import PageHeader from '../../../components/PageHeader.vue'
import GameContext from '../../../components/GameContext.vue'
import TableActions from '../../../components/TableActions.vue'
import TablePagination from '../../../components/TablePagination.vue'

const store = useMechaGameComputerOpponentsStore()
const gamesStore = useGamesStore()
const { selectedGame } = storeToRefs(gamesStore)

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'description', label: 'Description' },
  { key: 'aggression', label: 'Aggression' },
  { key: 'iq', label: 'IQ' },
]

const showModal = ref(false)
const modalMode = ref('create')
const modalForm = ref({ name: '', description: '', aggression: 5, iq: 5 })
const modalError = ref('')
const showDeleteModal = ref(false)
const toDelete = ref(null)

watch(() => selectedGame.value, (g) => { if (g) store.fetchMechaGameComputerOpponents(g.id) }, { immediate: true })

function openCreate() {
  modalMode.value = 'create'
  modalForm.value = { name: '', description: '', aggression: 5, iq: 5 }
  modalError.value = ''
  showModal.value = true
}

function openEdit(row) {
  modalMode.value = 'edit'
  modalForm.value = { ...row }
  modalError.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  modalError.value = ''
}

async function handleSubmit(formData) {
  modalError.value = ''
  const allowed = ['name', 'description', 'aggression', 'iq']
  const data = Object.fromEntries(allowed.map(k => [k, formData[k]]))
  try {
    if (modalMode.value === 'create') {
      await store.createMechaGameComputerOpponent(data)
    } else {
      await store.updateMechaGameComputerOpponent(modalForm.value.id, data)
    }
    closeModal()
  } catch (e) {
    modalError.value = e.message || 'Failed to save.'
  }
}

async function handleDelete() {
  try {
    await store.deleteMechaGameComputerOpponent(toDelete.value.id)
    showDeleteModal.value = false
    toDelete.value = null
  } catch (e) {
    console.error('Failed to delete computer opponent:', e)
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
.field-hint { font-size: var(--font-size-sm, 0.875rem); color: var(--color-text-muted); margin: var(--space-xs) 0 0; }
.required { color: var(--color-danger); }
.error { color: var(--color-warning-dark); background: var(--color-warning-light); padding: var(--space-sm) var(--space-md); border-radius: var(--radius-sm); border: 1px solid var(--color-warning); margin-top: var(--space-md); }
.error p { margin: 0; }
</style>
