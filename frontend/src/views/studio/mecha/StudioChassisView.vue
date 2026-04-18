<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage chassis.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Chassis" actionText="Create New Chassis" :showIcon="false" titleLevel="h2"
        @action="openCreate" />
      <ResourceTable :columns="columns" :rows="store.chassis" :loading="store.loading" :error="store.error">
        <template #cell-name="{ row }">
          <a href="#" class="edit-link" @click.prevent="openEdit(row)">{{ row.name }}</a>
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>
      <TablePagination :pageNumber="store.pageNumber" :hasMore="store.hasMore"
        @page-change="(p) => store.fetchMechaGameChassis(selectedGame.id, p)" />
    </div>

    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal">
          <h2>{{ modalMode === 'create' ? 'Create Chassis' : 'Edit Chassis' }}</h2>
          <form @submit.prevent="handleSubmit(modalForm)" class="modal-form">
            <div class="form-group">
              <label>Name <span class="required">*</span></label>
              <input v-model="modalForm.name" required maxlength="100" autocomplete="off" />
              <FieldHint>Display name for this chassis (e.g. "Raven Light", "Enforcer").</FieldHint>
            </div>
            <div class="form-group">
              <label>Description</label>
              <textarea v-model="modalForm.description" rows="3" maxlength="4096"></textarea>
              <FieldHint>Narrative description shown to designers and players.</FieldHint>
            </div>
            <div class="form-group">
              <label>Chassis Class <span class="required">*</span></label>
              <select v-model="modalForm.chassis_class" required>
                <option value="light">Light</option>
                <option value="medium">Medium</option>
                <option value="heavy">Heavy</option>
                <option value="assault">Assault</option>
              </select>
              <FieldHint>Weight class — determines general role. Light is fast and lightly armoured; assault is the slowest with maximum armour and firepower.</FieldHint>
            </div>
            <div class="form-row">
              <div class="form-group half">
                <label>Armor Points <span class="required">*</span></label>
                <select v-model.number="modalForm.armor_points" required>
                  <option v-for="v in armorPointsOptions" :key="v" :value="v">{{ v }}</option>
                </select>
                <FieldHint>Maximum armour (50–1000, step 50). Absorbs damage before structure takes any; regenerates partially each turn at depot sectors.</FieldHint>
              </div>
              <div class="form-group half">
                <label>Structure Points <span class="required">*</span></label>
                <select v-model.number="modalForm.structure_points" required>
                  <option v-for="v in structurePointsOptions" :key="v" :value="v">{{ v }}</option>
                </select>
                <FieldHint>Maximum structure (50–1000, step 50). The mech is destroyed when this reaches 0. Repaired via the Squad Management sheet at a depot sector.</FieldHint>
              </div>
            </div>
            <div class="form-row">
              <div class="form-group half">
                <label>Heat Capacity <span class="required">*</span></label>
                <select v-model.number="modalForm.heat_capacity" required>
                  <option v-for="v in heatCapacityOptions" :key="v" :value="v">{{ v }}</option>
                </select>
                <FieldHint>Maximum heat before shutdown (10–200, step 10). Weapons add heat per shot; roughly capacity / 3 dissipates each turn.</FieldHint>
              </div>
              <div class="form-group half">
                <label>Speed <span class="required">*</span></label>
                <select v-model.number="modalForm.speed" required>
                  <option v-for="v in speedOptions" :key="v" :value="v">{{ v }}</option>
                </select>
                <FieldHint>Sector hops allowed per turn (1–10). Long-range weapons reach 2 hops, so values above ~5 are only useful for rapid repositioning.</FieldHint>
              </div>
            </div>
            <div class="form-row">
              <div class="form-group third">
                <label>Small Slots <span class="required">*</span></label>
                <select v-model.number="modalForm.small_slots" required @change="markSlotsTouched">
                  <option v-for="v in smallSlotsOptions" :key="v" :value="v">{{ v }}</option>
                </select>
                <FieldHint>Hardpoints that natively fit small weapons. Small weapons can also spill up into medium or large slots when these are full.</FieldHint>
              </div>
              <div class="form-group third">
                <label>Medium Slots <span class="required">*</span></label>
                <select v-model.number="modalForm.medium_slots" required @change="markSlotsTouched">
                  <option v-for="v in mediumSlotsOptions" :key="v" :value="v">{{ v }}</option>
                </select>
                <FieldHint>Hardpoints that natively fit medium weapons. Medium weapons spill up into large slots only.</FieldHint>
              </div>
              <div class="form-group third">
                <label>Large Slots <span class="required">*</span></label>
                <select v-model.number="modalForm.large_slots" required @change="markSlotsTouched">
                  <option v-for="v in largeSlotsOptions" :key="v" :value="v">{{ v }}</option>
                </select>
                <FieldHint>Hardpoints that natively fit large weapons. Large weapons can only use large slots.</FieldHint>
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

    <ConfirmationModal :visible="showDeleteModal" title="Delete Chassis"
      :message="`Are you sure you want to delete '${toDelete?.name}'?`"
      @confirm="handleDelete" @cancel="showDeleteModal = false" />
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useMechaGameChassisStore } from '../../../stores/mechaGameChassis'
import { useGamesStore } from '../../../stores/games'
import ResourceTable from '../../../components/ResourceTable.vue'
import ConfirmationModal from '../../../components/ConfirmationModal.vue'
import PageHeader from '../../../components/PageHeader.vue'
import GameContext from '../../../components/GameContext.vue'
import TableActions from '../../../components/TableActions.vue'
import TablePagination from '../../../components/TablePagination.vue'
import FieldHint from '../../../components/FieldHint.vue'

const store = useMechaGameChassisStore()
const gamesStore = useGamesStore()
const { selectedGame } = storeToRefs(gamesStore)

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'chassis_class', label: 'Class' },
  { key: 'armor_points', label: 'Armor' },
  { key: 'structure_points', label: 'Structure' },
  { key: 'heat_capacity', label: 'Heat Cap.' },
  { key: 'speed', label: 'Speed' },
]

const showModal = ref(false)
const modalMode = ref('create')
// Per-class default slot allocation; mirrors the per-class backfill applied by
// the mecha_game_chassis_slots migration and the backend's
// DefaultSlotsForChassisClass helper so a designer sees the same starting
// shape regardless of which layer creates the chassis.
function defaultSlotsForClass(cls) {
  switch (cls) {
    case 'light': return { small_slots: 2, medium_slots: 1, large_slots: 0 }
    case 'heavy': return { small_slots: 2, medium_slots: 2, large_slots: 2 }
    case 'assault': return { small_slots: 2, medium_slots: 3, large_slots: 3 }
    case 'medium':
    default: return { small_slots: 2, medium_slots: 2, large_slots: 1 }
  }
}
function freshModalForm(cls = 'medium') {
  return { name: '', description: '', chassis_class: cls, armor_points: 100, structure_points: 50, heat_capacity: 30, speed: 4, ...defaultSlotsForClass(cls) }
}
const modalForm = ref(freshModalForm())
// Tracks whether the designer has manually changed any slot value in the
// current modal session. Used to decide whether a chassis-class change should
// reapply per-class slot defaults (only if the user has not customised them).
const slotsTouched = ref(false)
const modalError = ref('')
const showDeleteModal = ref(false)
const toDelete = ref(null)

function markSlotsTouched() {
  slotsTouched.value = true
}

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

const armorPointsOptions = computed(() => buildOptions(50, 1000, 50, modalForm.value.armor_points))
const structurePointsOptions = computed(() => buildOptions(50, 1000, 50, modalForm.value.structure_points))
const heatCapacityOptions = computed(() => buildOptions(10, 200, 10, modalForm.value.heat_capacity))
const speedOptions = computed(() => buildOptions(1, 10, 1, modalForm.value.speed))
const smallSlotsOptions = computed(() => buildOptions(0, 10, 1, modalForm.value.small_slots))
const mediumSlotsOptions = computed(() => buildOptions(0, 10, 1, modalForm.value.medium_slots))
const largeSlotsOptions = computed(() => buildOptions(0, 10, 1, modalForm.value.large_slots))

watch(() => selectedGame.value, (g) => { if (g) store.fetchMechaGameChassis(g.id) }, { immediate: true })

// When the designer changes chassis class in create mode, reapply per-class
// slot defaults — but only if they have not already hand-tuned the slots. In
// edit mode we never auto-overwrite the stored loadout capacity.
watch(() => modalForm.value.chassis_class, (cls) => {
  if (!showModal.value) return
  if (modalMode.value !== 'create') return
  if (slotsTouched.value) return
  const d = defaultSlotsForClass(cls)
  modalForm.value.small_slots = d.small_slots
  modalForm.value.medium_slots = d.medium_slots
  modalForm.value.large_slots = d.large_slots
})

function openCreate() {
  modalMode.value = 'create'
  modalForm.value = freshModalForm()
  slotsTouched.value = false
  modalError.value = ''
  showModal.value = true
}

function openEdit(row) {
  modalMode.value = 'edit'
  modalForm.value = { ...row }
  slotsTouched.value = false
  modalError.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  modalError.value = ''
}

async function handleSubmit(formData) {
  modalError.value = ''
  const allowed = ['name', 'description', 'chassis_class', 'armor_points', 'structure_points', 'heat_capacity', 'speed', 'small_slots', 'medium_slots', 'large_slots']
  const data = Object.fromEntries(allowed.map(k => [k, formData[k]]))
  try {
    if (modalMode.value === 'create') {
      await store.createMechaGameChassis(data)
    } else {
      await store.updateMechaGameChassis(modalForm.value.id, data)
    }
    closeModal()
  } catch (e) {
    modalError.value = e.message || 'Failed to save.'
  }
}

async function handleDelete() {
  try {
    await store.deleteMechaGameChassis(toDelete.value.id)
    showDeleteModal.value = false
    toDelete.value = null
  } catch (e) {
    console.error('Failed to delete chassis:', e)
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
.form-group input, .form-group select, .form-group textarea { width: 100%; padding: var(--space-sm); border: 1px solid var(--color-border); border-radius: var(--radius-sm); font-size: var(--font-size-base); }
.form-group textarea { resize: vertical; }
.form-row { display: flex; gap: var(--space-sm); }
.form-row .half { flex: 1; }
.form-row .third { flex: 1; }
.required { color: var(--color-danger); }
.error { color: var(--color-warning-dark); background: var(--color-warning-light); padding: var(--space-sm) var(--space-md); border-radius: var(--radius-sm); border: 1px solid var(--color-warning); margin-top: var(--space-md); }
.error p { margin: 0; }
</style>
