<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage squads.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Squads" actionText="Create New Squad" :showIcon="false" titleLevel="h2"
        @action="openCreate" />
      <ResourceTable :columns="columns" :rows="store.squads" :loading="store.loading" :error="store.error">
        <template #cell-name="{ row }">
          <span class="name-cell">
            <a href="#" class="expand-link" @click.prevent="toggleMechs(row)">
              <span class="expand-icon">{{ expandedSquadId === row.id ? '▾' : '▸' }}</span>
              {{ row.name }}
            </a>
            <span v-if="row.squad_type === 'starter'" class="starter-badge" title="Player Starter Squad">Starter</span>
          </span>
        </template>
        <template #cell-squad_type="{ row }">
          {{ row.squad_type === 'starter' ? 'Starter' : 'Opponent' }}
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
        <template #row-detail="{ row }">
          <div v-if="expandedSquadId === row.id" class="mech-panel">
            <div class="mech-panel-header">
              <span class="mech-panel-title">Mechs in {{ row.name }}</span>
              <button class="btn-small" @click="openMechCreate(row)">Add Mech</button>
            </div>
            <div v-if="mechsStore.loading" class="mech-loading">Loading mechs...</div>
            <div v-else-if="mechsStore.getMechsForSquad(row.id).length === 0" class="mech-empty">
              No mechs assigned to this squad.
            </div>
            <table v-else class="mech-table">
              <thead>
                <tr>
                  <th>Callsign</th>
                  <th>Chassis</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="mech in mechsStore.getMechsForSquad(row.id)" :key="mech.id">
                  <td>{{ mech.callsign }}</td>
                  <td>{{ chassisName(mech.mecha_game_chassis_id) }}</td>
                  <td>
                    <button class="btn-link" @click="openMechEdit(row, mech)">Edit</button>
                    <button class="btn-link btn-danger" @click="confirmMechDelete(row, mech)">Delete</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </template>
      </ResourceTable>
      <TablePagination :pageNumber="store.pageNumber" :hasMore="store.hasMore"
        @page-change="(p) => store.fetchMechaGameSquads(selectedGame.id, p)" />
    </div>

    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal">
          <h2>{{ modalMode === 'create' ? 'Create Squad' : 'Edit Squad' }}</h2>
          <form @submit.prevent="handleSubmit(modalForm)" class="modal-form">
            <div class="form-group">
              <label>Name <span class="required">*</span></label>
              <input v-model="modalForm.name" required maxlength="100" autocomplete="off" />
            </div>
            <div class="form-group">
              <label>Description</label>
              <textarea v-model="modalForm.description" rows="3" maxlength="4096"></textarea>
            </div>
            <div class="form-group">
              <label>Type <span class="required">*</span></label>
              <select v-model="modalForm.squad_type" required>
                <option value="opponent">Opponent</option>
                <option value="starter">Starter</option>
              </select>
              <FieldHint>
                <span v-if="modalForm.squad_type === 'starter'">This squad is cloned for every player who joins the game. Only one starter squad is allowed per game.</span>
                <span v-else>This squad template is randomly assigned to a computer opponent when the game starts.</span>
              </FieldHint>
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

    <ConfirmationModal :visible="showDeleteModal" title="Delete Squad"
      :message="`Are you sure you want to delete '${toDelete?.name}'?`"
      @confirm="handleDelete" @cancel="showDeleteModal = false" />

    <div v-if="showMechModal" class="modal-overlay" @click.self="closeMechModal">
      <div class="modal">
        <h2>{{ mechModalMode === 'create' ? 'Add Mech' : 'Edit Mech' }}</h2>
        <form @submit.prevent="handleMechSubmit(mechModalForm)" class="modal-form">
          <div class="form-group">
            <label>Callsign <span class="required">*</span></label>
            <input v-model="mechModalForm.callsign" required maxlength="100" autocomplete="off"
              placeholder="e.g. Iron Fist, Wraith" />
          </div>
          <div class="form-group">
            <label>Chassis <span class="required">*</span></label>
            <select v-model="mechModalForm.mecha_game_chassis_id" required>
              <option value="" disabled>Select chassis...</option>
              <option v-for="chassis in (chassisStore.chassis || [])" :key="chassis.id" :value="chassis.id">
                {{ chassis.name }}
              </option>
            </select>
          </div>
          <div class="modal-actions">
            <button type="submit">{{ mechModalMode === 'create' ? 'Add' : 'Save' }}</button>
            <button type="button" @click="closeMechModal">Cancel</button>
          </div>
        </form>
        <div v-if="mechModalError" class="error"><p>{{ mechModalError }}</p></div>
      </div>
    </div>

    <ConfirmationModal :visible="showMechDeleteModal" title="Delete Mech"
      :message="`Are you sure you want to remove '${mechToDelete?.callsign}' from this squad??`"
      @confirm="handleMechDelete" @cancel="showMechDeleteModal = false" />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useMechaGameSquadsStore } from '../../../stores/mechaGameSquads'
import { useMechaGameSquadMechsStore } from '../../../stores/mechaGameSquadMechs'
import { useMechaGameChassisStore } from '../../../stores/mechaGameChassis'
import { useGamesStore } from '../../../stores/games'
import ResourceTable from '../../../components/ResourceTable.vue'
import ConfirmationModal from '../../../components/ConfirmationModal.vue'
import PageHeader from '../../../components/PageHeader.vue'
import GameContext from '../../../components/GameContext.vue'
import TableActions from '../../../components/TableActions.vue'
import TablePagination from '../../../components/TablePagination.vue'
import FieldHint from '../../../components/FieldHint.vue'

const store = useMechaGameSquadsStore()
const mechsStore = useMechaGameSquadMechsStore()
const chassisStore = useMechaGameChassisStore()
const gamesStore = useGamesStore()
const { selectedGame } = storeToRefs(gamesStore)

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'squad_type', label: 'Type' },
  { key: 'description', label: 'Description' },
]

const showModal = ref(false)
const modalMode = ref('create')
const modalForm = ref({ name: '', description: '', squad_type: 'opponent' })
const modalError = ref('')
const showDeleteModal = ref(false)
const toDelete = ref(null)

// Squad mech sub-panel state
const expandedSquadId = ref(null)
const showMechModal = ref(false)
const mechModalMode = ref('create')
const mechModalForm = ref({ callsign: '', mecha_game_chassis_id: '', squad_id: '' })
const mechModalError = ref('')
const showMechDeleteModal = ref(false)
const mechToDelete = ref(null)
const activeSquadForMech = ref(null)

watch(() => selectedGame.value, (g) => {
  if (g) {
    store.fetchMechaGameSquads(g.id)
    chassisStore.fetchMechaGameChassis(g.id)
  }
}, { immediate: true })

function chassisName(chassisId) {
  const found = (chassisStore.chassis || []).find(c => c.id === chassisId)
  return found ? found.name : chassisId
}

function openCreate() {
  modalMode.value = 'create'
  modalForm.value = { name: '', description: '', squad_type: 'opponent' }
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
  const data = {
    name: formData.name,
    description: formData.description,
    squad_type: formData.squad_type,
  }
  try {
    if (modalMode.value === 'create') {
      await store.createMechaGameSquad(data)
    } else {
      await store.updateMechaGameSquad(modalForm.value.id, data)
    }
    closeModal()
  } catch (e) {
    modalError.value = e.message || 'Failed to save.'
  }
}

async function handleDelete() {
  try {
    await store.deleteMechaGameSquad(toDelete.value.id)
    showDeleteModal.value = false
    toDelete.value = null
  } catch (e) {
    console.error('Failed to delete squad:', e)
  }
}

function getActions(row) {
  return [
    { key: 'edit', label: 'Edit', handler: () => openEdit(row) },
    { key: 'delete', label: 'Delete', danger: true, handler: () => { toDelete.value = row; showDeleteModal.value = true } },
  ]
}

// Squad mech sub-panel functions
async function toggleMechs(row) {
  if (expandedSquadId.value === row.id) {
    expandedSquadId.value = null
  } else {
    expandedSquadId.value = row.id
    mechsStore.gameId = selectedGame.value.id
    await mechsStore.fetchMechaGameSquadMechs(selectedGame.value.id, row.id)
  }
}

function openMechCreate(squad) {
  activeSquadForMech.value = squad
  mechModalMode.value = 'create'
  mechModalForm.value = { callsign: '', mecha_game_chassis_id: '', squad_id: squad.id }
  mechModalError.value = ''
  showMechModal.value = true
}

function openMechEdit(squad, mech) {
  activeSquadForMech.value = squad
  mechModalMode.value = 'edit'
  mechModalForm.value = { ...mech, squad_id: squad.id }
  mechModalError.value = ''
  showMechModal.value = true
}

function closeMechModal() {
  showMechModal.value = false
  mechModalError.value = ''
}

async function handleMechSubmit(formData) {
  mechModalError.value = ''
  const squadId = activeSquadForMech.value.id
  const data = {
    callsign: formData.callsign,
    mecha_game_chassis_id: formData.mecha_game_chassis_id,
  }
  try {
    if (mechModalMode.value === 'create') {
      await mechsStore.createMechaGameSquadMech(squadId, data)
    } else {
      await mechsStore.updateMechaGameSquadMech(squadId, formData.id, data)
    }
    closeMechModal()
  } catch (e) {
    mechModalError.value = e.message || 'Failed to save mech.'
  }
}

function confirmMechDelete(squad, mech) {
  activeSquadForMech.value = squad
  mechToDelete.value = mech
  showMechDeleteModal.value = true
}

async function handleMechDelete() {
  try {
    await mechsStore.deleteMechaGameSquadMech(activeSquadForMech.value.id, mechToDelete.value.id)
    showMechDeleteModal.value = false
    mechToDelete.value = null
  } catch (e) {
    console.error('Failed to delete mech:', e)
  }
}
</script>

<style scoped>
.expand-link { color: var(--color-primary); text-decoration: none; display: inline-flex; align-items: center; gap: 0.25em; }
.expand-link:hover { text-decoration: underline; }
.expand-icon { font-size: 0.85em; }
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
.form-group input, .form-group textarea, .form-group select { width: 100%; padding: var(--space-sm); border: 1px solid var(--color-border); border-radius: var(--radius-sm); font-size: var(--font-size-base); }
.form-group textarea { resize: vertical; }
.required { color: var(--color-danger); }
.error { color: var(--color-warning-dark); background: var(--color-warning-light); padding: var(--space-sm) var(--space-md); border-radius: var(--radius-sm); border: 1px solid var(--color-warning); margin-top: var(--space-md); }
.error p { margin: 0; }

/* Mech sub-panel */
.mech-panel { background: var(--color-bg-subtle, #f9f9f9); border: 1px solid var(--color-border); border-radius: var(--radius-sm); padding: var(--space-sm) var(--space-md); margin: var(--space-xs) 0; }
.mech-panel-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: var(--space-sm); }
.mech-panel-title { font-weight: var(--font-weight-semibold); font-size: var(--font-size-sm, 0.875rem); }
.mech-table { width: 100%; border-collapse: collapse; font-size: var(--font-size-sm, 0.875rem); }
.mech-table th, .mech-table td { text-align: left; padding: var(--space-xs) var(--space-sm); border-bottom: 1px solid var(--color-border); }
.mech-table th { font-weight: var(--font-weight-semibold); color: var(--color-text-muted); }
.mech-loading, .mech-empty { color: var(--color-text-muted); font-size: var(--font-size-sm, 0.875rem); padding: var(--space-xs) 0; }
.btn-small { padding: 2px 10px; font-size: var(--font-size-sm, 0.875rem); border: 1px solid var(--color-border); border-radius: var(--radius-sm); background: var(--color-bg); cursor: pointer; }
.btn-small:hover { background: var(--color-bg-subtle, #f3f3f3); }
.btn-link { background: none; border: none; padding: 0 var(--space-xs); cursor: pointer; color: var(--color-primary); font-size: var(--font-size-sm, 0.875rem); }
.btn-link:hover { text-decoration: underline; }
.btn-danger { color: var(--color-danger); }
select { width: 100%; padding: var(--space-sm); border: 1px solid var(--color-border); border-radius: var(--radius-sm); font-size: var(--font-size-base); }
</style>
