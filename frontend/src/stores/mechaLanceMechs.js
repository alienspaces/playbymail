import { defineStore } from 'pinia'
import {
  fetchLanceMechs as apiFetchLanceMechs,
  createLanceMech as apiCreateLanceMech,
  updateLanceMech as apiUpdateLanceMech,
  deleteLanceMech as apiDeleteLanceMech,
} from '../api/mechaLanceMechs'

export const useMechaLanceMechsStore = defineStore('mechaLanceMechs', {
  state: () => ({
    mechsByLance: {},
    loading: false,
    error: null,
    gameId: null,
  }),
  actions: {
    async fetchLanceMechs(gameId, lanceId) {
      this.loading = true
      this.error = null
      this.gameId = gameId
      try {
        const result = await apiFetchLanceMechs(gameId, lanceId)
        this.mechsByLance = { ...this.mechsByLance, [lanceId]: result.data }
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    getMechsForLance(lanceId) {
      return this.mechsByLance[lanceId] || []
    },
    async createLanceMech(lanceId, data) {
      this.loading = true
      this.error = null
      try {
        const mech = await apiCreateLanceMech(this.gameId, lanceId, data)
        const existing = this.mechsByLance[lanceId] || []
        this.mechsByLance = { ...this.mechsByLance, [lanceId]: [...existing, mech] }
        return mech
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async updateLanceMech(lanceId, mechId, data) {
      this.loading = true
      this.error = null
      try {
        const updated = await apiUpdateLanceMech(this.gameId, lanceId, mechId, data)
        const mechs = this.mechsByLance[lanceId] || []
        const idx = mechs.findIndex(m => m.id === mechId)
        if (idx !== -1) {
          const updatedMechs = [...mechs]
          updatedMechs[idx] = updated
          this.mechsByLance = { ...this.mechsByLance, [lanceId]: updatedMechs }
        }
        return updated
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async deleteLanceMech(lanceId, mechId) {
      this.loading = true
      this.error = null
      try {
        await apiDeleteLanceMech(this.gameId, lanceId, mechId)
        const mechs = this.mechsByLance[lanceId] || []
        this.mechsByLance = { ...this.mechsByLance, [lanceId]: mechs.filter(m => m.id !== mechId) }
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
  },
})
