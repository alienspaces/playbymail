<!--
  StudioGameInstancesView.vue
  This component manages game instances for a selected game.
-->
<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage game instances.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <SectionHeader 
        title="Game Instances" 
        resourceName="Game Instance" 
        @create="openCreate"
      />
      <ResourceTable
        :columns="columns"
        :rows="gameInstancesStore.gameInstances"
        :loading="gameInstancesStore.loading"
        :error="gameInstancesStore.error"
      >
        <template #actions="{ row }">
          <button @click="openDetail(row)">View</button>
          <button @click="openEdit(row)">Edit</button>
          <button @click="confirmDelete(row)">Delete</button>
        </template>
      </ResourceTable>

      <!-- Create/Edit Modal -->
      <div v-if="showModal" class="modal-overlay">
        <div class="modal">
          <h2>{{ modalMode === 'create' ? 'Create Game Instance' : 'Edit Game Instance' }}</h2>
          <form @submit.prevent="handleSubmit(modalForm)">
            <div class="form-group">
              <label for="instance-name">Name:</label>
              <input v-model="modalForm.name" id="instance-name" required maxlength="1024" />
            </div>
            <div class="form-group">
              <label for="instance-description">Description:</label>
              <textarea v-model="modalForm.description" id="instance-description" rows="4" maxlength="4096" />
            </div>
            <div class="form-group">
              <label for="instance-max-turns">Max Turns (optional):</label>
              <input v-model.number="modalForm.max_turns" id="instance-max-turns" type="number" min="1" />
            </div>
            <div class="form-group">
              <label for="instance-turn-deadline">Turn Deadline (hours):</label>
              <input v-model.number="modalForm.turn_deadline_hours" id="instance-turn-deadline" type="number" min="1" required />
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
          <h2>Delete Game Instance</h2>
          <p>Are you sure you want to delete <b>{{ deleteTarget?.name }}</b>?</p>
          <div class="modal-actions">
            <button @click="deleteInstance">Delete</button>
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
import { useRouter } from 'vue-router';
import { useGameInstancesStore } from '../../../stores/gameInstances';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import GameContext from '../../../components/GameContext.vue';
import SectionHeader from '../../../components/SectionHeader.vue';

const router = useRouter();
const gameInstancesStore = useGameInstancesStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'status', label: 'Status' },
  { key: 'current_turn', label: 'Current Turn' },
  { key: 'max_turns', label: 'Max Turns' },
  { key: 'turn_deadline_hours', label: 'Turn Deadline (hrs)' },
  { key: 'created_at', label: 'Created' }
];

const showModal = ref(false);
const modalMode = ref('create');
const modalForm = ref({ 
  name: '', 
  description: '', 
  max_turns: null, 
  turn_deadline_hours: 168 
});
const modalError = ref('');
const showDeleteConfirm = ref(false);
const deleteTarget = ref(null);
const deleteError = ref('');

// Watch for game selection changes
watch(
  () => selectedGame.value,
  (newGame) => {
    if (newGame) {
      gameInstancesStore.fetchGameInstances(newGame.id);
    }
  },
  { immediate: true }
);

function openCreate() {
  modalMode.value = 'create';
  modalForm.value = { 
    name: '', 
    description: '', 
    max_turns: null, 
    turn_deadline_hours: 168 
  };
  modalError.value = '';
  showModal.value = true;
}

function openEdit(row) {
  modalMode.value = 'edit';
  modalForm.value = { 
    name: row.name, 
    description: row.description || '', 
    max_turns: row.max_turns, 
    turn_deadline_hours: row.turn_deadline_hours 
  };
  modalError.value = '';
  showModal.value = true;
}

function openDetail(row) {
  router.push(`/studio/${selectedGame.value.id}/instances/${row.id}`);
}

function closeModal() {
  showModal.value = false;
  modalError.value = '';
}

async function handleSubmit(form) {
  modalError.value = '';
  try {
    if (modalMode.value === 'create') {
      await gameInstancesStore.createGameInstance(selectedGame.value.id, form);
    } else {
      await gameInstancesStore.updateGameInstance(selectedGame.value.id, modalForm.value.id, form);
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

async function deleteInstance() {
  if (!deleteTarget.value) return;
  deleteError.value = '';
  try {
    await gameInstancesStore.deleteGameInstance(selectedGame.value.id, deleteTarget.value.id);
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

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.modal h2 {
  margin-bottom: 1rem;
  color: #333;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #333;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.1);
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

.modal-actions button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
}

.modal-actions button[type="submit"] {
  background: #007bff;
  color: white;
}

.modal-actions button[type="button"] {
  background: #6c757d;
  color: white;
}

.error {
  color: #dc3545;
  margin-top: 1rem;
  padding: 0.75rem;
  background: #f8d7da;
  border-radius: 4px;
  border: 1px solid #f5c6cb;
}
</style> 