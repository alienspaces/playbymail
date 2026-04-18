<!--
  StudioCreaturesView.vue
  This component follows the same pattern as StudioLocationsView.vue and StudioItemsView.vue.
  Added: Creature portrait image upload functionality.
-->
<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage creatures.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Creatures" actionText="Create New Creature" :showIcon="false" titleLevel="h2"
        @action="openCreate" />
      <ResourceTable :columns="columns" :rows="creaturesStore.creatures" :loading="creaturesStore.loading"
        :error="creaturesStore.error">
        <template #cell-name="{ row }">
          <a href="#" class="edit-link" @click.prevent="openEdit(row)">{{ row.name }}</a>
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>
      <TablePagination :pageNumber="creaturesStore.pageNumber" :hasMore="creaturesStore.hasMore"
        @page-change="(p) => creaturesStore.fetchAdventureGameCreatures(selectedGame.id, p)" />
    </div>

    <!-- Custom modal for create/edit with portrait upload support -->
    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal modal-wide">
          <h2>{{ modalMode === 'create' ? 'Create Creature' : 'Edit Creature' }}</h2>
          <form @submit.prevent="handleSubmit(modalForm)" class="modal-form">
            <div class="form-group">
              <label for="creature-name">Name <span class="required">*</span></label>
              <input v-model="modalForm.name" id="creature-name" required maxlength="1024" autocomplete="off" />
            </div>
            <div class="form-group">
              <label for="creature-description">Description <span class="required">*</span></label>
              <textarea v-model="modalForm.description" id="creature-description" required maxlength="4096"
                rows="4"></textarea>
            </div>
            <div class="form-group">
              <label for="creature-disposition">Disposition <span class="required">*</span></label>
              <select v-model="modalForm.disposition" id="creature-disposition" required>
                <option value="aggressive">Aggressive — hostile, blocks location items, free attack on flee</option>
                <option value="inquisitive">Inquisitive — curious, retaliates only if attacked</option>
                <option value="indifferent">Indifferent — ignores player, retaliates only if attacked</option>
              </select>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label for="creature-max-health">Max Health</label>
                <input v-model.number="modalForm.max_health" id="creature-max-health" type="number" min="1"
                  max="9999" />
                <span class="help-text">Maximum hit points</span>
              </div>
              <div class="form-group">
                <label for="creature-attack-damage">Attack Damage</label>
                <input v-model.number="modalForm.attack_damage" id="creature-attack-damage" type="number" min="0"
                  max="999" />
                <span class="help-text">Damage dealt per attack</span>
              </div>
              <div class="form-group">
                <label for="creature-defense">Defense</label>
                <input v-model.number="modalForm.defense" id="creature-defense" type="number" min="0" max="999" />
                <span class="help-text">Damage reduction against player attacks</span>
              </div>
            </div>

            <div class="form-group">
              <label for="creature-attack-method">Attack Method <span class="required">*</span></label>
              <select v-model="modalForm.attack_method" id="creature-attack-method" required>
                <option value="claws">Claws</option>
                <option value="bite">Bite</option>
                <option value="sting">Sting</option>
                <option value="weapon">Weapon</option>
                <option value="spell">Spell</option>
                <option value="slam">Slam</option>
                <option value="touch">Touch</option>
                <option value="breath">Breath</option>
                <option value="gaze">Gaze</option>
              </select>
              <span class="help-text">How the creature physically attacks</span>
            </div>
            <div class="form-group">
              <label for="creature-attack-description">Attack Description</label>
              <input v-model="modalForm.attack_description" id="creature-attack-description" maxlength="512"
                placeholder="e.g. reaches through you with a spectral hand" autocomplete="off" />
              <span class="help-text">Narrative fragment used in combat messages</span>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label for="creature-body-decay-turns">Body Decay Turns</label>
                <input v-model.number="modalForm.body_decay_turns" id="creature-body-decay-turns" type="number" min="0"
                  max="99" />
                <span class="help-text">Turns the body persists after death (0 = instant removal)</span>
              </div>
              <div class="form-group">
                <label for="creature-respawn-turns">Respawn Turns</label>
                <input v-model.number="modalForm.respawn_turns" id="creature-respawn-turns" type="number" min="0"
                  max="999" />
                <span class="help-text">Turns before creature respawns (0 = no respawn)</span>
              </div>
            </div>

            <!-- Portrait Image Upload (only in edit mode) -->
            <div v-if="modalMode === 'edit' && modalForm.id && selectedGame" class="form-section">
              <CreaturePortraitUpload :gameId="selectedGame.id" :creatureId="modalForm.id"
                @imageUpdated="onImageUpdated" @loadingChanged="onImageUploadLoadingChanged" />
            </div>

            <div class="modal-actions">
              <button type="submit" :disabled="imageUploadLoading">
                {{ modalMode === 'create' ? 'Create' : 'Save' }}
              </button>
              <button type="button" @click="closeModal" :disabled="imageUploadLoading">
                {{ imageUploadLoading ? 'Uploading...' : 'Cancel' }}
              </button>
            </div>
          </form>
          <div v-if="modalError" class="error">
            <p>{{ modalError }}</p>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Confirm Delete Dialog -->
    <ConfirmationModal :visible="showDeleteConfirm" title="Delete Creature"
      :message="`Are you sure you want to delete '${deleteTarget?.name}'?`" @confirm="deleteAdventureGameCreature"
      @cancel="closeDelete" />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { useAdventureGameCreaturesStore } from '../../../stores/adventureGameCreatures';
import { useGamesStore } from '../../../stores/games';
import ResourceTable from '../../../components/ResourceTable.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import TableActions from '../../../components/TableActions.vue';
import TablePagination from '../../../components/TablePagination.vue';
import CreaturePortraitUpload from '../../../components/CreaturePortraitUpload.vue';

const creaturesStore = useAdventureGameCreaturesStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'description', label: 'Description' },
  { key: 'disposition', label: 'Disposition' },
  { key: 'max_health', label: 'HP' },
  { key: 'attack_damage', label: 'ATK' },
  { key: 'defense', label: 'DEF' },
  { key: 'created_at', label: 'Created' }
];

const showModal = ref(false);
const modalMode = ref('create');
const defaultCreatureForm = () => ({
  name: '',
  description: '',
  disposition: 'aggressive',
  max_health: 50,
  attack_damage: 10,
  defense: 0,
  attack_method: 'claws',
  attack_description: '',
  body_decay_turns: 3,
  respawn_turns: 0,
});
const modalForm = ref(defaultCreatureForm());
const modalError = ref('');
const showDeleteConfirm = ref(false);
const deleteTarget = ref(null);
const imageUploadLoading = ref(false);

watch(
  () => selectedGame.value,
  (newGame) => {
    if (newGame) {
      creaturesStore.fetchAdventureGameCreatures(newGame.id);
    }
  },
  { immediate: true }
);

function openCreate() {
  modalMode.value = 'create';
  modalForm.value = defaultCreatureForm();
  modalError.value = '';
  showModal.value = true;
}

function openEdit(row) {
  modalMode.value = 'edit';
  modalForm.value = { ...row };
  modalError.value = '';
  showModal.value = true;
}

function closeModal() {
  if (imageUploadLoading.value) return;
  showModal.value = false;
  modalForm.value = defaultCreatureForm();
  modalError.value = '';
}

async function handleSubmit(form) {
  modalError.value = '';
  try {
    if (modalMode.value === 'create') {
      await creaturesStore.createAdventureGameCreature(form);
    } else {
      await creaturesStore.updateAdventureGameCreature(modalForm.value.id, form);
    }
    closeModal();
  } catch (err) {
    modalError.value = err.message || 'Failed to save.';
  }
}

function confirmDelete(row) {
  deleteTarget.value = row;
  showDeleteConfirm.value = true;
}

function closeDelete() {
  showDeleteConfirm.value = false;
  deleteTarget.value = null;
}

async function deleteAdventureGameCreature() {
  if (!deleteTarget.value) return;
  try {
    await creaturesStore.deleteAdventureGameCreature(deleteTarget.value.id);
    closeDelete();
  } catch (err) {
    console.error('Failed to delete creature:', err);
  }
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
      handler: () => confirmDelete(row)
    }
  ];
}

function onImageUpdated() {
  console.log('[StudioCreaturesView] Creature portrait image updated');
}

function onImageUploadLoadingChanged(isLoading) {
  imageUploadLoading.value = isLoading;
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

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: var(--color-bg);
  border-radius: var(--radius-md);
  padding: var(--space-lg);
  width: 500px;
  max-width: 95vw;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: var(--shadow-lg);
}

.modal-wide {
  width: 640px;
}

.modal h2 {
  margin: 0 0 var(--space-md) 0;
  font-size: var(--font-size-lg);
  color: var(--color-text);
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: var(--space-md);
}

.form-group label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-bold);
  color: var(--color-text);
}

.form-group input,
.form-group textarea,
.form-group select {
  padding: var(--space-xs) var(--space-sm);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
  background: var(--color-bg);
  color: var(--color-text);
}

.form-group textarea {
  resize: vertical;
}

.help-text {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}

.required {
  color: var(--color-danger);
}

.form-section {
  border-top: 1px solid var(--color-border);
  padding-top: var(--space-md);
}

.modal-actions {
  display: flex;
  gap: var(--space-sm);
  justify-content: flex-end;
  padding-top: var(--space-sm);
  border-top: 1px solid var(--color-border);
}

.modal-actions button {
  padding: var(--space-xs) var(--space-md);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-bold);
  cursor: pointer;
  border: none;
}

.modal-actions button[type="submit"] {
  background: var(--color-primary);
  color: white;
}

.modal-actions button[type="submit"]:hover:not(:disabled) {
  background: var(--color-primary-dark);
}

.modal-actions button[type="button"] {
  background: var(--color-bg-light);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}

.modal-actions button[type="button"]:hover:not(:disabled) {
  background: var(--color-bg-hover);
}

.modal-actions button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.error {
  margin-top: var(--space-sm);
  padding: var(--space-xs) var(--space-sm);
  background: var(--color-danger-light, #ffebee);
  color: var(--color-danger);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
}
</style>
