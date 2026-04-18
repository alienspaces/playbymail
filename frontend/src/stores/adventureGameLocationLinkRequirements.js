import { defineStore } from 'pinia';
import {
  fetchAdventureGameLocationLinkRequirements as apiFetchLocationLinkRequirements,
  createAdventureGameLocationLinkRequirement as apiCreateLocationLinkRequirement,
  updateAdventureGameLocationLinkRequirement as apiUpdateLocationLinkRequirement,
  deleteAdventureGameLocationLinkRequirement as apiDeleteLocationLinkRequirement
} from '../api/adventureGameLocationLinkRequirements';

export const useAdventureGameLocationLinkRequirementsStore = defineStore('adventureGameLocationLinkRequirements', {
  state: () => ({
    requirements: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchAdventureGameLocationLinkRequirements(gameId, pageNumber = 1) {
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

    async createAdventureGameLocationLinkRequirement(data) {
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

    async updateAdventureGameLocationLinkRequirement(requirementId, data) {
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

    async deleteAdventureGameLocationLinkRequirement(requirementId) {
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
