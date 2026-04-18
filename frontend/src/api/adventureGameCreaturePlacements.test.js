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
  fetchAdventureGameCreaturePlacements,
  createAdventureGameCreaturePlacement,
  updateAdventureGameCreaturePlacement,
  deleteAdventureGameCreaturePlacement,
} from './adventureGameCreaturePlacements'

describe('creaturePlacements API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  const mockJson = (data, paginationHeader = null) => ({
    ok: true,
    json: () => Promise.resolve(data),
    headers: { get: (name) => name === 'X-Pagination' ? paginationHeader : null },
  })

  describe('fetchAdventureGameCreaturePlacements', () => {
    it('calls GET /api/v1/adventure-games/:gameId/creature-placements and returns data with hasMore', async () => {
      const placements = [{ id: 'p1', creature_id: 'c1' }]
      mockApiFetch.mockResolvedValue(mockJson({ data: placements }, '{"has_more":false}'))

      const result = await fetchAdventureGameCreaturePlacements('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/adventure-games/game-1/creature-placements'),
        expect.objectContaining({
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual({ data: placements, hasMore: false })
    })

    it('returns empty data when data is null/undefined', async () => {
      mockApiFetch.mockResolvedValue(mockJson({}))

      const result = await fetchAdventureGameCreaturePlacements('game-1')
      expect(result).toEqual({ data: [], hasMore: false })
    })
  })

  describe('createAdventureGameCreaturePlacement', () => {
    it('calls POST with data and returns json.data', async () => {
      const data = { creature_id: 'c1', location_id: 'loc1' }
      const created = { id: 'p1', ...data }
      mockApiFetch.mockResolvedValue(mockJson({ data: created }))

      const result = await createAdventureGameCreaturePlacement('game-1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/creature-placements',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ 'Content-Type': 'application/json', Authorization: 'Bearer test-token' }),
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(created)
    })
  })

  describe('updateAdventureGameCreaturePlacement', () => {
    it('calls PUT with data and returns json.data', async () => {
      const data = { location_id: 'loc2' }
      const updated = { id: 'p1', creature_id: 'c1', ...data }
      mockApiFetch.mockResolvedValue(mockJson({ data: updated }))

      const result = await updateAdventureGameCreaturePlacement('game-1', 'p1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/creature-placements/p1',
        expect.objectContaining({
          method: 'PUT',
          headers: expect.objectContaining({ 'Content-Type': 'application/json', Authorization: 'Bearer test-token' }),
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(updated)
    })
  })

  describe('deleteAdventureGameCreaturePlacement', () => {
    it('calls DELETE and returns void', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteAdventureGameCreaturePlacement('game-1', 'p1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/creature-placements/p1',
        expect.objectContaining({
          method: 'DELETE',
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
    })
  })
})
