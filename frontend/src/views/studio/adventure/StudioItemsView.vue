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
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Items" actionText="Create New Item" :showIcon="false" titleLevel="h2" @action="openCreate" />
      <ResourceTable :columns="columns" :rows="formattedItems" :loading="itemsStore.loading" :error="itemsStore.error">
        <template #cell-name="{ row }">
          <a href="#" class="edit-link" @click.prevent="openEdit(row)">{{ row.name }}</a>
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>
      <TablePagination :pageNumber="itemsStore.pageNumber" :hasMore="itemsStore.hasMore"
        @page-change="(p) => itemsStore.fetchItems(selectedGame.id, p)" />
    </div>

    <ResourceModalForm :visible="showModal" :mode="modalMode" title="Item" :fields="fields" :modelValue="modalForm"
      :error="modalError" @submit="handleSubmit" @cancel="closeModal" />

    <ConfirmationModal :visible="showDeleteModal" title="Delete Item"
      :message="`Are you sure you want to delete '${itemToDelete?.name}'?`" @confirm="handleDelete"
      @cancel="closeDeleteModal" />
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { storeToRefs } from 'pinia';
import { useItemsStore } from '../../../stores/items';
import { useGamesStore } from '../../../stores/games';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import TableActions from '../../../components/TableActions.vue';
import TablePagination from '../../../components/TablePagination.vue';

const itemsStore = useItemsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

// Format items for table display
const formattedItems = computed(() => {
  return itemsStore.items.map(item => ({
    ...item,
    is_starting_item: item.is_starting_item ? 'Yes' : 'No',
    can_be_equipped: item.can_be_equipped ? 'Yes' : 'No',
  }));
});

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'description', label: 'Description' },
  { key: 'is_starting_item', label: 'Starting Item' },
  { key: 'can_be_equipped', label: 'Equippable' },
];

const fields = [
  { key: 'name', label: 'Name', required: true, maxlength: 1024 },
  { key: 'description', label: 'Description', required: true, maxlength: 4096, type: 'textarea' },
  { key: 'is_starting_item', label: 'Starting Item', type: 'checkbox', help: 'Automatically assigned to characters when they join the game' },
  { key: 'can_be_equipped', label: 'Can Be Equipped', type: 'checkbox', help: 'Player can equip this item to gain stats via item effects' },
  {
    key: 'item_category',
    label: 'Item Category',
    type: 'select',
    placeholder: 'Select category',
    options: [
      { value: '', label: '— None —' },
      { value: 'weapon', label: 'Weapon' },
      { value: 'armor', label: 'Armor' },
      { value: 'clothing', label: 'Clothing' },
      { value: 'jewelry', label: 'Jewelry' },
      { value: 'consumable', label: 'Consumable' },
      { value: 'misc', label: 'Miscellaneous' }
    ]
  },
  { key: 'equipment_slot', label: 'Equipment Slot', placeholder: 'e.g. weapon, armor_chest', help: 'Slot name used when equipping (leave blank for default)' },
];

const defaultItemForm = () => ({
  name: '',
  description: '',
  is_starting_item: false,
  can_be_equipped: false,
  item_category: '',
  equipment_slot: '',
});

const showModal = ref(false);
const modalMode = ref('create');
const modalForm = ref(defaultItemForm());
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
  modalForm.value = defaultItemForm();
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

function closeModal() {
  showModal.value = false;
  modalForm.value = defaultItemForm();
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
.edit-link {
  color: var(--color-primary);
  text-decoration: none;
}

.edit-link:hover {
  text-decoration: underline;
}
</style>
