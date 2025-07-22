// Standard Store Method Naming Conventions:
// - fetch<ResourcePlural>(gameId?)
// - create<ResourceSingular>(data)
// - update<ResourceSingular>(id, data)
// - delete<ResourceSingular>(id)
// Example: fetchLocations, createLocation, updateLocation, deleteLocation

import { defineStore } from 'pinia';
import { listGames, createGame as apiCreateGame, updateGame as apiUpdateGame, deleteGame as apiDeleteGame } from '../api/games';

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
    async fetchGames() {
      this.loading = true;
      this.error = null;
      try {
        const res = await listGames();
        this.games = res.data || [];
      } catch (err) {
        this.error = err.message;
      } finally {
        this.loading = false;
      }
    },
    async createGame({ name, game_type }) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiCreateGame({ name, game_type });
        if (res.data) {
          this.games.push(res.data);
        }
        return res.data;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async updateGame(id, { name, game_type }) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiUpdateGame(id, { name, game_type });
        if (res.data) {
          const idx = this.games.findIndex(g => g.id === id);
          if (idx !== -1) this.games[idx] = res.data;
        }
        return res.data;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async deleteGame(id) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteGame(id);
        this.games = this.games.filter(g => g.id !== id);
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },
  },
}); 