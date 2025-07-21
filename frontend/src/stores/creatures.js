import { defineStore } from 'pinia';
import { fetchCreatures, createCreature, updateCreature, deleteCreature } from '../api/creatures';

export const useCreaturesStore = defineStore('creatures', {
  state: () => ({
    creatures: [],
    loading: false,
    error: null,
    gameId: null,
  }),
  actions: {
    async loadCreatures(gameId) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        this.creatures = await fetchCreatures(gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async addCreature(data) {
      this.loading = true;
      this.error = null;
      try {
        await createCreature(this.gameId, data);
        await this.loadCreatures(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async editCreature(creatureId, data) {
      this.loading = true;
      this.error = null;
      try {
        await updateCreature(this.gameId, creatureId, data);
        await this.loadCreatures(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async removeCreature(creatureId) {
      this.loading = true;
      this.error = null;
      try {
        await deleteCreature(this.gameId, creatureId);
        await this.loadCreatures(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    }
  }
}); 