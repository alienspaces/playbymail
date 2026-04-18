import { defineStore } from 'pinia';
import {
  fetchAdventureGameItemEffects as apiFetchItemEffects,
  createAdventureGameItemEffect as apiCreateItemEffect,
  updateAdventureGameItemEffect as apiUpdateItemEffect,
  deleteAdventureGameItemEffect as apiDeleteItemEffect,
} from '../api/adventureGameItemEffects';

export const useAdventureGameItemEffectsStore = defineStore('adventureGameItemEffects', {
  state: () => ({
    itemEffects: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchAdventureGameItemEffects(gameId, pageNumber = 1) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        const result = await apiFetchItemEffects(gameId, { page_number: pageNumber });
        this.itemEffects = result.data;
        this.hasMore = result.hasMore;
        this.pageNumber = pageNumber;
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createAdventureGameItemEffect(data) {
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
    async updateAdventureGameItemEffect(itemEffectId, data) {
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
    async deleteAdventureGameItemEffect(itemEffectId) {
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
