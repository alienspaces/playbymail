import { defineStore } from 'pinia';
import { fetchLocations as apiFetchLocations, createLocation as apiCreateLocation, updateLocation as apiUpdateLocation, deleteLocation as apiDeleteLocation } from '../api/locations';

/**
 * Pinia store for managing game locations.
 * @typedef {import('../types').GameLocation} GameLocation
 */
export const useLocationsStore = defineStore('locations', {
  state: () => ({
    locations: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchLocations(gameId, pageNumber = 1) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        const result = await apiFetchLocations(gameId, { page_number: pageNumber });
        this.locations = result.data;
        this.hasMore = result.hasMore;
        this.pageNumber = pageNumber;
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createLocation(data) {
      this.loading = true;
      this.error = null;
      try {
        const location = await apiCreateLocation(this.gameId, data);
        this.locations.push(location);
        return location;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async updateLocation(locationId, data) {
      this.loading = true;
      this.error = null;
      try {
        const updated = await apiUpdateLocation(this.gameId, locationId, data);
        const idx = this.locations.findIndex(l => l.id === locationId);
        if (idx !== -1) this.locations[idx] = updated;
        return updated;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async deleteLocation(locationId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteLocation(this.gameId, locationId);
        this.locations = this.locations.filter(l => l.id !== locationId);
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    }
  }
}); 