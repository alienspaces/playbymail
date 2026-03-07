import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockFetch = vi.fn()
globalThis.fetch = mockFetch

const mockHandleApiError = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  handleApiError: (...args) => mockHandleApiError(...args),
}))

import { getJoinGameInfo, getJoinSheet, submitJoinGame } from './joinGame'

const SUB_ID = 'sub-abc-123'
const BASE = 'http://localhost:8080/api/v1/game-subscriptions'

describe('joinGame API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation(() => {})
  })

  describe('getJoinGameInfo', () => {
    it('calls GET /api/v1/game-subscriptions/:id/join', async () => {
      const data = { game_subscription_id: SUB_ID, game_name: 'Test Game' }
      mockFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data }),
      })

      const result = await getJoinGameInfo(SUB_ID)

      expect(mockFetch).toHaveBeenCalledWith(
        `${BASE}/${SUB_ID}/join`,
        expect.objectContaining({ headers: { 'Content-Type': 'application/json' } })
      )
      expect(result.data).toEqual(data)
    })

    it('calls handleApiError on failure', async () => {
      const errorRes = { ok: false, status: 404 }
      mockFetch.mockResolvedValue(errorRes)
      mockHandleApiError.mockRejectedValue(new Error('Failed to load game information'))

      await expect(getJoinGameInfo(SUB_ID)).rejects.toThrow('Failed to load game information')
      expect(mockHandleApiError).toHaveBeenCalledWith(errorRes, 'Failed to load game information')
    })
  })

  describe('getJoinSheet', () => {
    it('calls GET /api/v1/game-subscriptions/:id/join/sheet without auth', async () => {
      const html = '<html><body>Turn sheet</body></html>'
      mockFetch.mockResolvedValue({
        ok: true,
        text: () => Promise.resolve(html),
      })

      const result = await getJoinSheet(SUB_ID)

      expect(mockFetch).toHaveBeenCalledWith(`${BASE}/${SUB_ID}/join/sheet`)
      expect(result).toBe(html)
    })

    it('calls handleApiError on failure', async () => {
      const errorRes = { ok: false, status: 500 }
      mockFetch.mockResolvedValue(errorRes)
      mockHandleApiError.mockRejectedValue(new Error('Failed to load join game turn sheet'))

      await expect(getJoinSheet(SUB_ID)).rejects.toThrow('Failed to load join game turn sheet')
      expect(mockHandleApiError).toHaveBeenCalledWith(
        errorRes,
        'Failed to load join game turn sheet'
      )
    })
  })

  describe('submitJoinGame', () => {
    const submitData = {
      email: 'player@example.com',
      name: 'Test Player',
      postal_address_line1: '123 Main St',
      state_province: 'VIC',
      postal_code: '3000',
      country: 'Australia',
      delivery_email: true,
      delivery_physical_local: false,
      delivery_physical_post: false,
    }

    it('calls POST /api/v1/game-subscriptions/:id/join without auth headers', async () => {
      mockFetch.mockResolvedValue({
        ok: true,
        json: () =>
          Promise.resolve({
            data: {
              game_subscription_id: 'sub-1',
              game_instance_id: 'inst-1',
              game_id: 'g1',
            },
          }),
      })

      const result = await submitJoinGame(SUB_ID, submitData)

      expect(mockFetch).toHaveBeenCalledWith(
        `${BASE}/${SUB_ID}/join`,
        expect.objectContaining({
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(submitData),
        })
      )
      expect(result.data.game_subscription_id).toBe('sub-1')
    })

    it('calls handleApiError on failure', async () => {
      const errorRes = { ok: false, status: 400 }
      mockFetch.mockResolvedValue(errorRes)
      mockHandleApiError.mockRejectedValue(new Error('Failed to submit join game'))

      await expect(submitJoinGame(SUB_ID, submitData)).rejects.toThrow(
        'Failed to submit join game'
      )
    })
  })
})
