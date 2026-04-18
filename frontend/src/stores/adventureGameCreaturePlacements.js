import { defineStore } from 'pinia';
import { fetchAdventureGameCreaturePlacements as apiFetchCreaturePlacements, createAdventureGameCreaturePlacement as apiCreateCreaturePlacement, updateAdventureGameCreaturePlacement as apiUpdateCreaturePlacement, deleteAdventureGameCreaturePlacement as apiDeleteCreaturePlacement } from '../api/adventureGameCreaturePlacements';

export const useAdventureGameCreaturePlacementsStore = defineStore('adventureGameCreaturePlacements', {
  state: () => ({
    creaturePlacements: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchAdventureGameCreaturePlacements(gameId, pageNumber = 1) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        const result = await apiFetchCreaturePlacements(gameId, { page_number: pageNumber });
        this.creaturePlacements = result.data;
        this.hasMore = result.hasMore;
        this.pageNumber = pageNumber;
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createAdventureGameCreaturePlacement(data) {
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
    async updateAdventureGameCreaturePlacement(placementId, data) {
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
    async deleteAdventureGameCreaturePlacement(placementId) {
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