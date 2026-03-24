<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage sector links.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Sector Links" actionText="Create New Link" :showIcon="false" titleLevel="h2"
        @action="openCreate" />
      <ResourceTable :columns="columns" :rows="formattedLinks" :loading="store.loading" :error="store.error">
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>
      <TablePagination :pageNumber="store.pageNumber" :hasMore="store.hasMore"
        @page-change="(p) => store.fetchSectorLinks(selectedGame.id, p)" />
    </div>

    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal">
          <h2>{{ modalMode === 'create' ? 'Create Sector Link' : 'Edit Sector Link' }}</h2>
          <form @submit.prevent="handleSubmit(modalForm)" class="modal-form">
            <div class="form-group">
              <label>From Sector <span class="required">*</span></label>
              <select v-model="modalForm.from_mecha_sector_id" required>
                <option value="">-- Select sector --</option>
                <option v-for="s in sectorsStore.sectors" :key="s.id" :value="s.id">{{ s.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label>To Sector <span class="required">*</span></label>
              <select v-model="modalForm.to_mecha_sector_id" required>
                <option value="">-- Select sector --</option>
                <option v-for="s in sectorsStore.sectors" :key="s.id" :value="s.id">{{ s.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label>Cover Modifier</label>
              <input v-model.number="modalForm.cover_modifier" type="number" min="-10" max="10" />
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

    <ConfirmationModal :visible="showDeleteModal" title="Delete Sector Link"
      message="Are you sure you want to delete this sector link?"
      @confirm="handleDelete" @cancel="showDeleteModal = false" />
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useMechaSectorLinksStore } from '../../../stores/mechaSectorLinks'
import { useMechaSectorsStore } from '../../../stores/mechaSectors'
import { useGamesStore } from '../../../stores/games'
import ResourceTable from '../../../components/ResourceTable.vue'
import ConfirmationModal from '../../../components/ConfirmationModal.vue'
import PageHeader from '../../../components/PageHeader.vue'
import GameContext from '../../../components/GameContext.vue'
import TableActions from '../../../components/TableActions.vue'
import TablePagination from '../../../components/TablePagination.vue'

const store = useMechaSectorLinksStore()
const sectorsStore = useMechaSectorsStore()
const gamesStore = useGamesStore()
const { selectedGame } = storeToRefs(gamesStore)

function sectorName(id) {
  return sectorsStore.sectors.find(s => s.id === id)?.name ?? id
}

const formattedLinks = computed(() =>
  store.sectorLinks.map(sl => ({
    ...sl,
    from_sector_name: sectorName(sl.from_mecha_sector_id),
    to_sector_name: sectorName(sl.to_mecha_sector_id),
  }))
)

const columns = [
  { key: 'from_sector_name', label: 'From Sector' },
  { key: 'to_sector_name', label: 'To Sector' },
  { key: 'cover_modifier', label: 'Cover Mod.' },
]

const showModal = ref(false)
const modalMode = ref('create')
const modalForm = ref({ from_mecha_sector_id: '', to_mecha_sector_id: '', cover_modifier: 0 })
const modalError = ref('')
const showDeleteModal = ref(false)
const toDelete = ref(null)

watch(() => selectedGame.value, (g) => {
  if (g) {
    store.fetchSectorLinks(g.id)
    sectorsStore.fetchSectors(g.id)
  }
}, { immediate: true })

function openCreate() {
  modalMode.value = 'create'
  modalForm.value = { from_mecha_sector_id: '', to_mecha_sector_id: '', cover_modifier: 0 }
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
  const allowed = ['from_mecha_sector_id', 'to_mecha_sector_id', 'cover_modifier']
  const data = Object.fromEntries(allowed.map(k => [k, formData[k]]))
  try {
    if (modalMode.value === 'create') {
      await store.createSectorLink(data)
    } else {
      await store.updateSectorLink(modalForm.value.id, data)
    }
    closeModal()
  } catch (e) {
    modalError.value = e.message || 'Failed to save.'
  }
}

async function handleDelete() {
  try {
    await store.deleteSectorLink(toDelete.value.id)
    showDeleteModal.value = false
    toDelete.value = null
  } catch (e) {
    console.error('Failed to delete sector link:', e)
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
.modal-form { display: flex; flex-direction: column; gap: 0.25rem; }
.form-group { margin-bottom: var(--space-sm); }
.form-group label { display: block; margin-bottom: var(--space-xs); font-weight: var(--font-weight-semibold); }
.form-group input, .form-group select { width: 100%; padding: var(--space-sm); border: 1px solid var(--color-border); border-radius: var(--radius-sm); font-size: var(--font-size-base); }
.required { color: var(--color-danger); }
.error { color: var(--color-warning-dark); background: var(--color-warning-light); padding: var(--space-sm) var(--space-md); border-radius: var(--radius-sm); border: 1px solid var(--color-warning); margin-top: var(--space-md); }
.error p { margin: 0; }
</style>
