import { defineStore } from 'pinia';
import {
  fetchLocationLinkRequirements as apiFetchLocationLinkRequirements,
  createLocationLinkRequirement as apiCreateLocationLinkRequirement,
  updateLocationLinkRequirement as apiUpdateLocationLinkRequirement,
  deleteLocationLinkRequirement as apiDeleteLocationLinkRequirement
} from '../api/locationLinkRequirements';

export const useLocationLinkRequirementsStore = defineStore('locationLinkRequirements', {
  state: () => ({
    requirements: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchLocationLinkRequirements(gameId, pageNumber = 1) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      try {
        const result = await apiFetchLocationLinkRequirements(gameId, { page_number: pageNumber });
        this.requirements = result.data;
        this.hasMore = result.hasMore;
        this.pageNumber = pageNumber;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async createLocationLinkRequirement(data) {
      this.loading = true;
      this.error = null;
      try {
        const newRequirement = await apiCreateLocationLinkRequirement(this.gameId, data);
        this.requirements.push(newRequirement);
        return newRequirement;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async updateLocationLinkRequirement(requirementId, data) {
      this.loading = true;
      this.error = null;
      try {
        const updatedRequirement = await apiUpdateLocationLinkRequirement(this.gameId, requirementId, data);
        const index = this.requirements.findIndex(r => r.id === requirementId);
        if (index !== -1) {
          this.requirements[index] = updatedRequirement;
        }
        return updatedRequirement;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async deleteLocationLinkRequirement(requirementId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteLocationLinkRequirement(this.gameId, requirementId);
        this.requirements = this.requirements.filter(r => r.id !== requirementId);
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    }
  }
});
