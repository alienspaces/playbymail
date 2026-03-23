import { defineStore } from 'pinia'
import {
  fetchLances as apiFetchLances,
  createLance as apiCreateLance,
  updateLance as apiUpdateLance,
  deleteLance as apiDeleteLance,
} from '../api/mechWargameLances'

export const useMechWargameLancesStore = defineStore('mechWargameLances', {
  state: () => ({
    lances: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchLances(gameId, pageNumber = 1) {
      this.loading = true
      this.error = null
      this.gameId = gameId
      try {
        const result = await apiFetchLances(gameId, { page_number: pageNumber })
        this.lances = result.data
        this.hasMore = result.hasMore
        this.pageNumber = pageNumber
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    async createLance(data) {
      this.loading = true
      this.error = null
      try {
        const lance = await apiCreateLance(this.gameId, data)
        this.lances.push(lance)
        return lance
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async updateLance(lanceId, data) {
      this.loading = true
      this.error = null
      try {
        const updated = await apiUpdateLance(this.gameId, lanceId, data)
        const idx = this.lances.findIndex(l => l.id === lanceId)
        if (idx !== -1) this.lances[idx] = updated
        return updated
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async deleteLance(lanceId) {
      this.loading = true
      this.error = null
      try {
        await apiDeleteLance(this.gameId, lanceId)
        this.lances = this.lances.filter(l => l.id !== lanceId)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
  },
})
