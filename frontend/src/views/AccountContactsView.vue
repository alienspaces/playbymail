<template>
  <div class="account-contacts-view">
    <div class="view-header">
      <h2>Account Contacts</h2>
      <button @click="showCreateModal = true" class="btn btn-primary">
        <svg class="icon" viewBox="0 0 24 24" fill="currentColor">
          <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
        </svg>
        Add Contact
      </button>
    </div>

    <div v-if="loading" class="loading">
      <p>Loading contacts...</p>
    </div>

    <div v-else-if="error" class="error">
      <p>{{ error }}</p>
    </div>

    <div v-else-if="accountContacts && accountContacts.length > 0" class="contacts-list">
      <DataCard v-for="contact in accountContacts" :key="contact.id" :title="contact.name || 'Unnamed Contact'">
        <template #actions>
          <button @click="editContact(contact)" class="btn btn-secondary">Edit</button>
          <button @click="deleteContact(contact)" class="btn btn-danger">Delete</button>
        </template>

        <div class="contact-info">
          <div class="info-item">
            <label>Name:</label>
            <span>{{ contact.name || 'N/A' }}</span>
          </div>
          <div class="info-item">
            <label>Address Line 1:</label>
            <span>{{ contact.postal_address_line1 || 'N/A' }}</span>
          </div>
          <div v-if="contact.postal_address_line2" class="info-item">
            <label>Address Line 2:</label>
            <span>{{ contact.postal_address_line2 }}</span>
          </div>
          <div class="info-item">
            <label>State/Province:</label>
            <span>{{ contact.state_province || 'N/A' }}</span>
          </div>
          <div class="info-item">
            <label>Country:</label>
            <span>{{ contact.country || 'N/A' }}</span>
          </div>
          <div class="info-item">
            <label>Postal Code:</label>
            <span>{{ contact.postal_code || 'N/A' }}</span>
          </div>
        </div>
      </DataCard>
    </div>

    <div v-else class="empty-state">
      <p>No contacts found. Add your first contact to get started.</p>
    </div>

    <!-- Create/Edit Contact Modal -->
    <ContactModal v-if="showCreateModal || showEditModal" :visible="showCreateModal || showEditModal"
      :contact="editingContact" :account-id="accountId" @close="closeModal" @saved="handleContactSaved" />
  </div>
</template>

<script>
import { getMyAccount, getAccountContacts, deleteAccountContact } from '@/api/account'
import ContactModal from '@/components/ContactModal.vue'
import DataCard from '@/components/DataCard.vue'

export default {
  name: 'AccountContactsView',
  components: {
    ContactModal,
    DataCard
  },
  data() {
    return {
      accountId: null,
      accountContacts: [],
      loading: true,
      error: null,
      showCreateModal: false,
      showEditModal: false,
      editingContact: null
    }
  },
  async mounted() {
    await this.loadAccountAndContacts()
  },
  methods: {
    async loadAccountAndContacts() {
      try {
        this.loading = true
        this.error = null

        // Get account first
        const account = await getMyAccount()
        this.accountId = account.id

        // Load contacts
        if (this.accountId) {
          try {
            this.accountContacts = await getAccountContacts(this.accountId)
          } catch (err) {
            // Contacts might not exist yet, that's okay
            console.log('No account contacts found:', err)
            this.accountContacts = []
          }
        }
      } catch (err) {
        this.error = err.message || 'Failed to load account contacts'
        console.error('Error loading account contacts:', err)
      } finally {
        this.loading = false
      }
    },
    editContact(contact) {
      this.editingContact = contact
      this.showEditModal = true
    },
    async deleteContact(contact) {
      if (!confirm(`Are you sure you want to delete the contact "${contact.name || 'this contact'}"?`)) {
        return
      }

      try {
        await deleteAccountContact(this.accountId, contact.id)
        await this.loadAccountAndContacts()
      } catch (err) {
        this.error = err.message || 'Failed to delete contact'
        console.error('Error deleting contact:', err)
      }
    },
    closeModal() {
      this.showCreateModal = false
      this.showEditModal = false
      this.editingContact = null
    },
    async handleContactSaved() {
      this.closeModal()
      await this.loadAccountAndContacts()
    }
  }
}
</script>

<style scoped>
.account-contacts-view {
  max-width: 900px;
  margin: 0 auto;
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-xl);
}

.view-header h2 {
  margin: 0;
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

.view-header h2::after {
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

.btn-primary {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-sm) var(--space-md);
  background: transparent;
  color: var(--color-button);
  border: 2px solid var(--color-button);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--font-size-md);
  font-weight: var(--font-weight-bold);
  transition: all 0.2s ease-in-out;
}

.btn-primary:hover {
  background: var(--color-button);
  color: var(--color-text-light);
}

.icon {
  width: 16px;
  height: 16px;
}

.loading,
.error,
.empty-state {
  text-align: center;
  padding: var(--space-xl);
}

.error {
  color: var(--color-warning-dark);
  background: var(--color-warning-light);
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-warning);
}

.empty-state {
  text-align: center;
  padding: var(--space-xl);
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  color: var(--color-text-muted);
  font-size: var(--font-size-md);
}

.contacts-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: var(--space-lg);
}

.btn-secondary {
  padding: var(--space-xs) var(--space-sm);
  background: transparent;
  color: var(--color-button);
  border: 1px solid var(--color-button);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: var(--color-button);
  color: var(--color-text-light);
}

.btn-danger {
  padding: var(--space-xs) var(--space-sm);
  background: transparent;
  color: var(--color-danger);
  border: 1px solid var(--color-danger);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  transition: all 0.2s;
}

.btn-danger:hover {
  background: var(--color-danger);
  color: var(--color-text-light);
}

.contact-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  flex: 1;
}

.info-item {
  display: flex;
  align-items: center;
  padding: var(--space-xs) 0;
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

/* Mobile responsive breakpoint */
@media (max-width: 768px) {
  .view-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-md);
  }

  .contacts-list {
    grid-template-columns: 1fr;
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
