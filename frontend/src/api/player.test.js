import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockFetch = vi.fn()
globalThis.fetch = mockFetch

const mockHandleApiError = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  getAuthHeaders: () => ({ Authorization: 'Bearer test-token' }),
  handleApiError: (...args) => mockHandleApiError(...args),
}))

import { verifyGameSubscriptionToken } from './player'

describe('player API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('verifyGameSubscriptionToken', () => {
    it('calls POST /api/v1/player/game-subscription-instances/:id/verify-token with email and turn_sheet_token', async () => {
      mockFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ session_token: 'session-abc' }),
      })

      const result = await verifyGameSubscriptionToken('gsi-123', 'user@example.com', 'token-xyz')

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/player/game-subscription-instances/gsi-123/verify-token',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            Authorization: 'Bearer test-token',
          }),
          body: JSON.stringify({ email: 'user@example.com', turn_sheet_token: 'token-xyz' }),
        })
      )
      expect(result).toBe('session-abc')
    })
  })
})
