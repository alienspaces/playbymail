import { describe, it, expect, vi, beforeEach } from 'vitest'

vi.mock('../stores/auth', () => ({
  useAuthStore: vi.fn(() => ({
    sessionToken: 'test-token',
    logout: vi.fn(),
  })),
}))

const mockFetch = vi.fn()
globalThis.fetch = mockFetch

import { requestAuth, verifyAuth, refreshSession } from './auth'

describe('auth API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('requestAuth', () => {
    it('sends POST /api/v1/request-auth with email', async () => {
      mockFetch.mockResolvedValue({ ok: true })

      const result = await requestAuth('test@example.com')

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/request-auth',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ email: 'test@example.com' }),
        })
      )
      expect(result).toBe(true)
    })

    it('returns false when request fails', async () => {
      mockFetch.mockResolvedValue({ ok: false })

      const result = await requestAuth('bad@example.com')
      expect(result).toBe(false)
    })
  })

  describe('verifyAuth', () => {
    it('sends POST /api/v1/verify-auth with email and token, returns session_token', async () => {
      mockFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ session_token: 'new-session-token' }),
      })

      const result = await verifyAuth('test@example.com', '123456')

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/verify-auth',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ email: 'test@example.com', verification_token: '123456' }),
        })
      )
      expect(result).toBe('new-session-token')
    })
  })

  describe('refreshSession', () => {
    it('sends POST /api/v1/refresh-session and returns response', async () => {
      const responseData = { status: 'ok', expires_in_seconds: 900 }
      mockFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(responseData),
      })

      const result = await refreshSession()

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/refresh-session',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual(responseData)
    })
  })
})
