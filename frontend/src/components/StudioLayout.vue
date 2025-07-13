<template>
  <div class="studio-layout">
    <aside class="sidebar">
      <nav>
        <ul>
          <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/locations` : '#'" :class="{disabled: !selectedGame}" active-class="active">Locations</router-link></li>
          <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/items` : '#'" :class="{disabled: !selectedGame}" active-class="active">Items</router-link></li>
          <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/creatures` : '#'" :class="{disabled: !selectedGame}" active-class="active">Creatures</router-link></li>
          <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/placement` : '#'" :class="{disabled: !selectedGame}" active-class="active">Placement</router-link></li>
        </ul>
        <div v-if="!selectedGame" class="select-game-warning">Select a game to manage resources.</div>
      </nav>
    </aside>
    <div class="main-content">
      <header class="studio-header">
        <slot name="header"><h1>Game Designer Studio</h1></slot>
      </header>
      <section class="studio-body">
        <div v-if="!selectedGame" class="select-game-warning-main">
          <p>Please select or create a game to access management features.</p>
        </div>
        <router-view v-else />
      </section>
    </div>
  </div>
</template>

<script setup>
import { storeToRefs } from 'pinia';
import { useGamesStore } from '../stores/games';

const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);
</script>

<style scoped>
.studio-layout {
  display: flex;
  min-height: 100vh;
}
.sidebar {
  width: 200px;
  background: #f7f7f7;
  border-right: 1px solid #ddd;
  padding: 2rem 1rem;
}
.sidebar ul {
  list-style: none;
}
.sidebar li {
  margin-bottom: 1rem;
}
.sidebar a {
  color: #222;
  text-decoration: none;
  font-weight: 500;
}
.sidebar a.active {
  color: #1976d2;
}
.sidebar a.disabled {
  pointer-events: none;
  color: #aaa;
}
.select-game-warning {
  color: #b00;
  margin-top: 2rem;
  font-size: 0.95em;
}
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.studio-header {
  padding: 1.5rem 2rem 1rem 2rem;
  border-bottom: 1px solid #eee;
  background: #fff;
}
.studio-body {
  flex: 1;
  padding: 2rem;
  background: #fafbfc;
}
.select-game-warning-main {
  color: #b00;
  margin: 2rem auto;
  font-size: 1.1em;
  text-align: center;
}
</style> 