import { defineStore } from 'pinia';
import { fetchItemPlacements as apiFetchItemPlacements, createItemPlacement as apiCreateItemPlacement, updateItemPlacement as apiUpdateItemPlacement, deleteItemPlacement as apiDeleteItemPlacement } from '../api/itemPlacements';

export const useItemPlacementsStore = defineStore('itemPlacements', {
  state: () => ({
    itemPlacements: [],
    loading: false,
    error: null,
    gameId: null,
  }),
  actions: {
    async fetchItemPlacements(gameId) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        this.itemPlacements = await apiFetchItemPlacements(gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createItemPlacement(data) {
      this.loading = true;
      this.error = null;
      try {
        await apiCreateItemPlacement(this.gameId, data);
        await this.fetchItemPlacements(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async updateItemPlacement(placementId, data) {
      this.loading = true;
      this.error = null;
      try {
        await apiUpdateItemPlacement(this.gameId, placementId, data);
        await this.fetchItemPlacements(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async deleteItemPlacement(placementId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteItemPlacement(this.gameId, placementId);
        await this.fetchItemPlacements(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    }
  }
}); 