<template>
  <div>
    <div v-if="!selectedGame">
      <p>Please select or create a game to manage location objects.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader
        title="Location Objects"
        actionText="Create Location Object"
        :showIcon="false"
        titleLevel="h2"
        @action="openCreate"
      />
      <ResourceTable
        :columns="columns"
        :rows="enhancedObjects"
        :loading="locationObjectsStore.loading"
        :error="locationObjectsStore.error"
        data-testid="location-objects-table"
      >
        <template #cell-name="{ row }">
          <a href="#" class="edit-link" @click.prevent="openEdit(row)">{{ row.name }}</a>
        </template>
        <template #cell-initial_state_name="{ row }">
          <span>{{ row.initial_state_name || '—' }}</span>
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>

      <!-- Object create / edit modal -->
      <ResourceModalForm
        :visible="showModal"
        :mode="modalMode"
        title="Location Object"
        :fields="objectFields"
        :modelValue="modalForm"
        :error="modalError"
        :options="objectFieldOptions"
        data-testid="location-object-form"
        @submit="handleSubmit"
        @cancel="closeModal"
      />

      <!-- States management panel (shown when editing an existing object) -->
      <div v-if="showStatesPanel" class="states-panel">
        <div class="states-panel-header">
          <h3>States for "{{ statesPanelObject?.name }}"</h3>
          <button class="btn-secondary btn-sm" @click="openCreateState">Add State</button>
          <button
            v-if="locationObjectStatesStore.states.length >= 2"
            class="btn-secondary btn-sm"
            @click="openFlowChart"
          >State Flow</button>
          <button class="btn-text btn-sm" @click="closeStatesPanel">Close</button>
        </div>
        <div v-if="locationObjectStatesStore.loading" class="loading-text">Loading states…</div>
        <div v-else-if="locationObjectStatesStore.error" class="error-text">{{ locationObjectStatesStore.error }}</div>
        <table v-else class="states-table" data-testid="states-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Description</th>
              <th>Order</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="state in locationObjectStatesStore.states" :key="state.id">
              <td>{{ state.name }}</td>
              <td>{{ state.description }}</td>
              <td>{{ state.sort_order }}</td>
              <td>
                <button class="btn-text btn-sm" @click="openEditState(state)">Edit</button>
                <button class="btn-text btn-sm danger" @click="confirmDeleteState(state)">Delete</button>
              </td>
            </tr>
            <tr v-if="locationObjectStatesStore.states.length === 0">
              <td colspan="4" class="empty-row">No states defined yet.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- State create / edit modal -->
      <ResourceModalForm
        :visible="showStateModal"
        :mode="stateModalMode"
        title="Object State"
        :fields="stateFields"
        :modelValue="stateModalForm"
        :error="stateModalError"
        data-testid="location-object-state-form"
        @submit="handleStateSubmit"
        @cancel="closeStateModal"
      />

      <!-- Object delete confirmation -->
      <ConfirmationModal
        :visible="showDeleteConfirm"
        title="Delete Location Object"
        message="Are you sure you want to delete this location object?"
        @confirm="confirmDelete"
        @cancel="closeDeleteConfirm"
      />

      <!-- State delete confirmation -->
      <ConfirmationModal
        :visible="showDeleteStateConfirm"
        title="Delete Object State"
        message="Are you sure you want to delete this state? Effects that reference it will lose their state reference."
        @confirm="confirmDeleteStateAction"
        @cancel="closeDeleteStateConfirm"
      />

      <!-- State flow chart modal -->
      <ObjectStateFlowModal
        :visible="showFlowChart"
        :objectName="statesPanelObject?.name || ''"
        :states="locationObjectStatesStore.states"
        :effects="flowChartEffects"
        :initialStateId="statesPanelObject?.initial_adventure_game_location_object_state_id || null"
        @close="closeFlowChart"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { useLocationsStore } from '../../../stores/locations';
import { useLocationObjectsStore } from '../../../stores/locationObjects';
import { useLocationObjectStatesStore } from '../../../stores/locationObjectStates';
import { useLocationObjectEffectsStore } from '../../../stores/locationObjectEffects';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';
import ObjectStateFlowModal from '../../../components/ObjectStateFlowModal.vue';
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import TableActions from '../../../components/TableActions.vue';

const locationsStore = useLocationsStore();
const locationObjectsStore = useLocationObjectsStore();
const locationObjectStatesStore = useLocationObjectStatesStore();
const locationObjectEffectsStore = useLocationObjectEffectsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

// ── Enhanced rows ────────────────────────────────────────────────────────────

/**
 * Build a map of all states across all loaded objects so we can resolve
 * initial_adventure_game_location_object_state_id → state name.
 */
const allStatesById = computed(() => {
  const map = {};
  for (const state of locationObjectStatesStore.states) {
    map[state.id] = state;
  }
  return map;
});

const enhancedObjects = computed(() =>
  locationObjectsStore.locationObjects.map((obj) => {
    const location = locationsStore.locations.find((l) => l.id === obj.adventure_game_location_id);
    const initialState = obj.initial_adventure_game_location_object_state_id
      ? allStatesById.value[obj.initial_adventure_game_location_object_state_id]
      : null;
    return {
      ...obj,
      location_name: location?.name || 'Unknown Location',
      initial_state_name: initialState?.name || '',
    };
  })
);

// ── Table columns ────────────────────────────────────────────────────────────

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'location_name', label: 'Location' },
  { key: 'initial_state_name', label: 'Initial State' },
  { key: 'is_hidden', label: 'Hidden' },
  { key: 'created_at', label: 'Created' },
];

// ── Object form ──────────────────────────────────────────────────────────────

const objectFields = computed(() => [
  { key: 'adventure_game_location_id', label: 'Location', type: 'select', required: true, placeholder: 'Select a location…' },
  { key: 'name', label: 'Name', type: 'text', required: true, placeholder: 'Object name' },
  { key: 'description', label: 'Description', type: 'textarea', required: true, placeholder: 'Object description' },
  ...(modalMode.value === 'edit' && locationObjectStatesStore.states.length > 0
    ? [{ key: 'initial_adventure_game_location_object_state_id', label: 'Initial State', type: 'select', placeholder: '— none —' }]
    : []),
  { key: 'is_hidden', label: 'Hidden', type: 'checkbox' },
]);

const objectFieldOptions = computed(() => ({
  adventure_game_location_id: locationsStore.locations.map((l) => ({ value: l.id, label: l.name })),
  initial_adventure_game_location_object_state_id: [
    { value: '', label: '— none —' },
    ...locationObjectStatesStore.states.map((s) => ({ value: s.id, label: s.name })),
  ],
}));

// ── Object modal state ───────────────────────────────────────────────────────

const showModal = ref(false);
const modalMode = ref('create');
const modalForm = ref({
  adventure_game_location_id: '',
  name: '',
  description: '',
  initial_adventure_game_location_object_state_id: '',
  is_hidden: false,
});
const modalError = ref('');
const showDeleteConfirm = ref(false);
const deleteTarget = ref(null);

// ── States panel ─────────────────────────────────────────────────────────────

const showStatesPanel = ref(false);
const statesPanelObject = ref(null);

function openStatesPanel(obj) {
  statesPanelObject.value = obj;
  showStatesPanel.value = true;
  if (selectedGame.value) {
    locationObjectStatesStore.fetchStates(selectedGame.value.id, obj.id);
  }
}

function closeStatesPanel() {
  showStatesPanel.value = false;
  statesPanelObject.value = null;
  locationObjectStatesStore.clearStates();
}

// ── State flow chart ──────────────────────────────────────────────────────────

const showFlowChart = ref(false);

const flowChartEffects = computed(() => {
  if (!statesPanelObject.value) return [];
  return locationObjectEffectsStore.locationObjectEffects.filter(
    (e) => e.adventure_game_location_object_id === statesPanelObject.value.id
  );
});

function openFlowChart() {
  showFlowChart.value = true;
}

function closeFlowChart() {
  showFlowChart.value = false;
}

// ── State form ───────────────────────────────────────────────────────────────

const stateFields = [
  { key: 'name', label: 'Name', type: 'text', required: true, placeholder: 'e.g. intact' },
  { key: 'description', label: 'Description', type: 'textarea', placeholder: 'Describe this state' },
  { key: 'sort_order', label: 'Sort Order', type: 'number', placeholder: '0' },
];

const showStateModal = ref(false);
const stateModalMode = ref('create');
const stateModalForm = ref({ name: '', description: '', sort_order: 0 });
const stateModalError = ref('');
const showDeleteStateConfirm = ref(false);
const deleteStateTarget = ref(null);

function openCreateState() {
  stateModalMode.value = 'create';
  stateModalForm.value = { name: '', description: '', sort_order: 0 };
  stateModalError.value = '';
  showStateModal.value = true;
}

function openEditState(state) {
  stateModalMode.value = 'edit';
  stateModalForm.value = { ...state };
  stateModalError.value = '';
  showStateModal.value = true;
}

function closeStateModal() {
  showStateModal.value = false;
  stateModalError.value = '';
}

async function handleStateSubmit(form) {
  stateModalError.value = '';
  try {
    if (stateModalMode.value === 'create') {
      await locationObjectStatesStore.createState(form);
    } else {
      await locationObjectStatesStore.updateState(stateModalForm.value.id, form);
    }
    closeStateModal();
  } catch (err) {
    stateModalError.value = err.message || 'Failed to save state.';
  }
}

function confirmDeleteState(state) {
  deleteStateTarget.value = state;
  showDeleteStateConfirm.value = true;
}

function closeDeleteStateConfirm() {
  showDeleteStateConfirm.value = false;
  deleteStateTarget.value = null;
}

async function confirmDeleteStateAction() {
  if (!deleteStateTarget.value) return;
  try {
    await locationObjectStatesStore.deleteState(deleteStateTarget.value.id);
    closeDeleteStateConfirm();
  } catch (err) {
    console.error('Failed to delete state:', err);
  }
}

// ── Watchers ─────────────────────────────────────────────────────────────────

watch(
  () => selectedGame.value,
  (newGame) => {
    if (newGame) {
      locationsStore.fetchLocations(newGame.id);
      locationObjectsStore.fetchLocationObjects(newGame.id);
      locationObjectEffectsStore.fetchLocationObjectEffects(newGame.id);
    }
  },
  { immediate: true }
);

// ── Object modal actions ─────────────────────────────────────────────────────

function openCreate() {
  modalMode.value = 'create';
  modalForm.value = {
    adventure_game_location_id: '',
    name: '',
    description: '',
    initial_adventure_game_location_object_state_id: '',
    is_hidden: false,
  };
  modalError.value = '';
  locationObjectStatesStore.clearStates();
  showModal.value = true;
}

function openEdit(row) {
  modalMode.value = 'edit';
  modalForm.value = {
    adventure_game_location_id: row.adventure_game_location_id,
    name: row.name,
    description: row.description,
    initial_adventure_game_location_object_state_id: row.initial_adventure_game_location_object_state_id || '',
    is_hidden: row.is_hidden,
    id: row.id,
  };
  modalError.value = '';
  // Load states for this object so the initial-state dropdown is populated.
  if (selectedGame.value) {
    locationObjectStatesStore.fetchStates(selectedGame.value.id, row.id);
  }
  showModal.value = true;
}

function closeModal() {
  showModal.value = false;
  modalError.value = '';
}

async function handleSubmit(form) {
  modalError.value = '';
  const payload = { ...form };
  if (!payload.initial_adventure_game_location_object_state_id) {
    delete payload.initial_adventure_game_location_object_state_id;
  }
  try {
    if (modalMode.value === 'create') {
      await locationObjectsStore.createLocationObject(payload);
    } else {
      await locationObjectsStore.updateLocationObject(modalForm.value.id, payload);
    }
    closeModal();
  } catch (err) {
    modalError.value = err.message || 'Failed to save.';
  }
}

function confirmDeleteOpen(row) {
  deleteTarget.value = row;
  showDeleteConfirm.value = true;
}

function closeDeleteConfirm() {
  showDeleteConfirm.value = false;
  deleteTarget.value = null;
}

async function confirmDelete() {
  if (!deleteTarget.value) return;
  try {
    await locationObjectsStore.deleteLocationObject(deleteTarget.value.id);
    closeDeleteConfirm();
  } catch (err) {
    console.error('Failed to delete location object:', err);
  }
}

function getActions(row) {
  return [
    { key: 'states', label: 'States', handler: () => openStatesPanel(row) },
    { key: 'edit', label: 'Edit', handler: () => openEdit(row) },
    { key: 'delete', label: 'Delete', danger: true, handler: () => confirmDeleteOpen(row) },
  ];
}
</script>

<style scoped>
.edit-link {
  color: var(--color-primary);
  text-decoration: none;
}

.edit-link:hover {
  text-decoration: underline;
}

.states-panel {
  margin-top: 1.5rem;
  padding: 1rem;
  border: 1px solid var(--color-border, #e2e8f0);
  border-radius: 6px;
  background: var(--color-surface-alt, #f8fafc);
}

.states-panel-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.states-panel-header h3 {
  margin: 0;
  flex: 1;
  font-size: 1rem;
}

.states-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

.states-table th,
.states-table td {
  text-align: left;
  padding: 0.4rem 0.6rem;
  border-bottom: 1px solid var(--color-border, #e2e8f0);
}

.states-table th {
  font-weight: 600;
  color: var(--color-text-secondary, #64748b);
}

.empty-row {
  color: var(--color-text-secondary, #64748b);
  font-style: italic;
}

.btn-sm {
  padding: 0.25rem 0.6rem;
  font-size: 0.8rem;
}

.btn-text {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--color-primary);
  padding: 0.2rem 0.4rem;
}

.btn-text.danger {
  color: var(--color-danger, #ef4444);
}

.btn-secondary {
  background: var(--color-surface, #fff);
  border: 1px solid var(--color-border, #e2e8f0);
  border-radius: 4px;
  cursor: pointer;
  color: var(--color-text, #1e293b);
}

.loading-text,
.error-text {
  padding: 0.5rem;
  font-size: 0.9rem;
}

.error-text {
  color: var(--color-danger, #ef4444);
}
</style>
