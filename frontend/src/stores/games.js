import { defineStore } from 'pinia';

export const useGamesStore = defineStore('games', {
  state: () => ({
    games: [],
    loading: false,
    error: null,
    selectedGame: null, // Holds the currently selected game object
  }),
  actions: {
    setSelectedGame(game) {
      this.selectedGame = game;
    },
    clearSelectedGame() {
      this.selectedGame = null;
    },
    // ...other actions for loading/creating games...
  },
}); 