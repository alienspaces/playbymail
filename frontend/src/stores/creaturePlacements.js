import { defineStore } from 'pinia';
import { fetchCreaturePlacements as apiFetchCreaturePlacements, createCreaturePlacement as apiCreateCreaturePlacement, updateCreaturePlacement as apiUpdateCreaturePlacement, deleteCreaturePlacement as apiDeleteCreaturePlacement } from '../api/creaturePlacements';

export const useCreaturePlacementsStore = defineStore('creaturePlacements', {
  state: () => ({
    creaturePlacements: [],
    loading: false,
    error: null,
    gameId: null,
  }),
  actions: {
    async fetchCreaturePlacements(gameId) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        this.creaturePlacements = await apiFetchCreaturePlacements(gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createCreaturePlacement(data) {
      this.loading = true;
      this.error = null;
      try {
        await apiCreateCreaturePlacement(this.gameId, data);
        await this.fetchCreaturePlacements(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async updateCreaturePlacement(placementId, data) {
      this.loading = true;
      this.error = null;
      try {
        await apiUpdateCreaturePlacement(this.gameId, placementId, data);
        await this.fetchCreaturePlacements(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async deleteCreaturePlacement(placementId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteCreaturePlacement(this.gameId, placementId);
        await this.fetchCreaturePlacements(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    }
  }
}); 