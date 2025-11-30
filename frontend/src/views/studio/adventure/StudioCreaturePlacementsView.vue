<template>
  <div>
    <div v-if="!selectedGame">
      <p>Please select or create a game to manage creature placements.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Creature Placements" actionText="Create Creature Placement" :showIcon="false" titleLevel="h2"
        @action="openCreaturePlacementCreate" />
      <ResourceTable :columns="creaturePlacementColumns" :rows="enhancedCreaturePlacements"
        :loading="creaturePlacementsStore.loading" :error="creaturePlacementsStore.error">
        <template #cell-creature_name="{ row }">
          <a href="#" class="edit-link" @click.prevent="openCreaturePlacementEdit(row)">{{ row.creature_name }}</a>
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>

      <!-- Create/Edit Creature Placement Modal -->
      <ResourceModalForm :visible="showCreaturePlacementModal" :mode="creaturePlacementModalMode"
        title="Creature Placement" :fields="creaturePlacementFields" :modelValue="creaturePlacementModalForm"
        :error="creaturePlacementModalError" :options="creaturePlacementOptions" @submit="handleCreaturePlacementSubmit"
        @cancel="closeCreaturePlacementModal" />

      <!-- Confirm Delete Dialog -->
      <ConfirmationModal :visible="showCreaturePlacementDeleteConfirm" title="Delete Creature Placement"
        message="Are you sure you want to delete this creature placement?" @confirm="deleteCreaturePlacement"
        @cancel="closeCreaturePlacementDelete" />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { useCreaturesStore } from '../../../stores/creatures';
import { useLocationsStore } from '../../../stores/locations';
import { useCreaturePlacementsStore } from '../../../stores/creaturePlacements';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import ResourceModalForm from '../../../components/ResourceModalForm.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import TableActions from '../../../components/TableActions.vue';

const creaturesStore = useCreaturesStore();
const locationsStore = useLocationsStore();
const creaturePlacementsStore = useCreaturePlacementsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

// Enhance creature placements with names for display
const enhancedCreaturePlacements = computed(() => {
  return creaturePlacementsStore.creaturePlacements.map(placement => {
    const creature = creaturesStore.creatures.find(creature => creature.id === placement.adventure_game_creature_id);
    const location = locationsStore.locations.find(loc => loc.id === placement.adventure_game_location_id);
    return {
      ...placement,
      creature_name: creature?.name || 'Unknown Creature',
      location_name: location?.name || 'Unknown Location'
    };
  });
});

const creaturePlacementColumns = [
  { key: 'creature_name', label: 'Creature' },
  { key: 'location_name', label: 'Location' },
  { key: 'initial_count', label: 'Count' },
  { key: 'created_at', label: 'Created' }
];

// Field configuration for ResourceModalForm
const creaturePlacementFields = [
  { key: 'adventure_game_creature_id', label: 'Creature', type: 'select', required: true, placeholder: 'Select a creature...' },
  { key: 'adventure_game_location_id', label: 'Location', type: 'select', required: true, placeholder: 'Select a location...' },
  { key: 'initial_count', label: 'Initial Count', type: 'number', required: true, min: 1 }
];

// Options for select fields
const creaturePlacementOptions = computed(() => ({
  adventure_game_creature_id: creaturesStore.creatures.map(creature => ({
    value: creature.id,
    label: creature.name
  })),
  adventure_game_location_id: locationsStore.locations.map(location => ({
    value: location.id,
    label: location.name
  }))
}));

const showCreaturePlacementModal = ref(false);
const creaturePlacementModalMode = ref('create');
const creaturePlacementModalForm = ref({ adventure_game_creature_id: '', adventure_game_location_id: '', initial_count: 1 });
const creaturePlacementModalError = ref('');
const showCreaturePlacementDeleteConfirm = ref(false);
const creaturePlacementDeleteTarget = ref(null);

// Watch for game selection changes
watch(
  () => selectedGame.value,
  (newGame) => {
    if (newGame) {
      creaturesStore.fetchCreatures(newGame.id);
      locationsStore.fetchLocations(newGame.id);
      creaturePlacementsStore.fetchCreaturePlacements(newGame.id);
    }
  },
  { immediate: true }
);

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
  showCreaturePlacementDeleteConfirm.value = true;
}

function closeCreaturePlacementDelete() {
  showCreaturePlacementDeleteConfirm.value = false;
  creaturePlacementDeleteTarget.value = null;
}

async function deleteCreaturePlacement() {
  if (!creaturePlacementDeleteTarget.value) return;
  try {
    await creaturePlacementsStore.deleteCreaturePlacement(creaturePlacementDeleteTarget.value.id);
    closeCreaturePlacementDelete();
  } catch (err) {
    console.error('Failed to delete creature placement:', err);
  }
}

function getActions(row) {
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

<style scoped>
.edit-link {
  color: var(--color-primary);
  text-decoration: none;
}

.edit-link:hover {
  text-decoration: underline;
}
</style>
