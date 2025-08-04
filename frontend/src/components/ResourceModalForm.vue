<template>
  <div v-if="visible" class="modal-overlay">
    <div class="modal">
      <h2>{{ mode === 'create' ? `Create ${title}` : `Edit ${title}` }}</h2>
      <form @submit.prevent="handleSubmit">
        <div v-for="field in fields" :key="field.key" class="form-group">
          <label :for="field.key">{{ field.label }}{{ field.required ? ' *' : '' }}</label>
          <slot name="field" :field="field" :value="form[field.key]" :update="val => form[field.key] = val">
            <!-- Render textarea for textarea type -->
            <textarea
              v-if="field.type === 'textarea'"
              v-model="form[field.key]"
              :id="field.key"
              :required="field.required"
              :maxlength="field.maxlength"
              :placeholder="field.placeholder"
              :rows="field.rows || 4"
            />
            <!-- Render input for other types -->
            <input
              v-else
              v-model="form[field.key]"
              :id="field.key"
              :type="field.type || 'text'"
              :required="field.required"
              :maxlength="field.maxlength"
              :placeholder="field.placeholder"
              :min="field.min"
              :max="field.max"
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
/* Component-specific styles only */
</style> 