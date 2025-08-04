<!--
  StudioLocationsView.vue
  This component follows the same pattern as StudioItemsView.vue and StudioCreaturesView.vue.
-->
<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage locations.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <SectionHeader 
        title="Locations" 
        resourceName="Location" 
        @create="openCreate"
      />
      <ResourceTable
        :columns="columns"
        :rows="locationsStore.locations"
        :loading="locationsStore.loading"
        :error="locationsStore.error"
      >
        <template #actions="{ row }">
          <button @click="openEdit(row)">Edit</button>
          <button @click="confirmDelete(row)">Delete</button>
        </template>
      </ResourceTable>

      <!-- Create/Edit Modal -->
      <div v-if="showModal" class="modal-overlay">
        <div class="modal">
          <h2>{{ modalMode === 'create' ? 'Create Location' : 'Edit Location' }}</h2>
          <form @submit.prevent="handleSubmit(modalForm)">
            <div class="form-group">
              <label for="location-name">Name:</label>
              <input v-model="modalForm.name" id="location-name" required maxlength="1024" />
            </div>
            <div class="form-group">
              <label for="location-description">Description:</label>
              <textarea v-model="modalForm.description" id="location-description" rows="4" maxlength="4096" required />
            </div>
            <div class="modal-actions">
              <button type="submit">{{ modalMode === 'create' ? 'Create' : 'Save' }}</button>
              <button type="button" @click="closeModal">Cancel</button>
            </div>
          </form>
          <p v-if="modalError" class="error">{{ modalError }}</p>
        </div>
      </div>

      <!-- Confirm Delete Dialog -->
      <div v-if="showDeleteConfirm" class="modal-overlay">
        <div class="modal">
          <h2>Delete Location</h2>
          <p>Are you sure you want to delete <b>{{ deleteTarget?.name }}</b>?</p>
          <div class="modal-actions">
            <button @click="deleteLocation">Delete</button>
            <button @click="closeDelete">Cancel</button>
          </div>
          <p v-if="deleteError" class="error">{{ deleteError }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useLocationsStore } from '../../../stores/locations';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import GameContext from '../../../components/GameContext.vue';
import SectionHeader from '../../../components/SectionHeader.vue';

const locationsStore = useLocationsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

// Components are imported and used directly in template

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'description', label: 'Description' },
  { key: 'created_at', label: 'Created' }
];

const showModal = ref(false);
const modalMode = ref('create');
const modalForm = ref({ name: '', description: '' });
const modalError = ref('');
const showDeleteConfirm = ref(false);
const deleteTarget = ref(null);
const deleteError = ref('');

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
  modalForm.value = { name: '', description: '' };
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
      await locationsStore.createLocation(form);
    } else {
      await locationsStore.updateLocation(modalForm.value.id, form);
    }
    closeModal();
  } catch (err) {
    modalError.value = err.message || 'Failed to save.';
  }
}
function confirmDelete(row) {
  deleteTarget.value = row;
  deleteError.value = '';
  showDeleteConfirm.value = true;
}
function closeDelete() {
  showDeleteConfirm.value = false;
  deleteTarget.value = null;
  deleteError.value = '';
}
async function deleteLocation() {
  if (!deleteTarget.value) return;
  deleteError.value = '';
  try {
    await locationsStore.deleteLocation(deleteTarget.value.id);
    closeDelete();
  } catch (err) {
    deleteError.value = err.message || 'Failed to delete.';
  }
}
</script>

<style scoped>
.game-table-section {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}
button {
  margin-right: var(--space-sm);
}
.game-context-name {
  font-size: 1.1rem;
  font-weight: 400;
  color: #444;
  margin-bottom: var(--space-sm);
}
</style> 