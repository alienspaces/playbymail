<template>
  <div class="resource-table-section" data-testid="resource-table">
    <table v-if="rows && rows.length">
      <thead>
        <tr>
          <th v-for="col in columns" :key="col.key">{{ col.label }}</th>
          <th v-if="$slots.actions">Actions</th>
        </tr>
      </thead>
      <tbody>
        <template v-for="row in rows" :key="row.id">
          <tr :data-testid="`table-row-${row.id}`">
            <td v-for="col in columns" :key="col.key">
              <slot :name="`cell-${col.key}`" :row="row" :column="col">
                {{ row[col.key] }}
              </slot>
            </td>
            <td v-if="$slots.actions">
              <slot name="actions" :row="row" />
            </td>
          </tr>
          <tr v-if="$slots['row-detail']" :key="`${row.id}-detail`" class="row-detail-row">
            <td :colspan="columns.length + ($slots.actions ? 1 : 0)" class="row-detail-cell">
              <slot name="row-detail" :row="row" />
            </td>
          </tr>
        </template>
      </tbody>
    </table>
    <div v-else-if="loading" data-testid="table-loading">Loading...</div>
    <div v-else-if="error" class="error" data-testid="table-error">{{ error }}</div>
    <div v-else data-testid="table-empty">No records found.</div>
  </div>
</template>

<script setup>
defineProps({
  columns: Array,
  rows: Array,
  loading: Boolean,
  error: String
});
</script>

<style scoped>
.resource-table-section {
  width: 100%;
}

.row-detail-row > .row-detail-cell { padding: 0; background: transparent; border-bottom: none; }

/* Actions column styling - wider to accommodate inline buttons */
.resource-table-section th:last-child,
.resource-table-section td:last-child {
  text-align: center;
  width: auto;
  min-width: 120px;
  white-space: nowrap;
  vertical-align: middle;
  padding-left: var(--space-sm);
  padding-right: var(--space-sm);
}
</style> 