<template>
  <div class="page-header">
    <h1>{{ title }}</h1>
    <button v-if="actionText" @click="$emit('action')" class="action-btn">
      {{ actionText }}
    </button>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  actionText: {
    type: String,
    default: ''
  }
});

defineEmits(['action']);
</script>

<style scoped>
.page-header {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: var(--space-md);
  margin-bottom: var(--space-lg);
}

.page-header h1 {
  margin: 0;
}

.action-btn {
  background: var(--color-primary);
  color: var(--color-text-light);
  border: none;
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  font-weight: var(--font-weight-bold);
  cursor: pointer;
  transition: background 0.2s;
}

.action-btn:hover {
  background: var(--color-primary-dark);
}

@media (max-width: 768px) {
  .page-header {
    gap: var(--space-sm);
  }
}
</style> 