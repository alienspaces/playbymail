<template>
  <div class="page-header">
    <component :is="titleTag" class="hand-drawn-title" :class="titleClass">
      <HandDrawnIcon v-if="showIcon" :type="iconType" :color="iconColor" class="title-icon" />
      {{ title }}
    </component>
    <Button v-if="actionText" @click="$emit('action')" variant="primary" size="small" :disabled="actionDisabled" data-testid="page-header-action">
      {{ actionText }}
    </Button>
    <p v-if="subtitle" class="page-subtitle">
      {{ subtitle }}
    </p>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import HandDrawnIcon from './HandDrawnIcon.vue';
import Button from './Button.vue';

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  actionText: {
    type: String,
    default: ''
  },
  actionDisabled: {
    type: Boolean,
    default: false
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
    default: 'h1', // 'h1' for main titles, 'h2' for section titles, 'h3' for subsections
    validator: (value) => ['h1', 'h2', 'h3'].includes(value)
  },
  subtitle: {
    type: String,
    default: ''
  }
});

const titleTag = computed(() => props.titleLevel);
const titleClass = computed(() => {
  if (props.titleLevel === 'h2') return 'section-title';
  if (props.titleLevel === 'h3') return 'subsection-title';
  return '';
});

defineEmits(['action']);
</script>

<style scoped>
/* Default page header - for h1 and h2 (main page titles) */
.page-header {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: var(--space-md);
  margin-bottom: calc(var(--space-md) / 2);
}

.page-header h1,
.page-header h2,
.page-header h3 {
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

/* Subsection titles (h3) - smaller with tighter spacing */
.page-header h3.subsection-title {
  font-size: var(--font-size-md);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* Tighter spacing for h3 subsection headers */
.page-header:has(h3) {
  gap: var(--space-xs);
  margin-bottom: var(--space-sm);
}

.hand-drawn-title {
  position: relative;
}

.page-subtitle {
  margin: 0;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.title-icon {
  font-size: 0.8em;
  margin-right: 0.2em;
}

@media (max-width: 768px) {
  .page-header {
    gap: var(--space-sm);
    margin-bottom: calc(var(--space-xl) / 2);
  }
}
</style>