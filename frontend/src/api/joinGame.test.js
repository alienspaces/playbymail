import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockApiFetch = vi.fn()
const mockHandleApiError = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  apiFetch: (...args) => mockApiFetch(...args),
  handleApiError: (...args) => mockHandleApiError(...args),
}))

import { getJoinGameInfo, verifyJoinGameEmail, submitJoinGame } from './joinGame'

const INSTANCE_ID = 'inst-abc-123'
const BASE = 'http://localhost:8080/api/v1/player/game-instances'

describe('joinGame API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation(() => {})
  })

  describe('getJoinGameInfo', () => {
    it('calls GET /api/v1/player/game-instances/:id/join-game', async () => {
      const data = { game_id: 'g1', game_name: 'Test Game', instance: { id: INSTANCE_ID } }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data }),
      })

      const result = await getJoinGameInfo(INSTANCE_ID)

      expect(mockApiFetch).toHaveBeenCalledWith(
        `${BASE}/${INSTANCE_ID}/join-game`,
        expect.objectContaining({ headers: { 'Content-Type': 'application/json' } })
      )
      expect(result.data).toEqual(data)
    })

    it('calls handleApiError on failure', async () => {
      const errorRes = { ok: false, status: 404 }
      mockApiFetch.mockResolvedValue(errorRes)
      mockHandleApiError.mockRejectedValue(new Error('Failed to load game information'))

      await expect(getJoinGameInfo(INSTANCE_ID)).rejects.toThrow('Failed to load game information')
      expect(mockHandleApiError).toHaveBeenCalledWith(errorRes, 'Failed to load game information')
    })
  })

  describe('verifyJoinGameEmail', () => {
    it('calls POST /api/v1/player/game-instances/:id/join-game/verify-email', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: { has_account: false } }),
      })

      const result = await verifyJoinGameEmail(INSTANCE_ID, 'player@example.com')

      expect(mockApiFetch).toHaveBeenCalledWith(
        `${BASE}/${INSTANCE_ID}/join-game/verify-email`,
        expect.objectContaining({
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email: 'player@example.com' }),
        })
      )
      expect(result.data.has_account).toBe(false)
    })

    it('calls handleApiError on failure', async () => {
      const errorRes = { ok: false, status: 400 }
      mockApiFetch.mockResolvedValue(errorRes)
      mockHandleApiError.mockRejectedValue(new Error('Failed to verify email'))

      await expect(verifyJoinGameEmail(INSTANCE_ID, 'bad@example.com')).rejects.toThrow(
        'Failed to verify email'
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

    it('calls POST /api/v1/player/game-instances/:id/join-game with submit data', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () =>
          Promise.resolve({
            data: {
              game_subscription_id: 'sub-1',
              game_instance_id: INSTANCE_ID,
              game_id: 'g1',
            },
          }),
      })

      const result = await submitJoinGame(INSTANCE_ID, submitData)

      expect(mockApiFetch).toHaveBeenCalledWith(
        `${BASE}/${INSTANCE_ID}/join-game`,
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
      mockApiFetch.mockResolvedValue(errorRes)
      mockHandleApiError.mockRejectedValue(new Error('Failed to submit join game'))

      await expect(submitJoinGame(INSTANCE_ID, submitData)).rejects.toThrow(
        'Failed to submit join game'
      )
    })
  })
})
