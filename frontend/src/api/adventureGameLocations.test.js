import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockApiFetch = vi.fn()
const mockHandleApiError = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  getAuthHeaders: () => ({ Authorization: 'Bearer test-token' }),
  apiFetch: (...args) => mockApiFetch(...args),
  handleApiError: (...args) => mockHandleApiError(...args),
}))

import {
  fetchAdventureGameLocations,
  createAdventureGameLocation,
  updateAdventureGameLocation,
  deleteAdventureGameLocation,
} from './adventureGameLocations'

describe('locations API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('fetchAdventureGameLocations', () => {
    it('calls GET /api/v1/adventure-games/:gameId/locations and returns data with hasMore', async () => {
      const locations = [{ id: 'loc1', name: 'Location 1' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: locations }),
        headers: { get: (name) => name === 'X-Pagination' ? '{"has_more":false}' : null },
      })

      const result = await fetchAdventureGameLocations('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/adventure-games/game-1/locations'),
        expect.objectContaining({
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual({ data: locations, hasMore: false })
    })

    it('returns hasMore true when header indicates more pages', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: [{ id: 'loc1' }] }),
        headers: { get: (name) => name === 'X-Pagination' ? '{"has_more":true}' : null },
      })

      const result = await fetchAdventureGameLocations('game-1', { page_number: 1 })
      expect(result.hasMore).toBe(true)
    })

    it('returns empty data when response data is null', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: null }),
        headers: { get: () => null },
      })

      const result = await fetchAdventureGameLocations('game-1')
      expect(result).toEqual({ data: [], hasMore: false })
    })

    it('encodes gameId in URL', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
        headers: { get: () => null },
      })

      await fetchAdventureGameLocations('game/id')

      const calledUrl = mockApiFetch.mock.calls[0][0]
      expect(calledUrl).toContain('game%2Fid')
    })
  })

  describe('createAdventureGameLocation', () => {
    it('sends POST with data and returns json.data', async () => {
      const data = { name: 'New Location', description: 'A place' }
      const created = { id: 'loc-new', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: created }),
      })

      const result = await createAdventureGameLocation('game-1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/locations',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ 'Content-Type': 'application/json', Authorization: 'Bearer test-token' }),
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(created)
    })
  })

  describe('updateAdventureGameLocation', () => {
    it('sends PUT with data and returns json.data', async () => {
      const data = { name: 'Updated Location' }
      const updated = { id: 'loc1', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: updated }),
      })

      const result = await updateAdventureGameLocation('game-1', 'loc1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/locations/loc1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(updated)
    })
  })

  describe('deleteAdventureGameLocation', () => {
    it('sends DELETE and returns void', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteAdventureGameLocation('game-1', 'loc1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/locations/loc1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })
})
