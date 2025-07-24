import { defineStore } from 'pinia';
import { ref } from 'vue';
import * as locationLinksApi from '../api/locationLinks';

export const useLocationLinksStore = defineStore('locationLinks', () => {
  const locationLinks = ref([]);
  const loading = ref(false);
  const error = ref('');

  const currentGameId = ref(null);

  return {
    locationLinks,
    loading,
    error,
    currentGameId,

    async fetchLocationLinks(gameId) {
      loading.value = true;
      error.value = '';
      try {
        const data = await locationLinksApi.fetchLocationLinks(gameId);
        locationLinks.value = data;
        currentGameId.value = gameId;
      } catch (err) {
        error.value = err.message || 'Failed to fetch location links';
        throw err;
      } finally {
        loading.value = false;
      }
    },

    async createLocationLink(data) {
      if (!currentGameId.value) throw new Error('No game selected');
      loading.value = true;
      error.value = '';
      try {
        const newLink = await locationLinksApi.createLocationLink(currentGameId.value, data);
        locationLinks.value.push(newLink);
        return newLink;
      } catch (err) {
        error.value = err.message || 'Failed to create location link';
        throw err;
      } finally {
        loading.value = false;
      }
    },

    async updateLocationLink(locationLinkId, data) {
      if (!currentGameId.value) throw new Error('No game selected');
      loading.value = true;
      error.value = '';
      try {
        const updatedLink = await locationLinksApi.updateLocationLink(currentGameId.value, locationLinkId, data);
        const index = locationLinks.value.findIndex(link => link.id === locationLinkId);
        if (index !== -1) {
          locationLinks.value[index] = updatedLink;
        }
        return updatedLink;
      } catch (err) {
        error.value = err.message || 'Failed to update location link';
        throw err;
      } finally {
        loading.value = false;
      }
    },

    async deleteLocationLink(locationLinkId) {
      if (!currentGameId.value) throw new Error('No game selected');
      loading.value = true;
      error.value = '';
      try {
        await locationLinksApi.deleteLocationLink(currentGameId.value, locationLinkId);
        const index = locationLinks.value.findIndex(link => link.id === locationLinkId);
        if (index !== -1) {
          locationLinks.value.splice(index, 1);
        }
      } catch (err) {
        error.value = err.message || 'Failed to delete location link';
        throw err;
      } finally {
        loading.value = false;
      }
    },

    clearLocationLinks() {
      locationLinks.value = [];
      currentGameId.value = null;
      error.value = '';
    }
  };
}); 