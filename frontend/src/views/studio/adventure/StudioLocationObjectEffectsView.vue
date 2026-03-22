<template>
  <div>
    <div v-if="!selectedGame">
      <p>Please select or create a game to manage location object effects.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader
        title="Object Effects"
        actionText="Create Object Effect"
        :showIcon="false"
        titleLevel="h2"
        @action="openCreate"
      />
      <ResourceTable
        :columns="columns"
        :rows="enhancedEffects"
        :loading="locationObjectEffectsStore.loading"
        :error="locationObjectEffectsStore.error"
        data-testid="location-object-effects-table"
      >
        <template #cell-action_type="{ row }">
          <a href="#" class="edit-link" @click.prevent="openEdit(row)">{{ row.action_type }}</a>
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>
      <TablePagination :pageNumber="locationObjectEffectsStore.pageNumber" :hasMore="locationObjectEffectsStore.hasMore"
        @page-change="(p) => locationObjectEffectsStore.fetchLocationObjectEffects(selectedGame.id, p)" />

      <ResourceModalForm
        :visible="showModal"
        :mode="modalMode"
        title="Object Effect"
        :fields="effectFields"
        :modelValue="modalForm"
        :error="modalError"
        :options="effectFieldOptions"
        data-testid="location-object-effect-form"
        @submit="handleSubmit"
        @cancel="closeModal"
      />

      <ConfirmationModal
        :visible="showDeleteConfirm"
        title="Delete Object Effect"
        message="Are you sure you want to delete this object effect?"
        @confirm="confirmDelete"
        @cancel="closeDeleteConfirm"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { useLocationsStore } from '../../../stores/locations';
import { useLocationObjectsStore } from '../../../stores/locationObjects';
import { useLocationObjectEffectsStore } from '../../../stores/locationObjectEffects';
import { useLocationObjectStatesStore } from '../../../stores/locationObjectStates';
import { useItemsStore } from '../../../stores/items';
import { useCreaturesStore } from '../../../stores/creatures';
import { useLocationLinksStore } from '../../../stores/locationLinks';
import { useGamesStore } from '../../../stores/games';
import { fetchLocationObjectStates } from '../../../api/locationObjectStates';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import TableActions from '../../../components/TableActions.vue';
import TablePagination from '../../../components/TablePagination.vue';

const locationsStore = useLocationsStore();
const locationObjectsStore = useLocationObjectsStore();
const locationObjectEffectsStore = useLocationObjectEffectsStore();
const locationObjectStatesStore = useLocationObjectStatesStore();
const itemsStore = useItemsStore();
const creaturesStore = useCreaturesStore();
const locationLinksStore = useLocationLinksStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

// ── Per-effect states cache (keyed by objectId) ──────────────────────────────

/** @type {import('vue').Ref<Record<string, import('../../../types').GameLocationObjectState[]>>} */
const statesByObjectId = ref({});

async function loadStatesForObject(objectId) {
  if (!selectedGame.value || !objectId || statesByObjectId.value[objectId]) return;
  try {
    const states = await fetchLocationObjectStates(selectedGame.value.id, objectId);
    statesByObjectId.value = { ...statesByObjectId.value, [objectId]: states };
  } catch {
    // ignore fetch errors silently
  }
}

// ── Enhanced rows ─────────────────────────────────────────────────────────────

/**
 * Flat map of all known states from the per-effect cache.
 */
const allStatesById = computed(() => {
  const map = {};
  for (const states of Object.values(statesByObjectId.value)) {
    for (const s of states) {
      map[s.id] = s;
    }
  }
  // Also include states loaded via the shared store (e.g. when editing).
  for (const s of locationObjectStatesStore.states) {
    map[s.id] = s;
  }
  return map;
});

const enhancedEffects = computed(() =>
  locationObjectEffectsStore.locationObjectEffects.map((effect) => {
    const obj = locationObjectsStore.locationObjects.find((o) => o.id === effect.adventure_game_location_object_id);
    const reqState = effect.required_adventure_game_location_object_state_id
      ? allStatesById.value[effect.required_adventure_game_location_object_state_id]
      : null;
    return {
      ...effect,
      object_name: obj?.name || 'Unknown Object',
      required_state_name: reqState?.name || '',
    };
  })
);

// ── Table columns ─────────────────────────────────────────────────────────────

const columns = [
  { key: 'object_name', label: 'Object' },
  { key: 'action_type', label: 'Action' },
  { key: 'effect_type', label: 'Effect' },
  { key: 'required_state_name', label: 'Required State' },
  { key: 'result_description', label: 'Description' },
  { key: 'is_repeatable', label: 'Repeatable' },
];

// ── Field definitions ─────────────────────────────────────────────────────────

const ACTION_TYPES = [
  'inspect', 'touch', 'open', 'close', 'lock', 'unlock', 'search',
  'break', 'push', 'pull', 'move', 'burn', 'read', 'take', 'listen',
  'insert', 'pour', 'disarm', 'climb', 'use',
];

const EFFECT_TYPES = [
  'info', 'change_state', 'change_object_state', 'give_item', 'remove_item',
  'open_link', 'close_link', 'reveal_object', 'hide_object', 'damage',
  'heal', 'summon_creature', 'teleport', 'nothing', 'remove_object', 'place_item',
];

// Maps each effect_type to which result fields are shown and which are required.
const OBJECT_EFFECT_TYPE_FIELD_RULES = {
  info:                 { show: [], required: [] },
  nothing:              { show: [], required: [] },
  change_state:         { show: ['result_adventure_game_location_object_state_id'], required: ['result_adventure_game_location_object_state_id'] },
  change_object_state:  { show: ['result_adventure_game_location_object_id', 'result_adventure_game_location_object_state_id'], required: ['result_adventure_game_location_object_id', 'result_adventure_game_location_object_state_id'] },
  reveal_object:        { show: ['result_adventure_game_location_object_id'], required: ['result_adventure_game_location_object_id'] },
  hide_object:          { show: ['result_adventure_game_location_object_id'], required: ['result_adventure_game_location_object_id'] },
  remove_object:        { show: [], required: [] },
  give_item:            { show: ['result_adventure_game_item_id'], required: ['result_adventure_game_item_id'] },
  remove_item:          { show: ['result_adventure_game_item_id'], required: ['result_adventure_game_item_id'] },
  open_link:            { show: ['result_adventure_game_location_link_id'], required: ['result_adventure_game_location_link_id'] },
  close_link:           { show: ['result_adventure_game_location_link_id'], required: ['result_adventure_game_location_link_id'] },
  damage:               { show: ['result_value_min', 'result_value_max'], required: ['result_value_min', 'result_value_max'] },
  heal:                 { show: ['result_value_min', 'result_value_max'], required: ['result_value_min', 'result_value_max'] },
  summon_creature:      { show: ['result_adventure_game_creature_id'], required: ['result_adventure_game_creature_id'] },
  teleport:             { show: ['result_adventure_game_location_id'], required: ['result_adventure_game_location_id'] },
  place_item:           { show: ['result_adventure_game_item_id', 'result_adventure_game_location_id'], required: ['result_adventure_game_item_id', 'result_adventure_game_location_id'] },
};

// All possible result fields with their base definitions.
const OBJECT_RESULT_FIELDS = [
  { key: 'result_adventure_game_location_object_state_id', label: 'Result State', type: 'select', placeholder: '— no state change —' },
  { key: 'result_adventure_game_item_id', label: 'Result Item', type: 'select', placeholder: 'Item to give/remove…' },
  { key: 'result_adventure_game_location_link_id', label: 'Result Link', type: 'select', placeholder: 'Link to open/close…' },
  { key: 'result_adventure_game_creature_id', label: 'Result Creature', type: 'select', placeholder: 'Creature to summon…' },
  { key: 'result_adventure_game_location_object_id', label: 'Result Object', type: 'select', placeholder: 'Object to reveal/hide/change…' },
  { key: 'result_adventure_game_location_id', label: 'Result Location', type: 'select', placeholder: 'Location to teleport to…' },
  { key: 'result_value_min', label: 'Min Value', type: 'number', placeholder: 'e.g. 5' },
  { key: 'result_value_max', label: 'Max Value', type: 'number', placeholder: 'e.g. 10' },
];

// Base fields always shown regardless of effect type.
const OBJECT_BASE_FIELDS = [
  { key: 'adventure_game_location_object_id', label: 'Object', type: 'select', required: true, placeholder: 'Select an object…' },
  { key: 'action_type', label: 'Action Type', type: 'select', required: true, placeholder: 'Select action type…' },
  { key: 'effect_type', label: 'Effect Type', type: 'select', required: true, placeholder: 'Select effect type…' },
  { key: 'result_description', label: 'Result Description', type: 'textarea', required: true, placeholder: 'What the player sees' },
  { key: 'required_adventure_game_location_object_state_id', label: 'Required State', type: 'select', placeholder: '— any state —' },
  { key: 'required_adventure_game_item_id', label: 'Required Item', type: 'select', placeholder: 'Optional required item…' },
];

const OBJECT_TAIL_FIELDS = [
  { key: 'is_repeatable', label: 'Repeatable', type: 'checkbox' },
];

const effectFields = computed(() => {
  const effectType = modalForm.value.effect_type;
  const rules = OBJECT_EFFECT_TYPE_FIELD_RULES[effectType] || { show: [], required: [] };
  const showSet = new Set(rules.show);
  const requiredSet = new Set(rules.required);

  const resultFields = OBJECT_RESULT_FIELDS
    .filter((f) => showSet.has(f.key))
    .map((f) => ({ ...f, required: requiredSet.has(f.key) }));

  return [...OBJECT_BASE_FIELDS, ...resultFields, ...OBJECT_TAIL_FIELDS];
});

// ── States for the currently selected source and result objects ───────────────

const sourceObjectStates = computed(() => {
  const objectId = modalForm.value.adventure_game_location_object_id;
  return objectId ? (statesByObjectId.value[objectId] || locationObjectStatesStore.states) : [];
});

const resultObjectStates = computed(() => {
  const objectId = modalForm.value.result_adventure_game_location_object_id;
  return objectId ? (statesByObjectId.value[objectId] || []) : sourceObjectStates.value;
});

const stateSelectOptions = (states) => [
  { value: '', label: '— none —' },
  ...states.map((s) => ({ value: s.id, label: s.name })),
];

// ── Field options ─────────────────────────────────────────────────────────────

const effectFieldOptions = computed(() => ({
  adventure_game_location_object_id: locationObjectsStore.locationObjects.map((o) => ({ value: o.id, label: o.name })),
  action_type: ACTION_TYPES.map((t) => ({ value: t, label: t })),
  effect_type: EFFECT_TYPES.map((t) => ({ value: t, label: t })),
  required_adventure_game_location_object_state_id: stateSelectOptions(sourceObjectStates.value),
  required_adventure_game_item_id: [{ value: '', label: '— none —' }, ...itemsStore.items.map((i) => ({ value: i.id, label: i.name }))],
  result_adventure_game_location_object_state_id: stateSelectOptions(resultObjectStates.value),
  result_adventure_game_item_id: [{ value: '', label: '— none —' }, ...itemsStore.items.map((i) => ({ value: i.id, label: i.name }))],
  result_adventure_game_location_link_id: [{ value: '', label: '— none —' }, ...locationLinksStore.locationLinks.map((l) => ({ value: l.id, label: l.name }))],
  result_adventure_game_creature_id: [{ value: '', label: '— none —' }, ...creaturesStore.creatures.map((c) => ({ value: c.id, label: c.name }))],
  result_adventure_game_location_object_id: [{ value: '', label: '— none —' }, ...locationObjectsStore.locationObjects.map((o) => ({ value: o.id, label: o.name }))],
  result_adventure_game_location_id: [{ value: '', label: '— none —' }, ...locationsStore.locations.map((l) => ({ value: l.id, label: l.name }))],
}));

// ── Modal state ───────────────────────────────────────────────────────────────

const showModal = ref(false);
const modalMode = ref('create');

const defaultForm = () => ({
  adventure_game_location_object_id: '',
  action_type: '',
  effect_type: '',
  result_description: '',
  required_adventure_game_location_object_state_id: '',
  required_adventure_game_item_id: '',
  result_adventure_game_location_object_state_id: '',
  result_adventure_game_item_id: '',
  result_adventure_game_location_link_id: '',
  result_adventure_game_creature_id: '',
  result_adventure_game_location_object_id: '',
  result_adventure_game_location_id: '',
  result_value_min: null,
  result_value_max: null,
  is_repeatable: false,
});

const modalForm = ref(defaultForm());
const modalError = ref('');
const showDeleteConfirm = ref(false);
const deleteTarget = ref(null);

// ── Watchers ──────────────────────────────────────────────────────────────────

watch(
  () => selectedGame.value,
  (newGame) => {
    if (newGame) {
      locationsStore.fetchLocations(newGame.id);
      locationObjectsStore.fetchLocationObjects(newGame.id);
      locationObjectEffectsStore.fetchLocationObjectEffects(newGame.id);
      itemsStore.fetchItems(newGame.id);
      creaturesStore.fetchCreatures(newGame.id);
      locationLinksStore.fetchLocationLinks(newGame.id);
    }
  },
  { immediate: true }
);

// When objects finish loading, pre-fetch their states for the table display.
watch(
  () => locationObjectsStore.locationObjects,
  (objs) => {
    for (const obj of objs) {
      loadStatesForObject(obj.id);
    }
  }
);

// When the source object changes in the form, load its states.
watch(
  () => modalForm.value.adventure_game_location_object_id,
  (objectId) => {
    if (objectId && selectedGame.value) {
      loadStatesForObject(objectId);
      locationObjectStatesStore.fetchStates(selectedGame.value.id, objectId);
    } else {
      locationObjectStatesStore.clearStates();
    }
  }
);

// When the result object changes in the form, load its states.
watch(
  () => modalForm.value.result_adventure_game_location_object_id,
  (objectId) => {
    if (objectId && selectedGame.value) {
      loadStatesForObject(objectId);
    }
  }
);

// ── Modal actions ─────────────────────────────────────────────────────────────

function openCreate() {
  modalMode.value = 'create';
  modalForm.value = defaultForm();
  modalError.value = '';
  locationObjectStatesStore.clearStates();
  showModal.value = true;
}

function openEdit(row) {
  modalMode.value = 'edit';
  modalForm.value = {
    ...defaultForm(),
    ...row,
    required_adventure_game_location_object_state_id: row.required_adventure_game_location_object_state_id || '',
    result_adventure_game_location_object_state_id: row.result_adventure_game_location_object_state_id || '',
  };
  modalError.value = '';
  showModal.value = true;
}

function closeModal() {
  showModal.value = false;
  modalError.value = '';
}

async function handleSubmit(form) {
  modalError.value = '';
  const payload = { ...form };
  const optionalFields = [
    'required_adventure_game_location_object_state_id',
    'required_adventure_game_item_id',
    'result_adventure_game_location_object_state_id',
    'result_adventure_game_item_id',
    'result_adventure_game_location_link_id',
    'result_adventure_game_creature_id',
    'result_adventure_game_location_object_id',
    'result_adventure_game_location_id',
  ];
  optionalFields.forEach((f) => {
    if (payload[f] === '' || payload[f] === null) delete payload[f];
  });
  if (!payload.result_value_min && payload.result_value_min !== 0) delete payload.result_value_min;
  if (!payload.result_value_max && payload.result_value_max !== 0) delete payload.result_value_max;

  try {
    if (modalMode.value === 'create') {
      await locationObjectEffectsStore.createLocationObjectEffect(payload);
    } else {
      await locationObjectEffectsStore.updateLocationObjectEffect(modalForm.value.id, payload);
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
    await locationObjectEffectsStore.deleteLocationObjectEffect(deleteTarget.value.id);
    closeDeleteConfirm();
  } catch (err) {
    console.error('Failed to delete object effect:', err);
  }
}

function getActions(row) {
  return [
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
</style>
