import { defineStore } from 'pinia';
import { fetchLocationLinks as apiFetchLocationLinks, createLocationLink as apiCreateLocationLink, updateLocationLink as apiUpdateLocationLink, deleteLocationLink as apiDeleteLocationLink } from '../api/locationLinks';

export const useLocationLinksStore = defineStore('locationLinks', {
  state: () => ({
    locationLinks: [],
    loading: false,
    error: null,
    gameId: null,
  }),
  actions: {
    async fetchLocationLinks(gameId) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        this.locationLinks = await apiFetchLocationLinks(gameId);
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async createLocationLink(data) {
      this.loading = true;
      this.error = null;
      try {
        const newLink = await apiCreateLocationLink(this.gameId, data);
        this.locationLinks.push(newLink);
        return newLink;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async updateLocationLink(locationLinkId, data) {
      this.loading = true;
      this.error = null;
      try {
        const updatedLink = await apiUpdateLocationLink(this.gameId, locationLinkId, data);
        const index = this.locationLinks.findIndex(link => link.id === locationLinkId);
        if (index !== -1) {
          this.locationLinks[index] = updatedLink;
        }
        return updatedLink;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async deleteLocationLink(locationLinkId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteLocationLink(this.gameId, locationLinkId);
        this.locationLinks = this.locationLinks.filter(link => link.id !== locationLinkId);
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    }
  }
}); 