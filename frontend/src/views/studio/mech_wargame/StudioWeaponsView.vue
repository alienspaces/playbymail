<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage weapons.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Weapons" actionText="Create New Weapon" :showIcon="false" titleLevel="h2"
        @action="openCreate" />
      <ResourceTable :columns="columns" :rows="store.weapons" :loading="store.loading" :error="store.error">
        <template #cell-name="{ row }">
          <a href="#" class="edit-link" @click.prevent="openEdit(row)">{{ row.name }}</a>
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>
      <TablePagination :pageNumber="store.pageNumber" :hasMore="store.hasMore"
        @page-change="(p) => store.fetchWeapons(selectedGame.id, p)" />
    </div>

    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal">
          <h2>{{ modalMode === 'create' ? 'Create Weapon' : 'Edit Weapon' }}</h2>
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
                <label>Damage</label>
                <input v-model.number="modalForm.damage" type="number" min="0" />
              </div>
              <div class="form-group half">
                <label>Heat Generated</label>
                <input v-model.number="modalForm.heat_generated" type="number" min="0" />
              </div>
            </div>
            <div class="form-row">
              <div class="form-group half">
                <label>Range</label>
                <input v-model.number="modalForm.range" type="number" min="0" />
              </div>
              <div class="form-group half">
                <label>Tonnage</label>
                <input v-model.number="modalForm.tonnage" type="number" min="0" />
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

    <ConfirmationModal :visible="showDeleteModal" title="Delete Weapon"
      :message="`Are you sure you want to delete '${toDelete?.name}'?`"
      @confirm="handleDelete" @cancel="showDeleteModal = false" />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useMechWargameWeaponsStore } from '../../../stores/mechWargameWeapons'
import { useGamesStore } from '../../../stores/games'
import ResourceTable from '../../../components/ResourceTable.vue'
import ConfirmationModal from '../../../components/ConfirmationModal.vue'
import PageHeader from '../../../components/PageHeader.vue'
import GameContext from '../../../components/GameContext.vue'
import TableActions from '../../../components/TableActions.vue'
import TablePagination from '../../../components/TablePagination.vue'

const store = useMechWargameWeaponsStore()
const gamesStore = useGamesStore()
const { selectedGame } = storeToRefs(gamesStore)

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'damage', label: 'Damage' },
  { key: 'heat_generated', label: 'Heat' },
  { key: 'range', label: 'Range' },
  { key: 'tonnage', label: 'Tonnage' },
]

const showModal = ref(false)
const modalMode = ref('create')
const modalForm = ref({ name: '', description: '', damage: 0, heat_generated: 0, range: 1, tonnage: 0 })
const modalError = ref('')
const showDeleteModal = ref(false)
const toDelete = ref(null)

watch(() => selectedGame.value, (g) => { if (g) store.fetchWeapons(g.id) }, { immediate: true })

function openCreate() {
  modalMode.value = 'create'
  modalForm.value = { name: '', description: '', damage: 0, heat_generated: 0, range: 1, tonnage: 0 }
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
  const allowed = ['name', 'description', 'damage', 'heat_generated', 'range', 'tonnage']
  const data = Object.fromEntries(allowed.map(k => [k, formData[k]]))
  try {
    if (modalMode.value === 'create') {
      await store.createWeapon(data)
    } else {
      await store.updateWeapon(modalForm.value.id, data)
    }
    closeModal()
  } catch (e) {
    modalError.value = e.message || 'Failed to save.'
  }
}

async function handleDelete() {
  try {
    await store.deleteWeapon(toDelete.value.id)
    showDeleteModal.value = false
    toDelete.value = null
  } catch (e) {
    console.error('Failed to delete weapon:', e)
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
.required { color: var(--color-danger); }
.error { color: var(--color-warning-dark); background: var(--color-warning-light); padding: var(--space-sm) var(--space-md); border-radius: var(--radius-sm); border: 1px solid var(--color-warning); margin-top: var(--space-md); }
.error p { margin: 0; }
</style>
