<!--
  StudioLocationLinksView.vue
  This component manages location links between locations in a game.
-->
<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage location links.</p>
    </div>
    <div v-else class="game-table-section">
      <p class="game-context-name">Game: {{ selectedGame.name }}</p>
      <div class="section-header">
        <h2>Location Links</h2>
        <button @click="openCreate">Create New Location Link</button>
      </div>
      <ResourceTable
        :columns="columns"
        :rows="enhancedLocationLinks"
        :loading="locationLinksStore.loading"
        :error="locationLinksStore.error"
      >
        <template #actions="{ row }">
          <button @click="openEdit(row)">Edit</button>
          <button @click="confirmDelete(row)">Delete</button>
        </template>
      </ResourceTable>
      
      <!-- Custom Modal for Location Links -->
      <div v-if="showModal" class="modal-overlay">
        <div class="modal">
          <h2>{{ modalMode === 'create' ? 'Create' : 'Edit' }} Location Link</h2>
          <form @submit.prevent="handleSubmit" class="form">
            <div class="form-group">
              <label for="fromLocation">From Location *</label>
              <select 
                id="fromLocation" 
                v-model="modalForm.from_adventure_game_location_id" 
                required
                class="form-control"
              >
                <option value="">Select a location...</option>
                <option 
                  v-for="location in locationsStore.locations" 
                  :key="location.id" 
                  :value="location.id"
                >
                  {{ location.name }}
                </option>
              </select>
            </div>
            
            <div class="form-group">
              <label for="toLocation">To Location *</label>
              <select 
                id="toLocation" 
                v-model="modalForm.to_adventure_game_location_id" 
                required
                class="form-control"
              >
                <option value="">Select a location...</option>
                <option 
                  v-for="location in locationsStore.locations" 
                  :key="location.id" 
                  :value="location.id"
                >
                  {{ location.name }}
                </option>
              </select>
            </div>
            
            <div class="form-group">
              <label for="name">Link Name *</label>
              <input 
                id="name" 
                v-model="modalForm.name" 
                type="text" 
                required 
                maxlength="64"
                class="form-control"
                placeholder="e.g., North Path, Secret Door"
              />
            </div>
            
            <div class="form-group">
              <label for="description">Description *</label>
              <textarea 
                id="description" 
                v-model="modalForm.description" 
                maxlength="255"
                required
                class="form-control"
                placeholder="Describe the link between locations..."
                rows="3"
              ></textarea>
            </div>
            
            <div class="modal-actions">
              <button type="submit" :disabled="submitting">
                {{ submitting ? 'Saving...' : (modalMode === 'create' ? 'Create' : 'Update') }}
              </button>
              <button type="button" @click="closeModal">Cancel</button>
            </div>
            
            <p v-if="modalError" class="error">{{ modalError }}</p>
          </form>
        </div>
      </div>
      
      <div v-if="showDeleteConfirm" class="modal-overlay">
        <div class="modal">
          <h2>Delete Location Link</h2>
          <p>Are you sure you want to delete the link <b>"{{ deleteTarget?.name }}"</b>?</p>
          <div class="modal-actions">
            <button @click="deleteLocationLink">Delete</button>
            <button @click="closeDelete">Cancel</button>
          </div>
          <p v-if="deleteError" class="error">{{ deleteError }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useLocationLinksStore } from '../../../stores/locationLinks';
import { useLocationsStore } from '../../../stores/locations';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';

const locationLinksStore = useLocationLinksStore();
const locationsStore = useLocationsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const columns = [
  { key: 'name', label: 'Link Name' },
  { key: 'from_location_name', label: 'From Location' },
  { key: 'to_location_name', label: 'To Location' },
  { key: 'description', label: 'Description' },
  { key: 'created_at', label: 'Created' }
];

const showModal = ref(false);
const modalMode = ref('create');
const modalForm = ref({});
const modalError = ref('');
const submitting = ref(false);
const showDeleteConfirm = ref(false);
const deleteTarget = ref(null);
const deleteError = ref('');

// Watch for game selection changes
watch(
  () => selectedGame.value,
  (newGame) => {
    if (newGame) {
      locationLinksStore.fetchLocationLinks(newGame.id);
      locationsStore.fetchLocations(newGame.id);
    }
  },
  { immediate: true }
);

// Enhance location links with location names for display
const enhancedLocationLinks = computed(() => {
  return locationLinksStore.locationLinks.map(link => {
    const fromLocation = locationsStore.locations.find(loc => loc.id === link.from_adventure_game_location_id);
    const toLocation = locationsStore.locations.find(loc => loc.id === link.to_adventure_game_location_id);
    return {
      ...link,
      from_location_name: fromLocation?.name || 'Unknown',
      to_location_name: toLocation?.name || 'Unknown'
    };
  });
});

function openCreate() {
  modalMode.value = 'create';
  modalForm.value = { 
    from_adventure_game_location_id: '', 
    to_adventure_game_location_id: '', 
    name: '', 
    description: '' 
  };
  modalError.value = '';
  showModal.value = true;
}

function openEdit(row) {
  modalMode.value = 'edit';
  modalForm.value = { 
    id: row.id,
    from_adventure_game_location_id: row.from_adventure_game_location_id, 
    to_adventure_game_location_id: row.to_adventure_game_location_id, 
    name: row.name, 
    description: row.description || '' 
  };
  modalError.value = '';
  showModal.value = true;
}

function closeModal() {
  showModal.value = false;
  modalError.value = '';
  submitting.value = false;
}

async function handleSubmit() {
  modalError.value = '';
  submitting.value = true;
  
  try {
    if (modalMode.value === 'create') {
      await locationLinksStore.createLocationLink(modalForm.value);
    } else {
      await locationLinksStore.updateLocationLink(modalForm.value.id, modalForm.value);
    }
    closeModal();
  } catch (err) {
    modalError.value = err.message || 'Failed to save.';
  } finally {
    submitting.value = false;
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

async function deleteLocationLink() {
  if (!deleteTarget.value) return;
  deleteError.value = '';
  try {
    await locationLinksStore.deleteLocationLink(deleteTarget.value.id);
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

.form {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.form-group label {
  font-weight: 500;
  color: var(--color-text);
}

.form-control {
  padding: var(--space-sm);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-md);
  background: var(--color-bg);
  color: var(--color-text);
}

.form-control:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px rgba(var(--color-primary-rgb), 0.2);
}

select.form-control {
  cursor: pointer;
}

textarea.form-control {
  resize: vertical;
  min-height: 80px;
}
</style> 