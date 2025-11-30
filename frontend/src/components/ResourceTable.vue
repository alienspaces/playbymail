<template>
  <div class="resource-table-section">
    <table v-if="rows && rows.length">
      <thead>
        <tr>
          <th v-for="col in columns" :key="col.key">{{ col.label }}</th>
          <th v-if="$slots.actions">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in rows" :key="row.id">
          <td v-for="col in columns" :key="col.key">
            <slot :name="`cell-${col.key}`" :row="row" :column="col">
              {{ row[col.key] }}
            </slot>
          </td>
          <td v-if="$slots.actions">
            <slot name="actions" :row="row" />
          </td>
        </tr>
      </tbody>
    </table>
    <div v-else-if="loading">Loading...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else>No records found.</div>
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