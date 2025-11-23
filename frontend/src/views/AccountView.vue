<template>
  <div class="account-view">
    <h1>Account</h1>
    
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
          <label>Name:</label>
          <div v-if="!editingName" class="name-display">
            <span>{{ account.name }}</span>
            <button @click="startEditName" class="edit-btn">Edit</button>
          </div>
          <div v-else class="name-edit">
            <input 
              v-model="editingNameValue" 
              @keyup.enter="saveName"
              @keyup.esc="cancelEditName"
              ref="nameInput"
              class="name-input"
            />
            <div class="edit-actions">
              <button @click="saveName" class="save-btn">Save</button>
              <button @click="cancelEditName" class="cancel-btn">Cancel</button>
            </div>
          </div>
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
          <button @click="showDeleteConfirmation" class="delete-btn">Delete Account</button>
        </div>
      </section>
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
import { getMyAccount, updateMyAccount, deleteMyAccount } from '@/api/account'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import ConfirmationModal from '@/components/ConfirmationModal.vue'

export default {
  name: 'AccountView',
  components: {
    ConfirmationModal
  },
  data() {
    return {
      account: null,
      loading: true,
      error: null,
      editingName: false,
      editingNameValue: '',
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
    startEditName() {
      this.editingName = true
      this.editingNameValue = this.account.name
      this.$nextTick(() => {
        this.$refs.nameInput.focus()
      })
    },
    async saveName() {
      if (!this.editingNameValue.trim()) {
        this.error = 'Name cannot be empty'
        return
      }
      
      try {
        this.loading = true
        this.error = null
        const updatedAccount = await updateMyAccount({ name: this.editingNameValue.trim() })
        this.account = updatedAccount
        this.editingName = false
      } catch (err) {
        this.error = err.message || 'Failed to update name'
        console.error('Error updating name:', err)
      } finally {
        this.loading = false
      }
    },
    cancelEditName() {
      this.editingName = false
      this.editingNameValue = ''
      this.error = null
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
.account-view {
  max-width: 600px;
  margin: 2rem auto;
  padding: 2rem;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

h1 {
  margin-bottom: 1.5rem;
  color: #333;
}

h2 {
  margin-bottom: 1rem;
  color: #555;
  font-size: 1.2rem;
}

h3 {
  margin: 0 0 0.5rem 0;
  color: #333;
  font-size: 1rem;
}

.loading, .error {
  text-align: center;
  padding: 2rem;
}

.error {
  color: #d32f2f;
}

.account-info section {
  margin-bottom: 2rem;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
  padding: 0.75rem 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-item label {
  font-weight: 600;
  min-width: 120px;
  color: #666;
}

.info-item span {
  color: #333;
  flex: 1;
}

.account-id {
  font-family: monospace;
  background: #f5f5f5;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.9rem;
}

.name-display {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex: 1;
}

.edit-btn {
  background: #2196f3;
  color: white;
  border: none;
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8rem;
}

.edit-btn:hover {
  background: var(--color-primary);
}

.name-edit {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex: 1;
}

.name-input {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.name-input:focus {
  outline: none;
  border-color: #2196f3;
  box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.1);
}

.edit-actions {
  display: flex;
  gap: 0.5rem;
}

.save-btn, .cancel-btn {
  padding: 0.25rem 0.75rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8rem;
}

.save-btn {
  background: #4caf50;
  color: white;
}

.save-btn:hover {
  background: var(--color-success);
}

.cancel-btn {
  background: #f44336;
  color: white;
}

.cancel-btn:hover {
  background: #d32f2f;
}

/* Danger Zone Styles */
.danger-zone {
  border: 1px solid #ffcdd2;
  border-radius: 8px;
  padding: 1.5rem;
  background: #fff5f5;
}

.danger-zone h2 {
  color: #d32f2f;
  margin-bottom: 1rem;
}

.danger-item {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: wrap;
}

.danger-content {
  flex: 1;
  min-width: 0;
}

.danger-content h3 {
  color: #d32f2f;
  margin-bottom: 0.5rem;
}

.danger-content p {
  color: #666;
  margin: 0;
  font-size: 0.9rem;
  line-height: 1.4;
}

.delete-btn {
  background: #d32f2f;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  white-space: nowrap;
  flex-shrink: 0;
}

.delete-btn:hover {
  background: #b71c1c;
}

/* Mobile responsive breakpoint */
@media (max-width: 768px) {
  .danger-item {
    flex-direction: column;
    align-items: stretch;
  }
  
  .delete-btn {
    align-self: flex-start;
    margin-top: 1rem;
  }
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.modal-content h2 {
  color: #d32f2f;
  margin-bottom: 1rem;
}

.modal-content p {
  margin-bottom: 1rem;
  color: #333;
  line-height: 1.5;
}

.warning-text {
  color: #d32f2f;
  font-weight: 600;
  background: #ffebee;
  padding: 0.75rem;
  border-radius: 4px;
  border-left: 4px solid #d32f2f;
}

.confirmation-input {
  margin: 1.5rem 0;
}

.confirmation-input label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #333;
}

.confirm-input {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  text-transform: uppercase;
}

.confirm-input:focus {
  outline: none;
  border-color: #d32f2f;
  box-shadow: 0 0 0 2px rgba(211, 47, 47, 0.1);
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

.confirm-delete-btn {
  background: #d32f2f;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.confirm-delete-btn:hover:not(:disabled) {
  background: #b71c1c;
}

.confirm-delete-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}
</style> 