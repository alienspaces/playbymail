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
        await apiCreateItem(this.gameId, data);
        await this.fetchItems(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async updateItem(itemId, data) {
      this.loading = true;
      this.error = null;
      try {
        await apiUpdateItem(this.gameId, itemId, data);
        await this.fetchItems(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async deleteItem(itemId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteItem(this.gameId, itemId);
        await this.fetchItems(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    }
  }
}); 