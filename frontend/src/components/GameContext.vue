<template>
  <p v-if="gameName" class="game-context">
    <span class="game-context-label">Game:</span>&nbsp;
    <span class="game-context-name">{{ gameName }}</span>
  </p>
</template>

<script setup>
defineProps({
  gameName: {
    type: String,
    default: ''
  }
});
</script>

<style scoped>
.game-context {
  margin: 0 0 var(--space-md) 0;
  font-size: var(--font-size-base);
  color: var(--color-text);
  font-weight: var(--font-weight-normal);
}

.game-context-label {
  font-weight: var(--font-weight-semibold);
  margin-right: var(--space-xs);
}

.game-context-name {
  font-weight: var(--font-weight-semibold);
}
</style>
