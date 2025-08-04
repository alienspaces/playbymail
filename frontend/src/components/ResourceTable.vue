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
          <td v-for="col in columns" :key="col.key">{{ row[col.key] }}</td>
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

/* Add spacing between action buttons */
.resource-table-section td:last-child {
  display: flex;
  gap: 8px;
  align-items: center;
}

.resource-table-section td:last-child button {
  margin-right: 8px;
}

.resource-table-section td:last-child button:last-child {
  margin-right: 0;
}
</style> 