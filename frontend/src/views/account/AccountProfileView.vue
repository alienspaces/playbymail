<template>
  <div class="account-profile-view">
    <PageHeader 
      title="Account Profile" 
      :showIcon="false"
      titleLevel="h2"
    />

    <!-- Loading state -->
    <div v-if="loading" class="loading-state">
      <p>Loading account information...</p>
    </div>
    
    <!-- Error state -->
    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
      <AppButton @click="loadAccount" variant="primary" size="small">
        Retry
      </AppButton>
    </div>
    
    <!-- Account information -->
    <div v-else-if="account" class="account-grid">
      <!-- Profile Information Card -->
      <DataCard title="Profile Information" class="game-card">
        <div class="game-info">
          <DataItem label="Email" :value="account.email" />
          <DataItem label="Account ID" :value="account.id" />
          <DataItem label="Account Created" :value="formatDate(account.created_at)" />
          <DataItem v-if="account.updated_at" label="Last Updated" :value="formatDate(account.updated_at)" />
        </div>
      </DataCard>

      <!-- Danger Zone Card -->
      <DataCard title="Danger Zone" variant="danger" class="game-card">
        <div class="game-info">
          <p>This action cannot be undone. This will permanently delete your account and all associated data.</p>
        </div>
        <template #primary>
          <AppButton @click="showDeleteConfirmation" variant="danger" size="small">
            Delete Account
          </AppButton>
        </template>
      </DataCard>
    </div>

    <!-- Delete Confirmation Modal -->
    <ConfirmationModal
      :visible="showDeleteModal"
      title="Delete Account"
      message="Are you sure you want to delete your account? This action cannot be undone."
      warning="All your data, including games, characters, and settings will be permanently deleted."
      confirmText="Delete Account"
      :loading="deleting"
      loadingText="Deleting..."
      :requireConfirmation="true"
      confirmationText="DELETE"
      @confirm="confirmDeleteAccount"
      @cancel="hideDeleteConfirmation"
    />
  </div>
</template>

<script>
import { getMyAccount, deleteMyAccount } from '@/api/account'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import ConfirmationModal from '@/components/ConfirmationModal.vue'
import PageHeader from '@/components/PageHeader.vue'
import DataCard from '@/components/DataCard.vue'
import DataItem from '@/components/DataItem.vue'
import AppButton from '@/components/Button.vue'

export default {
  name: 'AccountProfileView',
  components: {
    ConfirmationModal,
    PageHeader,
    DataCard,
    DataItem,
    AppButton
  },
  data() {
    return {
      account: null,
      loading: true,
      error: null,
      showDeleteModal: false,
      deleteConfirmationText: '',
      deleting: false
    }
  },
  async mounted() {
    await this.loadAccount()
  },
  methods: {
    async loadAccount() {
      try {
        this.loading = true
        this.error = null
        this.account = await getMyAccount()
      } catch (err) {
        this.error = err.message || 'Failed to load account information'
        console.error('Error loading account:', err)
      } finally {
        this.loading = false
      }
    },
    showDeleteConfirmation() {
      this.showDeleteModal = true
      this.deleteConfirmationText = ''
      this.error = null
    },
    hideDeleteConfirmation() {
      this.showDeleteModal = false
      this.deleteConfirmationText = ''
      this.error = null
    },
    async confirmDeleteAccount() {
      if (this.deleteConfirmationText !== 'DELETE') {
        return
      }
      
      try {
        this.deleting = true
        this.error = null
        await deleteMyAccount()
        
        // Clear auth store and redirect to home
        const authStore = useAuthStore()
        authStore.logout()
        
        const router = useRouter()
        router.push('/')
      } catch (err) {
        this.error = err.message || 'Failed to delete account'
        console.error('Error deleting account:', err)
      } finally {
        this.deleting = false
      }
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
    }
  }
}
</script>

<style scoped>
.account-profile-view {
  width: 100%;
}

.loading-state,
.error-state {
  text-align: center;
  padding: var(--space-xl);
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.error-state button {
  margin-top: var(--space-md);
}

.account-grid {
  display: flex;
  flex-direction: column;
  gap: var(--space-lg);
}

.game-card {
  min-height: 280px;
}

.game-info {
  margin-bottom: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.game-info p {
  margin: 0;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  line-height: 1.4;
}

:deep(.data-card-danger) .game-info p {
  color: var(--color-text);
}

:deep(.data-card-danger.game-card) {
  min-height: auto;
}

:deep(.data-card-danger .card-content) {
  flex: 0 1 auto;
}
</style>

