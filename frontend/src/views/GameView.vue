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
            <th>Status</th>
            <th>Turn Duration</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="game in games" :key="game.id">
            <td><a href="#" class="edit-link" @click.prevent="openEdit(game)">{{ game.name }}</a></td>
            <td>{{ game.game_type }}</td>
            <td>{{ game.status || 'draft' }}</td>
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

    <!-- Confirm publish dialog -->
    <ConfirmationModal :visible="showPublishConfirm" title="Publish Game"
      :message="`Are you sure you want to publish '${publishTarget?.name}'? Once published, the game cannot be modified.`"
      warning="Published games are immutable. To make changes, you'll need to create a new version."
      :error="publishError" @confirm="publishGame" @cancel="closePublish" />
  </div>
</template>

<script>
import { useGamesStore } from '../stores/games';
import { publishGame as apiPublishGame } from '../api/games';
import { getMyAccountSubscriptions } from '../api/accountSubscriptions';
import { listGames } from '../api/games';
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
      showPublishConfirm: false,
      publishTarget: null,
      publishError: '',
      accountSubscriptions: [],
      gameCount: 0,
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
  async created() {
    this.gamesStore = useGamesStore();
    // Filter to only show games where the user has Designer subscription
    await this.gamesStore.fetchGames({ subscriptionType: 'Designer' });
    await this.loadSubscriptionInfo();
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
    async openCreate() {
      // Check subscription limits before allowing creation
      await this.loadSubscriptionInfo();
      const canCreate = this.checkCanCreateGame();
      if (!canCreate.canCreate) {
        this.modalError = canCreate.error;
        return;
      }

      this.modalMode = 'create'
      this.modalForm = { id: '', name: '', game_type: 'adventure', turn_duration_hours: 168, description: '' }
      this.modalError = ''
      this.showModal = true
    },
    async loadSubscriptionInfo() {
      try {
        const accountSubsResponse = await getMyAccountSubscriptions();
        this.accountSubscriptions = accountSubsResponse.data || [];
        const gamesResponse = await listGames();
        this.gameCount = (gamesResponse.data || []).length;
      } catch (err) {
        console.error('Failed to load subscription info:', err);
      }
    },
    checkCanCreateGame() {
      const basicSub = this.accountSubscriptions.find(sub =>
        sub.subscription_type === 'basic_game_designer' && sub.status === 'active'
      );
      const professionalSub = this.accountSubscriptions.find(sub =>
        sub.subscription_type === 'professional_game_designer' && sub.status === 'active'
      );

      if (professionalSub) {
        return { canCreate: true };
      }

      if (basicSub) {
        if (this.gameCount >= 10) {
          return {
            canCreate: false,
            error: `You have reached the limit of 10 games for Basic Game Designer subscription. You currently have ${this.gameCount} games. Upgrade to Professional Game Designer to create unlimited games.`
          };
        }
        return { canCreate: true };
      }

      return {
        canCreate: false,
        error: 'You need an active game designer subscription to create games.'
      };
    },
    openEdit(game) {
      // Prevent editing published games
      if (game.status === 'published') {
        this.modalError = 'Published games cannot be modified. Create a new version to make changes.';
        return;
      }

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
      const actions = [
        {
          key: 'manage',
          label: 'Manage',
          primary: true,
          handler: () => this.selectGame(game)
        }
      ];

      // Only allow edit/delete for draft games
      if (game.status !== 'published') {
        actions.push({
          key: 'edit',
          label: 'Edit',
          handler: () => this.openEdit(game)
        });
        actions.push({
          key: 'delete',
          label: 'Delete',
          danger: true,
          handler: () => this.confirmDelete(game)
        });
      }

      // Add publish button for draft games
      if (game.status === 'draft' || !game.status) {
        actions.push({
          key: 'publish',
          label: 'Publish',
          handler: () => this.confirmPublish(game)
        });
      }

      return actions;
    },
    confirmPublish(game) {
      this.publishTarget = game
      this.publishError = ''
      this.showPublishConfirm = true
    },
    closePublish() {
      this.showPublishConfirm = false
      this.publishTarget = null
      this.publishError = ''
    },
    async publishGame() {
      if (!this.publishTarget) return
      this.publishError = ''
      try {
        const res = await apiPublishGame(this.publishTarget.id)
        if (res.data) {
          // Update game in store
          const idx = this.games.findIndex(g => g.id === this.publishTarget.id)
          if (idx !== -1) {
            this.games[idx] = res.data
          }
        }
        this.closePublish()
      } catch (err) {
        this.publishError = err.message || 'Failed to publish game'
        console.error('Error publishing game:', err)
      }
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