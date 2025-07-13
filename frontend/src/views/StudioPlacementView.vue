<template>
  <div>
    <div v-if="!selectedGame">
      <p>Please select or create a game to manage placement.</p>
    </div>
    <div v-else>
      <h1>Placement</h1>
      <p>Assign items and creatures to locations here.</p>
      <!-- TODO: Add placement UI -->
    </div>
  </div>
</template>

<script setup>
import { useGamesStore } from '../stores/games';
import { storeToRefs } from 'pinia';

const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);
</script>

<style scoped>
h1 {
  margin-bottom: 1rem;
}
</style> 