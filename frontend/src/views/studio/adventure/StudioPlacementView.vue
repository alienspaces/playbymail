<template>
  <div>
    <div v-if="!gameId">
      <p>Please select or create a game to manage placement.</p>
    </div>
    <div v-else>
      <!-- Items Section -->
      <section class="studio-section">
        <h2>Items</h2>
        <ResourceTable :columns="itemColumns" :rows="itemsStore.items" :loading="itemsStore.loading"
          :error="itemsStore.error" />
        <TablePagination :pageNumber="itemsStore.pageNumber" :hasMore="itemsStore.hasMore"
          @page-change="(p) => itemsStore.fetchAdventureGameItems(gameId, p)" />
      </section>
      <section class="studio-section">
        <h3>Item Placements</h3>
        <div class="section-actions">
          <button @click="openItemPlacementCreate">Create Item Placement</button>
        </div>
        <ResourceTable :columns="itemPlacementColumns" :rows="itemPlacementsStore.itemPlacements"
          :loading="itemPlacementsStore.loading" :error="itemPlacementsStore.error">
          <template #actions="{ row }">
            <TableActions :actions="getItemPlacementActions(row)" />
          </template>
        </ResourceTable>
        <TablePagination :pageNumber="itemPlacementsStore.pageNumber" :hasMore="itemPlacementsStore.hasMore"
          @page-change="(p) => itemPlacementsStore.fetchAdventureGameItemPlacements(gameId, p)" />

        <!-- Create/Edit Modal using ResourceModalForm -->
        <ResourceModalForm :visible="showItemPlacementModal" :mode="itemPlacementModalMode" title="Item Placement"
          :fields="itemPlacementFields" :modelValue="itemPlacementModalForm" :error="itemPlacementModalError"
          @submit="handleItemPlacementSubmit" @cancel="closeItemPlacementModal">
          <template v-slot:field="{ field }">
            <select v-if="field.key === 'adventure_game_item_id'"
              v-model="itemPlacementModalForm.adventure_game_item_id" class="form-select">
              <option v-for="item in itemsStore.items" :key="item.id" :value="item.id">{{ item.name }}</option>
            </select>
            <select v-else-if="field.key === 'adventure_game_location_id'"
              v-model="itemPlacementModalForm.adventure_game_location_id" class="form-select">
              <option v-for="loc in locationsStore.locations" :key="loc.id" :value="loc.id">{{ loc.name }}</option>
            </select>
            <input v-else v-model="itemPlacementModalForm[field.key]" :type="field.type || 'text'"
              :required="field.required" :maxlength="field.maxlength" :placeholder="field.placeholder" />
          </template>
        </ResourceModalForm>

        <ConfirmationModal :visible="showItemPlacementDeleteConfirm" title="Delete Item Placement"
          message="Are you sure you want to delete this item placement?" confirmText="Delete"
          :error="itemPlacementDeleteError" @confirm="deleteAdventureGameItemPlacement" @cancel="closeItemPlacementDelete" />
      </section>

      <!-- Creatures Section -->
      <section class="studio-section">
        <h2>Creatures</h2>
        <ResourceTable :columns="creatureColumns" :rows="creaturesStore.creatures" :loading="creaturesStore.loading"
          :error="creaturesStore.error" />
        <TablePagination :pageNumber="creaturesStore.pageNumber" :hasMore="creaturesStore.hasMore"
          @page-change="(p) => creaturesStore.fetchAdventureGameCreatures(gameId, p)" />
      </section>
      <section class="studio-section">
        <h3>Creature Placements</h3>
        <div class="section-actions">
          <button @click="openCreaturePlacementCreate">Create Creature Placement</button>
        </div>
        <ResourceTable :columns="creaturePlacementColumns" :rows="creaturePlacementsStore.creaturePlacements"
          :loading="creaturePlacementsStore.loading" :error="creaturePlacementsStore.error">
          <template #actions="{ row }">
            <TableActions :actions="getCreaturePlacementActions(row)" />
          </template>
        </ResourceTable>
        <TablePagination :pageNumber="creaturePlacementsStore.pageNumber" :hasMore="creaturePlacementsStore.hasMore"
          @page-change="(p) => creaturePlacementsStore.fetchAdventureGameCreaturePlacements(gameId, p)" />
        <ResourceModalForm :visible="showCreaturePlacementModal" :mode="creaturePlacementModalMode"
          title="Creature Placement" :fields="creaturePlacementFields" :modelValue="creaturePlacementModalForm"
          :error="creaturePlacementModalError" @submit="handleCreaturePlacementSubmit"
          @cancel="closeCreaturePlacementModal">
          <template v-slot:field="{ field }">
            <select v-if="field.key === 'adventure_game_creature_id'"
              v-model="creaturePlacementModalForm.adventure_game_creature_id" class="form-select">
              <option v-for="creature in creaturesStore.creatures" :key="creature.id" :value="creature.id">{{
                creature.name }}</option>
            </select>
            <select v-else-if="field.key === 'adventure_game_location_id'"
              v-model="creaturePlacementModalForm.adventure_game_location_id" class="form-select">
              <option v-for="loc in locationsStore.locations" :key="loc.id" :value="loc.id">{{ loc.name }}</option>
            </select>
            <input v-else v-model="creaturePlacementModalForm[field.key]" :type="field.type || 'text'"
              :required="field.required" :maxlength="field.maxlength" :placeholder="field.placeholder" />
          </template>
        </ResourceModalForm>
        <ConfirmationModal :visible="showCreaturePlacementDeleteConfirm" title="Delete Creature Placement"
          message="Are you sure you want to delete this creature placement?" confirmText="Delete"
          :error="creaturePlacementDeleteError" @confirm="deleteAdventureGameCreaturePlacement"
          @cancel="closeCreaturePlacementDelete" />
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useAdventureGameItemsStore } from '../../../stores/adventureGameItems';
import { useAdventureGameCreaturesStore } from '../../../stores/adventureGameCreatures';
import { useAdventureGameLocationsStore } from '../../../stores/adventureGameLocations';
import { useAdventureGameItemPlacementsStore } from '../../../stores/adventureGameItemPlacements';
import { useAdventureGameCreaturePlacementsStore } from '../../../stores/adventureGameCreaturePlacements';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';
import TableActions from '../../../components/TableActions.vue';
import TablePagination from '../../../components/TablePagination.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';

const itemsStore = useAdventureGameItemsStore();
const creaturesStore = useAdventureGameCreaturesStore();
const locationsStore = useAdventureGameLocationsStore();
const itemPlacementsStore = useAdventureGameItemPlacementsStore();
const creaturePlacementsStore = useAdventureGameCreaturePlacementsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const gameId = computed(() => selectedGame.value ? selectedGame.value.id : null);

watch(
  () => gameId.value,
  (newGameId) => {
    if (newGameId) {
      itemsStore.fetchAdventureGameItems(newGameId);
      creaturesStore.fetchAdventureGameCreatures(newGameId);
      locationsStore.fetchAdventureGameLocations(newGameId);
      itemPlacementsStore.fetchAdventureGameItemPlacements(newGameId);
      creaturePlacementsStore.fetchAdventureGameCreaturePlacements(newGameId);
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
      await itemPlacementsStore.createAdventureGameItemPlacement(form);
    } else {
      await itemPlacementsStore.updateAdventureGameItemPlacement(itemPlacementModalForm.value.id, form);
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
async function deleteAdventureGameItemPlacement() {
  if (!itemPlacementDeleteTarget.value) return;
  itemPlacementDeleteError.value = '';
  try {
    await itemPlacementsStore.deleteAdventureGameItemPlacement(itemPlacementDeleteTarget.value.id);
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
      await creaturePlacementsStore.createAdventureGameCreaturePlacement(form);
    } else {
      await creaturePlacementsStore.updateAdventureGameCreaturePlacement(creaturePlacementModalForm.value.id, form);
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
async function deleteAdventureGameCreaturePlacement() {
  if (!creaturePlacementDeleteTarget.value) return;
  creaturePlacementDeleteError.value = '';
  try {
    await creaturePlacementsStore.deleteAdventureGameCreaturePlacement(creaturePlacementDeleteTarget.value.id);
    closeCreaturePlacementDelete();
  } catch (err) {
    creaturePlacementDeleteError.value = err.message || 'Failed to delete.';
  }
}

function getItemPlacementActions(row) {
  return [
    {
      key: 'edit',
      label: 'Edit',
      handler: () => openItemPlacementEdit(row)
    },
    {
      key: 'delete',
      label: 'Delete',
      danger: true,
      handler: () => confirmItemPlacementDelete(row)
    }
  ];
}

function getCreaturePlacementActions(row) {
  return [
    {
      key: 'edit',
      label: 'Edit',
      handler: () => openCreaturePlacementEdit(row)
    },
    {
      key: 'delete',
      label: 'Delete',
      danger: true,
      handler: () => confirmCreaturePlacementDelete(row)
    }
  ];
}
</script>
