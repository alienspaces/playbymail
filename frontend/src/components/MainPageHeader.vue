<template>
  <div class="main-page-header">
    <h1 class="hand-drawn-title">
      <HandDrawnIcon :type="iconType" :color="iconColor" class="title-icon" />
      {{ title }}
    </h1>
    <div v-if="$slots.actions" class="header-actions">
      <slot name="actions"></slot>
    </div>
  </div>
</template>

<script setup>
import HandDrawnIcon from './HandDrawnIcon.vue';

defineProps({
  title: {
    type: String,
    required: true
  },
  iconType: {
    type: String,
    required: true
  },
  iconColor: {
    type: String,
    default: 'blue'
  }
});
</script>

<style scoped>
.main-page-header {
  width: 100%;
  background: var(--color-bg);
  border-bottom: 1px solid var(--color-border);
  padding: var(--space-lg) var(--space-lg);
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  flex-shrink: 0;
}

.main-page-header h1 {
  margin: 0;
  font-size: var(--font-size-xl);
  font-family: var(--font-family-heading);
  font-weight: 700;
  color: var(--color-text);
  text-transform: uppercase;
  letter-spacing: 1px;
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

.header-actions {
  display: flex;
  gap: var(--space-md);
}

@media (max-width: 768px) {
  .main-page-header {
    padding: var(--space-md) var(--space-lg);
  }
  
  .main-page-header h1 {
    font-size: var(--font-size-lg);
  }
}
</style>

