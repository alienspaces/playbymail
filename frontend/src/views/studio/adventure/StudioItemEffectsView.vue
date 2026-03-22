<template>
  <div>
    <div v-if="!selectedGame">
      <p>Please select or create a game to manage item effects.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader
        title="Item Effects"
        actionText="Create Item Effect"
        :showIcon="false"
        titleLevel="h2"
        @action="openCreate"
      />
      <ResourceTable
        :columns="columns"
        :rows="enhancedEffects"
        :loading="itemEffectsStore.loading"
        :error="itemEffectsStore.error"
        data-testid="item-effects-table"
      >
        <template #cell-action_type="{ row }">
          <a href="#" class="edit-link" @click.prevent="openEdit(row)">{{ row.action_type }}</a>
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>
      <TablePagination :pageNumber="itemEffectsStore.pageNumber" :hasMore="itemEffectsStore.hasMore"
        @page-change="(p) => itemEffectsStore.fetchItemEffects(selectedGame.id, p)" />

      <ResourceModalForm
        :visible="showModal"
        :mode="modalMode"
        title="Item Effect"
        :fields="effectFields"
        :modelValue="modalForm"
        :error="modalError"
        :options="effectFieldOptions"
        data-testid="item-effect-form"
        @submit="handleSubmit"
        @cancel="closeModal"
      />

      <ConfirmationModal
        :visible="showDeleteConfirm"
        title="Delete Item Effect"
        message="Are you sure you want to delete this item effect?"
        @confirm="confirmDelete"
        @cancel="closeDeleteConfirm"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { useItemEffectsStore } from '../../../stores/itemEffects';
import { useItemsStore } from '../../../stores/items';
import { useLocationsStore } from '../../../stores/locations';
import { useLocationLinksStore } from '../../../stores/locationLinks';
import { useCreaturesStore } from '../../../stores/creatures';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import TableActions from '../../../components/TableActions.vue';
import TablePagination from '../../../components/TablePagination.vue';

const itemEffectsStore = useItemEffectsStore();
const itemsStore = useItemsStore();
const locationsStore = useLocationsStore();
const locationLinksStore = useLocationLinksStore();
const creaturesStore = useCreaturesStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

// ── Enhanced rows ─────────────────────────────────────────────────────────────

const enhancedEffects = computed(() =>
  itemEffectsStore.itemEffects.map((effect) => {
    const item = itemsStore.items.find((i) => i.id === effect.adventure_game_item_id);
    const requiredItem = effect.required_adventure_game_item_id
      ? itemsStore.items.find((i) => i.id === effect.required_adventure_game_item_id)
      : null;
    const requiredLocation = effect.required_adventure_game_location_id
      ? locationsStore.locations.find((l) => l.id === effect.required_adventure_game_location_id)
      : null;
    return {
      ...effect,
      item_name: item?.name || 'Unknown Item',
      required_item_name: requiredItem?.name || '',
      required_location_name: requiredLocation?.name || '',
    };
  })
);

// ── Table columns ─────────────────────────────────────────────────────────────

const columns = [
  { key: 'item_name', label: 'Item' },
  { key: 'action_type', label: 'Action' },
  { key: 'effect_type', label: 'Effect' },
  { key: 'required_item_name', label: 'Required Item' },
  { key: 'required_location_name', label: 'Required Location' },
  { key: 'result_description', label: 'Description' },
  { key: 'is_repeatable', label: 'Repeatable' },
];

// ── Field definitions ─────────────────────────────────────────────────────────

const ACTION_TYPES = ['use', 'equip', 'unequip', 'inspect', 'drop', 'pickup'];

const EFFECT_TYPES = [
  'info', 'damage_target', 'damage_wielder', 'heal_target', 'heal_wielder',
  'teleport', 'open_link', 'close_link', 'give_item', 'remove_item',
  'summon_creature', 'nothing', 'weapon_damage', 'armor_defense',
];

// Maps each effect_type to which result fields are shown and which are required.
const EFFECT_TYPE_FIELD_RULES = {
  damage_target:   { show: ['result_value_min', 'result_value_max'], required: ['result_value_min', 'result_value_max'] },
  damage_wielder:  { show: ['result_value_min', 'result_value_max'], required: ['result_value_min', 'result_value_max'] },
  heal_target:     { show: ['result_value_min', 'result_value_max'], required: ['result_value_min', 'result_value_max'] },
  heal_wielder:    { show: ['result_value_min', 'result_value_max'], required: ['result_value_min', 'result_value_max'] },
  weapon_damage:   { show: ['result_value_min', 'result_value_max'], required: ['result_value_min', 'result_value_max'] },
  armor_defense:   { show: ['result_value_min', 'result_value_max'], required: ['result_value_min', 'result_value_max'] },
  teleport:        { show: ['result_adventure_game_location_id'], required: ['result_adventure_game_location_id'] },
  open_link:       { show: ['result_adventure_game_location_link_id'], required: ['result_adventure_game_location_link_id'] },
  close_link:      { show: ['result_adventure_game_location_link_id'], required: ['result_adventure_game_location_link_id'] },
  give_item:       { show: ['result_adventure_game_item_id'], required: ['result_adventure_game_item_id'] },
  remove_item:     { show: ['result_adventure_game_item_id'], required: ['result_adventure_game_item_id'] },
  summon_creature: { show: ['result_adventure_game_creature_id'], required: ['result_adventure_game_creature_id'] },
  info:            { show: [], required: [] },
  nothing:         { show: [], required: [] },
};

// All possible result fields with their base definitions (never required by default).
const RESULT_FIELDS = [
  { key: 'result_adventure_game_item_id', label: 'Result Item', type: 'select', placeholder: 'Item to give/remove…' },
  { key: 'result_adventure_game_location_link_id', label: 'Result Link', type: 'select', placeholder: 'Link to open/close…' },
  { key: 'result_adventure_game_creature_id', label: 'Result Creature', type: 'select', placeholder: 'Creature to summon…' },
  { key: 'result_adventure_game_location_id', label: 'Result Location', type: 'select', placeholder: 'Location to teleport to…' },
  { key: 'result_value_min', label: 'Min Value', type: 'number', placeholder: 'e.g. 5' },
  { key: 'result_value_max', label: 'Max Value', type: 'number', placeholder: 'e.g. 10' },
];

// Base fields always shown regardless of effect type.
const BASE_FIELDS = [
  { key: 'adventure_game_item_id', label: 'Item', type: 'select', required: true, placeholder: 'Select an item…' },
  { key: 'action_type', label: 'Action Type', type: 'select', required: true, placeholder: 'Select action type…' },
  { key: 'effect_type', label: 'Effect Type', type: 'select', required: true, placeholder: 'Select effect type…' },
  { key: 'result_description', label: 'Result Description', type: 'textarea', required: true, placeholder: 'What the player sees' },
  { key: 'required_adventure_game_item_id', label: 'Required Item', type: 'select', placeholder: 'Optional — item that must be held…' },
  { key: 'required_adventure_game_location_id', label: 'Required Location', type: 'select', placeholder: 'Optional — location where effect triggers…' },
];

const TAIL_FIELDS = [
  { key: 'is_repeatable', label: 'Repeatable', type: 'checkbox' },
];

const effectFields = computed(() => {
  const effectType = modalForm.value.effect_type;
  const rules = EFFECT_TYPE_FIELD_RULES[effectType] || { show: [], required: [] };
  const showSet = new Set(rules.show);
  const requiredSet = new Set(rules.required);

  const resultFields = RESULT_FIELDS
    .filter((f) => showSet.has(f.key))
    .map((f) => ({ ...f, required: requiredSet.has(f.key) }));

  return [...BASE_FIELDS, ...resultFields, ...TAIL_FIELDS];
});

// ── Field options ─────────────────────────────────────────────────────────────

const effectFieldOptions = computed(() => ({
  adventure_game_item_id: itemsStore.items.map((i) => ({ value: i.id, label: i.name })),
  action_type: ACTION_TYPES.map((t) => ({ value: t, label: t })),
  effect_type: EFFECT_TYPES.map((t) => ({ value: t, label: t })),
  required_adventure_game_item_id: [
    { value: '', label: '— none —' },
    ...itemsStore.items.map((i) => ({ value: i.id, label: i.name })),
  ],
  required_adventure_game_location_id: [
    { value: '', label: '— none —' },
    ...locationsStore.locations.map((l) => ({ value: l.id, label: l.name })),
  ],
  result_adventure_game_item_id: [
    { value: '', label: '— none —' },
    ...itemsStore.items.map((i) => ({ value: i.id, label: i.name })),
  ],
  result_adventure_game_location_link_id: [
    { value: '', label: '— none —' },
    ...locationLinksStore.locationLinks.map((l) => ({ value: l.id, label: l.name })),
  ],
  result_adventure_game_creature_id: [
    { value: '', label: '— none —' },
    ...creaturesStore.creatures.map((c) => ({ value: c.id, label: c.name })),
  ],
  result_adventure_game_location_id: [
    { value: '', label: '— none —' },
    ...locationsStore.locations.map((l) => ({ value: l.id, label: l.name })),
  ],
}));

// ── Modal state ───────────────────────────────────────────────────────────────

const showModal = ref(false);
const modalMode = ref('create');

const defaultForm = () => ({
  adventure_game_item_id: '',
  action_type: '',
  effect_type: '',
  result_description: '',
  required_adventure_game_item_id: '',
  required_adventure_game_location_id: '',
  result_adventure_game_item_id: '',
  result_adventure_game_location_link_id: '',
  result_adventure_game_creature_id: '',
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
      itemEffectsStore.fetchItemEffects(newGame.id);
      itemsStore.fetchItems(newGame.id);
      locationsStore.fetchLocations(newGame.id);
      locationLinksStore.fetchLocationLinks(newGame.id);
      creaturesStore.fetchCreatures(newGame.id);
    }
  },
  { immediate: true }
);

// ── Modal actions ─────────────────────────────────────────────────────────────

function openCreate() {
  modalMode.value = 'create';
  modalForm.value = defaultForm();
  modalError.value = '';
  showModal.value = true;
}

function openEdit(row) {
  modalMode.value = 'edit';
  modalForm.value = {
    ...defaultForm(),
    ...row,
    required_adventure_game_item_id: row.required_adventure_game_item_id || '',
    required_adventure_game_location_id: row.required_adventure_game_location_id || '',
    result_adventure_game_item_id: row.result_adventure_game_item_id || '',
    result_adventure_game_location_link_id: row.result_adventure_game_location_link_id || '',
    result_adventure_game_creature_id: row.result_adventure_game_creature_id || '',
    result_adventure_game_location_id: row.result_adventure_game_location_id || '',
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
    'required_adventure_game_item_id',
    'required_adventure_game_location_id',
    'result_adventure_game_item_id',
    'result_adventure_game_location_link_id',
    'result_adventure_game_creature_id',
    'result_adventure_game_location_id',
  ];
  optionalFields.forEach((f) => {
    if (payload[f] === '' || payload[f] === null) delete payload[f];
  });
  if (!payload.result_value_min && payload.result_value_min !== 0) delete payload.result_value_min;
  if (!payload.result_value_max && payload.result_value_max !== 0) delete payload.result_value_max;

  try {
    if (modalMode.value === 'create') {
      await itemEffectsStore.createItemEffect(payload);
    } else {
      await itemEffectsStore.updateItemEffect(modalForm.value.id, payload);
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
    await itemEffectsStore.deleteItemEffect(deleteTarget.value.id);
    closeDeleteConfirm();
  } catch (err) {
    console.error('Failed to delete item effect:', err);
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
