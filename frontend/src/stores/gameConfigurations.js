import { defineStore } from 'pinia';
import { 
  listGameConfigurations, 
  getGameConfiguration, 
  createGameConfiguration, 
  updateGameConfiguration, 
  deleteGameConfiguration 
} from '../api/gameConfigurations';

export const useGameConfigurationsStore = defineStore('gameConfigurations', {
  state: () => ({
    configurations: [],
    loading: false,
    error: null,
    configurationsByGameType: {}, // Cache configurations by game type
  }),

  getters: {
    // Get configurations for a specific game type
    getConfigurationsByGameType: (state) => (gameType) => {
      return state.configurationsByGameType[gameType] || [];
    },

    // Get a specific configuration by key and game type
    getConfigurationByKey: (state) => (gameType, configKey) => {
      const configs = state.configurationsByGameType[gameType] || [];
      return configs.find(config => config.config_key === configKey);
    },

    // Get all configuration keys for a game type
    getConfigurationKeysByGameType: (state) => (gameType) => {
      const configs = state.configurationsByGameType[gameType] || [];
      return configs.map(config => config.config_key);
    },
  },

  actions: {
    // Fetch all game configurations
    async fetchGameConfigurations(params = {}) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await listGameConfigurations(params);
        this.configurations = response.data || [];
        
        // Group configurations by game type for easy access
        this.configurationsByGameType = {};
        this.configurations.forEach(config => {
          if (!this.configurationsByGameType[config.game_type]) {
            this.configurationsByGameType[config.game_type] = [];
          }
          this.configurationsByGameType[config.game_type].push(config);
        });
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    // Fetch configurations for a specific game type
    async fetchGameConfigurationsByGameType(gameType) {
      return this.fetchGameConfigurations({ gameType });
    },

    // Get a specific configuration
    async fetchGameConfiguration(id) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await getGameConfiguration(id);
        return response.data;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    // Create a new game configuration
    async createGameConfiguration(data) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await createGameConfiguration(data);
        const newConfig = response.data;
        
        // Add to local state
        this.configurations.push(newConfig);
        
        // Update game type cache
        if (!this.configurationsByGameType[newConfig.game_type]) {
          this.configurationsByGameType[newConfig.game_type] = [];
        }
        this.configurationsByGameType[newConfig.game_type].push(newConfig);
        
        return newConfig;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    // Update an existing game configuration
    async updateGameConfiguration(id, data) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await updateGameConfiguration(id, data);
        const updatedConfig = response.data;
        
        // Update local state
        const index = this.configurations.findIndex(config => config.id === id);
        if (index !== -1) {
          this.configurations[index] = updatedConfig;
        }
        
        // Update game type cache
        const gameTypeConfigs = this.configurationsByGameType[updatedConfig.game_type] || [];
        const gameTypeIndex = gameTypeConfigs.findIndex(config => config.id === id);
        if (gameTypeIndex !== -1) {
          gameTypeConfigs[gameTypeIndex] = updatedConfig;
        }
        
        return updatedConfig;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.loading = false;
      }
    },

    // Delete a game configuration
    async deleteGameConfiguration(id) {
      this.loading = true;
      this.error = null;
      
      try {
        await deleteGameConfiguration(id);
        
        // Remove from local state
        this.configurations = this.configurations.filter(config => config.id !== id);
        
        // Remove from game type cache
        Object.keys(this.configurationsByGameType).forEach(gameType => {
          this.configurationsByGameType[gameType] = this.configurationsByGameType[gameType].filter(
            config => config.id !== id
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