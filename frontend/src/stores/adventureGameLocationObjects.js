import { defineStore } from 'pinia';
import {
  fetchAdventureGameLocationObjects as apiFetchLocationObjects,
  createAdventureGameLocationObject as apiCreateLocationObject,
  updateAdventureGameLocationObject as apiUpdateLocationObject,
  deleteAdventureGameLocationObject as apiDeleteLocationObject,
} from '../api/adventureGameLocationObjects';

export const useAdventureGameLocationObjectsStore = defineStore('adventureGameLocationObjects', {
  state: () => ({
    locationObjects: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchAdventureGameLocationObjects(gameId, pageNumber = 1) {
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
    async createAdventureGameLocationObject(data) {
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
    async updateAdventureGameLocationObject(locationObjectId, data) {
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
    async deleteAdventureGameLocationObject(locationObjectId) {
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
