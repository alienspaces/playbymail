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
      <GameContext :gameName="selectedGame.name" />
      <PageHeader 
        title="Game Instances" 
        actionText="Create New Game Instance" 
        :showIcon="false"
        titleLevel="h2"
        @action="openCreate"
      />
      <ResourceTable
        :columns="columns"
        :rows="gameInstancesStore.gameInstances"
        :loading="gameInstancesStore.loading"
        :error="gameInstancesStore.error"
      >
        <template #actions="{ row }">
          <TableActionsMenu :actions="getActions(row)" />
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
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import TableActionsMenu from '../../../components/TableActionsMenu.vue';

const gameInstancesStore = useGameInstancesStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const columns = [
  { key: 'id', label: 'ID' },
  { key: 'game_id', label: 'Game ID' },
  { key: 'status', label: 'Status' },
  { key: 'current_turn', label: 'Current Turn' },
  { key: 'last_turn_processed_at', label: 'Last Turn Processed At' },
  { key: 'next_turn_due_at', label: 'Next Turn Due At' },
  { key: 'started_at', label: 'Started At' },
  { key: 'completed_at', label: 'Completed At' },
];

const fields = [      
  { key: 'game_id', label: 'Game ID', required: true, maxlength: 1024 },
  { key: 'status', label: 'Status', required: true, maxlength: 4096, type: 'select', options: [
    { label: 'Created', value: 'created' },
    { label: 'Started', value: 'started' },
    { label: 'Paused', value: 'paused' },
    { label: 'Completed', value: 'completed' },
  ] }
];

const showModal = ref(false);
const modalMode = ref('create');
const modalForm = ref({
  name: '',
  description: ''
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
    description: ''
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
    description: ''
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
.game-table-section {
  /* Consistent spacing - no top margin, handled by GameContext */
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
  background: transparent;
  color: var(--color-button);
  border: 2px solid var(--color-button);
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.section-header button:hover {
  background: var(--color-button);
  color: var(--color-text-light);
}
</style> 