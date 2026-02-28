import { defineStore } from 'pinia';
import { 
  listGameInstanceParameters,
  getGameInstanceParameter,
  createGameInstanceParameter as apiCreateGameInstanceParameter,
  updateGameInstanceParameter as apiUpdateGameInstanceParameter,
  deleteGameInstanceParameter as apiDeleteGameInstanceParameter,
} from '../api/gameInstanceParameters';

export const useGameInstanceParametersStore = defineStore('gameInstanceParameters', {
  state: () => ({
    gameInstanceParameters: [],
    selectedGameInstanceParameter: null,
    loading: false,
    error: null,
  }),
  
  getters: {
    getParametersByGameInstanceId: (state) => (gameInstanceId) => {
      return state.gameInstanceParameters.filter(param => param.game_instance_id === gameInstanceId);
    },
    getParameterByKey: (state) => (gameInstanceId, configKey) => {
      return state.gameInstanceParameters.find(
        param => param.game_instance_id === gameInstanceId && param.config_key === configKey
      );
    }
  },

  actions: {
    setSelectedGameInstanceParameter(parameter) {
      this.selectedGameInstanceParameter = parameter;
    },
    
    clearSelectedGameInstanceParameter() {
      this.selectedGameInstanceParameter = null;
    },

    async fetchGameInstanceParameters(gameId, gameInstanceId, params = {}) {
      this.loading = true;
      this.error = null;
      try {
        const res = await listGameInstanceParameters(gameId, gameInstanceId, params);
        this.gameInstanceParameters = res.data || [];
      } catch (err) {
        this.error = err.message;
      } finally {
        this.loading = false;
      }
    },

    async getGameInstanceParameter(gameId, gameInstanceId, parameterId) {
      this.loading = true;
      this.error = null;
      try {
        const res = await getGameInstanceParameter(gameId, gameInstanceId, parameterId);
        const parameter = res.data;
        
        const index = this.gameInstanceParameters.findIndex(p => p.id === parameter.id);
        if (index !== -1) {
          this.gameInstanceParameters[index] = parameter;
        } else {
          this.gameInstanceParameters.push(parameter);
        }
        
        return parameter;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async createGameInstanceParameter(gameId, gameInstanceId, data) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiCreateGameInstanceParameter(gameId, gameInstanceId, data);
        const parameter = res.data;
        this.gameInstanceParameters.push(parameter);
        return parameter;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async updateGameInstanceParameter(gameId, gameInstanceId, parameterId, data) {
      this.loading = true;
      this.error = null;
      try {
        const res = await apiUpdateGameInstanceParameter(gameId, gameInstanceId, parameterId, data);
        const parameter = res.data;
        
        const index = this.gameInstanceParameters.findIndex(p => p.id === parameter.id);
        if (index !== -1) {
          this.gameInstanceParameters[index] = parameter;
        }
        
        return parameter;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async deleteGameInstanceParameter(gameId, gameInstanceId, parameterId) {
      this.loading = true;
      this.error = null;
      try {
        await apiDeleteGameInstanceParameter(gameId, gameInstanceId, parameterId);
        this.gameInstanceParameters = this.gameInstanceParameters.filter(p => p.id !== parameterId);
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },
  }
});
