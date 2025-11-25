<template>
  <div class="account-profile-view">
    <div v-if="loading" class="loading">
      <p>Loading account information...</p>
    </div>

    <div v-else-if="error" class="error">
      <p>{{ error }}</p>
    </div>

    <div v-else-if="account" class="account-info">
      <section>
        <h2>Profile Information</h2>
        <div class="info-item">
          <label>Email:</label>
          <span>{{ account.email }}</span>
        </div>
        <div class="info-item">
          <label>Account Created:</label>
          <span>{{ formatDate(account.created_at) }}</span>
        </div>
        <div v-if="account.updated_at" class="info-item">
          <label>Last Updated:</label>
          <span>{{ formatDate(account.updated_at) }}</span>
        </div>
      </section>

      <section>
        <h2>Account ID</h2>
        <div class="info-item">
          <label>ID:</label>
          <span class="account-id">{{ account.id }}</span>
        </div>
      </section>

      <section class="danger-zone">
        <h2>Danger Zone</h2>
        <div class="danger-item">
          <div class="danger-content">
            <h3>Delete Account</h3>
            <p>This action cannot be undone. This will permanently delete your account and all associated data.</p>
          </div>
          <button @click="showDeleteConfirmation" class="btn delete-btn">Delete Account</button>
        </div>
      </section>
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
import { getMyAccount, deleteMyAccount } from '@/api/account'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import ConfirmationModal from '@/components/ConfirmationModal.vue'

export default {
  name: 'AccountProfileView',
  components: {
    ConfirmationModal
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
  /* Content is already aligned by AccountLayout */
}

h2 {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-lg);
  color: var(--color-text);
  font-family: var(--font-family-heading);
  text-transform: uppercase;
  letter-spacing: 0.8px;
  position: relative;
  display: inline-block;
  text-shadow: 0.3px 0.3px 0px rgba(0, 0, 0, 0.1);
  transform: rotate(0.2deg);
}

h2::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(to right,
      var(--color-pen-blue) 0%,
      transparent 100%);
  opacity: 0.3;
  transform: rotate(-0.2deg);
}

h3 {
  margin: 0 0 var(--space-sm) 0;
  color: var(--color-text);
  font-size: var(--font-size-md);
}

.loading,
.error {
  text-align: center;
  padding: var(--space-lg);
}

.error {
  color: var(--color-warning-dark);
  background: var(--color-warning-light);
  padding: var(--space-sm) var(--space-md);
  margin-top: var(--space-md);
  text-align: center;
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-warning);
}

.account-info section {
  margin-bottom: var(--space-xl);
}

.account-info section:last-child {
  margin-bottom: 0;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: var(--space-md);
  padding: var(--space-sm) 0;
  border-bottom: 1px solid var(--color-border-light);
}

.info-item:last-child {
  border-bottom: none;
}

.info-item label {
  font-weight: var(--font-weight-semibold);
  min-width: 140px;
  color: var(--color-text-muted);
}

.info-item span {
  color: var(--color-text);
  flex: 1;
}

.account-id {
  font-family: monospace;
  background: var(--color-bg-alt);
  padding: var(--space-xs) var(--space-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
}

/* Danger Zone Styles */
.danger-zone {
  border: 1px solid var(--color-danger);
  border-radius: var(--radius-md);
  padding: var(--space-lg);
  background: var(--color-danger-light);
}

.danger-zone h2 {
  color: var(--color-danger);
  margin-bottom: var(--space-md);
}

.danger-item {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--space-lg);
  flex-wrap: wrap;
}

.danger-content {
  flex: 1;
  min-width: 0;
}

.danger-content h3 {
  color: var(--color-danger);
  margin-bottom: var(--space-sm);
}

.danger-content p {
  color: var(--color-text-muted);
  margin: 0;
  font-size: var(--font-size-sm);
  line-height: 1.4;
}

.delete-btn {
  background: var(--color-danger);
  color: var(--color-text-light);
  border: 2px solid var(--color-danger);
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--font-size-md);
  font-weight: var(--font-weight-bold);
  transition: all 0.2s ease-in-out;
  white-space: nowrap;
  flex-shrink: 0;
}

.delete-btn:hover {
  background: var(--color-danger-dark);
  border-color: var(--color-danger-dark);
}

/* Mobile responsive breakpoint */
@media (max-width: 768px) {
  .danger-item {
    flex-direction: column;
    align-items: stretch;
  }

  .delete-btn {
    align-self: flex-start;
    margin-top: var(--space-md);
  }

  .info-item {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-xs);
  }

  .info-item label {
    min-width: auto;
  }
}
</style>
