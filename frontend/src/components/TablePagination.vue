<template>
  <div v-if="pageNumber > 1 || hasMore" class="table-pagination" data-testid="table-pagination">
    <Button
      variant="primary"
      size="small"
      :disabled="pageNumber <= 1"
      @click="$emit('page-change', pageNumber - 1)"
      data-testid="pagination-prev"
    >Previous</Button>
    <span class="page-indicator" data-testid="pagination-page">Page {{ pageNumber }}</span>
    <Button
      variant="primary"
      size="small"
      :disabled="!hasMore"
      @click="$emit('page-change', pageNumber + 1)"
      data-testid="pagination-next"
    >Next</Button>
  </div>
</template>

<script setup>
import Button from './Button.vue';

defineProps({
  pageNumber: {
    type: Number,
    default: 1,
  },
  hasMore: {
    type: Boolean,
    default: false,
  },
});

defineEmits(['page-change']);
</script>

<style scoped>
.table-pagination {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  margin-top: var(--space-md);
}

.page-indicator {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
  min-width: 4rem;
  text-align: center;
}
</style>
