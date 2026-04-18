import { defineStore } from 'pinia';
import { fetchAdventureGameCreatures as apiFetchCreatures, createAdventureGameCreature as apiCreateCreature, updateAdventureGameCreature as apiUpdateCreature, deleteAdventureGameCreature as apiDeleteCreature } from '../api/adventureGameCreatures';

export const useAdventureGameCreaturesStore = defineStore('adventureGameCreatures', {
  state: () => ({
    creatures: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchAdventureGameCreatures(gameId, pageNumber = 1) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        const result = await apiFetchCreatures(gameId, { page_number: pageNumber });
        this.creatures = result.data;
        this.hasMore = result.hasMore;
        this.pageNumber = pageNumber;
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createAdventureGameCreature(data) {
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
    async updateAdventureGameCreature(creatureId, data) {
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
    async deleteAdventureGameCreature(creatureId) {
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