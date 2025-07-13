import { defineStore } from 'pinia';
import { fetchLocations, createLocation, updateLocation, deleteLocation } from '../api/locations';

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
  }),
  actions: {
    /**
     * Load all locations for a game.
     * @param {string} gameId
     * @returns {Promise<void>}
     */
    async loadLocations(gameId) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        this.locations = await fetchLocations(gameId);
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    /**
     * Add a new location.
     * @param {Partial<GameLocation>} data
     * @returns {Promise<GameLocation>}
     */
    async addLocation(data) {
      this.loading = true;
      this.error = null;
      try {
        const location = await createLocation(this.gameId, data);
        this.locations.push(location);
        return location;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    /**
     * Update a location.
     * @param {string} locationId
     * @param {Partial<GameLocation>} data
     * @returns {Promise<GameLocation>}
     */
    async updateLocation(locationId, data) {
      this.loading = true;
      this.error = null;
      try {
        const updated = await updateLocation(locationId, data);
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
    /**
     * Delete a location.
     * @param {string} locationId
     * @returns {Promise<void>}
     */
    async removeLocation(locationId) {
      this.loading = true;
      this.error = null;
      try {
        await deleteLocation(locationId);
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