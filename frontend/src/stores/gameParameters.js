import { defineStore } from 'pinia';
import { 
  listGameParameters, 
  getGameParameter, 
  createGameParameter, 
  updateGameParameter, 
  deleteGameParameter 
} from '../api/gameParameters';

export const useGameParametersStore = defineStore('gameParameters', {
  state: () => ({
    parameters: [],
    loading: false,
    error: null,
    parametersByGameType: {}, // Cache parameters by game type
  }),

  getters: {
    // Get parameters for a specific game type
    getParametersByGameType: (state) => (gameType) => {
      return state.parametersByGameType[gameType] || [];
    },

    // Get a specific parameter by key and game type
    getParameterByKey: (state) => (gameType, configKey) => {
      const params = state.parametersByGameType[gameType] || [];
      return params.find(param => param.config_key === configKey);
    },

    // Get all parameter keys for a game type
    getParameterKeysByGameType: (state) => (gameType) => {
      const params = state.parametersByGameType[gameType] || [];
      return params.map(param => param.config_key);
    },
  },

  actions: {
    // Fetch all game parameters
    async fetchGameParameters(params = {}) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await listGameParameters(params);
        this.parameters = response.data || [];
        
        // Group parameters by game type for easy access
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

    // Fetch parameters for a specific game type
    async fetchGameParametersByGameType(gameType) {
      return this.fetchGameParameters({ gameType });
    },

    // Get a specific parameter
    async fetchGameParameter(id) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await getGameParameter(id);
        return response.data;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    // Create a new game parameter
    async createGameParameter(data) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await createGameParameter(data);
        const newParam = response.data;
        
        // Add to local state
        this.parameters.push(newParam);
        
        // Update game type cache
        if (!this.parametersByGameType[newParam.game_type]) {
          this.parametersByGameType[newParam.game_type] = [];
        }
        this.parametersByGameType[newParam.game_type].push(newParam);
        
        return newParam;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    // Update an existing game parameter
    async updateGameParameter(id, data) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await updateGameParameter(id, data);
        const updatedParam = response.data;
        
        // Update local state
        const index = this.parameters.findIndex(param => param.id === id);
        if (index !== -1) {
          this.parameters[index] = updatedParam;
        }
        
        // Update game type cache
        const gameTypeParams = this.parametersByGameType[updatedParam.game_type] || [];
        const gameTypeIndex = gameTypeParams.findIndex(param => param.id === id);
        if (gameTypeIndex !== -1) {
          gameTypeParams[gameTypeIndex] = updatedParam;
        }
        
        return updatedParam;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    // Delete a game parameter
    async deleteGameParameter(id) {
      this.loading = true;
      this.error = null;
      
      try {
        await deleteGameParameter(id);
        
        // Remove from local state
        this.parameters = this.parameters.filter(param => param.id !== id);
        
        // Remove from game type cache
        Object.keys(this.parametersByGameType).forEach(gameType => {
          this.parametersByGameType[gameType] = this.parametersByGameType[gameType].filter(
            param => param.id !== id
          );
        });
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    // Clear error state
    clearError() {
      this.error = null;
    },
  },
});
