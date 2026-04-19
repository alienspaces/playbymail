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
        @page-change="(p) => store.fetchMechaGameWeapons(selectedGame.id, p)" />
    </div>

    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal">
          <h2>{{ modalMode === 'create' ? 'Create Weapon' : 'Edit Weapon' }}</h2>
          <form @submit.prevent="handleSubmit(modalForm)" class="modal-form">
            <div class="form-group">
              <label>Name <span class="required">*</span></label>
              <input v-model="modalForm.name" required maxlength="100" autocomplete="off" />
              <FieldHint>Display name for this weapon (e.g. "AC/10", "Medium Laser", "LRM-20").</FieldHint>
            </div>
            <div class="form-group">
              <label>Description</label>
              <textarea v-model="modalForm.description" rows="3" maxlength="4096"></textarea>
              <FieldHint>Narrative description shown to designers and players.</FieldHint>
            </div>
            <div class="form-row">
              <div class="form-group half">
                <label>Damage <span class="required">*</span></label>
                <select v-model.number="modalForm.damage" required>
                  <option v-for="v in damageOptions" :key="v" :value="v">{{ v }}</option>
                </select>
                <FieldHint>Damage dealt per hit (1–20). Armour absorbs damage first; anything left over eats into structure.</FieldHint>
              </div>
              <div class="form-group half">
                <label>Heat Cost <span class="required">*</span></label>
                <select v-model.number="modalForm.heat_cost" required>
                  <option v-for="v in heatCostOptions" :key="v" :value="v">{{ v }}</option>
                </select>
                <FieldHint>Heat added to the firing mech each time this weapon fires (0–20), even on a miss. Overheating shuts the mech down.</FieldHint>
              </div>
            </div>
            <div class="form-row">
              <div class="form-group half">
                <label>Range Band <span class="required">*</span></label>
                <select v-model="modalForm.range_band" required>
                  <option value="short">Short (brawl)</option>
                  <option value="medium">Medium (versatile)</option>
                  <option value="long">Long (standoff)</option>
                </select>
                <FieldHint>Engagement distance. Short fires in the same sector only; medium fires same sector or 1 hop; long fires 1–2 hops but cannot fire in the same sector.</FieldHint>
              </div>
              <div class="form-group half">
                <label>Mount Size <span class="required">*</span></label>
                <select v-model="modalForm.mount_size" required>
                  <option value="small">Small</option>
                  <option value="medium">Medium</option>
                  <option value="large">Large</option>
                </select>
                <FieldHint>Mount size category (small / medium / large). Must fit an available slot on the chassis — small weapons can spill into medium or large slots, medium weapons into large, but large weapons only fit large slots.</FieldHint>
              </div>
            </div>
            <div class="form-row">
              <div class="form-group half">
                <label>Ammo Capacity</label>
                <input type="number" min="0" max="200" step="1" v-model.number="modalForm.ammo_capacity" />
                <FieldHint>Rounds each trigger-pull draws from the mech's shared ammo pool. Set to 0 for energy/beam weapons that never consume ammo. Mechs can carry ammo bins to add to the pool; when the pool hits 0, the weapon simply cannot fire until rearmed at a depot.</FieldHint>
              </div>
              <div class="form-group half"></div>
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
import { ref, computed, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useMechaGameWeaponsStore } from '../../../stores/mechaGameWeapons'
import { useGamesStore } from '../../../stores/games'
import ResourceTable from '../../../components/ResourceTable.vue'
import ConfirmationModal from '../../../components/ConfirmationModal.vue'
import PageHeader from '../../../components/PageHeader.vue'
import GameContext from '../../../components/GameContext.vue'
import TableActions from '../../../components/TableActions.vue'
import TablePagination from '../../../components/TablePagination.vue'
import FieldHint from '../../../components/FieldHint.vue'

const store = useMechaGameWeaponsStore()
const gamesStore = useGamesStore()
const { selectedGame } = storeToRefs(gamesStore)

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'damage', label: 'Damage' },
  { key: 'heat_cost', label: 'Heat' },
  { key: 'range_band', label: 'Range' },
  { key: 'mount_size', label: 'Mount' },
  { key: 'ammo_capacity', label: 'Ammo' },
]

const showModal = ref(false)
const modalMode = ref('create')
const modalForm = ref({ name: '', description: '', damage: 5, heat_cost: 3, range_band: 'medium', mount_size: 'medium', ammo_capacity: 0 })
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

const damageOptions = computed(() => buildOptions(1, 20, 1, modalForm.value.damage))
const heatCostOptions = computed(() => buildOptions(0, 20, 1, modalForm.value.heat_cost))

watch(() => selectedGame.value, (g) => { if (g) store.fetchMechaGameWeapons(g.id) }, { immediate: true })

function openCreate() {
  modalMode.value = 'create'
  modalForm.value = { name: '', description: '', damage: 5, heat_cost: 3, range_band: 'medium', mount_size: 'medium', ammo_capacity: 0 }
  modalError.value = ''
  showModal.value = true
}

function openEdit(row) {
  modalMode.value = 'edit'
  modalForm.value = { ammo_capacity: 0, ...row }
  modalError.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  modalError.value = ''
}

async function handleSubmit(formData) {
  modalError.value = ''
  const allowed = ['name', 'description', 'damage', 'heat_cost', 'range_band', 'mount_size', 'ammo_capacity']
  const data = Object.fromEntries(allowed.map(k => [k, formData[k]]))
  try {
    if (modalMode.value === 'create') {
      await store.createMechaGameWeapon(data)
    } else {
      await store.updateMechaGameWeapon(modalForm.value.id, data)
    }
    closeModal()
  } catch (e) {
    modalError.value = e.message || 'Failed to save.'
  }
}

async function handleDelete() {
  try {
    await store.deleteMechaGameWeapon(toDelete.value.id)
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
.form-group input, .form-group textarea, .form-group select { width: 100%; padding: var(--space-sm); border: 1px solid var(--color-border); border-radius: var(--radius-sm); font-size: var(--font-size-base); }
.form-group textarea { resize: vertical; }
.form-row { display: flex; gap: var(--space-sm); }
.form-row .half { flex: 1; }
.required { color: var(--color-danger); }
.error { color: var(--color-warning-dark); background: var(--color-warning-light); padding: var(--space-sm) var(--space-md); border-radius: var(--radius-sm); border: 1px solid var(--color-warning); margin-top: var(--space-md); }
.error p { margin: 0; }
</style>
