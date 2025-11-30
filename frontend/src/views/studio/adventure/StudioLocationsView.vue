<!--
  StudioLocationsView.vue
  This component follows the same pattern as StudioCreaturesView.vue and StudioItemsView.vue.
  Added: Turn sheet image upload and preview functionality.
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

    <!-- Custom modal for create/edit with image upload support -->
    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal modal-wide">
          <h2>{{ modalMode === 'create' ? 'Create Location' : 'Edit Location' }}</h2>
          <form @submit.prevent="handleSubmit(modalForm)" class="modal-form">
            <div class="form-group">
              <label for="location-name">Name <span class="required">*</span></label>
              <input v-model="modalForm.name" id="location-name" required maxlength="1024" autocomplete="off" />
            </div>
            <div class="form-group">
              <label for="location-description">Description <span class="required">*</span></label>
              <textarea v-model="modalForm.description" id="location-description" required maxlength="4096"
                rows="4"></textarea>
            </div>
            <div class="form-group checkbox-group">
              <label class="checkbox-label">
                <input type="checkbox" v-model="modalForm.is_starting_location" />
                This is a starting location for new players
              </label>
            </div>

            <!-- Turn Sheet Image Upload (only in edit mode) -->
            <div v-if="modalMode === 'edit' && modalForm.id && selectedGame" class="form-section">
              <LocationTurnSheetImageUpload :gameId="selectedGame.id" :locationId="modalForm.id"
                @imagesUpdated="onImagesUpdated" @loadingChanged="onImageUploadLoadingChanged" />
            </div>

            <div class="modal-actions">
              <button type="submit" :disabled="imageUploadLoading">
                {{ modalMode === 'create' ? 'Create' : 'Save' }}
              </button>
              <button type="button" @click="closeModal" :disabled="imageUploadLoading">
                {{ imageUploadLoading ? 'Uploading...' : 'Cancel' }}
              </button>
            </div>
          </form>
          <div v-if="modalError" class="error">
            <p>{{ modalError }}</p>
          </div>
        </div>
      </div>
    </Teleport>

    <ConfirmationModal :visible="showDeleteModal" title="Delete Location"
      :message="`Are you sure you want to delete '${locationToDelete?.name}'?`" @confirm="handleDelete"
      @cancel="closeDeleteModal" />

    <!-- Turn Sheet Preview Modal -->
    <LocationTurnSheetPreviewModal :visible="showPreviewModal" :gameId="selectedGame?.id || ''"
      :locationId="previewLocationId" :locationName="previewLocationName" title="Location Choice Turn Sheet Preview"
      @close="closePreviewModal" />
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { storeToRefs } from 'pinia';
import { useLocationsStore } from '../../../stores/locations';
import { useGamesStore } from '../../../stores/games';
import ResourceTable from '../../../components/ResourceTable.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import TableActions from '../../../components/TableActions.vue';
import LocationTurnSheetImageUpload from '../../../components/LocationTurnSheetImageUpload.vue';
import LocationTurnSheetPreviewModal from '../../../components/LocationTurnSheetPreviewModal.vue';

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

const showModal = ref(false);
const modalMode = ref('create');
const modalForm = ref({ name: '', description: '', is_starting_location: false });
const modalError = ref('');
const showDeleteModal = ref(false);
const locationToDelete = ref(null);

// Image upload loading state
const imageUploadLoading = ref(false);

// Preview modal state
const showPreviewModal = ref(false);
const previewLocationId = ref('');
const previewLocationName = ref('');

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

function openPreview(location) {
  previewLocationId.value = location.id;
  previewLocationName.value = location.name;
  showPreviewModal.value = true;
}

function closePreviewModal() {
  showPreviewModal.value = false;
  previewLocationId.value = '';
  previewLocationName.value = '';
}

function closeModal() {
  // Don't close if images are being uploaded
  if (imageUploadLoading.value) {
    console.log('[StudioLocationsView] Preventing modal close - upload in progress');
    return;
  }
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

function onImagesUpdated() {
  // Called when turn sheet images are updated
  console.log('[StudioLocationsView] Turn sheet images updated');
}

function onImageUploadLoadingChanged(isLoading) {
  console.log(`[StudioLocationsView] Image upload loading changed: ${isLoading}`);
  imageUploadLoading.value = isLoading;
}

function getActions(row) {
  return [
    {
      key: 'preview',
      label: 'Preview Turn Sheet',
      handler: () => openPreview(row)
    },
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

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.form-group {
  margin-bottom: var(--space-sm);
}

.form-group label {
  display: block;
  margin-bottom: var(--space-xs);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: var(--space-sm);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-base);
}

.form-group textarea {
  resize: vertical;
}

.checkbox-group {
  margin-top: var(--space-sm);
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  cursor: pointer;
  font-weight: normal;
}

.checkbox-label input[type="checkbox"] {
  width: auto;
}

.required {
  color: var(--color-danger);
}

.error {
  color: var(--color-warning-dark);
  background: var(--color-warning-light);
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-warning);
  margin-top: var(--space-md);
}

.error p {
  margin: 0;
}

.modal-wide {
  max-width: 600px;
  width: 100%;
}

.form-section {
  margin-top: var(--space-md);
  padding-top: var(--space-md);
  border-top: 1px solid var(--color-border);
}

@media (max-width: 768px) {
  .modal-actions {
    flex-direction: column-reverse;
  }

  .modal-actions button {
    width: 100%;
  }
}
</style>
