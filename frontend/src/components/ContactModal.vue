<template>
  <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal">
      <h2>{{ contact ? 'Edit Contact' : 'Add Contact' }}</h2>
      <form @submit.prevent="handleSubmit" class="modal-form">
        <div class="form-group">
          <label for="name">Name <span class="required">*</span></label>
          <input
            id="name"
            v-model="formData.name"
            type="text"
            required
            maxlength="255"
            placeholder="Enter contact name"
          />
        </div>

        <div class="form-group">
          <label for="postal_address_line1">Address Line 1 <span class="required">*</span></label>
          <input
            id="postal_address_line1"
            v-model="formData.postal_address_line1"
            type="text"
            required
            maxlength="255"
            placeholder="Street address"
          />
        </div>

        <div class="form-group">
          <label for="postal_address_line2">Address Line 2</label>
          <input
            id="postal_address_line2"
            v-model="formData.postal_address_line2"
            type="text"
            maxlength="255"
            placeholder="Apartment, suite, etc. (optional)"
          />
        </div>

        <div class="form-group">
          <label for="state_province">State/Province <span class="required">*</span></label>
          <input
            id="state_province"
            v-model="formData.state_province"
            type="text"
            required
            maxlength="100"
            placeholder="State or province"
          />
        </div>

        <div class="form-group">
          <label for="country">Country <span class="required">*</span></label>
          <input
            id="country"
            v-model="formData.country"
            type="text"
            required
            maxlength="100"
            placeholder="Country"
          />
        </div>

        <div class="form-group">
          <label for="postal_code">Postal Code <span class="required">*</span></label>
          <input
            id="postal_code"
            v-model="formData.postal_code"
            type="text"
            required
            maxlength="20"
            placeholder="Postal/ZIP code"
          />
        </div>

        <div class="modal-actions">
          <button type="submit" :disabled="loading">
            {{ loading ? 'Saving...' : (contact ? 'Update' : 'Create') }}
          </button>
          <button type="button" @click="$emit('close')">Cancel</button>
        </div>
      </form>
      <div v-if="error" class="error">
        <p>{{ error }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import { createAccountContact, updateAccountContact } from '@/api/account';

const props = defineProps({
  visible: Boolean,
  contact: Object,
  accountId: String,
  accountUserId: String
});

const emit = defineEmits(['close', 'saved']);

const formData = ref({
  name: '',
  postal_address_line1: '',
  postal_address_line2: '',
  state_province: '',
  country: '',
  postal_code: ''
});

const loading = ref(false);
const error = ref(null);

// Watch for contact changes to populate form
watch(() => props.contact, (newContact) => {
  if (newContact) {
    formData.value = {
      name: newContact.name || '',
      postal_address_line1: newContact.postal_address_line1 || '',
      postal_address_line2: newContact.postal_address_line2 || '',
      state_province: newContact.state_province || '',
      country: newContact.country || '',
      postal_code: newContact.postal_code || ''
    };
  } else {
    resetForm();
  }
}, { immediate: true });

// Watch for visibility to reset form when closed
watch(() => props.visible, (isVisible) => {
  if (!isVisible) {
    resetForm();
    error.value = null;
  } else if (props.contact) {
    // Populate form when opening with existing contact
    formData.value = {
      name: props.contact.name || '',
      postal_address_line1: props.contact.postal_address_line1 || '',
      postal_address_line2: props.contact.postal_address_line2 || '',
      state_province: props.contact.state_province || '',
      country: props.contact.country || '',
      postal_code: props.contact.postal_code || ''
    };
  }
});

function resetForm() {
  formData.value = {
    name: '',
    postal_address_line1: '',
    postal_address_line2: '',
    state_province: '',
    country: '',
    postal_code: ''
  };
}

async function handleSubmit() {
  if (!props.accountId || !props.accountUserId) {
    error.value = 'Account ID and Account User ID are required';
    return;
  }

  try {
    loading.value = true;
    error.value = null;

    const contactData = {
      name: formData.value.name.trim(),
      postal_address_line1: formData.value.postal_address_line1.trim(),
      postal_address_line2: formData.value.postal_address_line2.trim() || undefined,
      state_province: formData.value.state_province.trim(),
      country: formData.value.country.trim(),
      postal_code: formData.value.postal_code.trim()
    };

    if (props.contact && props.contact.id) {
      await updateAccountContact(props.accountId, props.accountUserId, props.contact.id, contactData);
    } else {
      await createAccountContact(props.accountId, props.accountUserId, contactData);
    }

    emit('saved');
  } catch (err) {
    error.value = err.message || 'Failed to save contact';
    console.error('Error saving contact:', err);
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.modal-form {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.required {
  color: var(--color-danger);
}

.error {
  color: var(--color-warning-dark);
  background: var(--color-warning-light);
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-warning);
  margin-top: var(--space-md);
}

.error p {
  margin: 0;
}

@media (max-width: 768px) {
  .modal-actions {
    flex-direction: column-reverse;
  }

  .modal-actions button {
    width: 100%;
  }
}
</style>

