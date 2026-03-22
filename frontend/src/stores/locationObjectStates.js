import { defineStore } from 'pinia';
import {
  fetchLocationObjectStates as apiFetchStates,
  createLocationObjectState as apiCreateState,
  updateLocationObjectState as apiUpdateState,
  deleteLocationObjectState as apiDeleteState,
} from '../api/locationObjectStates';

export const useLocationObjectStatesStore = defineStore('locationObjectStates', {
  state: () => ({
    /** @type {GameLocationObjectState[]} */
    states: [],
    loading: false,
    error: null,
    gameId: null,
    locationObjectId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchStates(gameId, locationObjectId, pageNumber = 1) {
      this.loading = true;
      this.error = null;
      this.gameId = gameId;
      this.locationObjectId = locationObjectId;
      try {
        const result = await apiFetchStates(gameId, locationObjectId, { page_number: pageNumber });
        this.states = result.data;
        this.hasMore = result.hasMore;
        this.pageNumber = pageNumber;
      } catch (e) {
        this.error = e.message;
      } finally {
        this.loading = false;
      }
    },
    async createState(data) {
      this.loading = true;
      this.error = null;
      try {
        const created = await apiCreateState(this.gameId, this.locationObjectId, data);
        this.states.push(created);
        return created;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async updateState(stateId, data) {
      this.loading = true;
      this.error = null;
      try {
        const updated = await apiUpdateState(this.gameId, this.locationObjectId, stateId, data);
        const idx = this.states.findIndex((s) => s.id === stateId);
        if (idx !== -1) this.states[idx] = updated;
        return updated;
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    async deleteState(stateId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteState(this.gameId, this.locationObjectId, stateId);
        this.states = this.states.filter((s) => s.id !== stateId);
      } catch (e) {
        this.error = e.message;
        throw e;
      } finally {
        this.loading = false;
      }
    },
    clearStates() {
      this.states = [];
      this.locationObjectId = null;
    },
  },
});
