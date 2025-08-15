<template>
  <div class="game-list">
    <div class="game-table-section">
      <PageHeader 
        title="Games" 
        actionText="Create New Game" 
        @action="openCreate"
      />
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
              <button @click="openEdit(game)">Edit</button>
              <button @click="confirmDelete(game)">Delete</button>
              <button @click="selectGame(game)">Manage</button>
            </td>
          </tr>
        </tbody>
      </table>
      <p v-else>No games found.</p>
    </div>

    <!-- Modal for create/edit -->
    <div v-if="showModal" class="modal-overlay">
      <div class="modal">
        <h2>{{ modalMode === 'create' ? 'Create Game' : 'Edit Game' }}</h2>
        <form @submit.prevent="modalMode === 'create' ? createGame() : updateGame()">
          <div class="form-group">
            <label for="game-name">Name:</label>
            <input v-model="modalForm.name" id="game-name" required maxlength="1024" autocomplete="off" />
          </div>
          <div class="form-group">
            <label for="game-type">Type:</label>
            <select v-model="modalForm.game_type" id="game-type" required>
              <option value="adventure">Adventure</option>
            </select>
          </div>
          <div class="form-group">
            <label for="turn-duration">Turn Duration (hours):</label>
            <input 
              v-model.number="modalForm.turn_duration_hours" 
              id="turn-duration" 
              type="number" 
              min="1" 
              required 
              placeholder="168 (1 week)"
            />
          </div>
          <div class="modal-actions">
            <button type="submit">{{ modalMode === 'create' ? 'Create' : 'Save' }}</button>
            <button type="button" @click="closeModal">Cancel</button>
          </div>
        </form>
        <p v-if="modalError" class="error">{{ modalError }}</p>
      </div>
    </div>

    <!-- Confirm delete dialog -->
    <div v-if="showDeleteConfirm" class="modal-overlay">
      <div class="modal">
        <h2>Delete Game</h2>
        <p>Are you sure you want to delete <b>{{ deleteTarget?.name }}</b>?</p>
        <div class="modal-actions">
          <button @click="deleteGame">Delete</button>
          <button @click="closeDelete">Cancel</button>
        </div>
        <p v-if="deleteError" class="error">{{ deleteError }}</p>
      </div>
    </div>
  </div>
</template>

<script>
import { useGamesStore } from '../stores/games';
import PageHeader from '../components/PageHeader.vue';

export default {
  name: 'GameView',
  components: {
    PageHeader
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
button {
  margin-right: var(--space-sm);
}
</style> 