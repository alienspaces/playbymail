import { defineStore } from 'pinia';
import {
  fetchItemEffects as apiFetchItemEffects,
  createItemEffect as apiCreateItemEffect,
  updateItemEffect as apiUpdateItemEffect,
  deleteItemEffect as apiDeleteItemEffect,
} from '../api/itemEffects';

export const useItemEffectsStore = defineStore('itemEffects', {
  state: () => ({
    itemEffects: [],
    loading: false,
    error: null,
    gameId: null,
  }),
  actions: {
    async fetchItemEffects(gameId) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        this.itemEffects = await apiFetchItemEffects(gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createItemEffect(data) {
      this.loading = true;
      this.error = null;
      try {
        const created = await apiCreateItemEffect(this.gameId, data);
        this.itemEffects.push(created);
        return created;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async updateItemEffect(itemEffectId, data) {
      this.loading = true;
      this.error = null;
      try {
        const updated = await apiUpdateItemEffect(this.gameId, itemEffectId, data);
        const idx = this.itemEffects.findIndex((e) => e.id === itemEffectId);
        if (idx !== -1) this.itemEffects[idx] = updated;
        return updated;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async deleteItemEffect(itemEffectId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteItemEffect(this.gameId, itemEffectId);
        this.itemEffects = this.itemEffects.filter((e) => e.id !== itemEffectId);
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
  },
});
