<!--
  StudioGameInstancesView.vue
  This component follows the same pattern as other studio views.
-->
<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage game instances.</p>
    </div>
    <div v-else class="game-table-section">
      <p class="game-context-name">Game: {{ selectedGame.name }}</p>
      <div class="section-header">
        <h2>Game Instances</h2>
        <button @click="openCreate">Create New Game Instance</button>
      </div>
      <ResourceTable
        :columns="columns"
        :rows="gameInstancesStore.gameInstances"
        :loading="gameInstancesStore.loading"
        :error="gameInstancesStore.error"
      >
        <template #actions="{ row }">
          <button @click="openEdit(row)">Edit</button>
          <button @click="openDelete(row)">Delete</button>
        </template>
      </ResourceTable>
    </div>

    <ResourceModalForm
      :visible="showModal"
      :mode="modalMode"
      title="Game Instance"
      :fields="fields"
      :modelValue="modalForm"
      :error="modalError"
      @submit="handleSubmit"
      @cancel="closeModal"
    />

    <ConfirmationModal
      :visible="showDeleteModal"
      title="Delete Game Instance"
      :message="`Are you sure you want to delete '${instanceToDelete?.name}'?`"
      @confirm="handleDelete"
      @cancel="closeDeleteModal"
    />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { useGameInstancesStore } from '../../../stores/gameInstances';
import { useGamesStore } from '../../../stores/games';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';

const gameInstancesStore = useGameInstancesStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'description', label: 'Description' },
  { key: 'max_turns', label: 'Max Turns' },
  { key: 'turn_deadline_hours', label: 'Turn Deadline (hours)' }
];

const fields = [
  { key: 'name', label: 'Name', required: true, maxlength: 1024 },
  { key: 'description', label: 'Description', required: false, maxlength: 4096, type: 'textarea' },
  { key: 'max_turns', label: 'Max Turns (optional)', required: false, type: 'number', min: 1 },
  { key: 'turn_deadline_hours', label: 'Turn Deadline (hours)', required: true, type: 'number', min: 1 }
];

const showModal = ref(false);
const modalMode = ref('create');
const modalForm = ref({
  name: '',
  description: '',
  max_turns: null,
  turn_deadline_hours: null
});
const modalError = ref('');
const showDeleteModal = ref(false);
const instanceToDelete = ref(null);

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
    turn_deadline_hours: null
  };
  modalError.value = '';
  showModal.value = true;
}

function openEdit(instance) {
  modalMode.value = 'edit';
  modalForm.value = { ...instance };
  modalError.value = '';
  showModal.value = true;
}

function openDelete(instance) {
  instanceToDelete.value = instance;
  showDeleteModal.value = true;
}

function closeModal() {
  showModal.value = false;
  modalForm.value = {
    name: '',
    description: '',
    max_turns: null,
    turn_deadline_hours: null
  };
  modalError.value = '';
}

function closeDeleteModal() {
  showDeleteModal.value = false;
  instanceToDelete.value = null;
}

async function handleSubmit(formData) {
  modalError.value = '';
  try {
    if (modalMode.value === 'create') {
      await gameInstancesStore.createGameInstance(formData);
    } else {
      await gameInstancesStore.updateGameInstance(modalForm.value.id, formData);
    }
    closeModal();
  } catch (error) {
    modalError.value = error.message || 'Failed to save.';
  }
}

async function handleDelete() {
  try {
    await gameInstancesStore.deleteGameInstance(instanceToDelete.value.id);
    closeDeleteModal();
  } catch (error) {
    console.error('Failed to delete game instance:', error);
  }
}
</script>

<style scoped>
.game-table-section {
  margin-top: 20px;
}

.game-context-name {
  font-weight: bold;
  margin-bottom: 10px;
}

.section-header {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  margin-bottom: 20px;
}

.section-header h2 {
  margin: 0 0 10px 0;
}

.section-header button {
  padding: 8px 16px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.section-header button:hover {
  background-color: #0056b3;
}
</style> 