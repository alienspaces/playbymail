<template>
  <div v-if="visible" class="modal-overlay" @click="$emit('cancel')">
    <div class="modal-content" @click.stop>
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
        <button @click="$emit('cancel')" class="cancel-btn">Cancel</button>
        <button 
          @click="$emit('confirm')" 
          :disabled="isConfirmDisabled || loading"
          class="confirm-btn"
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

const emit = defineEmits(['confirm', 'cancel']);

const confirmationValue = ref('');
const confirmationId = ref(`confirm-${Date.now()}`);

const isConfirmDisabled = computed(() => {
  if (!props.requireConfirmation) return false;
  return confirmationValue.value !== props.confirmationText;
});
</script>

<style scoped>
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

.cancel-btn {
  background: #f44336;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.cancel-btn:hover {
  background: #d32f2f;
}

.confirm-btn {
  background: #d32f2f;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.confirm-btn:hover:not(:disabled) {
  background: #b71c1c;
}

.confirm-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}
</style> 