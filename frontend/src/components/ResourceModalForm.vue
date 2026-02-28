<template>
  <!-- Teleport modal to body to avoid z-index stacking context issues -->
  <Teleport to="body">
    <div v-if="visible" class="modal-overlay" @click.self="$emit('cancel')">
      <div class="modal">
        <h2>{{ mode === 'create' ? `Create ${title}` : `Edit ${title}` }}</h2>
        <form @submit.prevent="handleSubmit" class="modal-form">
          <div v-for="field in fields" :key="field.key" class="form-group">
            <label v-if="field.type !== 'checkbox' && field.type !== 'info'" :for="field.key">
              {{ field.label }}<span v-if="field.required" class="required"> *</span>
            </label>
            <slot name="field" :field="field" :value="form[field.key]" :update="val => form[field.key] = val">
              <!-- Render select for select type -->
              <select v-if="field.type === 'select'" v-model="form[field.key]" :id="field.key"
                :required="field.required" :placeholder="field.placeholder" class="form-select">
                <option v-if="field.placeholder" value="" disabled>{{ field.placeholder }}</option>
                <option v-for="option in getFieldOptions(field)" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
              <!-- Render textarea for textarea type -->
              <TextareaWithCounter v-else-if="field.type === 'textarea'" :id="field.key" v-model="form[field.key]"
                :max-length="field.maxlength" :required="field.required" :placeholder="field.placeholder"
                :rows="field.rows || 4" />
              <!-- Render checkbox for checkbox type -->
              <div v-else-if="field.type === 'checkbox'" class="checkbox-group">
                <input v-model="form[field.key]" :id="field.key" type="checkbox" :required="field.required" />
                <label :for="field.key" class="checkbox-label">{{ field.checkboxLabel || field.label }}{{ field.required
                  ?
                  ' *' : '' }}</label>
              </div>
              <!-- Render informational notice -->
              <div v-else-if="field.type === 'info'" class="info-notice">
                <p>{{ field.text }}</p>
              </div>
              <!-- Render input for other types -->
              <input v-else v-model="form[field.key]" :id="field.key" :type="field.type || 'text'"
                :required="field.required" :maxlength="field.maxlength" :placeholder="field.placeholder"
                :min="field.min" :max="field.max" autocomplete="off" />
            </slot>
          </div>
          <div class="modal-actions">
            <button type="submit" data-testid="modal-submit">{{ mode === 'create' ? 'Create' : 'Save' }}</button>
            <button type="button" @click="$emit('cancel')" data-testid="modal-cancel">Cancel</button>
          </div>
        </form>
        <div v-if="error" class="error" data-testid="modal-error">
          <p>{{ error }}</p>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { reactive, watch } from 'vue';
import TextareaWithCounter from './TextareaWithCounter.vue';
const props = defineProps({
  visible: Boolean,
  mode: String, // 'create' or 'edit'
  title: String,
  fields: Array,
  modelValue: Object,
  error: String,
  options: Object // New prop for field options
});
const emit = defineEmits(['submit', 'cancel']);
const form = reactive({});

// Watch for modelValue changes to sync form data
watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      Object.assign(form, val);
    }
  },
  { immediate: true, deep: true }
);

// Watch for modal visibility to reset form when it opens
watch(
  () => props.visible,
  (isVisible) => {
    if (isVisible && props.modelValue) {
      // Reset form with modelValue when modal opens
      Object.keys(form).forEach(key => delete form[key]);
      Object.assign(form, props.modelValue);
    }
  }
);

function getFieldOptions(field) {
  if (field.options) {
    return field.options;
  }
  if (props.options && props.options[field.key]) {
    return props.options[field.key];
  }
  return [];
}

function handleSubmit() {
  emit('submit', { ...form });
}
</script>

<style scoped>
.modal-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
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

.info-notice {
  background: var(--color-bg-light);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: var(--space-sm) var(--space-md);
}

.info-notice p {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
  line-height: 1.5;
}

.checkbox-group {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-xs) 0;
}

.checkbox-group input[type="checkbox"] {
  width: 1.5em;
  height: 1.5em;
  margin: 0;
  flex-shrink: 0;
  cursor: pointer;
  min-width: 1.5em;
  min-height: 1.5em;
}

.checkbox-label {
  margin: 0;
  font-weight: var(--font-weight-normal);
  color: var(--color-text);
  line-height: 1.5;
}

@media (max-width: 768px) {
  .modal {
    padding: var(--space-lg);
    width: 95%;
  }

  .modal-actions {
    flex-direction: column-reverse;
  }

  .modal-actions button {
    width: 100%;
  }
}
</style>