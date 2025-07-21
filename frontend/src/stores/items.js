import { defineStore } from 'pinia';
import { fetchItems, createItem, updateItem, deleteItem } from '../api/items';

export const useItemsStore = defineStore('items', {
  state: () => ({
    items: [],
    loading: false,
    error: null,
    gameId: null,
  }),
  actions: {
    async loadItems(gameId) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        this.items = await fetchItems(gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async addItem(data) {
      this.loading = true;
      this.error = null;
      try {
        await createItem(this.gameId, data);
        await this.loadItems(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async editItem(itemId, data) {
      this.loading = true;
      this.error = null;
      try {
        await updateItem(this.gameId, itemId, data);
        await this.loadItems(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async removeItem(itemId) {
      this.loading = true;
      this.error = null;
      try {
        await deleteItem(this.gameId, itemId);
        await this.loadItems(this.gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    }
  }
}); 