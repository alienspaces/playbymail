import { defineStore } from 'pinia'
import {
  fetchChassis as apiFetchChassis,
  createChassis as apiCreateChassis,
  updateChassis as apiUpdateChassis,
  deleteChassis as apiDeleteChassis,
} from '../api/mechWargameChassis'

export const useMechWargameChassisStore = defineStore('mechWargameChassis', {
  state: () => ({
    chassis: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchChassis(gameId, pageNumber = 1) {
      this.loading = true
      this.error = null
      this.gameId = gameId
      try {
        const result = await apiFetchChassis(gameId, { page_number: pageNumber })
        this.chassis = result.data
        this.hasMore = result.hasMore
        this.pageNumber = pageNumber
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    async createChassis(data) {
      this.loading = true
      this.error = null
      try {
        const chassis = await apiCreateChassis(this.gameId, data)
        this.chassis.push(chassis)
        return chassis
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async updateChassis(chassisId, data) {
      this.loading = true
      this.error = null
      try {
        const updated = await apiUpdateChassis(this.gameId, chassisId, data)
        const idx = this.chassis.findIndex(c => c.id === chassisId)
        if (idx !== -1) this.chassis[idx] = updated
        return updated
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async deleteChassis(chassisId) {
      this.loading = true
      this.error = null
      try {
        await apiDeleteChassis(this.gameId, chassisId)
        this.chassis = this.chassis.filter(c => c.id !== chassisId)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
  },
})
