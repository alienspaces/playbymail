<template>
  <div v-if="visible" class="modal-overlay" @click.self="$emit('cancel')">
    <div class="modal">
      <h2>{{ title }}</h2>
      <p>{{ message }}</p>
      <p v-if="warning" class="warning-text">{{ warning }}</p>
      
      <div v-if="requireConfirmation" class="confirmation-input">
        <label :for="confirmationId">Type "{{ confirmationText }}" to confirm:</label>
        <input 
          :id="confirmationId"
          v-model="confirmationValue" 
          :placeholder="confirmationText"
          class="confirm-input"
        />
      </div>
      
      <div class="modal-actions">
        <button type="button" @click="$emit('cancel')">Cancel</button>
        <button 
          type="button"
          @click="$emit('confirm')" 
          :disabled="isConfirmDisabled || loading"
          class="danger-btn"
        >
          {{ loading ? loadingText : confirmText }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';

const props = defineProps({
  visible: Boolean,
  title: String,
  message: String,
  warning: String,
  confirmText: {
    type: String,
    default: 'Confirm'
  },
  loading: Boolean,
  loadingText: {
    type: String,
    default: 'Loading...'
  },
  requireConfirmation: Boolean,
  confirmationText: {
    type: String,
    default: 'DELETE'
  }
});

defineEmits(['confirm', 'cancel']);

const confirmationValue = ref('');
const confirmationId = ref(`confirm-${Date.now()}`);

const isConfirmDisabled = computed(() => {
  if (!props.requireConfirmation) return false;
  return confirmationValue.value !== props.confirmationText;
});
</script>

<style scoped>
.modal h2 {
  color: var(--color-danger);
}

.modal p {
  margin-bottom: var(--space-md);
  color: var(--color-text);
  line-height: 1.5;
}

.warning-text {
  color: var(--color-danger);
  font-weight: var(--font-weight-semibold);
  background: var(--color-danger-light);
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  border-left: 4px solid var(--color-danger);
  margin-bottom: var(--space-md);
}

.confirmation-input {
  margin: var(--space-lg) 0;
}

.confirmation-input label {
  display: block;
  margin-bottom: var(--space-xs);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
}

.confirm-input {
  width: 100%;
  padding: var(--space-sm) var(--space-md);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-md);
  text-transform: uppercase;
  font-family: inherit;
}

.confirm-input:focus {
  outline: none;
  border-color: var(--color-danger);
  box-shadow: 0 0 0 2px rgba(220, 38, 38, 0.1);
}

.danger-btn {
  background: var(--color-danger) !important;
  color: var(--color-text-light) !important;
  border-color: var(--color-danger) !important;
}

.danger-btn:hover:not(:disabled) {
  background: var(--color-danger-dark) !important;
  border-color: var(--color-danger-dark) !important;
}

.danger-btn:disabled {
  background: var(--color-text-muted) !important;
  border-color: var(--color-text-muted) !important;
  cursor: not-allowed;
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