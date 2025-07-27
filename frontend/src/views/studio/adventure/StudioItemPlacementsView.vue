<template>
  <div>
    <div v-if="!selectedGame">
      <p>Please select or create a game to manage item placements.</p>
    </div>
    <div v-else class="game-table-section">
      <p class="game-context-name">Game: {{ selectedGame.name }}</p>
      <div class="section-header">
        <h2>Item Placements</h2>
        <button @click="openItemPlacementCreate">Create Item Placement</button>
      </div>
      <ResourceTable
        :columns="itemPlacementColumns"
        :rows="itemPlacementsStore.itemPlacements"
        :loading="itemPlacementsStore.loading"
        :error="itemPlacementsStore.error"
      >
        <template #actions="{ row }">
          <div style="display: flex; gap: var(--space-sm);">
            <button @click="openItemPlacementEdit(row)">Edit</button>
            <button @click="confirmItemPlacementDelete(row)">Delete</button>
          </div>
        </template>
      </ResourceTable>

      <!-- Create/Edit Item Placement Modal -->
      <div v-if="showItemPlacementModal" class="modal-overlay">
        <div class="modal">
          <h2>{{ itemPlacementModalMode === 'create' ? 'Create Item Placement' : 'Edit Item Placement' }}</h2>
          <form @submit.prevent="handleItemPlacementSubmit(itemPlacementModalForm)">
            <div class="form-group">
              <label for="adventure_game_item_id">Item:</label>
              <select
                id="adventure_game_item_id"
                v-model="itemPlacementModalForm.adventure_game_item_id"
                required
                class="form-select"
              >
                <option value="" disabled>Select an item...</option>
                <option v-for="item in itemsStore.items" :key="item.id" :value="item.id">{{ item.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label for="adventure_game_location_id">Location:</label>
              <select
                id="adventure_game_location_id"
                v-model="itemPlacementModalForm.adventure_game_location_id"
                required
                class="form-select"
              >
                <option value="" disabled>Select a location...</option>
                <option v-for="loc in locationsStore.locations" :key="loc.id" :value="loc.id">{{ loc.name }}</option>
              </select>
            </div>
            <div class="modal-actions">
              <button type="submit">{{ itemPlacementModalMode === 'create' ? 'Create' : 'Save' }}</button>
              <button type="button" @click="closeItemPlacementModal">Cancel</button>
            </div>
          </form>
          <p v-if="itemPlacementModalError" class="error">{{ itemPlacementModalError }}</p>
        </div>
      </div>

      <!-- Confirm Delete Dialog -->
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
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useItemsStore } from '../../../stores/items';
import { useLocationsStore } from '../../../stores/locations';
import { useItemPlacementsStore } from '../../../stores/itemPlacements';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';

const itemsStore = useItemsStore();
const locationsStore = useLocationsStore();
const itemPlacementsStore = useItemPlacementsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const itemPlacementColumns = [
  { key: 'adventure_game_item_id', label: 'Item ID' },
  { key: 'adventure_game_location_id', label: 'Location ID' },
  { key: 'created_at', label: 'Created' }
];

const showItemPlacementModal = ref(false);
const itemPlacementModalMode = ref('create');
const itemPlacementModalForm = ref({ adventure_game_item_id: '', adventure_game_location_id: '' });
const itemPlacementModalError = ref('');
const showItemPlacementDeleteConfirm = ref(false);
const itemPlacementDeleteTarget = ref(null);
const itemPlacementDeleteError = ref('');

// Watch for game selection changes
watch(
  () => selectedGame.value,
  (newGame) => {
    if (newGame) {
      itemsStore.fetchItems(newGame.id);
      locationsStore.fetchLocations(newGame.id);
      itemPlacementsStore.fetchItemPlacements(newGame.id);
    }
  },
  { immediate: true }
);

function openItemPlacementCreate() {
  itemPlacementModalMode.value = 'create';
  itemPlacementModalForm.value = { adventure_game_item_id: '', adventure_game_location_id: '' };
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