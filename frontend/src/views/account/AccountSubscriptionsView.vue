<template>
  <div class="account-subscriptions-view">
    <PageHeader title="Subscriptions" :showIcon="false" titleLevel="h2"
      subtitle="Manage your account and game subscriptions" />

    <!-- Loading state -->
    <div v-if="loading" class="loading-state">
      <p>Loading subscriptions...</p>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
      <AppButton @click="loadSubscriptions" variant="primary" size="small">
        Retry
      </AppButton>
    </div>

    <!-- Subscriptions content -->
    <div v-else class="subscriptions-content">
      <!-- Account Subscriptions Section -->
      <div class="subscriptions-section">
        <h3>Account Subscriptions</h3>
        <p class="section-description">Subscriptions that grant permissions for game design capabilities.</p>

        <div v-if="accountSubscriptions.length > 0" class="subscriptions-grid">
          <DataCard v-for="subscription in accountSubscriptions" :key="subscription.id"
            :title="formatSubscriptionType(subscription.subscription_type)" class="subscription-card">
            <div class="subscription-info">
              <DataItem label="Type" :value="formatSubscriptionType(subscription.subscription_type)" />
              <DataItem label="Status" :value="subscription.status" />
              <DataItem label="Auto Renew" :value="subscription.auto_renew ? 'Yes' : 'No'" />
              <DataItem label="Created" :value="formatDate(subscription.created_at)" />
              <DataItem v-if="subscription.subscription_type === 'basic_game_designer'" label="Game Limit"
                :value="`${gameCount} / 10 games`" />
              <DataItem v-else-if="subscription.subscription_type === 'professional_game_designer'" label="Game Limit"
                value="Unlimited" />
            </div>
          </DataCard>
        </div>

        <div v-else class="empty-state">
          <p>No account subscriptions found.</p>
        </div>
      </div>

      <!-- Game Subscriptions Section -->
      <div class="subscriptions-section">
        <h3>Game Subscriptions</h3>
        <p class="section-description">Subscriptions to specific games for playing or managing.</p>

        <div v-if="gameSubscriptions.length > 0" class="subscriptions-grid">
          <DataCard v-for="subscription in gameSubscriptions" :key="subscription.id"
            :title="subscription.game_name || 'Unknown Game'" class="subscription-card">
            <div class="subscription-info">
              <DataItem label="Game" :value="subscription.game_name || 'Unknown Game'" />
              <DataItem label="Type" :value="formatGameSubscriptionType(subscription.subscription_type)" />
              <DataItem label="Status" :value="subscription.status" />
              <DataItem v-if="subscription.instance_limit !== null && subscription.instance_limit !== undefined"
                label="Instance Limit"
                :value="subscription.instance_limit === null ? 'Unlimited' : subscription.instance_limit" />
              <DataItem v-if="subscription.game_instance_ids && subscription.game_instance_ids.length > 0"
                label="Instances"
                :value="`${subscription.game_instance_ids.length} instance${subscription.game_instance_ids.length !== 1 ? 's' : ''}`" />
              <DataItem label="Created" :value="formatDate(subscription.created_at)" />
            </div>

            <template #actions>
              <TableActions :actions="getGameSubscriptionActions(subscription)" />
            </template>
          </DataCard>
        </div>

        <div v-else class="empty-state">
          <p>No game subscriptions found.</p>
        </div>
      </div>
    </div>

    <!-- Confirm Cancel Dialog -->
    <ConfirmationModal :visible="showCancelModal" title="Cancel Game Subscription"
      :message="`Are you sure you want to cancel your ${formatGameSubscriptionType(subscriptionToCancel?.subscription_type)} subscription to '${subscriptionToCancel?.game_name || 'this game'}'?`"
      @confirm="handleCancelSubscription" @cancel="closeCancelModal" />
  </div>
</template>

<script>
import { getMyAccountSubscriptions } from '@/api/accountSubscriptions'
import { getMyGameSubscriptions, cancelGameSubscription } from '@/api/gameSubscriptions'
import { listGames } from '@/api/games'
import PageHeader from '@/components/PageHeader.vue'
import DataCard from '@/components/DataCard.vue'
import DataItem from '@/components/DataItem.vue'
import TableActions from '@/components/TableActions.vue'
import ConfirmationModal from '@/components/ConfirmationModal.vue'
import AppButton from '@/components/Button.vue'

export default {
  name: 'AccountSubscriptionsView',
  components: {
    PageHeader,
    DataCard,
    DataItem,
    TableActions,
    ConfirmationModal,
    AppButton
  },
  data() {
    return {
      accountSubscriptions: [],
      gameSubscriptions: [],
      games: [],
      gameCount: 0,
      loading: true,
      error: null,
      showCancelModal: false,
      subscriptionToCancel: null
    }
  },
  async mounted() {
    await this.loadSubscriptions()
  },
  methods: {
    async loadSubscriptions() {
      try {
        this.loading = true
        this.error = null

        // Load account subscriptions
        const accountSubsResponse = await getMyAccountSubscriptions()
        this.accountSubscriptions = accountSubsResponse.data || []

        // Load game subscriptions
        const gameSubsResponse = await getMyGameSubscriptions()
        this.gameSubscriptions = gameSubsResponse.data || []

        // Load games to get names and count
        const gamesResponse = await listGames()
        this.games = gamesResponse.data || []

        // Count games owned by user
        this.gameCount = this.games.length

        // Enrich game subscriptions with game names
        this.gameSubscriptions = this.gameSubscriptions.map(sub => {
          const game = this.games.find(g => g.id === sub.game_id)
          return {
            ...sub,
            game_name: game ? game.name : 'Unknown Game'
          }
        })
      } catch (err) {
        this.error = err.message || 'Failed to load subscriptions'
        console.error('Error loading subscriptions:', err)
      } finally {
        this.loading = false
      }
    },
    formatSubscriptionType(type) {
      const types = {
        basic_game_designer: 'Basic Game Designer',
        professional_game_designer: 'Professional Game Designer'
      }
      return types[type] || type
    },
    formatGameSubscriptionType(type) {
      const types = {
        Player: 'Player',
        Manager: 'Manager',
        Designer: 'Designer'
      }
      return types[type] || type
    },
    formatDate(dateString) {
      if (!dateString) return 'N/A'
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    },
    openCancel(subscription) {
      this.subscriptionToCancel = subscription
      this.showCancelModal = true
    },
    closeCancelModal() {
      this.showCancelModal = false
      this.subscriptionToCancel = null
    },
    async handleCancelSubscription() {
      if (!this.subscriptionToCancel) return

      try {
        await cancelGameSubscription(this.subscriptionToCancel.id)
        this.closeCancelModal()
        await this.loadSubscriptions()
      } catch (err) {
        this.error = err.message || 'Failed to cancel subscription'
        console.error('Error cancelling subscription:', err)
        this.closeCancelModal()
      }
    },
    getGameSubscriptionActions(subscription) {
      const actions = []
      // Only allow cancelling Player and Manager subscriptions
      if (subscription.subscription_type === 'Player' || subscription.subscription_type === 'Manager') {
        actions.push({
          key: 'cancel',
          label: 'Cancel',
          danger: true,
          handler: () => this.openCancel(subscription)
        })
      }
      return actions
    }
  }
}
</script>

<style scoped>
.account-subscriptions-view {
  width: 100%;
}

.loading-state,
.error-state,
.empty-state {
  text-align: center;
  padding: var(--space-xl);
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.error-state button {
  margin-top: var(--space-md);
}

.subscriptions-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-xl);
}

.subscriptions-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.subscriptions-section h3 {
  margin: 0;
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text);
}

.section-description {
  margin: 0;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.subscriptions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(600px, 1fr));
  gap: var(--space-lg);
}

.subscription-card {
  min-height: 200px;
}

.subscription-info {
  margin-bottom: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}
</style>
