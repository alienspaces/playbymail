<template>
  <div v-if="visible" class="modal-overlay">
    <div class="modal">
      <h2>{{ mode === 'create' ? `Create ${title}` : `Edit ${title}` }}</h2>
      <form @submit.prevent="handleSubmit">
        <div v-for="field in fields" :key="field.key" class="form-group">
          <label :for="field.key">{{ field.label }}</label>
          <slot name="field" :field="field" :value="form[field.key]" :update="val => form[field.key] = val">
            <input
              v-model="form[field.key]"
              :id="field.key"
              :type="field.type || 'text'"
              :required="field.required"
              :maxlength="field.maxlength"
              :placeholder="field.placeholder"
            />
          </slot>
        </div>
        <div class="modal-actions">
          <button type="submit">{{ mode === 'create' ? 'Create' : 'Save' }}</button>
          <button type="button" @click="$emit('cancel')">Cancel</button>
        </div>
      </form>
      <p v-if="error" class="error">{{ error }}</p>
    </div>
  </div>
</template>

<script setup>
import { reactive, watch } from 'vue';
const props = defineProps({
  visible: Boolean,
  mode: String, // 'create' or 'edit'
  title: String,
  fields: Array,
  modelValue: Object,
  error: String
});
const emit = defineEmits(['submit', 'cancel']);
const form = reactive({});
watch(
  () => props.modelValue,
  (val) => {
    if (val) Object.assign(form, val);
  },
  { immediate: true, deep: true }
);
function handleSubmit() {
  emit('submit', { ...form });
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal {
  background: var(--color-bg);
  padding: var(--space-lg);
  border-radius: var(--radius-md);
  min-width: 300px;
  max-width: 90vw;
  box-shadow: 0 2px 16px rgba(0,0,0,0.2);
}
.form-group {
  display: flex;
  flex-direction: column;
  margin-bottom: var(--space-md);
}
label {
  margin-bottom: var(--space-xs);
  font-weight: 500;
}
input {
  padding: var(--space-sm);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-base);
}
.modal-actions {
  margin-top: var(--space-md);
  display: flex;
  gap: var(--space-md);
  justify-content: flex-start;
}
.error {
  color: var(--color-error);
  margin-top: var(--space-md);
}
</style> 