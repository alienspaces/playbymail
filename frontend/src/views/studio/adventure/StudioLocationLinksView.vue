<!--
  StudioLocationLinksView.vue
  This component manages location links between locations in a game.
-->
<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage location links.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Location Links" actionText="Create New Location Link" :showIcon="false" titleLevel="h2"
        @action="openCreate" />
      <ResourceTable :columns="columns" :rows="enhancedLocationLinks" :loading="locationLinksStore.loading"
        :error="locationLinksStore.error">
        <template #actions="{ row }">
          <TableActionsMenu :actions="getActions(row)" />
        </template>
      </ResourceTable>

      <!-- Create/Edit Location Link Modal -->
      <ResourceModalForm :visible="showModal" :mode="modalMode" title="Location Link" :fields="locationLinkFields"
        :modelValue="modalForm" :error="modalError" :options="locationLinkOptions" @submit="handleSubmit"
        @cancel="closeModal" />

      <!-- Confirm Delete Dialog -->
      <ConfirmationModal :visible="showDeleteConfirm" title="Delete Location Link"
        :message="`Are you sure you want to delete the link '${deleteTarget?.name}'?`" @confirm="deleteLocationLink"
        @cancel="closeDelete" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useLocationLinksStore } from '../../../stores/locationLinks';
import { useLocationsStore } from '../../../stores/locations';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import TableActionsMenu from '../../../components/TableActionsMenu.vue';

const locationLinksStore = useLocationLinksStore();
const locationsStore = useLocationsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const columns = [
  { key: 'name', label: 'Link Name' },
  { key: 'from_location_name', label: 'From Location' },
  { key: 'to_location_name', label: 'To Location' },
  { key: 'description', label: 'Description' },
  { key: 'created_at', label: 'Created' }
];

// Field configuration for ResourceModalForm
const locationLinkFields = [
  { key: 'from_adventure_game_location_id', label: 'From Location', type: 'select', required: true, placeholder: 'Select a location...' },
  { key: 'to_adventure_game_location_id', label: 'To Location', type: 'select', required: true, placeholder: 'Select a location...' },
  { key: 'name', label: 'Link Name', type: 'text', required: true, maxlength: 64, placeholder: 'e.g., North Path, Secret Door' },
  { key: 'description', label: 'Description', type: 'textarea', required: true, maxlength: 255, placeholder: 'Describe the link between locations...', rows: 3 }
];

// Options for select fields
const locationLinkOptions = computed(() => ({
  from_adventure_game_location_id: locationsStore.locations.map(location => ({
    value: location.id,
    label: location.name
  })),
  to_adventure_game_location_id: locationsStore.locations.map(location => ({
    value: location.id,
    label: location.name
  }))
}));

const showModal = ref(false);
const modalMode = ref('create');
const modalForm = ref({});
const modalError = ref('');
const showDeleteConfirm = ref(false);
const deleteTarget = ref(null);

// Watch for game selection changes
watch(
  () => selectedGame.value,
  (newGame) => {
    if (newGame) {
      locationLinksStore.fetchLocationLinks(newGame.id);
      locationsStore.fetchLocations(newGame.id);
    }
  },
  { immediate: true }
);

// Enhance location links with location names for display
const enhancedLocationLinks = computed(() => {
  return locationLinksStore.locationLinks.map(link => {
    const fromLocation = locationsStore.locations.find(loc => loc.id === link.from_adventure_game_location_id);
    const toLocation = locationsStore.locations.find(loc => loc.id === link.to_adventure_game_location_id);
    return {
      ...link,
      from_location_name: fromLocation?.name || 'Unknown',
      to_location_name: toLocation?.name || 'Unknown'
    };
  });
});

function openCreate() {
  modalMode.value = 'create';
  modalForm.value = {
    from_adventure_game_location_id: '',
    to_adventure_game_location_id: '',
    name: '',
    description: ''
  };
  modalError.value = '';
  showModal.value = true;
}

function openEdit(row) {
  modalMode.value = 'edit';
  modalForm.value = {
    id: row.id,
    from_adventure_game_location_id: row.from_adventure_game_location_id,
    to_adventure_game_location_id: row.to_adventure_game_location_id,
    name: row.name,
    description: row.description || ''
  };
  modalError.value = '';
  showModal.value = true;
}

function closeModal() {
  showModal.value = false;
  modalError.value = '';
}

async function handleSubmit(formData) {
  modalError.value = '';
  try {
    if (modalMode.value === 'create') {
      await locationLinksStore.createLocationLink(formData);
    } else {
      await locationLinksStore.updateLocationLink(modalForm.value.id, formData);
    }
    closeModal();
  } catch (err) {
    modalError.value = err.message || 'Failed to save.';
  }
}

function confirmDelete(row) {
  deleteTarget.value = row;
  showDeleteConfirm.value = true;
}

function closeDelete() {
  showDeleteConfirm.value = false;
  deleteTarget.value = null;
}

async function deleteLocationLink() {
  if (!deleteTarget.value) return;
  try {
    await locationLinksStore.deleteLocationLink(deleteTarget.value.id);
    closeDelete();
  } catch (err) {
    console.error('Failed to delete location link:', err);
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
      handler: () => confirmDelete(row)
    }
  ];
}
</script>
