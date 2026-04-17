import { defineStore } from 'pinia'
import {
  fetchSquads as apiFetchSquads,
  createSquad as apiCreateSquad,
  updateSquad as apiUpdateSquad,
  deleteSquad as apiDeleteSquad,
} from '../api/mechaSquads'

export const useMechaSquadsStore = defineStore('mechaSquads', {
  state: () => ({
    squads: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchSquads(gameId, pageNumber = 1) {
      this.loading = true
      this.error = null
      this.gameId = gameId
      try {
        const result = await apiFetchSquads(gameId, { page_number: pageNumber })
        this.squads = result.data
        this.hasMore = result.hasMore
        this.pageNumber = pageNumber
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    async createSquad(data) {
      this.loading = true
      this.error = null
      try {
        const squad = await apiCreateSquad(this.gameId, data)
        this.squads.push(squad)
        return squad
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async updateSquad(squadId, data) {
      this.loading = true
      this.error = null
      try {
        const updated = await apiUpdateSquad(this.gameId, squadId, data)
        const idx = this.squads.findIndex(l => l.id === squadId)
        if (idx !== -1) this.squads[idx] = updated
        return updated
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async deleteSquad(squadId) {
      this.loading = true
      this.error = null
      try {
        await apiDeleteSquad(this.gameId, squadId)
        this.squads = this.squads.filter(s => s.id !== squadId)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
  },
})
