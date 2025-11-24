<template>
  <div class="section-header">
    <h2 class="hand-drawn-title">
      <HandDrawnIcon :type="iconType" :color="iconColor" class="title-icon" />
      {{ title }}
    </h2>
    <button @click="$emit('create')" class="create-btn">
      Create New {{ resourceName }}
    </button>
  </div>
</template>

<script setup>
import HandDrawnIcon from './HandDrawnIcon.vue';

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  resourceName: {
    type: String,
    required: true
  },
  iconType: {
    type: String,
    default: 'circle' // 'star', 'circle', 'checkbox-checked', etc.
  },
  iconColor: {
    type: String,
    default: 'blue' // 'blue', 'black', 'red'
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
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  position: relative;
}

.hand-drawn-title {
  position: relative;
}

.title-icon {
  font-size: 0.8em;
  margin-right: 0.2em;
}

.create-btn {
  background: transparent;
  color: var(--color-button);
  border: 2px solid var(--color-button);
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  font-weight: var(--font-weight-bold);
  cursor: pointer;
  transition: all 0.2s;
}

.create-btn:hover {
  background: var(--color-button);
  color: var(--color-text-light);
}

@media (max-width: 768px) {
  .section-header {
    gap: var(--space-sm);
  }
}
</style> 