<!--
  DataItem.vue
  Shared component for displaying label/value pairs in data panels.
-->
<template>
  <div class="data-item">
    <label class="data-item-label">{{ label }}</label>
    <span class="data-item-value">{{ value || 'N/A' }}</span>
  </div>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    required: true
  },
  value: {
    type: [String, Number],
    default: null
  }
});
</script>

<style scoped>
.data-item {
  display: flex;
  align-items: center;
  padding: var(--space-xs) 0;
}

.data-item-label {
  font-weight: var(--font-weight-semibold);
  min-width: 140px;
  color: var(--color-text-muted);
}

.data-item-value {
  color: var(--color-text);
  flex: 1;
}

/* Mobile responsive */
@media (max-width: 768px) {
  .data-item {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-xs);
  }

  .data-item-label {
    min-width: auto;
  }
}
</style>

