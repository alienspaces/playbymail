import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockApiFetch = vi.fn()
const mockHandleApiError = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  getAuthHeaders: () => ({ Authorization: 'Bearer test-token' }),
  apiFetch: (...args) => mockApiFetch(...args),
  handleApiError: (...args) => mockHandleApiError(...args),
}))

import { listGameParameters } from './gameParameters'

describe('gameParameters API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  const mockJson = (data) => ({
    ok: true,
    json: () => Promise.resolve(data),
  })

  describe('listGameParameters', () => {
    it('calls GET /api/v1/game-parameters without params', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: [] }))
      await listGameParameters()
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/game-parameters',
        expect.any(Object)
      )
    })

    it('appends game_type and config_key query params', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: [] }))
      await listGameParameters({ gameType: 'adventure', configKey: 'foo' })
      const calledUrl = mockApiFetch.mock.calls[0][0]
      expect(calledUrl).toContain('game_type=adventure')
      expect(calledUrl).toContain('config_key=foo')
    })
  })
})
