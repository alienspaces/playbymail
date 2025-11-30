<template>
  <div class="account-contacts-view">
    <PageHeader title="Account Contacts" actionText="Create New Contact" :showIcon="false" titleLevel="h2"
      subtitle="Manage postal addresses for shipments and correspondence" @action="openCreate" />

    <!-- Loading state -->
    <div v-if="loading" class="loading-state">
      <p>Loading contacts...</p>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
      <AppButton @click="loadAccountAndContacts" variant="primary" size="small">
        Retry
      </AppButton>
    </div>

    <!-- Contacts grid -->
    <div v-else-if="accountContacts.length > 0" class="contacts-grid">
      <DataCard v-for="contact in accountContacts" :key="contact.id" :title="contact.name || 'Unnamed Contact'"
        class="contact-card">
        <div class="contact-info">
          <DataItem v-if="formatAddress(contact)" label="Address" :value="formatAddress(contact)" />
          <DataItem v-if="contact.state_province" label="State/Province" :value="contact.state_province" />
          <DataItem v-if="contact.country" label="Country" :value="contact.country" />
          <DataItem v-if="contact.postal_code" label="Postal Code" :value="contact.postal_code" />
        </div>

        <template #actions>
          <TableActions :actions="getContactActions(contact)" />
        </template>
      </DataCard>
    </div>

    <!-- Empty state -->
    <div v-else class="empty-state">
      <h3>No Contacts</h3>
      <p>You don't have any contacts yet. Create your first contact to get started.</p>
    </div>

    <!-- Create/Edit Contact Modal -->
    <ContactModal v-if="showCreateModal || showEditModal" :visible="showCreateModal || showEditModal"
      :contact="editingContact" :account-id="accountId" @close="closeModal" @saved="handleContactSaved" />

    <!-- Confirm Delete Dialog -->
    <ConfirmationModal :visible="showDeleteModal" title="Delete Contact"
      :message="`Are you sure you want to delete '${contactToDelete?.name || 'this contact'}'?`" @confirm="handleDelete"
      @cancel="closeDeleteModal" />
  </div>
</template>

<script>
import { getMyAccount, getAccountContacts, deleteAccountContact } from '@/api/account'
import ContactModal from '@/components/ContactModal.vue'
import PageHeader from '@/components/PageHeader.vue'
import DataCard from '@/components/DataCard.vue'
import DataItem from '@/components/DataItem.vue'
import TableActions from '@/components/TableActions.vue'
import ConfirmationModal from '@/components/ConfirmationModal.vue'

export default {
  name: 'AccountContactsView',
  components: {
    ContactModal,
    PageHeader,
    DataCard,
    DataItem,
    TableActions,
    ConfirmationModal
  },
  data() {
    return {
      accountId: null,
      accountContacts: [],
      loading: true,
      error: null,
      showCreateModal: false,
      showEditModal: false,
      editingContact: null,
      showDeleteModal: false,
      contactToDelete: null
    }
  },
  async mounted() {
    await this.loadAccountAndContacts()
  },
  methods: {
    formatAddress(contact) {
      const parts = []
      if (contact.postal_address_line1) parts.push(contact.postal_address_line1)
      if (contact.postal_address_line2) parts.push(contact.postal_address_line2)
      return parts.join(', ') || ''
    },
    openCreate() {
      this.showCreateModal = true
    },
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
    openDelete(contact) {
      this.contactToDelete = contact
      this.showDeleteModal = true
    },
    closeDeleteModal() {
      this.showDeleteModal = false
      this.contactToDelete = null
    },
    async handleDelete() {
      if (!this.contactToDelete) return

      try {
        await deleteAccountContact(this.accountId, this.contactToDelete.id)
        this.closeDeleteModal()
        await this.loadAccountAndContacts()
      } catch (err) {
        this.error = err.message || 'Failed to delete contact'
        console.error('Error deleting contact:', err)
        this.closeDeleteModal()
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
    },
    getContactActions(contact) {
      return [
        {
          key: 'edit',
          label: 'Edit',
          handler: () => this.editContact(contact)
        },
        {
          key: 'delete',
          label: 'Delete',
          danger: true,
          handler: () => this.openDelete(contact)
        }
      ]
    }
  }
}
</script>

<style scoped>
.account-contacts-view {
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

.contacts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(600px, 1fr));
  gap: var(--space-lg);
}

.contact-card {
  min-height: 200px;
}

.contact-info {
  margin-bottom: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}
</style>
