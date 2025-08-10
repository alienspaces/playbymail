import { defineStore } from 'pinia';
import { fetchItems as apiFetchItems, createItem as apiCreateItem, updateItem as apiUpdateItem, deleteItem as apiDeleteItem } from '../api/items';

export const useItemsStore = defineStore('items', {
  state: () => ({
    items: [],
    loading: false,
    error: null,
    gameId: null,
  }),
  actions: {
    async fetchItems(gameId) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        this.items = await apiFetchItems(gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createItem(data) {
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
    async updateItem(itemId, data) {
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
    async deleteItem(itemId) {
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