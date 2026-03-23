<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage lances.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Lances" actionText="Create New Lance" :showIcon="false" titleLevel="h2"
        @action="openCreate" />
      <ResourceTable :columns="columns" :rows="store.lances" :loading="store.loading" :error="store.error">
        <template #cell-name="{ row }">
          <span class="name-cell">
            <a href="#" class="edit-link" @click.prevent="openEdit(row)">{{ row.name }}</a>
            <span v-if="row.is_player_starter" class="starter-badge" title="Player Starter Lance">Starter</span>
          </span>
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>
      <TablePagination :pageNumber="store.pageNumber" :hasMore="store.hasMore"
        @page-change="(p) => store.fetchLances(selectedGame.id, p)" />
    </div>

    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal">
          <h2>{{ modalMode === 'create' ? 'Create Lance' : 'Edit Lance' }}</h2>
          <form @submit.prevent="handleSubmit(modalForm)" class="modal-form">
            <div class="form-group">
              <label>Name <span class="required">*</span></label>
              <input v-model="modalForm.name" required maxlength="100" autocomplete="off" />
            </div>
            <div class="form-group">
              <label>Description</label>
              <textarea v-model="modalForm.description" rows="3" maxlength="4096"></textarea>
            </div>
            <div class="form-group form-group--checkbox">
              <label class="checkbox-label">
                <input type="checkbox" v-model="modalForm.is_player_starter" />
                Player Starter Lance
              </label>
              <p class="field-hint">When enabled this lance serves as the template cloned for every new player who joins the game. Only one starter lance is allowed per game.</p>
            </div>
            <div v-if="!modalForm.is_player_starter" class="form-group">
              <label>Account User ID <span class="required">*</span></label>
              <input v-model="modalForm.account_user_id" :required="!modalForm.is_player_starter" autocomplete="off"
                placeholder="Account user ID of the lance commander" />
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

    <ConfirmationModal :visible="showDeleteModal" title="Delete Lance"
      :message="`Are you sure you want to delete '${toDelete?.name}'?`"
      @confirm="handleDelete" @cancel="showDeleteModal = false" />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useMechaLancesStore } from '../../../stores/mechaLances'
import { useGamesStore } from '../../../stores/games'
import ResourceTable from '../../../components/ResourceTable.vue'
import ConfirmationModal from '../../../components/ConfirmationModal.vue'
import PageHeader from '../../../components/PageHeader.vue'
import GameContext from '../../../components/GameContext.vue'
import TableActions from '../../../components/TableActions.vue'
import TablePagination from '../../../components/TablePagination.vue'

const store = useMechaLancesStore()
const gamesStore = useGamesStore()
const { selectedGame } = storeToRefs(gamesStore)

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'description', label: 'Description' },
  { key: 'account_user_id', label: 'Account User ID' },
]

const showModal = ref(false)
const modalMode = ref('create')
const modalForm = ref({ name: '', description: '', account_user_id: '', is_player_starter: false })
const modalError = ref('')
const showDeleteModal = ref(false)
const toDelete = ref(null)

watch(() => selectedGame.value, (g) => { if (g) store.fetchLances(g.id) }, { immediate: true })

function openCreate() {
  modalMode.value = 'create'
  modalForm.value = { name: '', description: '', account_user_id: '', is_player_starter: false }
  modalError.value = ''
  showModal.value = true
}

function openEdit(row) {
  modalMode.value = 'edit'
  modalForm.value = { ...row, is_player_starter: !!row.is_player_starter }
  modalError.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  modalError.value = ''
}

async function handleSubmit(formData) {
  modalError.value = ''
  const allowed = ['name', 'description', 'is_player_starter', 'account_user_id']
  const data = Object.fromEntries(allowed.map(k => [k, formData[k]]))
  if (data.is_player_starter) {
    delete data.account_user_id
  }
  try {
    if (modalMode.value === 'create') {
      await store.createLance(data)
    } else {
      await store.updateLance(modalForm.value.id, data)
    }
    closeModal()
  } catch (e) {
    modalError.value = e.message || 'Failed to save.'
  }
}

async function handleDelete() {
  try {
    await store.deleteLance(toDelete.value.id)
    showDeleteModal.value = false
    toDelete.value = null
  } catch (e) {
    console.error('Failed to delete lance:', e)
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
.name-cell { display: flex; align-items: center; gap: var(--space-xs); }
.starter-badge {
  display: inline-block;
  padding: 2px 6px;
  font-size: var(--font-size-xs, 0.75rem);
  font-weight: var(--font-weight-semibold);
  background: var(--color-primary);
  color: #fff;
  border-radius: var(--radius-sm);
  line-height: 1.4;
}
.modal-form { display: flex; flex-direction: column; gap: 0.25rem; }
.form-group { margin-bottom: var(--space-sm); }
.form-group label { display: block; margin-bottom: var(--space-xs); font-weight: var(--font-weight-semibold); }
.form-group input, .form-group textarea { width: 100%; padding: var(--space-sm); border: 1px solid var(--color-border); border-radius: var(--radius-sm); font-size: var(--font-size-base); }
.form-group textarea { resize: vertical; }
.form-group--checkbox label { display: flex; align-items: center; gap: var(--space-sm); cursor: pointer; }
.form-group--checkbox input[type="checkbox"] { width: auto; margin: 0; }
.checkbox-label { font-weight: var(--font-weight-semibold); }
.field-hint { font-size: var(--font-size-sm, 0.875rem); color: var(--color-text-muted); margin: var(--space-xs) 0 0; }
.required { color: var(--color-danger); }
.error { color: var(--color-warning-dark); background: var(--color-warning-light); padding: var(--space-sm) var(--space-md); border-radius: var(--radius-sm); border: 1px solid var(--color-warning); margin-top: var(--space-md); }
.error p { margin: 0; }
</style>
