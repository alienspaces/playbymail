<template>
  <div>
    <div v-if="!gameId">
      <p>Please select or create a game to manage creature placements.</p>
    </div>
    <div v-else>
      <h2>Creature Placements</h2>
      <button @click="openCreaturePlacementCreate">Create Creature Placement</button>
      <ResourceTable
        :columns="creaturePlacementColumns"
        :rows="creaturePlacementsStore.creaturePlacements"
        :loading="creaturePlacementsStore.loading"
        :error="creaturePlacementsStore.error"
      >
        <template #actions="{ row }">
          <div style="display: flex; gap: var(--space-sm);">
            <button @click="openCreaturePlacementEdit(row)">Edit</button>
            <button @click="confirmCreaturePlacementDelete(row)">Delete</button>
          </div>
        </template>
      </ResourceTable>
      <ResourceModalForm
        :visible="showCreaturePlacementModal"
        :mode="creaturePlacementModalMode"
        title="Creature Placement"
        :fields="creaturePlacementFields"
        :modelValue="creaturePlacementModalForm"
        :error="creaturePlacementModalError"
        @submit="handleCreaturePlacementSubmit"
        @cancel="closeCreaturePlacementModal"
      >
        <template v-slot:field="{ field }">
          <select
            v-if="field.key === 'adventure_game_creature_id'"
            :id="field.key"
            v-model="creaturePlacementModalForm.adventure_game_creature_id"
            required
            class="form-select"
          >
            <option value="" disabled>Select a creature...</option>
            <option v-for="creature in creaturesStore.creatures" :key="creature.id" :value="creature.id">{{ creature.name }}</option>
          </select>
          <select
            v-else-if="field.key === 'adventure_game_location_id'"
            :id="field.key"
            v-model="creaturePlacementModalForm.adventure_game_location_id"
            required
            class="form-select"
          >
            <option value="" disabled>Select a location...</option>
            <option v-for="loc in locationsStore.locations" :key="loc.id" :value="loc.id">{{ loc.name }}</option>
          </select>
          <input
            v-else
            v-model="creaturePlacementModalForm[field.key]"
            :id="field.key"
            :type="field.type || 'text'"
            :required="field.required"
            :maxlength="field.maxlength"
            :placeholder="field.placeholder"
          />
        </template>
      </ResourceModalForm>
      <div v-if="showCreaturePlacementDeleteConfirm" class="modal-overlay">
        <div class="modal">
          <h2>Delete Creature Placement</h2>
          <p>Are you sure you want to delete this creature placement?</p>
          <div class="modal-actions">
            <button @click="deleteCreaturePlacement">Delete</button>
            <button @click="closeCreaturePlacementDelete">Cancel</button>
          </div>
          <p v-if="creaturePlacementDeleteError" class="error">{{ creaturePlacementDeleteError }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useCreaturesStore } from '../../../stores/creatures';
import { useLocationsStore } from '../../../stores/locations';
import { useCreaturePlacementsStore } from '../../../stores/creaturePlacements';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';

const creaturesStore = useCreaturesStore();
const locationsStore = useLocationsStore();
const creaturePlacementsStore = useCreaturePlacementsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const gameId = computed(() => selectedGame.value ? selectedGame.value.id : null);

watch(
  () => gameId.value,
  (newGameId) => {
    if (newGameId) {
      creaturesStore.fetchCreatures(newGameId);
      locationsStore.fetchLocations(newGameId);
      creaturePlacementsStore.fetchCreaturePlacements(newGameId);
    }
  },
  { immediate: true }
);

const creaturePlacementColumns = [
  { key: 'adventure_game_creature_id', label: 'Creature' },
  { key: 'adventure_game_location_id', label: 'Location' },
  { key: 'initial_count', label: 'Count' },
  { key: 'created_at', label: 'Created' }
  // Do NOT include { key: 'actions', label: 'Actions' }
];
const creaturePlacementFields = [
  { key: 'adventure_game_creature_id', label: 'Creature', required: true },
  { key: 'adventure_game_location_id', label: 'Location', required: true },
  { key: 'initial_count', label: 'Count', type: 'number', required: true, min: 1 }
];

const showCreaturePlacementModal = ref(false);
const creaturePlacementModalMode = ref('create');
const creaturePlacementModalForm = ref({});
const creaturePlacementModalError = ref('');
const showCreaturePlacementDeleteConfirm = ref(false);
const creaturePlacementDeleteTarget = ref(null);
const creaturePlacementDeleteError = ref('');

function openCreaturePlacementCreate() {
  creaturePlacementModalMode.value = 'create';
  creaturePlacementModalForm.value = { adventure_game_creature_id: '', adventure_game_location_id: '', initial_count: 1 };
  creaturePlacementModalError.value = '';
  showCreaturePlacementModal.value = true;
}
function openCreaturePlacementEdit(row) {
  creaturePlacementModalMode.value = 'edit';
  creaturePlacementModalForm.value = { ...row };
  creaturePlacementModalError.value = '';
  showCreaturePlacementModal.value = true;
}
function closeCreaturePlacementModal() {
  showCreaturePlacementModal.value = false;
  creaturePlacementModalError.value = '';
}
async function handleCreaturePlacementSubmit(form) {
  creaturePlacementModalError.value = '';
  try {
    if (creaturePlacementModalMode.value === 'create') {
      await creaturePlacementsStore.createCreaturePlacement(form);
    } else {
      await creaturePlacementsStore.updateCreaturePlacement(creaturePlacementModalForm.value.id, form);
    }
    closeCreaturePlacementModal();
  } catch (err) {
    creaturePlacementModalError.value = err.message || 'Failed to save.';
  }
}
function confirmCreaturePlacementDelete(row) {
  creaturePlacementDeleteTarget.value = row;
  creaturePlacementDeleteError.value = '';
  showCreaturePlacementDeleteConfirm.value = true;
}
function closeCreaturePlacementDelete() {
  showCreaturePlacementDeleteConfirm.value = false;
  creaturePlacementDeleteTarget.value = null;
  creaturePlacementDeleteError.value = '';
}
async function deleteCreaturePlacement() {
  if (!creaturePlacementDeleteTarget.value) return;
  creaturePlacementDeleteError.value = '';
  try {
    await creaturePlacementsStore.deleteCreaturePlacement(creaturePlacementDeleteTarget.value.id);
    closeCreaturePlacementDelete();
  } catch (err) {
    creaturePlacementDeleteError.value = err.message || 'Failed to delete.';
  }
}
</script>

<style scoped>
h2 {
  margin-top: 2rem;
  margin-bottom: 1rem;
}
button {
  margin-bottom: 1rem;
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
</style> 