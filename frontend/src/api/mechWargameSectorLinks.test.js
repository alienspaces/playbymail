import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockApiFetch = vi.fn()
const mockHandleApiError = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  getAuthHeaders: () => ({ Authorization: 'Bearer test-token' }),
  apiFetch: (...args) => mockApiFetch(...args),
  handleApiError: (...args) => mockHandleApiError(...args),
}))

import { fetchSectorLinks, createSectorLink, updateSectorLink, deleteSectorLink } from './mechWargameSectorLinks'

describe('mechWargameSectorLinks API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('fetchSectorLinks', () => {
    it('calls GET /api/v1/mech-wargame-games/:gameId/sector-links and returns data with hasMore', async () => {
      const links = [{ id: 'sl1', from_sector_id: 's1', to_sector_id: 's2' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: links }),
        headers: { get: (name) => name === 'X-Pagination' ? '{"has_more":false}' : null },
      })

      const result = await fetchSectorLinks('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/mech-wargame-games/game-1/sector-links'),
        expect.objectContaining({ headers: expect.objectContaining({ Authorization: 'Bearer test-token' }) })
      )
      expect(result).toEqual({ data: links, hasMore: false })
    })

    it('returns empty array when response data is null', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: null }),
        headers: { get: () => null },
      })

      const result = await fetchSectorLinks('game-1')
      expect(result).toEqual({ data: [], hasMore: false })
    })

    it('encodes gameId in URL', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
        headers: { get: () => null },
      })

      await fetchSectorLinks('game/id')
      expect(mockApiFetch.mock.calls[0][0]).toContain('game%2Fid')
    })
  })

  describe('createSectorLink', () => {
    it('sends POST with data and returns json.data', async () => {
      const data = { from_sector_id: 's1', to_sector_id: 's2', cover_modifier: 1 }
      const created = { id: 'sl-new', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: created }),
      })

      const result = await createSectorLink('game-1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/mech-wargame-games/game-1/sector-links',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ 'Content-Type': 'application/json' }),
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(created)
    })
  })

  describe('updateSectorLink', () => {
    it('sends PUT with data and returns json.data', async () => {
      const data = { cover_modifier: 2 }
      const updated = { id: 'sl1', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: updated }),
      })

      const result = await updateSectorLink('game-1', 'sl1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/mech-wargame-games/game-1/sector-links/sl1',
        expect.objectContaining({ method: 'PUT', body: JSON.stringify(data) })
      )
      expect(result).toEqual(updated)
    })
  })

  describe('deleteSectorLink', () => {
    it('sends DELETE request', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteSectorLink('game-1', 'sl1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/mech-wargame-games/game-1/sector-links/sl1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })
})
