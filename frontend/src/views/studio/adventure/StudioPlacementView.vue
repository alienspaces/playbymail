<template>
  <div>
    <div v-if="!gameId">
      <p>Please select or create a game to manage placement.</p>
    </div>
    <div v-else>
      <!-- Items Section -->
      <section>
        <h2>Items</h2>
        <ResourceTable
          :columns="itemColumns"
          :rows="itemsStore.items"
          :loading="itemsStore.loading"
          :error="itemsStore.error"
        />
      </section>
      <section style="margin-bottom: 2rem;">
        <h3>Item Placements</h3>
        <button @click="openItemPlacementCreate">Create Item Placement</button>
        <ResourceTable
          :columns="itemPlacementColumns"
          :rows="itemPlacementsStore.itemPlacements"
          :loading="itemPlacementsStore.loading"
          :error="itemPlacementsStore.error"
        >
          <template #actions="{ row }">
            <button @click="openItemPlacementEdit(row)">Edit</button>
            <button @click="confirmItemPlacementDelete(row)">Delete</button>
          </template>
        </ResourceTable>
        <ResourceModalForm
          :visible="showItemPlacementModal"
          :mode="itemPlacementModalMode"
          title="Item Placement"
          :fields="itemPlacementFields"
          :modelValue="itemPlacementModalForm"
          :error="itemPlacementModalError"
          @submit="handleItemPlacementSubmit"
          @cancel="closeItemPlacementModal"
        >
          <template v-slot:field="{ field, value, update }">
            <select v-if="field.key === 'adventure_game_item_id'" v-model="itemPlacementModalForm.adventure_game_item_id">
              <option v-for="item in itemsStore.items" :key="item.id" :value="item.id">{{ item.name }}</option>
            </select>
            <select v-else-if="field.key === 'adventure_game_location_id'" v-model="itemPlacementModalForm.adventure_game_location_id">
              <option v-for="loc in locationsStore.locations" :key="loc.id" :value="loc.id">{{ loc.name }}</option>
            </select>
            <input v-else v-model="itemPlacementModalForm[field.key]" :type="field.type || 'text'" :required="field.required" :maxlength="field.maxlength" :placeholder="field.placeholder" />
          </template>
        </ResourceModalForm>
        <div v-if="showItemPlacementDeleteConfirm" class="modal-overlay">
          <div class="modal">
            <h2>Delete Item Placement</h2>
            <p>Are you sure you want to delete this item placement?</p>
            <div class="modal-actions">
              <button @click="deleteItemPlacement">Delete</button>
              <button @click="closeItemPlacementDelete">Cancel</button>
            </div>
            <p v-if="itemPlacementDeleteError" class="error">{{ itemPlacementDeleteError }}</p>
          </div>
        </div>
      </section>

      <!-- Creatures Section -->
      <section>
        <h2>Creatures</h2>
        <ResourceTable
          :columns="creatureColumns"
          :rows="creaturesStore.creatures"
          :loading="creaturesStore.loading"
          :error="creaturesStore.error"
        />
      </section>
      <section>
        <h3>Creature Placements</h3>
        <button @click="openCreaturePlacementCreate">Create Creature Placement</button>
        <ResourceTable
          :columns="creaturePlacementColumns"
          :rows="creaturePlacementsStore.creaturePlacements"
          :loading="creaturePlacementsStore.loading"
          :error="creaturePlacementsStore.error"
        >
          <template #actions="{ row }">
            <button @click="openCreaturePlacementEdit(row)">Edit</button>
            <button @click="confirmCreaturePlacementDelete(row)">Delete</button>
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
          <template v-slot:field="{ field, value, update }">
            <select v-if="field.key === 'adventure_game_creature_id'" v-model="creaturePlacementModalForm.adventure_game_creature_id">
              <option v-for="creature in creaturesStore.creatures" :key="creature.id" :value="creature.id">{{ creature.name }}</option>
            </select>
            <select v-else-if="field.key === 'adventure_game_location_id'" v-model="creaturePlacementModalForm.adventure_game_location_id">
              <option v-for="loc in locationsStore.locations" :key="loc.id" :value="loc.id">{{ loc.name }}</option>
            </select>
            <input v-else v-model="creaturePlacementModalForm[field.key]" :type="field.type || 'text'" :required="field.required" :maxlength="field.maxlength" :placeholder="field.placeholder" />
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
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useItemsStore } from '../../../stores/items';
import { useCreaturesStore } from '../../../stores/creatures';
import { useLocationsStore } from '../../../stores/locations';
import { useItemPlacementsStore } from '../../../stores/itemPlacements';
import { useCreaturePlacementsStore } from '../../../stores/creaturePlacements';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';

const itemsStore = useItemsStore();
const creaturesStore = useCreaturesStore();
const locationsStore = useLocationsStore();
const itemPlacementsStore = useItemPlacementsStore();
const creaturePlacementsStore = useCreaturePlacementsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const gameId = computed(() => selectedGame.value ? selectedGame.value.id : null);

watch(
  () => gameId.value,
  (newGameId) => {
    if (newGameId) {
      itemsStore.fetchItems(newGameId);
      creaturesStore.fetchCreatures(newGameId);
      locationsStore.fetchLocations(newGameId);
      itemPlacementsStore.fetchItemPlacements(newGameId);
      creaturePlacementsStore.fetchCreaturePlacements(newGameId);
    }
  },
  { immediate: true }
);

// Item columns
const itemColumns = [
  { key: 'name', label: 'Name' },
  { key: 'description', label: 'Description' },
  { key: 'created_at', label: 'Created' }
];
const itemPlacementColumns = [
  { key: 'adventure_game_item_id', label: 'Item' },
  { key: 'adventure_game_location_id', label: 'Location' },
  { key: 'initial_count', label: 'Count' },
  { key: 'created_at', label: 'Created' },
  { key: 'actions', label: 'Actions' }
];
const itemPlacementFields = [
  { key: 'adventure_game_item_id', label: 'Item', required: true },
  { key: 'adventure_game_location_id', label: 'Location', required: true },
  { key: 'initial_count', label: 'Count', type: 'number', required: true, min: 1 }
];

// Creature columns
const creatureColumns = [
  { key: 'name', label: 'Name' },
  { key: 'description', label: 'Description' },
  { key: 'created_at', label: 'Created' }
];
const creaturePlacementColumns = [
  { key: 'adventure_game_creature_id', label: 'Creature' },
  { key: 'adventure_game_location_id', label: 'Location' },
  { key: 'initial_count', label: 'Count' },
  { key: 'created_at', label: 'Created' },
  { key: 'actions', label: 'Actions' }
];
const creaturePlacementFields = [
  { key: 'adventure_game_creature_id', label: 'Creature', required: true },
  { key: 'adventure_game_location_id', label: 'Location', required: true },
  { key: 'initial_count', label: 'Count', type: 'number', required: true, min: 1 }
];

// Item Placement Modal State
const showItemPlacementModal = ref(false);
const itemPlacementModalMode = ref('create');
const itemPlacementModalForm = ref({});
const itemPlacementModalError = ref('');
const showItemPlacementDeleteConfirm = ref(false);
const itemPlacementDeleteTarget = ref(null);
const itemPlacementDeleteError = ref('');

function openItemPlacementCreate() {
  itemPlacementModalMode.value = 'create';
  itemPlacementModalForm.value = { adventure_game_item_id: '', adventure_game_location_id: '', initial_count: 1 };
  itemPlacementModalError.value = '';
  showItemPlacementModal.value = true;
}
function openItemPlacementEdit(row) {
  itemPlacementModalMode.value = 'edit';
  itemPlacementModalForm.value = { ...row };
  itemPlacementModalError.value = '';
  showItemPlacementModal.value = true;
}
function closeItemPlacementModal() {
  showItemPlacementModal.value = false;
  itemPlacementModalError.value = '';
}
async function handleItemPlacementSubmit(form) {
  itemPlacementModalError.value = '';
  try {
    if (itemPlacementModalMode.value === 'create') {
      await itemPlacementsStore.createItemPlacement(form);
    } else {
      await itemPlacementsStore.updateItemPlacement(itemPlacementModalForm.value.id, form);
    }
    closeItemPlacementModal();
  } catch (err) {
    itemPlacementModalError.value = err.message || 'Failed to save.';
  }
}
function confirmItemPlacementDelete(row) {
  itemPlacementDeleteTarget.value = row;
  itemPlacementDeleteError.value = '';
  showItemPlacementDeleteConfirm.value = true;
}
function closeItemPlacementDelete() {
  showItemPlacementDeleteConfirm.value = false;
  itemPlacementDeleteTarget.value = null;
  itemPlacementDeleteError.value = '';
}
async function deleteItemPlacement() {
  if (!itemPlacementDeleteTarget.value) return;
  itemPlacementDeleteError.value = '';
  try {
    await itemPlacementsStore.deleteItemPlacement(itemPlacementDeleteTarget.value.id);
    closeItemPlacementDelete();
  } catch (err) {
    itemPlacementDeleteError.value = err.message || 'Failed to delete.';
  }
}

// Creature Placement Modal State
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
h2, h3 {
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