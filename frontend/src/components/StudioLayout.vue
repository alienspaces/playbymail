<template>
  <div class="studio-layout">
    <aside class="sidebar">
      <button class="create-game-btn" @click="openCreate">+ Create Game</button>
      <nav>
        <ul>
          <li><router-link to="/studio" active-class="active">Games</router-link></li>
          <template v-if="showAdventureMenu">
            <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/locations` : '#'" :class="{disabled: !selectedGame}" active-class="active">Locations</router-link></li>
            <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/items` : '#'" :class="{disabled: !selectedGame}" active-class="active">Items</router-link></li>
            <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/creatures` : '#'" :class="{disabled: !selectedGame}" active-class="active">Creatures</router-link></li>
            <li><router-link :to="selectedGame ? `/studio/${selectedGame.id}/placement` : '#'" :class="{disabled: !selectedGame}" active-class="active">Placement</router-link></li>
          </template>
        </ul>
        <div v-if="!selectedGame" class="select-game-warning">Select a game to manage resources.</div>
      </nav>
      <button class="help-btn" @click="showHelp = true">Help</button>
    </aside>
    <div class="main-content">
      <header class="studio-header">
        <slot name="header"><h1>Game Designer Studio</h1></slot>
      </header>
      <section class="studio-body">
        <router-view />
      </section>
    </div>
    <div v-if="showHelp" class="help-panel-overlay" @click.self="showHelp = false">
      <div class="help-panel">
        <button class="close-help" @click="showHelp = false">&times;</button>
        <h2>Studio Help</h2>
        <p>This is context-sensitive help for the current section. (Stub)</p>
      </div>
    </div>
    <!-- Create Game Modal -->
    <div v-if="showModal" class="modal-overlay">
      <div class="modal">
        <h2>Create Game</h2>
        <form @submit.prevent="createGame">
          <label>
            Name:
            <input v-model="modalForm.name" required maxlength="1024" />
          </label>
          <label>
            Type:
            <select v-model="modalForm.game_type" required>
              <option value="adventure">Adventure</option>
            </select>
          </label>
          <div class="modal-actions">
            <button type="submit">Create</button>
            <button type="button" @click="closeModal">Cancel</button>
          </div>
        </form>
        <p v-if="modalError" class="error">{{ modalError }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { storeToRefs } from 'pinia';
import { useGamesStore } from '../stores/games';
import { useRoute } from 'vue-router';

const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);
const showHelp = ref(false);
const route = useRoute();

// Show adventure menu only on /studio/:gameId/*
const showAdventureMenu = computed(() => {
  // Matches /studio/:gameId/locations, /studio/:gameId/items, etc.
  return (
    selectedGame.value &&
    selectedGame.value.game_type === 'adventure' &&
    /^\/studio\/[^/]+\//.test(route.path)
  );
});

// Create Game modal logic
const showModal = ref(false);
const modalForm = ref({ name: '', game_type: 'adventure' });
const modalError = ref('');

function openCreate() {
  modalForm.value = { name: '', game_type: 'adventure' };
  modalError.value = '';
  showModal.value = true;
}
function closeModal() {
  showModal.value = false;
  modalError.value = '';
}
async function createGame() {
  modalError.value = '';
  try {
    const created = await gamesStore.createGame({ name: modalForm.value.name, game_type: modalForm.value.game_type });
    closeModal();
    if (created && created.id) {
      gamesStore.setSelectedGame(created);
    }
  } catch (err) {
    modalError.value = err.message;
  }
}
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
.help-btn {
  margin-top: 2rem;
  width: 100%;
  padding: 0.75rem 0;
  background: #1976d2;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-weight: 600;
  cursor: pointer;
  font-size: 1rem;
}
.help-panel-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.25);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}
.help-panel {
  background: #fff;
  padding: 2rem 2.5rem;
  border-radius: 10px;
  min-width: 320px;
  max-width: 90vw;
  box-shadow: 0 2px 16px rgba(0,0,0,0.18);
  position: relative;
}
.close-help {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background: none;
  border: none;
  font-size: 2rem;
  color: #888;
  cursor: pointer;
}
.create-game-btn {
  width: 100%;
  margin-bottom: 1.5rem;
  padding: 0.75rem 0;
  background: #1976d2;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-weight: 600;
  cursor: pointer;
  font-size: 1rem;
}
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 3000;
}
.modal {
  background: #fff;
  padding: 2rem 2.5rem;
  border-radius: 10px;
  min-width: 320px;
  max-width: 90vw;
  box-shadow: 0 2px 16px rgba(0,0,0,0.18);
  position: relative;
}
.modal h2 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  text-align: center;
}
.modal label {
  display: block;
  margin-bottom: 1rem;
  font-weight: 500;
}
.modal input, .modal select {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 1rem;
}
.modal-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 2rem;
}
.modal-actions button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-weight: 600;
  cursor: pointer;
  font-size: 1rem;
}
.modal-actions button:first-child {
  background: #1976d2;
  color: #fff;
}
.modal-actions button:last-child {
  background: #f44336;
  color: #fff;
}
.error {
  color: #f44336;
  margin-top: 1rem;
  font-size: 0.9em;
  text-align: center;
}
</style> 