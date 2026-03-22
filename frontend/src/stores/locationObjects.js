import { defineStore } from 'pinia';
import {
  fetchLocationObjects as apiFetchLocationObjects,
  createLocationObject as apiCreateLocationObject,
  updateLocationObject as apiUpdateLocationObject,
  deleteLocationObject as apiDeleteLocationObject,
} from '../api/locationObjects';

export const useLocationObjectsStore = defineStore('locationObjects', {
  state: () => ({
    locationObjects: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchLocationObjects(gameId, pageNumber = 1) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        const result = await apiFetchLocationObjects(gameId, { page_number: pageNumber });
        this.locationObjects = result.data;
        this.hasMore = result.hasMore;
        this.pageNumber = pageNumber;
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createLocationObject(data) {
      this.loading = true;
      this.error = null;
      try {
        const created = await apiCreateLocationObject(this.gameId, data);
        this.locationObjects.push(created);
        return created;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async updateLocationObject(locationObjectId, data) {
      this.loading = true;
      this.error = null;
      try {
        const updated = await apiUpdateLocationObject(this.gameId, locationObjectId, data);
        const idx = this.locationObjects.findIndex((o) => o.id === locationObjectId);
        if (idx !== -1) this.locationObjects[idx] = updated;
        return updated;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async deleteLocationObject(locationObjectId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteLocationObject(this.gameId, locationObjectId);
        this.locationObjects = this.locationObjects.filter((o) => o.id !== locationObjectId);
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
  },
});
