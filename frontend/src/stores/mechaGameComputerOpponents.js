import { defineStore } from 'pinia'
import {
  fetchMechaGameComputerOpponents as apiFetchComputerOpponents,
  createMechaGameComputerOpponent as apiCreateComputerOpponent,
  updateMechaGameComputerOpponent as apiUpdateComputerOpponent,
  deleteMechaGameComputerOpponent as apiDeleteComputerOpponent,
} from '../api/mechaGameComputerOpponents'

export const useMechaGameComputerOpponentsStore = defineStore('mechaGameComputerOpponents', {
  state: () => ({
    opponents: [],
    loading: false,
    error: null,
    gameId: null,
    pageNumber: 1,
    hasMore: false,
  }),
  actions: {
    async fetchMechaGameComputerOpponents(gameId, pageNumber = 1) {
      this.loading = true
      this.error = null
      this.gameId = gameId
      try {
        const result = await apiFetchComputerOpponents(gameId, { page_number: pageNumber })
        this.opponents = result.data
        this.hasMore = result.hasMore
        this.pageNumber = pageNumber
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    async createMechaGameComputerOpponent(data) {
      this.loading = true
      this.error = null
      try {
        const opponent = await apiCreateComputerOpponent(this.gameId, data)
        this.opponents.push(opponent)
        return opponent
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async updateMechaGameComputerOpponent(opponentId, data) {
      this.loading = true
      this.error = null
      try {
        const updated = await apiUpdateComputerOpponent(this.gameId, opponentId, data)
        const idx = this.opponents.findIndex(o => o.id === opponentId)
        if (idx !== -1) this.opponents[idx] = updated
        return updated
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
    async deleteMechaGameComputerOpponent(opponentId) {
      this.loading = true
      this.error = null
      try {
        await apiDeleteComputerOpponent(this.gameId, opponentId)
        this.opponents = this.opponents.filter(o => o.id !== opponentId)
      } catch (e) {
        this.error = e.message
        throw e
      } finally {
        this.loading = false
      }
    },
  },
})
