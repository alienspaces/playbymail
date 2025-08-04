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
      <div v-if="showCreaturePlacementModal" class="modal-overlay">
        <div class="modal">
          <h2>{{ creaturePlacementModalMode === 'create' ? 'Create Creature Placement' : 'Edit Creature Placement' }}</h2>
          <form @submit.prevent="handleCreaturePlacementSubmit(creaturePlacementModalForm)">
            <div class="form-group">
              <label for="adventure_game_creature_id">Creature:</label>
              <select
                id="adventure_game_creature_id"
                v-model="creaturePlacementModalForm.adventure_game_creature_id"
                required
                class="form-select"
              >
                <option value="" disabled>Select a creature...</option>
                <option v-for="creature in creaturesStore.creatures" :key="creature.id" :value="creature.id">{{ creature.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label for="adventure_game_location_id">Location:</label>
              <select
                id="adventure_game_location_id"
                v-model="creaturePlacementModalForm.adventure_game_location_id"
                required
                class="form-select"
              >
                <option value="" disabled>Select a location...</option>
                <option v-for="loc in locationsStore.locations" :key="loc.id" :value="loc.id">{{ loc.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label for="initial_count">Initial Count:</label>
              <input
                id="initial_count"
                v-model="creaturePlacementModalForm.initial_count"
                type="number"
                min="1"
                required
                class="form-input"
                placeholder="How many creatures to place"
              />
            </div>
            <div class="modal-actions">
              <button type="submit">{{ creaturePlacementModalMode === 'create' ? 'Create' : 'Save' }}</button>
              <button type="button" @click="closeCreaturePlacementModal">Cancel</button>
            </div>
          </form>
          <p v-if="creaturePlacementModalError" class="error">{{ creaturePlacementModalError }}</p>
        </div>
      </div>

      <!-- Confirm Delete Dialog -->
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
import { ref, watch, computed } from 'vue';
import { useCreaturesStore } from '../../../stores/creatures';
import { useLocationsStore } from '../../../stores/locations';
import { useCreaturePlacementsStore } from '../../../stores/creaturePlacements';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';

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

const showCreaturePlacementModal = ref(false);
const creaturePlacementModalMode = ref('create');
const creaturePlacementModalForm = ref({ adventure_game_creature_id: '', adventure_game_location_id: '', initial_count: 1 });
const creaturePlacementModalError = ref('');
const showCreaturePlacementDeleteConfirm = ref(false);
const creaturePlacementDeleteTarget = ref(null);
const creaturePlacementDeleteError = ref('');

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