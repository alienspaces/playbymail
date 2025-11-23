<!--
  StudioItemsView.vue
  This component follows the same pattern as StudioLocationsView.vue and StudioCreaturesView.vue.
-->
<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage items.</p>
    </div>
    <div v-else class="game-table-section">
      <p class="game-context-name">Game: {{ selectedGame.name }}</p>
      <div class="section-header">
        <h2>Items</h2>
        <button @click="openCreate">Create New Item</button>
      </div>
      <ResourceTable
        :columns="columns"
        :rows="itemsStore.items"
        :loading="itemsStore.loading"
        :error="itemsStore.error"
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
      title="Item"
      :fields="fields"
      :modelValue="modalForm"
      :error="modalError"
      @submit="handleSubmit"
      @cancel="closeModal"
    />

    <ConfirmationModal
      :visible="showDeleteModal"
      title="Delete Item"
      :message="`Are you sure you want to delete '${itemToDelete?.name}'?`"
      @confirm="handleDelete"
      @cancel="closeDeleteModal"
    />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { useItemsStore } from '../../../stores/items';
import { useGamesStore } from '../../../stores/games';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';

const itemsStore = useItemsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'description', label: 'Description' }
];

const fields = [
  { key: 'name', label: 'Name', required: true, maxlength: 1024 },
  { key: 'description', label: 'Description', required: true, maxlength: 4096, type: 'textarea' }
];

const showModal = ref(false);
const modalMode = ref('create');
const modalForm = ref({ name: '', description: '' });
const modalError = ref('');
const showDeleteModal = ref(false);
const itemToDelete = ref(null);

// Watch for game selection changes
watch(
  () => selectedGame.value,
  (newGame) => {
    if (newGame) {
      itemsStore.fetchItems(newGame.id);
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

function openEdit(item) {
  modalMode.value = 'edit';
  modalForm.value = { ...item };
  modalError.value = '';
  showModal.value = true;
}

function openDelete(item) {
  itemToDelete.value = item;
  showDeleteModal.value = true;
}

function closeModal() {
  showModal.value = false;
  modalForm.value = { name: '', description: '' };
  modalError.value = '';
}

function closeDeleteModal() {
  showDeleteModal.value = false;
  itemToDelete.value = null;
}

async function handleSubmit(formData) {
  modalError.value = '';
  try {
    if (modalMode.value === 'create') {
      await itemsStore.createItem(formData);
    } else {
      await itemsStore.updateItem(modalForm.value.id, formData);
    }
    closeModal();
  } catch (error) {
    modalError.value = error.message || 'Failed to save.';
  }
}

async function handleDelete() {
  try {
    await itemsStore.deleteItem(itemToDelete.value.id);
    closeDeleteModal();
  } catch (error) {
    console.error('Failed to delete item:', error);
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