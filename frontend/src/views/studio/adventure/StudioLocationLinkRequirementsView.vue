<!--
  StudioLocationLinkRequirementsView.vue
  Manages link requirements (visibility and traversal conditions) for location links.
-->
<template>
  <div>
    <div v-if="!selectedGame">
      <p>Select a game to manage location link requirements.</p>
    </div>
    <div v-else class="game-table-section">
      <GameContext :gameName="selectedGame.name" />
      <PageHeader title="Location Link Requirements" actionText="Create New Requirement" :showIcon="false"
        titleLevel="h2" @action="openCreate" />
      <ResourceTable :columns="columns" :rows="enhancedRequirements" :loading="requirementsStore.loading"
        :error="requirementsStore.error">
        <template #cell-link_name="{ row }">
          <a href="#" class="edit-link" @click.prevent="openEdit(row)">{{ row.link_name }}</a>
        </template>
        <template #actions="{ row }">
          <TableActions :actions="getActions(row)" />
        </template>
      </ResourceTable>
      <TablePagination :pageNumber="requirementsStore.pageNumber" :hasMore="requirementsStore.hasMore"
        @page-change="(p) => requirementsStore.fetchAdventureGameLocationLinkRequirements(selectedGame.id, p)" />

      <!-- Create/Edit Requirement Modal (custom — needs conditional fields) -->
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal-panel">
          <h3>{{ modalMode === 'create' ? 'Create' : 'Edit' }} Location Link Requirement</h3>
          <form @submit.prevent="handleSubmit">
            <div class="form-field">
              <label>Location Link</label>
              <select v-model="form.game_location_link_id" required>
                <option value="">Select a link...</option>
                <option v-for="link in locationLinksStore.locationLinks" :key="link.id" :value="link.id">
                  {{ getLinkLabel(link) }}
                </option>
              </select>
            </div>

            <div class="form-field">
              <label>Purpose</label>
              <select v-model="form.purpose" required>
                <option value="traverse">Traverse — required to walk through</option>
                <option value="visible">Visible — required to see the link at all</option>
              </select>
            </div>

            <div class="form-field">
              <label>Target Type</label>
              <select v-model="form.target_type" required @change="onTargetTypeChange">
                <option value="item">Item</option>
                <option value="creature">Creature</option>
              </select>
            </div>

            <div v-if="form.target_type === 'item'" class="form-field">
              <label>Item</label>
              <select v-model="form.game_item_id" required>
                <option value="">Select an item...</option>
                <option v-for="item in itemsStore.items" :key="item.id" :value="item.id">
                  {{ item.name }}
                </option>
              </select>
            </div>

            <div v-if="form.target_type === 'creature'" class="form-field">
              <label>Creature</label>
              <select v-model="form.game_creature_id" required>
                <option value="">Select a creature...</option>
                <option v-for="creature in creaturesStore.creatures" :key="creature.id" :value="creature.id">
                  {{ creature.name }}
                </option>
              </select>
            </div>

            <div class="form-field">
              <label>Condition</label>
              <select v-model="form.condition" required>
                <template v-if="form.target_type === 'item'">
                  <option value="in_inventory">In inventory (character holds required quantity)</option>
                  <option value="equipped">Equipped (character has item equipped)</option>
                </template>
                <template v-else>
                  <option value="dead_at_location">Dead at location (N creatures dead here)</option>
                  <option value="none_alive_at_location">None alive at location (no living instances here)</option>
                  <option value="none_alive_in_game">None alive in game (all instances dead)</option>
                </template>
              </select>
            </div>

            <div class="form-field">
              <label>Quantity</label>
              <input type="number" v-model.number="form.quantity" min="1" required />
              <span class="form-field-hint">Minimum number required (used for in_inventory and dead_at_location conditions)</span>
            </div>

            <div v-if="modalError" class="modal-error">{{ modalError }}</div>

            <div class="modal-actions">
              <button type="button" class="btn-secondary" @click="closeModal">Cancel</button>
              <button type="submit" class="btn-primary">{{ modalMode === 'create' ? 'Create' : 'Save' }}</button>
            </div>
          </form>
        </div>
      </div>

      <!-- Confirm Delete Dialog -->
      <ConfirmationModal :visible="showDeleteConfirm" title="Delete Requirement"
        message="Are you sure you want to delete this link requirement?" @confirm="deleteRequirement"
        @cancel="closeDelete" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useAdventureGameLocationLinkRequirementsStore } from '../../../stores/adventureGameLocationLinkRequirements';
import { useAdventureGameLocationLinksStore } from '../../../stores/adventureGameLocationLinks';
import { useAdventureGameItemsStore } from '../../../stores/adventureGameItems';
import { useAdventureGameCreaturesStore } from '../../../stores/adventureGameCreatures';
import { useAdventureGameLocationsStore } from '../../../stores/adventureGameLocations';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';
import ResourceTable from '../../../components/ResourceTable.vue';
import ConfirmationModal from '../../../components/ConfirmationModal.vue';
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import TableActions from '../../../components/TableActions.vue';
import TablePagination from '../../../components/TablePagination.vue';

const requirementsStore = useAdventureGameLocationLinkRequirementsStore();
const locationLinksStore = useAdventureGameLocationLinksStore();
const itemsStore = useAdventureGameItemsStore();
const creaturesStore = useAdventureGameCreaturesStore();
const locationsStore = useAdventureGameLocationsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const columns = [
  { key: 'link_name', label: 'Location Link' },
  { key: 'purpose', label: 'Purpose' },
  { key: 'target_name', label: 'Target' },
  { key: 'condition', label: 'Condition' },
  { key: 'quantity', label: 'Quantity' },
  { key: 'created_at', label: 'Created' }
];

const showModal = ref(false);
const modalMode = ref('create');
const modalError = ref('');
const showDeleteConfirm = ref(false);
const deleteTarget = ref(null);

const defaultForm = () => ({
  id: '',
  game_location_link_id: '',
  purpose: 'traverse',
  target_type: 'item',
  game_item_id: '',
  game_creature_id: '',
  condition: 'in_inventory',
  quantity: 1
});

const form = ref(defaultForm());

watch(
  () => selectedGame.value,
  (newGame) => {
    if (newGame) {
      requirementsStore.fetchAdventureGameLocationLinkRequirements(newGame.id);
      locationLinksStore.fetchAdventureGameLocationLinks(newGame.id);
      itemsStore.fetchAdventureGameItems(newGame.id);
      creaturesStore.fetchAdventureGameCreatures(newGame.id);
      locationsStore.fetchAdventureGameLocations(newGame.id);
    }
  },
  { immediate: true }
);

function getLinkLabel(link) {
  const from = locationsStore.locations.find(l => l.id === link.from_adventure_game_location_id);
  const to = locationsStore.locations.find(l => l.id === link.to_adventure_game_location_id);
  const fromName = from?.name || '?';
  const toName = to?.name || '?';
  return `${link.name} (${fromName} → ${toName})`;
}

const enhancedRequirements = computed(() => {
  return requirementsStore.requirements.map(req => {
    const link = locationLinksStore.locationLinks.find(l => l.id === req.game_location_link_id);
    const item = itemsStore.items.find(i => i.id === req.game_item_id);
    const creature = creaturesStore.creatures.find(c => c.id === req.game_creature_id);
    return {
      ...req,
      link_name: link ? getLinkLabel(link) : (req.game_location_link_id || 'Unknown'),
      target_name: item?.name || creature?.name || '—'
    };
  });
});

function onTargetTypeChange() {
  form.value.game_item_id = '';
  form.value.game_creature_id = '';
  form.value.condition = form.value.target_type === 'item' ? 'in_inventory' : 'dead_at_location';
}

function openCreate() {
  modalMode.value = 'create';
  form.value = defaultForm();
  modalError.value = '';
  showModal.value = true;
}

function openEdit(row) {
  modalMode.value = 'edit';
  const targetType = row.game_creature_id ? 'creature' : 'item';
  form.value = {
    id: row.id,
    game_location_link_id: row.game_location_link_id,
    purpose: row.purpose,
    target_type: targetType,
    game_item_id: row.game_item_id || '',
    game_creature_id: row.game_creature_id || '',
    condition: row.condition,
    quantity: row.quantity
  };
  modalError.value = '';
  showModal.value = true;
}

function closeModal() {
  showModal.value = false;
  modalError.value = '';
}

async function handleSubmit() {
  modalError.value = '';
  const payload = {
    game_location_link_id: form.value.game_location_link_id,
    purpose: form.value.purpose,
    condition: form.value.condition,
    quantity: form.value.quantity
  };
  if (form.value.target_type === 'item') {
    payload.game_item_id = form.value.game_item_id;
  } else {
    payload.game_creature_id = form.value.game_creature_id;
  }

  try {
    if (modalMode.value === 'create') {
      await requirementsStore.createAdventureGameLocationLinkRequirement(payload);
    } else {
      await requirementsStore.updateAdventureGameLocationLinkRequirement(form.value.id, payload);
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

async function deleteRequirement() {
  if (!deleteTarget.value) return;
  try {
    await requirementsStore.deleteAdventureGameLocationLinkRequirement(deleteTarget.value.id);
    closeDelete();
  } catch (err) {
    console.error('Failed to delete location link requirement:', err);
  }
}

function getActions(row) {
  return [
    { key: 'edit', label: 'Edit', handler: () => openEdit(row) },
    { key: 'delete', label: 'Delete', danger: true, handler: () => confirmDelete(row) }
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

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-panel {
  background: var(--color-bg);
  border-radius: var(--radius-md);
  padding: var(--space-lg);
  width: 100%;
  max-width: 540px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-panel h3 {
  margin: 0 0 var(--space-md);
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
  margin-bottom: var(--space-md);
}

.form-field label {
  font-weight: 500;
  font-size: var(--font-size-sm);
}

.form-field select,
.form-field input[type="number"] {
  padding: var(--space-sm);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg);
  color: var(--color-text);
  font-size: var(--font-size-base);
}

.form-field-hint {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}

.modal-error {
  color: var(--color-danger);
  font-size: var(--font-size-sm);
  margin-bottom: var(--space-sm);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-sm);
  margin-top: var(--space-md);
}

.btn-primary {
  background: var(--color-primary);
  color: white;
  border: none;
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--font-size-base);
}

.btn-secondary {
  background: var(--color-bg-light);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--font-size-base);
}
</style>
