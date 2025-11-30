<!--
  StudioLocationsView.vue
  This component follows the same pattern as StudioCreaturesView.vue and StudioItemsView.vue.
-->
<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage locations.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Locations" actionText="Create New Location" :showIcon="false" titleLevel="h2"
        @action="openCreate" />
      <ResourceTable :columns="columns" :rows="formattedLocations" :loading="locationsStore.loading"
        :error="locationsStore.error">
        <template #cell-name="{ row }">
          <a href="#" class="edit-link" @click.prevent="openEdit(row)">{{ row.name }}</a>
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>
    </div>

    <ResourceModalForm :visible="showModal" :mode="modalMode" title="Location" :fields="fields" :modelValue="modalForm"
      :error="modalError" @submit="handleSubmit" @cancel="closeModal" />

    <ConfirmationModal :visible="showDeleteModal" title="Delete Location"
      :message="`Are you sure you want to delete '${locationToDelete?.name}'?`" @confirm="handleDelete"
      @cancel="closeDeleteModal" />
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { storeToRefs } from 'pinia';
import { useLocationsStore } from '../../../stores/locations';
import { useGamesStore } from '../../../stores/games';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import TableActions from '../../../components/TableActions.vue';

const locationsStore = useLocationsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

// Format locations for display with formatted starting location
const formattedLocations = computed(() => {
  return locationsStore.locations.map(location => ({
    ...location,
    is_starting_location: location.is_starting_location ? 'Yes' : 'No'
  }));
});

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'description', label: 'Description' },
  { key: 'is_starting_location', label: 'Starting Location' }
];

const fields = [
  { key: 'name', label: 'Name', required: true, maxlength: 1024 },
  { key: 'description', label: 'Description', required: true, maxlength: 4096, type: 'textarea' },
  { key: 'is_starting_location', label: 'Starting Location', type: 'checkbox', checkboxLabel: 'This is a starting location for new players' }
];

const showModal = ref(false);
const modalMode = ref('create');
const modalForm = ref({ name: '', description: '', is_starting_location: false });
const modalError = ref('');
const showDeleteModal = ref(false);
const locationToDelete = ref(null);

// Watch for game selection changes
watch(
  () => selectedGame.value,
  (newGame) => {
    if (newGame) {
      locationsStore.fetchLocations(newGame.id);
    }
  },
  { immediate: true }
);

function openCreate() {
  modalMode.value = 'create';
  modalForm.value = { name: '', description: '', is_starting_location: false };
  modalError.value = '';
  showModal.value = true;
}

function openEdit(location) {
  modalMode.value = 'edit';
  // Get the original location from the store (not the formatted one)
  const originalLocation = locationsStore.locations.find(l => l.id === location.id);
  modalForm.value = { ...originalLocation };
  modalError.value = '';
  showModal.value = true;
}

function openDelete(location) {
  locationToDelete.value = location;
  showDeleteModal.value = true;
}

function closeModal() {
  showModal.value = false;
  modalForm.value = { name: '', description: '', is_starting_location: false };
  modalError.value = '';
}

function closeDeleteModal() {
  showDeleteModal.value = false;
  locationToDelete.value = null;
}

async function handleSubmit(formData) {
  modalError.value = '';
  try {
    // Only send allowed fields (exclude id, game_id, created_at, etc.)
    const allowedFields = ['name', 'description', 'is_starting_location'];
    const requestData = {};
    for (const field of allowedFields) {
      if (field in formData) {
        requestData[field] = formData[field];
      }
    }

    if (modalMode.value === 'create') {
      await locationsStore.createLocation(requestData);
    } else {
      await locationsStore.updateLocation(modalForm.value.id, requestData);
    }
    closeModal();
  } catch (error) {
    modalError.value = error.message || 'Failed to save.';
  }
}

async function handleDelete() {
  try {
    await locationsStore.deleteLocation(locationToDelete.value.id);
    closeDeleteModal();
  } catch (error) {
    console.error('Failed to delete location:', error);
  }
}

function getActions(row) {
  return [
    {
      key: 'edit',
      label: 'Edit',
      handler: () => openEdit(row)
    },
    {
      key: 'delete',
      label: 'Delete',
      danger: true,
      handler: () => openDelete(row)
    }
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
