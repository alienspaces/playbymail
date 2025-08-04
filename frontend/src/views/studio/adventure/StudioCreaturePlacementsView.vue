<template>
  <div>
    <div v-if="!selectedGame">
      <p>Please select or create a game to manage creature placements.</p>
    </div>
    <div v-else class="game-table-section">
      <p class="game-context-name">Game: {{ selectedGame.name }}</p>
      <div class="section-header">
        <h2>Creature Placements</h2>
        <button @click="openCreaturePlacementCreate">Create Creature Placement</button>
      </div>
      <ResourceTable
        :columns="creaturePlacementColumns"
        :rows="enhancedCreaturePlacements"
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

      <!-- Create/Edit Creature Placement Modal -->
      <ResourceModalForm
        :visible="showCreaturePlacementModal"
        :mode="creaturePlacementModalMode"
        title="Creature Placement"
        :fields="creaturePlacementFields"
        :modelValue="creaturePlacementModalForm"
        :error="creaturePlacementModalError"
        :options="creaturePlacementOptions"
        @submit="handleCreaturePlacementSubmit"
        @cancel="closeCreaturePlacementModal"
      />

      <!-- Confirm Delete Dialog -->
      <ConfirmationModal
        :visible="showCreaturePlacementDeleteConfirm"
        title="Delete Creature Placement"
        message="Are you sure you want to delete this creature placement?"
        @confirm="deleteCreaturePlacement"
        @cancel="closeCreaturePlacementDelete"
      />
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