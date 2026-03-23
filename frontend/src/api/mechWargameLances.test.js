import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockApiFetch = vi.fn()
const mockHandleApiError = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  getAuthHeaders: () => ({ Authorization: 'Bearer test-token' }),
  apiFetch: (...args) => mockApiFetch(...args),
  handleApiError: (...args) => mockHandleApiError(...args),
}))

import { fetchLances, createLance, updateLance, deleteLance } from './mechWargameLances'

describe('mechWargameLances API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('fetchLances', () => {
    it('calls GET /api/v1/mech-wargame-games/:gameId/lances and returns data with hasMore', async () => {
      const lances = [{ id: 'l1', name: 'Alpha Lance' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: lances }),
        headers: { get: (name) => name === 'X-Pagination' ? '{"has_more":false}' : null },
      })

      const result = await fetchLances('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/mech-wargame-games/game-1/lances'),
        expect.objectContaining({ headers: expect.objectContaining({ Authorization: 'Bearer test-token' }) })
      )
      expect(result).toEqual({ data: lances, hasMore: false })
    })

    it('returns empty array when response data is null', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: null }),
        headers: { get: () => null },
      })

      const result = await fetchLances('game-1')
      expect(result).toEqual({ data: [], hasMore: false })
    })

    it('encodes gameId in URL', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
        headers: { get: () => null },
      })

      await fetchLances('game/id')
      expect(mockApiFetch.mock.calls[0][0]).toContain('game%2Fid')
    })
  })

  describe('createLance', () => {
    it('sends POST with data and returns json.data', async () => {
      const data = { name: 'Alpha Lance', description: 'First lance' }
      const created = { id: 'l-new', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: created }),
      })

      const result = await createLance('game-1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/mech-wargame-games/game-1/lances',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ 'Content-Type': 'application/json' }),
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(created)
    })
  })

  describe('updateLance', () => {
    it('sends PUT with data and returns json.data', async () => {
      const data = { name: 'Alpha Lance Updated' }
      const updated = { id: 'l1', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: updated }),
      })

      const result = await updateLance('game-1', 'l1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/mech-wargame-games/game-1/lances/l1',
        expect.objectContaining({ method: 'PUT', body: JSON.stringify(data) })
      )
      expect(result).toEqual(updated)
    })
  })

  describe('deleteLance', () => {
    it('sends DELETE request', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteLance('game-1', 'l1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/mech-wargame-games/game-1/lances/l1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })
})
