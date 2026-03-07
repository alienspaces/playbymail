import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockFetch = vi.fn()
globalThis.fetch = mockFetch

const mockHandleApiError = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  apiFetch: (...args) => mockFetch(...args),
  handleApiError: (...args) => mockHandleApiError(...args),
}))

import { verifyGameSubscriptionToken, requestNewTurnSheetToken } from './player'

const GSI_BASE = 'http://localhost:8080/api/v1/player/game-subscription-instances'

describe('player API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('verifyGameSubscriptionToken', () => {
    it('calls POST verify-token with turn_sheet_token only (no email)', async () => {
      mockFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ session_token: 'session-abc' }),
      })

      const result = await verifyGameSubscriptionToken('gsi-123', 'token-xyz')

      expect(mockFetch).toHaveBeenCalledWith(
        `${GSI_BASE}/gsi-123/verify-token`,
        expect.objectContaining({
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ turn_sheet_token: 'token-xyz' }),
        })
      )
      expect(result).toBe('session-abc')
    })

    it('includes email when provided', async () => {
      mockFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ session_token: 'session-def' }),
      })

      const result = await verifyGameSubscriptionToken('gsi-123', 'token-xyz', 'user@example.com')

      expect(mockFetch).toHaveBeenCalledWith(
        `${GSI_BASE}/gsi-123/verify-token`,
        expect.objectContaining({
          body: JSON.stringify({ turn_sheet_token: 'token-xyz', email: 'user@example.com' }),
        })
      )
      expect(result).toBe('session-def')
    })
  })

  describe('requestNewTurnSheetToken', () => {
    it('calls POST request-token with email', async () => {
      mockFetch.mockResolvedValue({ ok: true })

      await requestNewTurnSheetToken('gsi-456', 'player@example.com')

      expect(mockFetch).toHaveBeenCalledWith(
        `${GSI_BASE}/gsi-456/request-token`,
        expect.objectContaining({
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email: 'player@example.com' }),
        })
      )
    })

    it('throws on failure', async () => {
      const errorRes = { ok: false, status: 404 }
      mockFetch.mockResolvedValue(errorRes)
      mockHandleApiError.mockRejectedValue(new Error('Failed to request new link'))

      await expect(requestNewTurnSheetToken('gsi-456', 'player@example.com'))
        .rejects.toThrow('Failed to request new link')
    })
  })
})
