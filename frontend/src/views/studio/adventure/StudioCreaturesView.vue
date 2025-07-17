<!--
  StudioCreaturesView.vue
  This component follows the same pattern as StudioLocationsView.vue and StudioItemsView.vue.
-->
<template>
  <div>
    <div v-if="!gameId">
      <p>Select a game to manage creatures.</p>
    </div>
    <div v-else>
      <h2>Game Creatures</h2>
      <div v-if="creaturesStore.loading">Loading...</div>
      <div v-else-if="creaturesStore.error">Error: {{ creaturesStore.error }}</div>
      <div v-else>
        <table v-if="creaturesStore.creatures.length">
          <thead>
            <tr>
              <th>Name</th>
              <th>Description</th>
              <th>Created</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="creature in creaturesStore.creatures" :key="creature.id">
              <td>{{ creature.name }}</td>
              <td>{{ creature.description }}</td>
              <td>{{ creature.created_at }}</td>
              <td>
                <!-- Edit/Delete actions to be implemented -->
                <button>Edit</button>
                <button>Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-else>No creatures found.</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { watch, computed } from 'vue';
import { useCreaturesStore } from '../../../stores/creatures';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';

const creaturesStore = useCreaturesStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const gameId = computed(() => selectedGame.value ? selectedGame.value.id : null);

watch(
  () => gameId.value,
  (newGameId) => {
    if (newGameId) creaturesStore.loadCreatures(newGameId);
  },
  { immediate: true }
);
</script>

<style scoped>
table {
  width: 100%;
  border-collapse: collapse;
}
th, td {
  border: 1px solid var(--color-border);
  padding: var(--space-sm);
}
</style> 