<template>
  <div class="section-header">
    <h2>{{ title }}</h2>
    <button @click="$emit('create')" class="create-btn">
      Create New {{ resourceName }}
    </button>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  resourceName: {
    type: String,
    required: true
  }
});

defineEmits(['create']);
</script>

<style scoped>
.section-header {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: var(--space-md);
  margin-bottom: var(--space-lg);
}

.section-header h2 {
  margin: 0;
}

.create-btn {
  background: var(--color-primary);
  color: var(--color-text-light);
  border: none;
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  font-weight: var(--font-weight-bold);
  cursor: pointer;
  transition: background 0.2s;
}

.create-btn:hover {
  background: var(--color-primary-dark);
}

@media (max-width: 768px) {
  .section-header {
    gap: var(--space-sm);
  }
}
</style> 