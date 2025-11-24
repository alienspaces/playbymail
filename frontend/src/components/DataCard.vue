<!--
  DataCard.vue
  Reusable card component with title, data/content, and standardized buttons.
-->
<template>
  <div class="data-card" :class="{ 
    'data-card-completed': completed,
    'data-card-danger': variant === 'danger'
  }">
    <div v-if="$slots.header || $slots.title || title || $slots['header-extra'] || $slots.actions" class="card-header">
      <div class="card-header-left">
        <slot name="title">
          <h3 v-if="title">{{ title }}</h3>
        </slot>
        <slot name="header-extra"></slot>
      </div>
      <div v-if="$slots.actions" class="card-actions">
        <slot name="actions"></slot>
      </div>
    </div>

    <div v-if="$slots.default" class="card-content">
      <slot></slot>
    </div>

    <div v-if="$slots.footer || $slots.primary || $slots.secondary || $slots.tertiary" class="card-footer">
      <slot name="footer"></slot>
      <div v-if="$slots.primary || $slots.secondary || $slots.tertiary" class="card-footer-actions">
        <slot name="primary"></slot>
        <slot name="secondary"></slot>
        <slot name="tertiary"></slot>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: ''
  },
  completed: {
    type: Boolean,
    default: false
  },
  variant: {
    type: String,
    default: '',
    validator: (value) => ['', 'danger', 'warning', 'success', 'info'].includes(value)
  }
});
</script>

<style scoped>
.data-card {
  background: var(--color-bg-light);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-lg) var(--space-lg) calc(var(--space-lg) / 2);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.data-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.data-card-completed {
  opacity: 0.7;
}

.data-card-danger {
  border-color: var(--color-danger);
  background: var(--color-danger-light);
}

.data-card-danger .card-header {
  border-bottom-color: var(--color-danger);
}

.data-card-danger .card-header h3 {
  color: var(--color-danger);
}

.data-card-danger .card-footer {
  border-top-color: var(--color-danger);
}

.data-card-danger .card-content {
  color: var(--color-text);
}

.data-card-danger .card-content p {
  color: var(--color-text);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-md);
  padding-bottom: var(--space-md);
  border-bottom: 1px solid var(--color-border);
}

.card-header-left {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  flex: 1;
  min-width: 0;
}

.card-header h3 {
  margin: 0;
  font-size: var(--font-size-lg);
  color: var(--color-text);
  font-weight: var(--font-weight-bold);
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-actions {
  display: flex;
  gap: var(--space-sm);
  flex-wrap: wrap;
  flex-shrink: 0;
}

.card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  min-height: 0;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-top: var(--space-md);
  padding-top: var(--space-md);
  border-top: 1px solid var(--color-border);
}

.card-footer-actions {
  display: flex;
  gap: var(--space-sm);
  flex-wrap: wrap;
  flex-shrink: 0;
}

/* Mobile responsive */
@media (max-width: 768px) {
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-md);
  }

  .card-actions {
    width: 100%;
    justify-content: flex-start;
  }

  .card-footer {
    flex-direction: column;
    align-items: stretch;
    gap: var(--space-md);
  }

  .card-footer-actions {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
