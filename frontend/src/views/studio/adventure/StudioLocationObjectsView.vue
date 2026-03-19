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
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>

      <ResourceModalForm
        :visible="showModal"
        :mode="modalMode"
        title="Location Object"
        :fields="fields"
        :modelValue="modalForm"
        :error="modalError"
        :options="fieldOptions"
        data-testid="location-object-form"
        @submit="handleSubmit"
        @cancel="closeModal"
      />

      <ConfirmationModal
        :visible="showDeleteConfirm"
        title="Delete Location Object"
        message="Are you sure you want to delete this location object?"
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
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import TableActions from '../../../components/TableActions.vue';

const locationsStore = useLocationsStore();
const locationObjectsStore = useLocationObjectsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const enhancedObjects = computed(() =>
  locationObjectsStore.locationObjects.map((obj) => {
    const location = locationsStore.locations.find((l) => l.id === obj.adventure_game_location_id);
    return {
      ...obj,
      location_name: location?.name || 'Unknown Location',
    };
  })
);

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'location_name', label: 'Location' },
  { key: 'initial_state', label: 'Initial State' },
  { key: 'is_hidden', label: 'Hidden' },
  { key: 'created_at', label: 'Created' },
];

const fields = [
  { key: 'adventure_game_location_id', label: 'Location', type: 'select', required: true, placeholder: 'Select a location...' },
  { key: 'name', label: 'Name', type: 'text', required: true, placeholder: 'Object name' },
  { key: 'description', label: 'Description', type: 'textarea', required: true, placeholder: 'Object description' },
  { key: 'initial_state', label: 'Initial State', type: 'text', placeholder: 'e.g. intact' },
  { key: 'is_hidden', label: 'Hidden', type: 'checkbox' },
];

const fieldOptions = computed(() => ({
  adventure_game_location_id: locationsStore.locations.map((l) => ({ value: l.id, label: l.name })),
}));

const showModal = ref(false);
const modalMode = ref('create');
const modalForm = ref({ adventure_game_location_id: '', name: '', description: '', initial_state: 'intact', is_hidden: false });
const modalError = ref('');
const showDeleteConfirm = ref(false);
const deleteTarget = ref(null);

watch(
  () => selectedGame.value,
  (newGame) => {
    if (newGame) {
      locationsStore.fetchLocations(newGame.id);
      locationObjectsStore.fetchLocationObjects(newGame.id);
    }
  },
  { immediate: true }
);

function openCreate() {
  modalMode.value = 'create';
  modalForm.value = { adventure_game_location_id: '', name: '', description: '', initial_state: 'intact', is_hidden: false };
  modalError.value = '';
  showModal.value = true;
}

function openEdit(row) {
  modalMode.value = 'edit';
  modalForm.value = { ...row };
  modalError.value = '';
  showModal.value = true;
}

function closeModal() {
  showModal.value = false;
  modalError.value = '';
}

async function handleSubmit(form) {
  modalError.value = '';
  try {
    if (modalMode.value === 'create') {
      await locationObjectsStore.createLocationObject(form);
    } else {
      await locationObjectsStore.updateLocationObject(modalForm.value.id, form);
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
