import { defineStore } from 'pinia';
import { fetchAdventureGameLocationLinks as apiFetchLocationLinks, createAdventureGameLocationLink as apiCreateLocationLink, updateAdventureGameLocationLink as apiUpdateLocationLink, deleteAdventureGameLocationLink as apiDeleteLocationLink } from '../api/adventureGameLocationLinks';

export const useAdventureGameLocationLinksStore = defineStore('adventureGameLocationLinks', {
  state: () => ({
    locationLinks: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchAdventureGameLocationLinks(gameId, pageNumber = 1) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        const result = await apiFetchLocationLinks(gameId, { page_number: pageNumber });
        this.locationLinks = result.data;
        this.hasMore = result.hasMore;
        this.pageNumber = pageNumber;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async createAdventureGameLocationLink(data) {
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

    async updateAdventureGameLocationLink(locationLinkId, data) {
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

    async deleteAdventureGameLocationLink(locationLinkId) {
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