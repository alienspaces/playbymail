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
        const res = await apiCreateCreaturePlacement(this.gameId, data);
        const placement = res.data || res;
        this.creaturePlacements.push(placement);
        return placement;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async updateCreaturePlacement(placementId, data) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiUpdateCreaturePlacement(this.gameId, placementId, data);
        const updated = res.data || res;
        const idx = this.creaturePlacements.findIndex(p => p.id === placementId);
        if (idx !== -1) this.creaturePlacements[idx] = updated;
        return updated;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async deleteCreaturePlacement(placementId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteCreaturePlacement(this.gameId, placementId);
        this.creaturePlacements = this.creaturePlacements.filter(p => p.id !== placementId);
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    }
  }
}); 