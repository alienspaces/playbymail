import { defineStore } from 'pinia'
import {
  fetchMechaGameSquadMechs as apiFetchSquadMechs,
  createMechaGameSquadMech as apiCreateSquadMech,
  updateMechaGameSquadMech as apiUpdateSquadMech,
  deleteMechaGameSquadMech as apiDeleteSquadMech,
} from '../api/mechaGameSquadMechs'

export const useMechaGameSquadMechsStore = defineStore('mechaGameSquadMechs', {
  state: () => ({
    mechsBySquad: {},
    loading: false,
    error: null,
    gameId: null,
  }),
  actions: {
    async fetchMechaGameSquadMechs(gameId, squadId) {
      this.loading = true
      this.error = null
      this.gameId = gameId
      try {
        const result = await apiFetchSquadMechs(gameId, squadId)
        this.mechsBySquad = { ...this.mechsBySquad, [squadId]: result.data }
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    getMechsForSquad(squadId) {
      return this.mechsBySquad[squadId] || []
    },
    async createMechaGameSquadMech(squadId, data) {
      this.loading = true
      this.error = null
      try {
        const mech = await apiCreateSquadMech(this.gameId, squadId, data)
        const existing = this.mechsBySquad[squadId] || []
        this.mechsBySquad = { ...this.mechsBySquad, [squadId]: [...existing, mech] }
        return mech
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async updateMechaGameSquadMech(squadId, mechId, data) {
      this.loading = true
      this.error = null
      try {
        const updated = await apiUpdateSquadMech(this.gameId, squadId, mechId, data)
        const mechs = this.mechsBySquad[squadId] || []
        const idx = mechs.findIndex(m => m.id === mechId)
        if (idx !== -1) {
          const updatedMechs = [...mechs]
          updatedMechs[idx] = updated
          this.mechsBySquad = { ...this.mechsBySquad, [squadId]: updatedMechs }
        }
        return updated
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async deleteMechaGameSquadMech(squadId, mechId) {
      this.loading = true
      this.error = null
      try {
        await apiDeleteSquadMech(this.gameId, squadId, mechId)
        const mechs = this.mechsBySquad[squadId] || []
        this.mechsBySquad = { ...this.mechsBySquad, [squadId]: mechs.filter(m => m.id !== mechId) }
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
  },
})
