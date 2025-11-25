<template>
  <div>
    <div v-if="!selectedGame">
      <p>Please select or create a game to manage item placements.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Item Placements" actionText="Create Item Placement" :showIcon="false" titleLevel="h2"
        @action="openItemPlacementCreate" />
      <ResourceTable :columns="itemPlacementColumns" :rows="enhancedItemPlacements"
        :loading="itemPlacementsStore.loading" :error="itemPlacementsStore.error">
        <template #actions="{ row }">
          <TableActionsMenu :actions="getActions(row)" />
        </template>
      </ResourceTable>

      <!-- Create/Edit Item Placement Modal -->
      <ResourceModalForm :visible="showItemPlacementModal" :mode="itemPlacementModalMode" title="Item Placement"
        :fields="itemPlacementFields" :modelValue="itemPlacementModalForm" :error="itemPlacementModalError"
        :options="itemPlacementOptions" @submit="handleItemPlacementSubmit" @cancel="closeItemPlacementModal" />

      <!-- Confirm Delete Dialog -->
      <ConfirmationModal :visible="showItemPlacementDeleteConfirm" title="Delete Item Placement"
        message="Are you sure you want to delete this item placement?" @confirm="deleteItemPlacement"
        @cancel="closeItemPlacementDelete" />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { useItemsStore } from '../../../stores/items';
import { useLocationsStore } from '../../../stores/locations';
import { useItemPlacementsStore } from '../../../stores/itemPlacements';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import TableActionsMenu from '../../../components/TableActionsMenu.vue';

const itemsStore = useItemsStore();
const locationsStore = useLocationsStore();
const itemPlacementsStore = useItemPlacementsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

// Enhance item placements with names for display
const enhancedItemPlacements = computed(() => {
  return itemPlacementsStore.itemPlacements.map(placement => {
    const item = itemsStore.items.find(item => item.id === placement.adventure_game_item_id);
    const location = locationsStore.locations.find(loc => loc.id === placement.adventure_game_location_id);
    return {
      ...placement,
      item_name: item?.name || 'Unknown Item',
      location_name: location?.name || 'Unknown Location'
    };
  });
});

const itemPlacementColumns = [
  { key: 'item_name', label: 'Item' },
  { key: 'location_name', label: 'Location' },
  { key: 'initial_count', label: 'Count' },
  { key: 'created_at', label: 'Created' }
];

// Field configuration for ResourceModalForm
const itemPlacementFields = [
  { key: 'adventure_game_item_id', label: 'Item', type: 'select', required: true, placeholder: 'Select an item...' },
  { key: 'adventure_game_location_id', label: 'Location', type: 'select', required: true, placeholder: 'Select a location...' },
  { key: 'initial_count', label: 'Initial Count', type: 'number', required: true, min: 1 }
];

// Options for select fields
const itemPlacementOptions = computed(() => ({
  adventure_game_item_id: itemsStore.items.map(item => ({
    value: item.id,
    label: item.name
  })),
  adventure_game_location_id: locationsStore.locations.map(location => ({
    value: location.id,
    label: location.name
  }))
}));

const showItemPlacementModal = ref(false);
const itemPlacementModalMode = ref('create');
const itemPlacementModalForm = ref({ adventure_game_item_id: '', adventure_game_location_id: '', initial_count: 1 });
const itemPlacementModalError = ref('');
const showItemPlacementDeleteConfirm = ref(false);
const itemPlacementDeleteTarget = ref(null);

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
  showItemPlacementDeleteConfirm.value = true;
}

function closeItemPlacementDelete() {
  showItemPlacementDeleteConfirm.value = false;
  itemPlacementDeleteTarget.value = null;
}

async function deleteItemPlacement() {
  if (!itemPlacementDeleteTarget.value) return;
  try {
    await itemPlacementsStore.deleteItemPlacement(itemPlacementDeleteTarget.value.id);
    closeItemPlacementDelete();
  } catch (err) {
    console.error('Failed to delete item placement:', err);
  }
}

function getActions(row) {
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
</script>
