import { defineStore } from 'pinia'
import {
  fetchSectorLinks as apiFetchSectorLinks,
  createSectorLink as apiCreateSectorLink,
  updateSectorLink as apiUpdateSectorLink,
  deleteSectorLink as apiDeleteSectorLink,
} from '../api/mechaSectorLinks'

export const useMechaSectorLinksStore = defineStore('mechaSectorLinks', {
  state: () => ({
    sectorLinks: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchSectorLinks(gameId, pageNumber = 1) {
      this.loading = true
      this.error = null
      this.gameId = gameId
      try {
        const result = await apiFetchSectorLinks(gameId, { page_number: pageNumber })
        this.sectorLinks = result.data
        this.hasMore = result.hasMore
        this.pageNumber = pageNumber
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    async createSectorLink(data) {
      this.loading = true
      this.error = null
      try {
        const sectorLink = await apiCreateSectorLink(this.gameId, data)
        this.sectorLinks.push(sectorLink)
        return sectorLink
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async updateSectorLink(sectorLinkId, data) {
      this.loading = true
      this.error = null
      try {
        const updated = await apiUpdateSectorLink(this.gameId, sectorLinkId, data)
        const idx = this.sectorLinks.findIndex(sl => sl.id === sectorLinkId)
        if (idx !== -1) this.sectorLinks[idx] = updated
        return updated
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async deleteSectorLink(sectorLinkId) {
      this.loading = true
      this.error = null
      try {
        await apiDeleteSectorLink(this.gameId, sectorLinkId)
        this.sectorLinks = this.sectorLinks.filter(sl => sl.id !== sectorLinkId)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
  },
})
