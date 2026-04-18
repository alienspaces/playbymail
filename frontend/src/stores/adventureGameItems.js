import { defineStore } from 'pinia';
import { fetchAdventureGameItems as apiFetchItems, createAdventureGameItem as apiCreateItem, updateAdventureGameItem as apiUpdateItem, deleteAdventureGameItem as apiDeleteItem } from '../api/adventureGameItems';

export const useAdventureGameItemsStore = defineStore('adventureGameItems', {
  state: () => ({
    items: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchAdventureGameItems(gameId, pageNumber = 1) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        const result = await apiFetchItems(gameId, { page_number: pageNumber });
        this.items = result.data;
        this.hasMore = result.hasMore;
        this.pageNumber = pageNumber;
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createAdventureGameItem(data) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiCreateItem(this.gameId, data);
        const item = res.data || res;
        this.items.push(item);
        return item;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async updateAdventureGameItem(itemId, data) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiUpdateItem(this.gameId, itemId, data);
        const updated = res.data || res;
        const idx = this.items.findIndex(i => i.id === itemId);
        if (idx !== -1) this.items[idx] = updated;
        return updated;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async deleteAdventureGameItem(itemId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteItem(this.gameId, itemId);
        this.items = this.items.filter(i => i.id !== itemId);
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    }
  }
}); 