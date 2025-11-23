<template>
  <button 
    :class="buttonClasses" 
    :disabled="disabled"
    @click="$emit('click', $event)"
    v-bind="$attrs"
  >
    <slot />
  </button>
</template>

<script setup>
import { computed } from 'vue';

// Component name for linting
defineOptions({
  name: 'BaseButton'
});

const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'success', 'warning', 'danger', 'info'].includes(value)
  },
  size: {
    type: String,
    default: 'medium',
    validator: (value) => ['small', 'medium', 'large'].includes(value)
  },
  disabled: {
    type: Boolean,
    default: false
  }
});

defineEmits(['click']);

const buttonClasses = computed(() => [
  'btn',
  `btn-${props.variant}`,
  `btn-${props.size}`
]);
</script>

<style scoped>
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-xs);
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-family: inherit;
  font-weight: var(--font-weight-bold);
  text-decoration: none;
  transition: all 0.2s ease;
  white-space: nowrap;
  user-select: none;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn:not(:disabled):hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.btn:not(:disabled):active {
  transform: translateY(0);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}

/* Size variants */
.btn-small {
  padding: var(--space-xs) var(--space-sm);
  font-size: var(--font-size-sm);
  min-height: 32px;
}

.btn-medium {
  padding: var(--space-sm) var(--space-md);
  font-size: var(--font-size-base);
  min-height: 40px;
}

.btn-large {
  padding: var(--space-md) var(--space-lg);
  font-size: var(--font-size-md);
  min-height: 48px;
}

/* Color variants */
.btn-primary {
  background: transparent;
  color: var(--color-button);
  border: 2px solid var(--color-button);
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-button);
  color: var(--color-text-light);
}

.btn-secondary {
  background: var(--color-bg-light);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}

.btn-secondary:hover:not(:disabled) {
  background: var(--color-border);
}

.btn-success {
  background: var(--color-success);
  color: var(--color-text-light);
}

.btn-success:hover:not(:disabled) {
  background: var(--color-success-dark);
}

.btn-warning {
  background: var(--color-warning);
  color: var(--color-text-light);
}

.btn-warning:hover:not(:disabled) {
  background: var(--color-warning-dark);
}

.btn-danger {
  background: var(--color-danger);
  color: var(--color-text-light);
}

.btn-danger:hover:not(:disabled) {
  background: var(--color-danger-dark);
}

.btn-info {
  background: var(--color-info);
  color: var(--color-text-light);
}

.btn-info:hover:not(:disabled) {
  background: var(--color-info-dark);
}

/* Full width option */
.btn-full {
  width: 100%;
}

/* Icon button */
.btn-icon {
  padding: var(--space-xs);
  min-width: 40px;
  min-height: 40px;
}

.btn-icon.btn-small {
  min-width: 32px;
  min-height: 32px;
}

.btn-icon.btn-large {
  min-width: 48px;
  min-height: 48px;
}
</style>
