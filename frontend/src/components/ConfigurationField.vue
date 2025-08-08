<template>
  <div class="configuration-field">
    <label :for="fieldId" class="field-label">
      {{ config.description || config.config_key }}
      <span v-if="config.is_required" class="required">*</span>
    </label>
    
    <!-- Boolean field (checkbox) -->
    <div v-if="config.value_type === 'boolean'" class="field-input">
      <input
        :id="fieldId"
        v-model="fieldValue"
        type="checkbox"
        :required="config.is_required"
        @change="handleChange"
      />
      <span class="checkbox-label">{{ config.description || config.config_key }}</span>
    </div>

    <!-- Number field -->
    <div v-else-if="config.value_type === 'integer'" class="field-input">
      <input
        :id="fieldId"
        v-model.number="fieldValue"
        type="number"
        :required="config.is_required"
        :min="getMinValue()"
        :max="getMaxValue()"
        :step="getStepValue()"
        @input="handleChange"
      />
      <span v-if="config.ui_hint === 'percentage'" class="input-suffix">%</span>
    </div>

    <!-- Select field -->
    <div v-else-if="config.ui_hint === 'select'" class="field-input">
      <select
        :id="fieldId"
        v-model="fieldValue"
        :required="config.is_required"
        @change="handleChange"
      >
        <option value="" disabled>Select an option</option>
        <option v-for="option in getSelectOptions()" :key="option.value" :value="option.value">
          {{ option.label }}
        </option>
      </select>
    </div>

    <!-- Text area for long text -->
    <div v-else-if="config.ui_hint === 'textarea'" class="field-input">
      <textarea
        :id="fieldId"
        v-model="fieldValue"
        :required="config.is_required"
        :rows="4"
        :maxlength="getMaxLength()"
        @input="handleChange"
      />
    </div>

    <!-- Default text input -->
    <div v-else class="field-input">
      <input
        :id="fieldId"
        v-model="fieldValue"
        type="text"
        :required="config.is_required"
        :maxlength="getMaxLength()"
        :placeholder="getPlaceholder()"
        @input="handleChange"
      />
    </div>

    <!-- Help text -->
    <p v-if="config.description" class="help-text">
      {{ config.description }}
    </p>

    <!-- Validation error -->
    <p v-if="validationError" class="error-text">
      {{ validationError }}
    </p>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';

const props = defineProps({
  config: {
    type: Object,
    required: true
  },
  modelValue: {
    type: [String, Number, Boolean],
    default: null
  },
  fieldId: {
    type: String,
    required: true
  }
});

const emit = defineEmits(['update:modelValue', 'validation-error']);

const fieldValue = ref(props.modelValue);
const validationError = ref('');

// Watch for external value changes
watch(() => props.modelValue, (newValue) => {
  fieldValue.value = newValue;
});

// Watch for internal value changes
watch(fieldValue, (newValue) => {
  emit('update:modelValue', newValue);
  validateField(newValue);
});

// Initialize with default value if no value provided
if (fieldValue.value === null && props.config.default_value) {
  fieldValue.value = parseValue(props.config.default_value, props.config.value_type);
}

function parseValue(value, type) {
  if (value === null || value === undefined) return null;
  
  switch (type) {
    case 'boolean':
      return value === 'true' || value === true;
    case 'integer':
      return parseInt(value, 10);
    default:
      return value;
  }
}

function handleChange() {
  validateField(fieldValue.value);
}

function validateField(value) {
  validationError.value = '';
  
  // Required field validation
  if (props.config.is_required && (value === null || value === undefined || value === '')) {
    validationError.value = 'This field is required';
    emit('validation-error', validationError.value);
    return;
  }

  // Type-specific validation
  if (value !== null && value !== undefined && value !== '') {
    switch (props.config.value_type) {
      case 'integer':
        if (isNaN(value) || !Number.isInteger(Number(value))) {
          validationError.value = 'Must be a whole number';
        } else {
          const numValue = Number(value);
          if (props.config.ui_hint === 'percentage' && (numValue < 0 || numValue > 100)) {
            validationError.value = 'Percentage must be between 0 and 100';
          }
        }
        break;
      case 'boolean':
        if (typeof value !== 'boolean') {
          validationError.value = 'Must be true or false';
        }
        break;
    }
  }

  emit('validation-error', validationError.value);
}

function getMinValue() {
  if (props.config.ui_hint === 'percentage') return 0;
  if (props.config.ui_hint === 'number') return 1;
  return undefined;
}

function getMaxValue() {
  if (props.config.ui_hint === 'percentage') return 100;
  return undefined;
}

function getStepValue() {
  if (props.config.ui_hint === 'percentage') return 1;
  return undefined;
}

function getMaxLength() {
  // Default max length for text fields
  return 255;
}

function getPlaceholder() {
  if (props.config.ui_hint === 'percentage') return 'Enter percentage (0-100)';
  if (props.config.ui_hint === 'number') return 'Enter number';
  return `Enter ${props.config.config_key}`;
}

function getSelectOptions() {
  // Parse validation rules for select options
  if (props.config.validation_rules) {
    try {
      const rules = JSON.parse(props.config.validation_rules);
      if (rules.options) {
        return rules.options;
      }
          } catch {
        // Invalid JSON, ignore
      }
  }
  
  // Default options for common select fields
  switch (props.config.config_key) {
    case 'combat_difficulty':
      return [
        { value: 'easy', label: 'Easy' },
        { value: 'normal', label: 'Normal' },
        { value: 'hard', label: 'Hard' },
        { value: 'expert', label: 'Expert' }
      ];
    case 'magic_enabled':
      return [
        { value: true, label: 'Enabled' },
        { value: false, label: 'Disabled' }
      ];
    default:
      return [];
  }
}
</script>

<style scoped>
.configuration-field {
  margin-bottom: var(--space-md);
}

.field-label {
  display: block;
  margin-bottom: var(--space-sm);
  font-weight: 500;
  color: var(--color-text);
}

.required {
  color: var(--color-error);
  margin-left: var(--space-xs);
}

.field-input {
  position: relative;
  display: flex;
  align-items: center;
}

.field-input input[type="text"],
.field-input input[type="number"],
.field-input select,
.field-input textarea {
  flex: 1;
  padding: var(--space-sm);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-md);
  background: var(--color-bg);
  color: var(--color-text);
}

.field-input input[type="checkbox"] {
  margin-right: var(--space-sm);
}

.checkbox-label {
  margin-left: var(--space-sm);
  color: var(--color-text);
}

.input-suffix {
  margin-left: var(--space-sm);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.help-text {
  margin-top: var(--space-xs);
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}

.error-text {
  margin-top: var(--space-xs);
  font-size: var(--font-size-sm);
  color: var(--color-error);
}

.field-input input:focus,
.field-input select:focus,
.field-input textarea:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px rgba(var(--color-primary-rgb), 0.2);
}
</style> 