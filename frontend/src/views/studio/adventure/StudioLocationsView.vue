<!--
  StudioLocationsView.vue
  This component follows the same pattern as StudioItemsView.vue and StudioCreaturesView.vue.
-->
<template>
  <div>
    <div v-if="!gameId">
      <p>Select a game to manage locations.</p>
    </div>
    <div v-else>
      <h3 v-if="selectedGame && selectedGame.name" class="game-context-name">
        Game: {{ selectedGame.name }}
      </h3>
      <div class="game-table-section">
        <h2>Game Locations</h2>
        <button @click="openCreate">Create New Location</button>
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
        <ResourceModalForm
          :visible="showModal"
          :mode="modalMode"
          title="Location"
          :fields="fields"
          :modelValue="modalForm"
          :error="modalError"
          @submit="handleSubmit"
          @cancel="closeModal"
        />
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
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useLocationsStore } from '../../../stores/locations';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';

const locationsStore = useLocationsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const gameId = computed(() => selectedGame.value ? selectedGame.value.id : null);

watch(
  () => gameId.value,
  (newGameId) => {
    if (newGameId) locationsStore.fetchLocations(newGameId);
  },
  { immediate: true }
);

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'description', label: 'Description' },
  { key: 'created_at', label: 'Created' }
];
const fields = [
  { key: 'name', label: 'Name', required: true, maxlength: 128 },
  { key: 'description', label: 'Description', required: false, maxlength: 1024 }
];

const showModal = ref(false);
const modalMode = ref('create');
const modalForm = ref({});
const modalError = ref('');
const showDeleteConfirm = ref(false);
const deleteTarget = ref(null);
const deleteError = ref('');

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
.game-table-section h2 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  font-size: 2rem;
}
button {
  margin-right: var(--space-sm);
}
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal {
  background: var(--color-bg);
  padding: var(--space-lg);
  border-radius: var(--radius-md);
  min-width: 300px;
  max-width: 90vw;
  box-shadow: 0 2px 16px rgba(0,0,0,0.2);
}
.modal-actions {
  margin-top: var(--space-md);
  display: flex;
  gap: var(--space-md);
  justify-content: flex-start;
}
.error {
  color: var(--color-error);
  margin-top: var(--space-md);
}
.game-context-name {
  font-size: 1.1rem;
  font-weight: 400;
  color: #444;
  margin-bottom: 0.5rem;
}
</style> 