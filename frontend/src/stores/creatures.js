import { defineStore } from 'pinia';
import { fetchCreatures as apiFetchCreatures, createCreature as apiCreateCreature, updateCreature as apiUpdateCreature, deleteCreature as apiDeleteCreature } from '../api/creatures';

export const useCreaturesStore = defineStore('creatures', {
  state: () => ({
    creatures: [],
    loading: false,
    error: null,
    gameId: null,
  }),
  actions: {
    async fetchCreatures(gameId) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        this.creatures = await apiFetchCreatures(gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createCreature(data) {
      this.loading = true;
      this.error = null;
      try {
        await apiCreateCreature(this.gameId, data);
        await this.fetchCreatures(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async updateCreature(creatureId, data) {
      this.loading = true;
      this.error = null;
      try {
        await apiUpdateCreature(this.gameId, creatureId, data);
        await this.fetchCreatures(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async deleteCreature(creatureId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteCreature(this.gameId, creatureId);
        await this.fetchCreatures(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    }
  }
}); 