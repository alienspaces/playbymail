import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockApiFetch = vi.fn()
const mockHandleApiError = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  getAuthHeaders: () => ({ Authorization: 'Bearer test-token' }),
  apiFetch: (...args) => mockApiFetch(...args),
  handleApiError: (...args) => mockHandleApiError(...args),
}))

import { fetchLanceMechs, createLanceMech, updateLanceMech, deleteLanceMech } from './mechaLanceMechs'

describe('mechaLanceMechs API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('fetchLanceMechs', () => {
    it('calls GET /api/v1/mecha-games/:gameId/lances/:lanceId/mechs and returns data with hasMore', async () => {
      const mechs = [{ id: 'm1', chassis_id: 'ch1' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: mechs }),
        headers: { get: (name) => name === 'X-Pagination' ? '{"has_more":false}' : null },
      })

      const result = await fetchLanceMechs('game-1', 'lance-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/mecha-games/game-1/lances/lance-1/mechs'),
        expect.objectContaining({ headers: expect.objectContaining({ Authorization: 'Bearer test-token' }) })
      )
      expect(result).toEqual({ data: mechs, hasMore: false })
    })

    it('returns empty array when response data is null', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: null }),
        headers: { get: () => null },
      })

      const result = await fetchLanceMechs('game-1', 'lance-1')
      expect(result).toEqual({ data: [], hasMore: false })
    })

    it('encodes both gameId and lanceId in URL', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
        headers: { get: () => null },
      })

      await fetchLanceMechs('game/id', 'lance/id')
      const calledUrl = mockApiFetch.mock.calls[0][0]
      expect(calledUrl).toContain('game%2Fid')
      expect(calledUrl).toContain('lance%2Fid')
    })
  })

  describe('createLanceMech', () => {
    it('sends POST with data and returns json.data', async () => {
      const data = { chassis_id: 'ch1', weapon_config: [] }
      const created = { id: 'm-new', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: created }),
      })

      const result = await createLanceMech('game-1', 'lance-1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/mecha-games/game-1/lances/lance-1/mechs',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ 'Content-Type': 'application/json' }),
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(created)
    })
  })

  describe('updateLanceMech', () => {
    it('sends PUT with data and returns json.data', async () => {
      const data = { chassis_id: 'ch2' }
      const updated = { id: 'm1', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: updated }),
      })

      const result = await updateLanceMech('game-1', 'lance-1', 'm1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/mecha-games/game-1/lances/lance-1/mechs/m1',
        expect.objectContaining({ method: 'PUT', body: JSON.stringify(data) })
      )
      expect(result).toEqual(updated)
    })
  })

  describe('deleteLanceMech', () => {
    it('sends DELETE request', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteLanceMech('game-1', 'lance-1', 'm1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/mecha-games/game-1/lances/lance-1/mechs/m1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })
})
