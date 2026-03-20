<template>
  <div class="account-profile-view">
    <PageHeader title="Account Profile" :showIcon="false" titleLevel="h2"
      subtitle="Review your profile details and account status" />

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
          <!-- Account Name with inline edit -->
          <div class="account-name-row">
            <template v-if="editingName">
              <span class="data-item-label">Account Name</span>
              <div class="name-edit-controls">
                <input
                  v-model="nameInput"
                  type="text"
                  class="name-input"
                  placeholder="Enter account name"
                  maxlength="255"
                  @keyup.enter="saveAccountName"
                  @keyup.escape="cancelEditName"
                />
                <div class="name-edit-buttons">
                  <AppButton @click="saveAccountName" variant="primary" size="small" :disabled="savingName">
                    {{ savingName ? 'Saving…' : 'Save' }}
                  </AppButton>
                  <AppButton @click="cancelEditName" variant="secondary" size="small" :disabled="savingName">
                    Cancel
                  </AppButton>
                </div>
                <p v-if="nameError" class="name-error">{{ nameError }}</p>
              </div>
            </template>
            <template v-else>
              <DataItem label="Account Name" :value="accountData ? (accountData.name || 'Not set') : '…'" />
              <AppButton @click="startEditName" variant="secondary" size="small" class="edit-name-btn">
                Edit
              </AppButton>
            </template>
          </div>
          <DataItem label="Email" :value="account.email" />
          <DataItem label="Account ID" :value="account.account_id" />
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
    <ConfirmationModal :visible="showDeleteModal" title="Delete Account"
      message="Are you sure you want to delete your account? This action cannot be undone."
      warning="All your data, including games, characters, and settings will be permanently deleted."
      confirmText="Delete Account" :loading="deleting" loadingText="Deleting..." :requireConfirmation="true"
      confirmationText="DELETE" @confirm="confirmDeleteAccount" @cancel="hideDeleteConfirmation" />
  </div>
</template>

<script>
import { getMe, getAccount, updateAccount, deleteAccountUser } from '@/api/account'
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
      accountData: null,
      loading: true,
      error: null,
      showDeleteModal: false,
      deleteConfirmationText: '',
      deleting: false,
      editingName: false,
      nameInput: '',
      savingName: false,
      nameError: null
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
        this.account = await getMe()
        this.accountData = await getAccount(this.account.account_id)
      } catch (err) {
        this.error = err.message || 'Failed to load account information'
        console.error('Error loading account:', err)
      } finally {
        this.loading = false
      }
    },
    startEditName() {
      this.nameInput = this.accountData ? this.accountData.name : ''
      this.nameError = null
      this.editingName = true
    },
    cancelEditName() {
      this.editingName = false
      this.nameError = null
    },
    async saveAccountName() {
      try {
        this.savingName = true
        this.nameError = null
        this.accountData = await updateAccount(this.account.account_id, { name: this.nameInput })
        this.editingName = false
      } catch (err) {
        this.nameError = err.message || 'Failed to update account name'
      } finally {
        this.savingName = false
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
        await deleteAccountUser(this.account.account_id, this.account.id)

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
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.error-state button {
  margin-top: var(--space-md);
}

.account-grid {
  display: flex;
  flex-direction: column;
  gap: var(--space-lg);
  max-width: 600px;
  width: 100%;
  margin: 0;
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

.account-name-row {
  display: flex;
  align-items: flex-start;
  gap: var(--space-sm);
}

.data-item-label {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
  font-weight: var(--font-weight-bold);
  min-width: 120px;
  flex-shrink: 0;
}

.name-edit-controls {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
  flex: 1;
}

.name-input {
  padding: var(--space-xs) var(--space-sm);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
  background: var(--color-bg);
  color: var(--color-text);
  width: 100%;
  box-sizing: border-box;
}

.name-input:focus {
  outline: none;
  border-color: var(--color-primary, #3b82f6);
}

.name-edit-buttons {
  display: flex;
  gap: var(--space-xs);
}

.name-error {
  color: var(--color-danger, #dc2626);
  font-size: var(--font-size-xs, 0.75rem);
  margin: 0;
}

.edit-name-btn {
  flex-shrink: 0;
  align-self: center;
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
