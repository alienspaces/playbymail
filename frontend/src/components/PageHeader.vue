<template>
  <div class="page-header">
    <component :is="titleTag" class="hand-drawn-title" :class="titleClass">
      <HandDrawnIcon v-if="showIcon" :type="iconType" :color="iconColor" class="title-icon" />
      {{ title }}
    </component>
    <button v-if="actionText" @click="$emit('action')" class="action-btn">
      {{ actionText }}
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import HandDrawnIcon from './HandDrawnIcon.vue';

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  actionText: {
    type: String,
    default: ''
  },
  iconType: {
    type: String,
    default: 'star' // 'star', 'circle', 'checkbox-checked', etc.
  },
  iconColor: {
    type: String,
    default: 'blue' // 'blue', 'black', 'red'
  },
  showIcon: {
    type: Boolean,
    default: true // Show icon by default
  },
  titleLevel: {
    type: String,
    default: 'h1', // 'h1' for main titles, 'h2' for section titles
    validator: (value) => ['h1', 'h2'].includes(value)
  }
});

const titleTag = computed(() => props.titleLevel);
const titleClass = computed(() => {
  return props.titleLevel === 'h2' ? 'section-title' : '';
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

.page-header h1,
.page-header h2 {
  margin: 0;
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  position: relative;
}

/* Section titles (h2) should be smaller than main titles (h1) */
.page-header h2.section-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
  text-transform: uppercase;
  letter-spacing: 0.8px;
  margin-bottom: var(--space-sm);
}

.hand-drawn-title {
  position: relative;
}

.title-icon {
  font-size: 0.8em;
  margin-right: 0.2em;
}

.action-btn {
  background: transparent;
  color: var(--color-button);
  border: 2px solid var(--color-button);
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  font-weight: var(--font-weight-bold);
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background: var(--color-button);
  color: var(--color-text-light);
}

@media (max-width: 768px) {
  .page-header {
    gap: var(--space-sm);
  }
}
</style> 