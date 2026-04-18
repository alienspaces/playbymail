import { defineStore } from 'pinia'
import {
  fetchMechaGameSectors as apiFetchSectors,
  createMechaGameSector as apiCreateSector,
  updateMechaGameSector as apiUpdateSector,
  deleteMechaGameSector as apiDeleteSector,
} from '../api/mechaGameSectors'

export const useMechaGameSectorsStore = defineStore('mechaGameSectors', {
  state: () => ({
    sectors: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchMechaGameSectors(gameId, pageNumber = 1) {
      this.loading = true
      this.error = null
      this.gameId = gameId
      try {
        const result = await apiFetchSectors(gameId, { page_number: pageNumber })
        this.sectors = result.data
        this.hasMore = result.hasMore
        this.pageNumber = pageNumber
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    async createMechaGameSector(data) {
      this.loading = true
      this.error = null
      try {
        const sector = await apiCreateSector(this.gameId, data)
        this.sectors.push(sector)
        return sector
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async updateMechaGameSector(sectorId, data) {
      this.loading = true
      this.error = null
      try {
        const updated = await apiUpdateSector(this.gameId, sectorId, data)
        const idx = this.sectors.findIndex(s => s.id === sectorId)
        if (idx !== -1) this.sectors[idx] = updated
        return updated
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async deleteMechaGameSector(sectorId) {
      this.loading = true
      this.error = null
      try {
        await apiDeleteSector(this.gameId, sectorId)
        this.sectors = this.sectors.filter(s => s.id !== sectorId)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
  },
})
