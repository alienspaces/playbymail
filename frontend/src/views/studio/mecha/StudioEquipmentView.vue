<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage equipment.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Equipment" actionText="Create New Equipment" :showIcon="false" titleLevel="h2"
        @action="openCreate" />
      <ResourceTable :columns="columns" :rows="store.equipment" :loading="store.loading" :error="store.error">
        <template #cell-name="{ row }">
          <a href="#" class="edit-link" @click.prevent="openEdit(row)">{{ row.name }}</a>
        </template>
        <template #cell-effect_kind="{ row }">
          {{ formatEffectKind(row.effect_kind) }}
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>
      <TablePagination :pageNumber="store.pageNumber" :hasMore="store.hasMore"
        @page-change="(p) => store.fetchMechaGameEquipment(selectedGame.id, p)" />
    </div>

    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal">
          <h2>{{ modalMode === 'create' ? 'Create Equipment' : 'Edit Equipment' }}</h2>
          <form @submit.prevent="handleSubmit(modalForm)" class="modal-form">
            <div class="form-group">
              <label>Name <span class="required">*</span></label>
              <input v-model="modalForm.name" required maxlength="100" autocomplete="off" />
              <FieldHint>Display name for this equipment (e.g. "Double Heat Sink", "Targeting Computer Mk II", "Jump Jets", "Ammo Bin (Standard)").</FieldHint>
            </div>
            <div class="form-group">
              <label>Description</label>
              <textarea v-model="modalForm.description" rows="3" maxlength="4096"></textarea>
              <FieldHint>Narrative description shown to designers and players.</FieldHint>
            </div>
            <div class="form-row">
              <div class="form-group half">
                <label>Effect Kind <span class="required">*</span></label>
                <select v-model="modalForm.effect_kind" required @change="onEffectKindChange">
                  <option value="heat_sink">Heat Sink</option>
                  <option value="targeting_computer">Targeting Computer</option>
                  <option value="armor_upgrade">Armor Upgrade</option>
                  <option value="jump_jets">Jump Jets</option>
                  <option value="ecm">ECM</option>
                  <option value="ammo_bin">Ammo Bin</option>
                </select>
                <FieldHint>What this piece of equipment does. Effects are strictly additive bonuses — they enhance the mech without gating existing capabilities. Heat sinks add dissipation, targeting computers add hit chance, armor upgrades add max armor, jump jets add move hops, ECM adds a cover penalty for incoming attacks, and ammo bins add rounds to the mech's shared ammo pool.</FieldHint>
              </div>
              <div class="form-group half">
                <label>Mount Size <span class="required">*</span></label>
                <select v-model="modalForm.mount_size" required>
                  <option value="small">Small</option>
                  <option value="medium">Medium</option>
                  <option value="large">Large</option>
                </select>
                <FieldHint>Slot category consumed on a chassis. Small fits small / medium / large slots, medium fits medium / large, large only fits large. Shares the slot budget with weapons.</FieldHint>
              </div>
            </div>
            <div class="form-row">
              <div class="form-group half">
                <label>Magnitude <span class="required">*</span></label>
                <input type="number" min="1" :max="magnitudeMax" step="1" v-model.number="modalForm.magnitude" required />
                <FieldHint>{{ magnitudeHint }}</FieldHint>
              </div>
              <div class="form-group half">
                <label>Heat Cost</label>
                <input type="number" min="0" max="20" step="1" v-model.number="modalForm.heat_cost" />
                <FieldHint>{{ heatCostHint }}</FieldHint>
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

    <ConfirmationModal :visible="showDeleteModal" title="Delete Equipment"
      :message="`Are you sure you want to delete '${toDelete?.name}'?`"
      @confirm="handleDelete" @cancel="showDeleteModal = false" />
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useMechaGameEquipmentStore } from '../../../stores/mechaGameEquipment'
import { useGamesStore } from '../../../stores/games'
import ResourceTable from '../../../components/ResourceTable.vue'
import ConfirmationModal from '../../../components/ConfirmationModal.vue'
import PageHeader from '../../../components/PageHeader.vue'
import GameContext from '../../../components/GameContext.vue'
import TableActions from '../../../components/TableActions.vue'
import TablePagination from '../../../components/TablePagination.vue'
import FieldHint from '../../../components/FieldHint.vue'

const store = useMechaGameEquipmentStore()
const gamesStore = useGamesStore()
const { selectedGame } = storeToRefs(gamesStore)

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'effect_kind', label: 'Effect' },
  { key: 'magnitude', label: 'Magnitude' },
  { key: 'heat_cost', label: 'Heat' },
  { key: 'mount_size', label: 'Mount' },
]

// Per-kind magnitude caps must match backend
// record/mecha_game_record/mecha_game_equipment.go MagnitudeMaxForEffectKind.
const MAGNITUDE_MAX = {
  heat_sink: 20,
  targeting_computer: 30,
  armor_upgrade: 200,
  jump_jets: 5,
  ecm: 50,
  ammo_bin: 200,
}

const EFFECT_KIND_LABELS = {
  heat_sink: 'Heat Sink',
  targeting_computer: 'Targeting Computer',
  armor_upgrade: 'Armor Upgrade',
  jump_jets: 'Jump Jets',
  ecm: 'ECM',
  ammo_bin: 'Ammo Bin',
}

const MAGNITUDE_HINTS = {
  heat_sink: 'Heat dissipated per turn on top of the chassis baseline (1–20).',
  targeting_computer: 'Percentage points added to the attacker\'s hit chance (1–30). Final hit chance is still capped at 95%.',
  armor_upgrade: 'Extra maximum armour points added to the chassis base (1–200). Raises both starting armour and the auto-repair ceiling.',
  jump_jets: 'Extra movement hops added to chassis base speed (1–5). Used by orders and AI.',
  ecm: 'Percentage points of cover added against incoming attacks on this mech (1–50). Stacks with sector cover.',
  ammo_bin: 'Extra rounds added to the mech\'s shared ammo pool (1–200). Purely additive — weapons with ammo_capacity > 0 draw from the combined pool.',
}

const HEAT_COST_HINTS = {
  heat_sink: 'Heat added each turn while the mech is powered (0–20). Typically 0; positive values model inefficient or experimental heat sinks.',
  targeting_computer: 'Heat added on any turn this mech declares an attack (0–20). Does not apply on quiet turns.',
  armor_upgrade: 'Heat added each turn while powered (0–20). Typically 0; positive values model reactive/powered armour.',
  jump_jets: 'Heat added on any turn this mech moves more hops than the chassis base speed (0–20). Normal-speed moves are free.',
  ecm: 'Heat added each turn while powered (0–20). ECM is always-on when the mech is not refitting.',
  ammo_bin: 'Heat added on any turn this mech fires an ammo-consuming weapon (0–20). Typically 0 for passive bins.',
}

const showModal = ref(false)
const modalMode = ref('create')
const defaultForm = () => ({ name: '', description: '', effect_kind: 'heat_sink', mount_size: 'small', magnitude: 1, heat_cost: 0 })
const modalForm = ref(defaultForm())
const modalError = ref('')
const showDeleteModal = ref(false)
const toDelete = ref(null)

const magnitudeMax = computed(() => MAGNITUDE_MAX[modalForm.value.effect_kind] || 1)
const magnitudeHint = computed(() => MAGNITUDE_HINTS[modalForm.value.effect_kind] || '')
const heatCostHint = computed(() => HEAT_COST_HINTS[modalForm.value.effect_kind] || '')

function formatEffectKind(kind) {
  return EFFECT_KIND_LABELS[kind] || kind
}

function onEffectKindChange() {
  const max = magnitudeMax.value
  if (modalForm.value.magnitude > max) modalForm.value.magnitude = max
  if (modalForm.value.magnitude < 1) modalForm.value.magnitude = 1
}

watch(() => selectedGame.value, (g) => { if (g) store.fetchMechaGameEquipment(g.id) }, { immediate: true })

function openCreate() {
  modalMode.value = 'create'
  modalForm.value = defaultForm()
  modalError.value = ''
  showModal.value = true
}

function openEdit(row) {
  modalMode.value = 'edit'
  modalForm.value = { ...defaultForm(), ...row }
  modalError.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  modalError.value = ''
}

async function handleSubmit(formData) {
  modalError.value = ''
  const allowed = ['name', 'description', 'effect_kind', 'mount_size', 'magnitude', 'heat_cost']
  const data = Object.fromEntries(allowed.map(k => [k, formData[k]]))
  try {
    if (modalMode.value === 'create') {
      await store.createMechaGameEquipment(data)
    } else {
      await store.updateMechaGameEquipment(modalForm.value.id, data)
    }
    closeModal()
  } catch (e) {
    modalError.value = e.message || 'Failed to save.'
  }
}

async function handleDelete() {
  try {
    await store.deleteMechaGameEquipment(toDelete.value.id)
    showDeleteModal.value = false
    toDelete.value = null
  } catch (e) {
    console.error('Failed to delete equipment:', e)
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
