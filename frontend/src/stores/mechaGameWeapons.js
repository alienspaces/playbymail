import { defineStore } from 'pinia'
import {
  fetchMechaGameWeapons as apiFetchWeapons,
  createMechaGameWeapon as apiCreateWeapon,
  updateMechaGameWeapon as apiUpdateWeapon,
  deleteMechaGameWeapon as apiDeleteWeapon,
} from '../api/mechaGameWeapons'

export const useMechaGameWeaponsStore = defineStore('mechaGameWeapons', {
  state: () => ({
    weapons: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchMechaGameWeapons(gameId, pageNumber = 1) {
      this.loading = true
      this.error = null
      this.gameId = gameId
      try {
        const result = await apiFetchWeapons(gameId, { page_number: pageNumber })
        this.weapons = result.data
        this.hasMore = result.hasMore
        this.pageNumber = pageNumber
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    async createMechaGameWeapon(data) {
      this.loading = true
      this.error = null
      try {
        const weapon = await apiCreateWeapon(this.gameId, data)
        this.weapons.push(weapon)
        return weapon
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async updateMechaGameWeapon(weaponId, data) {
      this.loading = true
      this.error = null
      try {
        const updated = await apiUpdateWeapon(this.gameId, weaponId, data)
        const idx = this.weapons.findIndex(w => w.id === weaponId)
        if (idx !== -1) this.weapons[idx] = updated
        return updated
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async deleteMechaGameWeapon(weaponId) {
      this.loading = true
      this.error = null
      try {
        await apiDeleteWeapon(this.gameId, weaponId)
        this.weapons = this.weapons.filter(w => w.id !== weaponId)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
  },
})
