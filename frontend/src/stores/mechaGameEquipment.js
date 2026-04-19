import { defineStore } from 'pinia'
import {
  fetchMechaGameEquipment as apiFetchEquipment,
  createMechaGameEquipment as apiCreateEquipment,
  updateMechaGameEquipment as apiUpdateEquipment,
  deleteMechaGameEquipment as apiDeleteEquipment,
} from '../api/mechaGameEquipment'

export const useMechaGameEquipmentStore = defineStore('mechaGameEquipment', {
  state: () => ({
    equipment: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchMechaGameEquipment(gameId, pageNumber = 1) {
      this.loading = true
      this.error = null
      this.gameId = gameId
      try {
        const result = await apiFetchEquipment(gameId, { page_number: pageNumber })
        this.equipment = result.data
        this.hasMore = result.hasMore
        this.pageNumber = pageNumber
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    async createMechaGameEquipment(data) {
      this.loading = true
      this.error = null
      try {
        const eq = await apiCreateEquipment(this.gameId, data)
        this.equipment.push(eq)
        return eq
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async updateMechaGameEquipment(equipmentId, data) {
      this.loading = true
      this.error = null
      try {
        const updated = await apiUpdateEquipment(this.gameId, equipmentId, data)
        const idx = this.equipment.findIndex(e => e.id === equipmentId)
        if (idx !== -1) this.equipment[idx] = updated
        return updated
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async deleteMechaGameEquipment(equipmentId) {
      this.loading = true
      this.error = null
      try {
        await apiDeleteEquipment(this.gameId, equipmentId)
        this.equipment = this.equipment.filter(e => e.id !== equipmentId)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
  },
})
