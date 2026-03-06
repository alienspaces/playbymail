import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockApiFetch = vi.fn()
const mockHandleApiError = vi.fn()
const mockGetAuthHeaders = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  apiFetch: (...args) => mockApiFetch(...args),
  handleApiError: (...args) => mockHandleApiError(...args),
  getAuthHeaders: (...args) => mockGetAuthHeaders(...args),
}))

import { getJoinGameInfo, getJoinSheet, submitJoinGame } from './joinGame'

const SUB_ID = 'sub-abc-123'
const BASE = 'http://localhost:8080/api/v1/game-subscriptions'

describe('joinGame API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation(() => {})
    mockGetAuthHeaders.mockReturnValue({ Authorization: 'Bearer test-token' })
  })

  describe('getJoinGameInfo', () => {
    it('calls GET /api/v1/game-subscriptions/:id/join', async () => {
      const data = { game_subscription_id: SUB_ID, game_name: 'Test Game' }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data }),
      })

      const result = await getJoinGameInfo(SUB_ID)

      expect(mockApiFetch).toHaveBeenCalledWith(
        `${BASE}/${SUB_ID}/join`,
        expect.objectContaining({ headers: { 'Content-Type': 'application/json' } })
      )
      expect(result.data).toEqual(data)
    })

    it('calls handleApiError on failure', async () => {
      const errorRes = { ok: false, status: 404 }
      mockApiFetch.mockResolvedValue(errorRes)
      mockHandleApiError.mockRejectedValue(new Error('Failed to load game information'))

      await expect(getJoinGameInfo(SUB_ID)).rejects.toThrow('Failed to load game information')
      expect(mockHandleApiError).toHaveBeenCalledWith(errorRes, 'Failed to load game information')
    })
  })

  describe('getJoinSheet', () => {
    it('calls GET /api/v1/game-subscriptions/:id/join/sheet with auth headers', async () => {
      const html = '<html><body>Turn sheet</body></html>'
      mockApiFetch.mockResolvedValue({
        ok: true,
        text: () => Promise.resolve(html),
      })

      const result = await getJoinSheet(SUB_ID)

      expect(mockApiFetch).toHaveBeenCalledWith(
        `${BASE}/${SUB_ID}/join/sheet`,
        expect.objectContaining({
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toBe(html)
    })

    it('calls handleApiError on failure', async () => {
      const errorRes = { ok: false, status: 401 }
      mockApiFetch.mockResolvedValue(errorRes)
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

    it('calls POST /api/v1/game-subscriptions/:id/join with auth headers and submit data', async () => {
      mockApiFetch.mockResolvedValue({
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

      expect(mockApiFetch).toHaveBeenCalledWith(
        `${BASE}/${SUB_ID}/join`,
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            Authorization: 'Bearer test-token',
          }),
          body: JSON.stringify(submitData),
        })
      )
      expect(result.data.game_subscription_id).toBe('sub-1')
    })

    it('calls handleApiError on failure', async () => {
      const errorRes = { ok: false, status: 400 }
      mockApiFetch.mockResolvedValue(errorRes)
      mockHandleApiError.mockRejectedValue(new Error('Failed to submit join game'))

      await expect(submitJoinGame(SUB_ID, submitData)).rejects.toThrow(
        'Failed to submit join game'
      )
    })
  })
})
