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
            <td><a href="#" class="edit-link" @click.prevent="openEdit(game)">{{ game.name }}</a></td>
            <td>{{ game.game_type }}</td>
            <td>{{ formatTurnDuration(game.turn_duration_hours) }}</td>
            <td>{{ formatDate(game.created_at) }}</td>
            <td>
              <TableActions :actions="getActions(game)" />
            </td>
          </tr>
        </tbody>
      </table>
      <p v-else>No games found.</p>
    </div>

    <!-- Modal for create/edit using ResourceModalForm -->
    <ResourceModalForm :visible="showModal" :mode="modalMode" title="Game" :fields="gameFields" :modelValue="modalForm"
      :error="modalError" @submit="handleSubmit" @cancel="closeModal" />

    <!-- Confirm delete dialog -->
    <ConfirmationModal :visible="showDeleteConfirm" title="Delete Game"
      :message="`Are you sure you want to delete '${deleteTarget?.name}'?`" :error="deleteError" @confirm="deleteGame"
      @cancel="closeDelete" />
  </div>
</template>

<script>
import { useGamesStore } from '../stores/games';
import PageHeader from '../components/PageHeader.vue';
import TableActions from '../components/TableActions.vue';
import ResourceModalForm from '../components/ResourceModalForm.vue';
import ConfirmationModal from '../components/ConfirmationModal.vue';

export default {
  name: 'GameView',
  components: {
    PageHeader,
    TableActions,
    ResourceModalForm,
    ConfirmationModal
  },
  data() {
    return {
      showModal: false,
      modalMode: 'create', // 'create' or 'edit'
      modalForm: {
        id: '',
        name: '',
        game_type: 'adventure',
        turn_duration_hours: 168, // Default to 1 week
        description: ''
      },
      modalError: '',
      showDeleteConfirm: false,
      deleteTarget: null,
      deleteError: '',
      gameFields: [
        { key: 'name', label: 'Name', required: true, maxlength: 1024 },
        { key: 'game_type', label: 'Type', required: true, type: 'select', options: [{ value: 'adventure', label: 'Adventure' }] },
        { key: 'turn_duration_hours', label: 'Turn Duration (hours)', required: true, type: 'number', min: 1, placeholder: '168 (1 week)' },
        { key: 'description', label: 'Description', required: true, maxlength: 512, type: 'textarea', rows: 4, placeholder: 'Game description that appears on the join game turn sheet' }
      ]
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
    // Filter to only show games where the user has Designer subscription
    this.gamesStore.fetchGames({ subscriptionType: 'Designer' });
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
      this.modalForm = { id: '', name: '', game_type: 'adventure', turn_duration_hours: 168, description: '' }
      this.modalError = ''
      this.showModal = true
    },
    openEdit(game) {
      this.modalMode = 'edit'
      this.modalForm = {
        id: game.id,
        name: game.name,
        game_type: game.game_type,
        turn_duration_hours: game.turn_duration_hours || 168,
        description: game.description || ''
      }
      this.modalError = ''
      this.showModal = true
    },
    closeModal() {
      this.showModal = false
      this.modalError = ''
    },
    async handleSubmit(formData) {
      this.modalError = ''
      try {
        const data = {
          name: formData.name,
          game_type: formData.game_type,
          turn_duration_hours: typeof formData.turn_duration_hours === 'string' ? parseInt(formData.turn_duration_hours, 10) : formData.turn_duration_hours,
          description: formData.description.trim()
        }
        if (this.modalMode === 'create') {
          const created = await this.gamesStore.createGame(data)
          this.closeModal()
          if (created && created.id) {
            this.selectGame(created)
          }
        } else {
          await this.gamesStore.updateGame(this.modalForm.id, data)
          this.closeModal()
        }
      } catch (err) {
        this.modalError = err.message || 'Failed to save.'
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
      this.$router.push(`/studio/${game.id}/turn-sheet-backgrounds`)
    },
    getActions(game) {
      return [
        {
          key: 'manage',
          label: 'Manage',
          primary: true,
          handler: () => this.selectGame(game)
        },
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
  width: auto;
  min-width: 180px;
  text-align: center;
  padding-left: var(--space-sm);
  padding-right: var(--space-sm);
  vertical-align: middle;
}

.edit-link {
  color: var(--color-primary);
  text-decoration: none;
}

.edit-link:hover {
  text-decoration: underline;
}
</style>