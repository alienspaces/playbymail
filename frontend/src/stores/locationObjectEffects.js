import { defineStore } from 'pinia';
import {
  fetchLocationObjectEffects as apiFetchLocationObjectEffects,
  createLocationObjectEffect as apiCreateLocationObjectEffect,
  updateLocationObjectEffect as apiUpdateLocationObjectEffect,
  deleteLocationObjectEffect as apiDeleteLocationObjectEffect,
} from '../api/locationObjectEffects';

export const useLocationObjectEffectsStore = defineStore('locationObjectEffects', {
  state: () => ({
    locationObjectEffects: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchLocationObjectEffects(gameId, pageNumber = 1) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        const result = await apiFetchLocationObjectEffects(gameId, { page_number: pageNumber });
        this.locationObjectEffects = result.data;
        this.hasMore = result.hasMore;
        this.pageNumber = pageNumber;
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createLocationObjectEffect(data) {
      this.loading = true;
      this.error = null;
      try {
        const created = await apiCreateLocationObjectEffect(this.gameId, data);
        this.locationObjectEffects.push(created);
        return created;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async updateLocationObjectEffect(locationObjectEffectId, data) {
      this.loading = true;
      this.error = null;
      try {
        const updated = await apiUpdateLocationObjectEffect(this.gameId, locationObjectEffectId, data);
        const idx = this.locationObjectEffects.findIndex((e) => e.id === locationObjectEffectId);
        if (idx !== -1) this.locationObjectEffects[idx] = updated;
        return updated;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async deleteLocationObjectEffect(locationObjectEffectId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteLocationObjectEffect(this.gameId, locationObjectEffectId);
        this.locationObjectEffects = this.locationObjectEffects.filter((e) => e.id !== locationObjectEffectId);
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
  },
});
