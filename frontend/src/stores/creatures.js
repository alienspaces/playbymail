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
        const res = await apiCreateCreature(this.gameId, data);
        const creature = res.data || res;
        this.creatures.push(creature);
        return creature;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async updateCreature(creatureId, data) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiUpdateCreature(this.gameId, creatureId, data);
        const updated = res.data || res;
        const idx = this.creatures.findIndex(c => c.id === creatureId);
        if (idx !== -1) this.creatures[idx] = updated;
        return updated;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async deleteCreature(creatureId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteCreature(this.gameId, creatureId);
        this.creatures = this.creatures.filter(c => c.id !== creatureId);
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    }
  }
}); 