<template>
  <p v-if="gameName" class="game-context-name">
    Game: {{ gameName }}
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
.game-context-name {
  font-size: 1.1rem;
  font-weight: 400;
  color: #444;
  margin-bottom: var(--space-sm);
}
</style> 