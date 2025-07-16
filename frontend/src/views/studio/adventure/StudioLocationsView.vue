<!--
  StudioLocationsView.vue
  This component follows the same pattern as StudioItemsView.vue and StudioCreaturesView.vue.
-->
<template>
  <div>
    <div v-if="!gameId">
      <p>Select a game to manage locations.</p>
    </div>
    <div v-else>
      <h2>Game Locations</h2>
      <div v-if="locationsStore.loading">Loading...</div>
      <div v-else-if="locationsStore.error">Error: {{ locationsStore.error }}</div>
      <div v-else>
        <table v-if="locationsStore.locations.length">
          <thead>
            <tr>
              <th>Name</th>
              <th>Description</th>
              <th>Created</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="location in locationsStore.locations" :key="location.id">
              <td>{{ location.name }}</td>
              <td>{{ location.description }}</td>
              <td>{{ location.created_at }}</td>
              <td>
                <!-- Edit/Delete actions to be implemented -->
                <button>Edit</button>
                <button>Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-else>No locations found.</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { watch, computed } from 'vue';
import { useLocationsStore } from '../../../stores/locations';
import { useGamesStore } from '../../../stores/games';
import { storeToRefs } from 'pinia';

const locationsStore = useLocationsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const gameId = computed(() => selectedGame.value ? selectedGame.value.id : null);

watch(
  () => gameId.value,
  (newGameId) => {
    if (newGameId) locationsStore.loadLocations(newGameId);
  },
  { immediate: true }
);
</script>

<style scoped>
.select-game-warning {
  color: #b00;
  margin-top: 2rem;
  font-size: 1.1em;
}
h1 {
  margin-bottom: 1rem;
}
table {
  width: 100%;
  border-collapse: collapse;
}
th, td {
  border: 1px solid #ccc;
  padding: 0.5em;
}
</style> 