<!--
  StudioCreaturesView.vue
  This component follows the same pattern as StudioLocationsView.vue and StudioItemsView.vue.
-->
<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage creatures.</p>
    </div>
    <div v-else class="game-table-section">
      <p class="game-context-name">Game: {{ selectedGame.name }}</p>
      <div class="section-header">
        <h2>Creatures</h2>
        <button @click="openCreate">Create New Creature</button>
      </div>
      <ResourceTable
        :columns="columns"
        :rows="creaturesStore.creatures"
        :loading="creaturesStore.loading"
        :error="creaturesStore.error"
      >
        <template #actions="{ row }">
          <button @click="openEdit(row)">Edit</button>
          <button @click="confirmDelete(row)">Delete</button>
        </template>
      </ResourceTable>

      <!-- Create/Edit Modal using ResourceModalForm -->
      <ResourceModalForm
        :visible="showModal"
        :mode="modalMode"
        title="Creature"
        :fields="fields"
        :modelValue="modalForm"
        :error="modalError"
        @submit="handleSubmit"
        @cancel="closeModal"
      >
        <template v-slot:field="{ field }">
          <textarea 
            v-if="field.key === 'description'" 
            v-model="modalForm[field.key]" 
            :rows="4" 
            :maxlength="field.maxlength" 
            :required="field.required"
          />
          <input 
            v-else 
            v-model="modalForm[field.key]" 
            :type="field.type || 'text'" 
            :required="field.required" 
            :maxlength="field.maxlength" 
          />
        </template>
      </ResourceModalForm>

      <!-- Confirm Delete Dialog -->
      <div v-if="showDeleteConfirm" class="modal-overlay">
        <div class="modal">
          <h2>Delete Creature</h2>
          <p>Are you sure you want to delete <b>{{ deleteTarget?.name }}</b>?</p>
          <div class="modal-actions">
            <button @click="deleteCreature">Delete</button>
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
import { storeToRefs } from 'pinia';
import { useCreaturesStore } from '../../../stores/creatures';
import { useGamesStore } from '../../../stores/games';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';

const creaturesStore = useCreaturesStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'description', label: 'Description' },
  { key: 'created_at', label: 'Created' }
];

// Field configuration for ResourceModalForm
const fields = [
  { key: 'name', label: 'Name', required: true, maxlength: 1024 },
  { key: 'description', label: 'Description', required: true, maxlength: 4096, type: 'textarea' }
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
      creaturesStore.fetchCreatures(newGame.id);
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
      await creaturesStore.createCreature(form);
    } else {
      await creaturesStore.updateCreature(modalForm.value.id, form);
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
async function deleteCreature() {
  if (!deleteTarget.value) return;
  deleteError.value = '';
  try {
    await creaturesStore.deleteCreature(deleteTarget.value.id);
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