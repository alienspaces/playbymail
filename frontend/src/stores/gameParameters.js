import { defineStore } from 'pinia';
import { listGameParameters } from '../api/gameParameters';

export const useGameParametersStore = defineStore('gameParameters', {
  state: () => ({
    parameters: [],
    loading: false,
    error: null,
    parametersByGameType: {},
  }),

  getters: {
    getParametersByGameType: (state) => (gameType) => {
      return state.parametersByGameType[gameType] || [];
    },

    getParameterByKey: (state) => (gameType, configKey) => {
      const params = state.parametersByGameType[gameType] || [];
      return params.find(param => param.config_key === configKey);
    },

    getParameterKeysByGameType: (state) => (gameType) => {
      const params = state.parametersByGameType[gameType] || [];
      return params.map(param => param.config_key);
    },
  },

  actions: {
    async fetchGameParameters(params = {}) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await listGameParameters(params);
        this.parameters = response.data || [];
        
        this.parametersByGameType = {};
        this.parameters.forEach(param => {
          if (!this.parametersByGameType[param.game_type]) {
            this.parametersByGameType[param.game_type] = [];
          }
          this.parametersByGameType[param.game_type].push(param);
        });
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async fetchGameParametersByGameType(gameType) {
      return this.fetchGameParameters({ gameType });
    },

    clearError() {
      this.error = null;
    },
  },
});
