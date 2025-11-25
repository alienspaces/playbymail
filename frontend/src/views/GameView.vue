<template>
  <div class="game-list">
    <div class="game-table-section">
      <PageHeader title="Games" actionText="Create New Game" :showIcon="false" titleLevel="h2"
        subtitle="Create and manage your adventure games" @action="openCreate" />
      <table v-if="games.length">
        <thead>
          <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Turn Duration</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="game in games" :key="game.id">
            <td>{{ game.name }}</td>
            <td>{{ game.game_type }}</td>
            <td>{{ formatTurnDuration(game.turn_duration_hours) }}</td>
            <td>{{ formatDate(game.created_at) }}</td>
            <td>
              <TableActionsMenu :actions="getActions(game)" />
            </td>
          </tr>
        </tbody>
      </table>
      <p v-else>No games found.</p>
    </div>

    <!-- Modal for create/edit -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <h2>{{ modalMode === 'create' ? 'Create Game' : 'Edit Game' }}</h2>
        <form @submit.prevent="modalMode === 'create' ? createGame() : updateGame()" class="modal-form">
          <div class="form-group">
            <label for="game-name">Name <span class="required">*</span></label>
            <input v-model="modalForm.name" id="game-name" required maxlength="1024" autocomplete="off" />
          </div>
          <div class="form-group">
            <label for="game-type">Type <span class="required">*</span></label>
            <select v-model="modalForm.game_type" id="game-type" required>
              <option value="adventure">Adventure</option>
            </select>
          </div>
          <div class="form-group">
            <label for="turn-duration">Turn Duration (hours) <span class="required">*</span></label>
            <input v-model.number="modalForm.turn_duration_hours" id="turn-duration" type="number" min="1" required
              placeholder="168 (1 week)" />
          </div>
          <div class="modal-actions">
            <button type="submit">{{ modalMode === 'create' ? 'Create' : 'Save' }}</button>
            <button type="button" @click="closeModal">Cancel</button>
          </div>
        </form>
        <div v-if="modalError" class="error">
          <p>{{ modalError }}</p>
        </div>
      </div>
    </div>

    <!-- Confirm delete dialog -->
    <div v-if="showDeleteConfirm" class="modal-overlay" @click.self="closeDelete">
      <div class="modal">
        <h2>Delete Game</h2>
        <p>Are you sure you want to delete <b>{{ deleteTarget?.name }}</b>?</p>
        <div class="modal-actions">
          <button type="button" @click="deleteGame" class="danger-btn">Delete</button>
          <button type="button" @click="closeDelete">Cancel</button>
        </div>
        <div v-if="deleteError" class="error">
          <p>{{ deleteError }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useGamesStore } from '../stores/games';
import PageHeader from '../components/PageHeader.vue';
import TableActionsMenu from '../components/TableActionsMenu.vue';

export default {
  name: 'GameView',
  components: {
    PageHeader,
    TableActionsMenu
  },
  data() {
    return {
      showModal: false,
      modalMode: 'create', // 'create' or 'edit'
      modalForm: {
        id: '',
        name: '',
        game_type: 'adventure',
        turn_duration_hours: 168 // Default to 1 week
      },
      modalError: '',
      showDeleteConfirm: false,
      deleteTarget: null,
      deleteError: ''
    }
  },
  computed: {
    games() {
      return this.gamesStore.games;
    },
    loading() {
      return this.gamesStore.loading;
    },
    error() {
      return this.gamesStore.error;
    }
  },
  created() {
    this.gamesStore = useGamesStore();
    this.gamesStore.fetchGames();
  },
  methods: {
    formatDate(dateStr) {
      if (!dateStr) return ''
      const d = new Date(dateStr)
      return d.toLocaleDateString()
    },
    formatTurnDuration(hours) {
      if (!hours) return 'Not set'
      if (hours % (24 * 7) === 0) {
        const weeks = hours / (24 * 7)
        return `${weeks} week${weeks === 1 ? '' : 's'}`
      }
      if (hours % 24 === 0) {
        const days = hours / 24
        return `${days} day${days === 1 ? '' : 's'}`
      }
      return `${hours} hour${hours === 1 ? '' : 's'}`
    },
    openCreate() {
      this.modalMode = 'create'
      this.modalForm = { id: '', name: '', game_type: 'adventure', turn_duration_hours: 168 }
      this.modalError = ''
      this.showModal = true
    },
    openEdit(game) {
      this.modalMode = 'edit'
      this.modalForm = {
        id: game.id,
        name: game.name,
        game_type: game.game_type,
        turn_duration_hours: game.turn_duration_hours || 168
      }
      this.modalError = ''
      this.showModal = true
    },
    closeModal() {
      this.showModal = false
      this.modalError = ''
    },
    async createGame() {
      this.modalError = ''
      try {
        const created = await this.gamesStore.createGame({
          name: this.modalForm.name,
          game_type: this.modalForm.game_type,
          turn_duration_hours: this.modalForm.turn_duration_hours
        });
        this.closeModal();
        if (created && created.id) {
          this.selectGame(created);
        }
      } catch (err) {
        this.modalError = err.message;
      }
    },
    async updateGame() {
      this.modalError = ''
      try {
        await this.gamesStore.updateGame(this.modalForm.id, {
          name: this.modalForm.name,
          game_type: this.modalForm.game_type,
          turn_duration_hours: this.modalForm.turn_duration_hours
        });
        this.closeModal();
      } catch (err) {
        this.modalError = err.message;
      }
    },
    confirmDelete(game) {
      this.deleteTarget = game
      this.deleteError = ''
      this.showDeleteConfirm = true
    },
    closeDelete() {
      this.showDeleteConfirm = false
      this.deleteTarget = null
      this.deleteError = ''
    },
    async deleteGame() {
      if (!this.deleteTarget) return
      this.deleteError = ''
      try {
        await this.gamesStore.deleteGame(this.deleteTarget.id);
        this.closeDelete();
      } catch (err) {
        this.deleteError = err.message;
      }
    },
    selectGame(game) {
      this.gamesStore.setSelectedGame(game)
      this.$router.push(`/studio/${game.id}/locations`)
    },
    getActions(game) {
      return [
        {
          key: 'edit',
          label: 'Edit',
          handler: () => this.openEdit(game)
        },
        {
          key: 'delete',
          label: 'Delete',
          danger: true,
          handler: () => this.confirmDelete(game)
        },
        {
          key: 'manage',
          label: 'Manage',
          handler: () => this.selectGame(game)
        }
      ];
    }
  }
}
</script>

<style scoped>
.game-list {
  max-width: 1000px;
  margin: 0;
  padding: 0;
}

.game-table-section {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.game-table-section table {
  margin-top: 0;
  /* Remove default table margin-top to match ResourceTable spacing */
}

.game-table-section table th:last-child,
.game-table-section table td:last-child {
  width: 60px;
  text-align: center;
  padding-left: var(--space-sm);
  padding-right: var(--space-sm);
  vertical-align: middle;
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.required {
  color: var(--color-danger);
}

.error {
  color: var(--color-warning-dark);
  background: var(--color-warning-light);
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-warning);
  margin-top: var(--space-md);
}

.error p {
  margin: 0;
}

.danger-btn {
  background: var(--color-danger) !important;
  color: var(--color-text-light) !important;
  border-color: var(--color-danger) !important;
}

.danger-btn:hover {
  background: var(--color-danger-dark) !important;
  border-color: var(--color-danger-dark) !important;
}

@media (max-width: 768px) {
  .modal-actions {
    flex-direction: column-reverse;
  }

  .modal-actions button {
    width: 100%;
  }
}
</style>