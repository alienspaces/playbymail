<!--
  AccountView.vue
  Account information and settings page.
-->
<template>
  <div class="account-dashboard">
    <div class="dashboard-header">
      <h2>Account</h2>
      <p>Manage your account information and settings</p>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="loading-state">
      <p>Loading account information...</p>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
      <button @click="loadAccount">Retry</button>
    </div>

    <!-- Account information -->
    <div v-else-if="account" class="account-grid">
      <!-- Profile Information Card -->
      <DataCard title="Profile Information" class="game-card">
        <div class="game-info">
          <DataItem label="Email" :value="account.email" />
          <DataItem label="Account Created" :value="formatDate(account.created_at)" />
          <DataItem v-if="account.updated_at" label="Last Updated" :value="formatDate(account.updated_at)" />
        </div>
      </DataCard>

      <!-- Contact Information Cards -->
      <DataCard 
        v-for="contact in accountContacts" 
        :key="contact.id"
        title="Contact Information"
      >
        <div class="account-info">
          <DataItem label="Name" :value="contact.name" />
          <DataItem label="Address Line 1" :value="contact.postal_address_line1" />
          <DataItem v-if="contact.postal_address_line2" label="Address Line 2" :value="contact.postal_address_line2" />
          <DataItem label="State/Province" :value="contact.state_province" />
          <DataItem label="Country" :value="contact.country" />
          <DataItem label="Postal Code" :value="contact.postal_code" />
        </div>
      </DataCard>
      
      <!-- Account ID Card -->
      <DataCard title="Account ID">
        <div class="account-info">
          <DataItem label="ID" :value="account.id" />
        </div>
      </DataCard>

      <!-- Danger Zone Card -->
      <DataCard title="Danger Zone">
        <div class="account-info">
          <p>This action cannot be undone. This will permanently delete your account and all associated data.</p>
        </div>
        <template #primary>
          <Button @click="showDeleteConfirmation" variant="danger" size="small">
            Delete Account
          </Button>
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

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { getMyAccount, getAccountContacts, deleteMyAccount } from '@/api/account';
import Button from '@/components/Button.vue';
import DataCard from '@/components/DataCard.vue';
import DataItem from '@/components/DataItem.vue';
import ConfirmationModal from '@/components/ConfirmationModal.vue';

const router = useRouter();
const authStore = useAuthStore();

const account = ref(null);
const accountContacts = ref([]);
const loading = ref(true);
const error = ref(null);
const showDeleteModal = ref(false);
const deleting = ref(false);

onMounted(async () => {
  await loadAccount();
});

const loadAccount = async () => {
  try {
    loading.value = true;
    error.value = null;
    account.value = await getMyAccount();
    
    // Load account contacts if we have an account
    if (account.value && account.value.id) {
      try {
        accountContacts.value = await getAccountContacts(account.value.id);
      } catch (err) {
        // Contacts might not exist yet, that's okay
        console.log('No account contacts found:', err);
        accountContacts.value = [];
      }
    }
  } catch (err) {
    error.value = err.message || 'Failed to load account information';
    console.error('Error loading account:', err);
  } finally {
    loading.value = false;
  }
};

const showDeleteConfirmation = () => {
  showDeleteModal.value = true;
  error.value = null;
};

const hideDeleteConfirmation = () => {
  showDeleteModal.value = false;
  error.value = null;
};

const confirmDeleteAccount = async () => {
  try {
    deleting.value = true;
    error.value = null;
    await deleteMyAccount();
    
    // Clear auth store and redirect to home
    authStore.logout();
    router.push('/');
  } catch (err) {
    error.value = err.message || 'Failed to delete account';
    console.error('Error deleting account:', err);
  } finally {
    deleting.value = false;
  }
};

const formatDate = (dateString) => {
  if (!dateString) return 'N/A';
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};
</script>

<style scoped>
.account-dashboard {
  max-width: 1200px;
  margin: 0 auto;
}

.dashboard-header {
  margin-bottom: var(--space-xl);
  text-align: center;
}

.dashboard-header h2 {
  margin: 0 0 var(--space-sm) 0;
  font-size: var(--font-size-xl);
  color: var(--color-text);
}

.dashboard-header p {
  margin: 0;
  color: var(--color-text-muted);
  font-size: var(--font-size-md);
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
  padding: var(--space-sm) var(--space-md);
  background: var(--color-primary);
  color: var(--color-text-light);
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
}

.account-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
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

.account-info {
  margin-bottom: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.account-info p {
  margin: 0;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  line-height: 1.4;
}
</style>
