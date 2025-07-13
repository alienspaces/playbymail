<template>
  <div class="game-list">
    <h1>Games</h1>
    <button @click="openCreate">Create New Game</button>
    <table v-if="games.length">
      <thead>
        <tr>
          <th>Name</th>
          <th>Type</th>
          <th>Created</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="game in games" :key="game.id">
          <td>{{ game.name }}</td>
          <td>{{ game.game_type }}</td>
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

    <!-- Modal for create/edit -->
    <div v-if="showModal" class="modal-overlay">
      <div class="modal">
        <h2>{{ modalMode === 'create' ? 'Create Game' : 'Edit Game' }}</h2>
        <form @submit.prevent="modalMode === 'create' ? createGame() : updateGame()">
          <label>
            Name:
            <input v-model="modalForm.name" required maxlength="1024" />
          </label>
          <label>
            Type:
            <select v-model="modalForm.game_type" required>
              <option value="adventure">Adventure</option>
            </select>
          </label>
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

export default {
  name: 'GameView',
  data() {
    return {
      games: [],
      loading: false,
      error: null,
      showModal: false,
      modalMode: 'create', // 'create' or 'edit'
      modalForm: {
        id: '',
        name: '',
        game_type: 'adventure'
      },
      modalError: '',
      showDeleteConfirm: false,
      deleteTarget: null,
      deleteError: ''
    }
  },
  mounted() {
    this.fetchGames()
  },
  methods: {
    async fetchGames() {
      this.loading = true
      this.error = null
      try {
        const res = await fetch('/v1/games', {
          headers: {
            'Content-Type': 'application/json',
            // TODO: Add authentication header if needed
          }
        })
        if (!res.ok) throw new Error('Failed to fetch games')
        const data = await res.json()
        this.games = data.data || []
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },
    formatDate(dateStr) {
      if (!dateStr) return ''
      const d = new Date(dateStr)
      return d.toLocaleDateString()
    },
    openCreate() {
      this.modalMode = 'create'
      this.modalForm = { id: '', name: '', game_type: 'adventure' }
      this.modalError = ''
      this.showModal = true
    },
    openEdit(game) {
      this.modalMode = 'edit'
      this.modalForm = { id: game.id, name: game.name, game_type: game.game_type }
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
        const res = await fetch('/v1/games', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ name: this.modalForm.name, game_type: this.modalForm.game_type })
        })
        if (!res.ok) throw new Error('Failed to create game')
        const created = await res.json()
        this.closeModal()
        await this.fetchGames()
        // Automatically select and redirect to the new game
        if (created && created.data && created.data.id) {
          this.selectGame(created.data)
        }
      } catch (err) {
        this.modalError = err.message
      }
    },
    async updateGame() {
      this.modalError = ''
      try {
        const res = await fetch(`/v1/games/${this.modalForm.id}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ name: this.modalForm.name, game_type: this.modalForm.game_type })
        })
        if (!res.ok) throw new Error('Failed to update game')
        this.closeModal()
        await this.fetchGames()
      } catch (err) {
        this.modalError = err.message
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
        const res = await fetch(`/v1/games/${this.deleteTarget.id}`, {
          method: 'DELETE',
          headers: { 'Content-Type': 'application/json' }
        })
        if (!res.ok) throw new Error('Failed to delete game')
        this.closeDelete()
        await this.fetchGames()
      } catch (err) {
        this.deleteError = err.message
      }
    },
    selectGame(game) {
      // Set selected game in Pinia store and redirect
      const gamesStore = useGamesStore()
      gamesStore.setSelectedGame(game)
      this.$router.push(`/studio/${game.id}/locations`)
    }
  }
}
</script>

<style scoped>
.game-list {
  max-width: 800px;
  margin: 2rem auto;
}
table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
}
th, td {
  border: 1px solid #ccc;
  padding: 0.5rem 1rem;
  text-align: left;
}
th {
  background: #f8f8f8;
}
button {
  margin-right: 0.5rem;
}
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal {
  background: #fff;
  padding: 2rem;
  border-radius: 8px;
  min-width: 300px;
  max-width: 90vw;
  box-shadow: 0 2px 16px rgba(0,0,0,0.2);
}
.modal-actions {
  margin-top: 1rem;
}
.error {
  color: #b00;
  margin-top: 1rem;
}
</style> 