<template>
  <div class="table-actions">
    <button
      v-for="action in actions"
      :key="action.key"
      @click="action.handler"
      class="action-btn"
      :class="{
        'action-btn-primary': action.primary,
        'action-btn-danger': action.danger
      }"
      :title="action.label"
      :aria-label="action.label"
      type="button"
    >
      <component :is="getIcon(action)" class="action-icon" />
      <span v-if="showLabels" class="action-label">{{ action.label }}</span>
    </button>
  </div>
</template>

<script setup>
import { h } from 'vue';

defineProps({
  actions: {
    type: Array,
    required: true,
    validator: (actions) => {
      return actions.every(action => 
        action.key && action.label && typeof action.handler === 'function'
      );
    }
  },
  showLabels: {
    type: Boolean,
    default: false
  }
});

// Icon components as render functions
const icons = {
  select: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
    h('path', { d: 'M9 5l7 7-7 7' })
  ]),
  preview: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
    h('path', { d: 'M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z' }),
    h('circle', { cx: '12', cy: '12', r: '3' })
  ]),
  view: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
    h('path', { d: 'M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z' }),
    h('circle', { cx: '12', cy: '12', r: '3' })
  ]),
  edit: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
    h('path', { d: 'M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7' }),
    h('path', { d: 'M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z' })
  ]),
  delete: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
    h('polyline', { points: '3 6 5 6 21 6' }),
    h('path', { d: 'M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2' }),
    h('line', { x1: '10', y1: '11', x2: '10', y2: '17' }),
    h('line', { x1: '14', y1: '11', x2: '14', y2: '17' })
  ]),
  start: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
    h('polygon', { points: '5 3 19 12 5 21 5 3' })
  ]),
  pause: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
    h('rect', { x: '6', y: '4', width: '4', height: '16' }),
    h('rect', { x: '14', y: '4', width: '4', height: '16' })
  ]),
  resume: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
    h('polygon', { points: '5 3 19 12 5 21 5 3' })
  ]),
  cancel: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
    h('circle', { cx: '12', cy: '12', r: '10' }),
    h('line', { x1: '15', y1: '9', x2: '9', y2: '15' }),
    h('line', { x1: '9', y1: '9', x2: '15', y2: '15' })
  ]),
  set: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
    h('line', { x1: '12', y1: '5', x2: '12', y2: '19' }),
    h('line', { x1: '5', y1: '12', x2: '19', y2: '12' })
  ]),
  remove: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
    h('line', { x1: '5', y1: '12', x2: '19', y2: '12' })
  ]),
  manage: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
    h('circle', { cx: '12', cy: '12', r: '3' }),
    h('path', { d: 'M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z' })
  ]),
  // Default fallback icon (circle with dot)
  default: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
    h('circle', { cx: '12', cy: '12', r: '10' }),
    h('circle', { cx: '12', cy: '12', r: '2', fill: 'currentColor' })
  ])
};

function getIcon(action) {
  // Use action.icon if specified, otherwise use action.key
  const iconKey = action.icon || action.key;
  return icons[iconKey] || icons.default;
}
</script>

<style scoped>
.table-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-xs);
  flex-wrap: nowrap;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-xs);
  padding: 6px;
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  color: var(--color-text-muted);
  transition: all 0.15s ease;
  min-width: 32px;
  height: 32px;
}

.action-btn:hover {
  background: var(--color-bg-alt);
  color: var(--color-button);
  border-color: var(--color-button);
}

.action-btn:focus {
  outline: 2px solid var(--color-primary);
  outline-offset: 1px;
}

/* Primary action style (e.g., Select) */
.action-btn-primary {
  background: var(--color-button);
  color: var(--color-text-light);
  border-color: var(--color-button);
}

.action-btn-primary:hover {
  background: var(--color-button-dark);
  border-color: var(--color-button-dark);
  color: var(--color-text-light);
}

/* Danger action style (e.g., Delete) */
.action-btn-danger {
  color: var(--color-danger);
  border-color: var(--color-danger-light);
}

.action-btn-danger:hover {
  background: var(--color-danger);
  color: var(--color-text-light);
  border-color: var(--color-danger);
}

.action-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.action-label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  white-space: nowrap;
}

/* When showing labels, add more padding */
.table-actions:has(.action-label) .action-btn {
  padding: 6px 10px;
}
</style>

