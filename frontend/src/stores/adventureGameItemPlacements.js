import { defineStore } from 'pinia';
import { fetchAdventureGameItemPlacements as apiFetchItemPlacements, createAdventureGameItemPlacement as apiCreateItemPlacement, updateAdventureGameItemPlacement as apiUpdateItemPlacement, deleteAdventureGameItemPlacement as apiDeleteItemPlacement } from '../api/adventureGameItemPlacements';

export const useAdventureGameItemPlacementsStore = defineStore('adventureGameItemPlacements', {
  state: () => ({
    itemPlacements: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchAdventureGameItemPlacements(gameId, pageNumber = 1) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        const result = await apiFetchItemPlacements(gameId, { page_number: pageNumber });
        this.itemPlacements = result.data;
        this.hasMore = result.hasMore;
        this.pageNumber = pageNumber;
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createAdventureGameItemPlacement(data) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiCreateItemPlacement(this.gameId, data);
        const placement = res.data || res;
        this.itemPlacements.push(placement);
        return placement;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async updateAdventureGameItemPlacement(placementId, data) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiUpdateItemPlacement(this.gameId, placementId, data);
        const updated = res.data || res;
        const idx = this.itemPlacements.findIndex(p => p.id === placementId);
        if (idx !== -1) this.itemPlacements[idx] = updated;
        return updated;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async deleteAdventureGameItemPlacement(placementId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteItemPlacement(this.gameId, placementId);
        this.itemPlacements = this.itemPlacements.filter(p => p.id !== placementId);
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    }
  }
}); 