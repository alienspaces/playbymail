<!--
  StudioItemsView.vue
  This component follows the same pattern as StudioLocationsView.vue and StudioCreaturesView.vue.
-->
<template>
  <div>
    <div v-if="!gameId">
      <p>Select a game to manage items.</p>
    </div>
    <div v-else>
      <h2>Game Items</h2>
      <div v-if="itemsStore.loading">Loading...</div>
      <div v-else-if="itemsStore.error">Error: {{ itemsStore.error }}</div>
      <div v-else>
        <table v-if="itemsStore.items.length">
          <thead>
            <tr>
              <th>Name</th>
              <th>Description</th>
              <th>Created</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in itemsStore.items" :key="item.id">
              <td>{{ item.name }}</td>
              <td>{{ item.description }}</td>
              <td>{{ item.created_at }}</td>
              <td>
                <!-- Edit/Delete actions to be implemented -->
                <button>Edit</button>
                <button>Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-else>No items found.</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { watch, computed } from 'vue';
import { useItemsStore } from '../stores/items';
import { useGamesStore } from '../stores/games';
import { storeToRefs } from 'pinia';

const itemsStore = useItemsStore();
const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

const gameId = computed(() => selectedGame.value ? selectedGame.value.id : null);

watch(
  () => gameId.value,
  (newGameId) => {
    if (newGameId) itemsStore.loadItems(newGameId);
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
  border: 1px solid #ccc;
  padding: 0.5em;
}
</style> 