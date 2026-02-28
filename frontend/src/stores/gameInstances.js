// Standard Store Method Naming Conventions:
// - fetch<ResourcePlural>(contextId?) - contextId when resources belong to parent
// - create<ResourceSingular>(contextId?, data) - contextId when resources belong to parent
// - update<ResourceSingular>(contextId?, id, data) - contextId when resources belong to parent  
// - delete<ResourceSingular>(contextId?, id) - contextId when resources belong to parent
// - For game-scoped resources: fetch*, create*, update*, delete* with gameId

import { defineStore } from 'pinia';
import { 
  listAllGameInstances,
  listGameInstances, 
  getGameInstance,
  createGameInstance as apiCreateGameInstance, 
  updateGameInstance as apiUpdateGameInstance, 
  deleteGameInstance as apiDeleteGameInstance,
  startGameInstance as apiStartGameInstance,
  pauseGameInstance as apiPauseGameInstance,
  resumeGameInstance as apiResumeGameInstance,
  cancelGameInstance as apiCancelGameInstance,
  resetGameInstance as apiResetGameInstance
} from '../api/gameInstances';

export const useGameInstancesStore = defineStore('gameInstances', {
  state: () => ({
    gameInstances: [],
    selectedGameInstance: null,
    loading: false,
    error: null,
  }),
  actions: {
    setSelectedGameInstance(instance) {
      this.selectedGameInstance = instance;
    },
    clearSelectedGameInstance() {
      this.selectedGameInstance = null;
    },
    async fetchGameInstances(gameId) {
      this.loading = true;
      this.error = null;
      try {
        const res = await listGameInstances(gameId);
        this.gameInstances = res.data || [];
      } catch (err) {
        this.error = err.message;
      } finally {
        this.loading = false;
      }
    },
    async fetchAllGameInstances() {
      this.loading = true;
      this.error = null;
      try {
        const res = await listAllGameInstances();
        this.gameInstances = res.data || [];
      } catch (err) {
        this.error = err.message;
      } finally {
        this.loading = false;
      }
    },
    async getGameInstance(gameId, instanceId) {
      this.loading = true;
      this.error = null;
      try {
        const res = await getGameInstance(gameId, instanceId);
        return res.data;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async createGameInstance(gameId, data) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiCreateGameInstance(gameId, data);
        if (res.data) {
          this.gameInstances.push(res.data);
        }
        return res.data;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async updateGameInstance(gameId, instanceId, data) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiUpdateGameInstance(gameId, instanceId, data);
        if (res.data) {
          const idx = this.gameInstances.findIndex(i => i.id === instanceId);
          if (idx !== -1) this.gameInstances[idx] = res.data;
          if (this.selectedGameInstance?.id === instanceId) {
            this.selectedGameInstance = res.data;
          }
        }
        return res.data;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async deleteGameInstance(gameId, instanceId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteGameInstance(gameId, instanceId);
        this.gameInstances = this.gameInstances.filter(i => i.id !== instanceId);
        if (this.selectedGameInstance?.id === instanceId) {
          this.selectedGameInstance = null;
        }
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },
    // Runtime management actions
    async startGameInstance(gameId, instanceId) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiStartGameInstance(gameId, instanceId);
        if (res.data) {
          const idx = this.gameInstances.findIndex(i => i.id === instanceId);
          if (idx !== -1) this.gameInstances[idx] = res.data;
          if (this.selectedGameInstance?.id === instanceId) {
            this.selectedGameInstance = res.data;
          }
        }
        return res.data;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async pauseGameInstance(gameId, instanceId) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiPauseGameInstance(gameId, instanceId);
        if (res.data) {
          const idx = this.gameInstances.findIndex(i => i.id === instanceId);
          if (idx !== -1) this.gameInstances[idx] = res.data;
          if (this.selectedGameInstance?.id === instanceId) {
            this.selectedGameInstance = res.data;
          }
        }
        return res.data;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async resumeGameInstance(gameId, instanceId) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiResumeGameInstance(gameId, instanceId);
        if (res.data) {
          const idx = this.gameInstances.findIndex(i => i.id === instanceId);
          if (idx !== -1) this.gameInstances[idx] = res.data;
          if (this.selectedGameInstance?.id === instanceId) {
            this.selectedGameInstance = res.data;
          }
        }
        return res.data;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async cancelGameInstance(gameId, instanceId) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiCancelGameInstance(gameId, instanceId);
        if (res.data) {
          const idx = this.gameInstances.findIndex(i => i.id === instanceId);
          if (idx !== -1) this.gameInstances[idx] = res.data;
          if (this.selectedGameInstance?.id === instanceId) {
            this.selectedGameInstance = res.data;
          }
        }
        return res.data;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async resetGameInstance(gameId, instanceId) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiResetGameInstance(gameId, instanceId);
        if (res.data) {
          const idx = this.gameInstances.findIndex(i => i.id === instanceId);
          if (idx !== -1) this.gameInstances[idx] = res.data;
          if (this.selectedGameInstance?.id === instanceId) {
            this.selectedGameInstance = res.data;
          }
        }
        return res.data;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },
  },
}); 