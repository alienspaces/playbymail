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
  fetchCreaturePlacements,
  createCreaturePlacement,
  updateCreaturePlacement,
  deleteCreaturePlacement,
} from './creaturePlacements'

describe('creaturePlacements API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  const mockJson = (data) => ({
    ok: true,
    json: () => Promise.resolve(data),
  })

  describe('fetchCreaturePlacements', () => {
    it('calls GET /api/v1/adventure-games/:gameId/creature-placements and returns json.data', async () => {
      const placements = [{ id: 'p1', creature_id: 'c1' }]
      mockApiFetch.mockResolvedValue(mockJson({ data: placements }))

      const result = await fetchCreaturePlacements('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/creature-placements',
        expect.objectContaining({
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual(placements)
    })

    it('returns empty array when data is null/undefined', async () => {
      mockApiFetch.mockResolvedValue(mockJson({}))

      const result = await fetchCreaturePlacements('game-1')
      expect(result).toEqual([])
    })
  })

  describe('createCreaturePlacement', () => {
    it('calls POST with data and returns json.data', async () => {
      const data = { creature_id: 'c1', location_id: 'loc1' }
      const created = { id: 'p1', ...data }
      mockApiFetch.mockResolvedValue(mockJson({ data: created }))

      const result = await createCreaturePlacement('game-1', data)

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

  describe('updateCreaturePlacement', () => {
    it('calls PUT with data and returns json.data', async () => {
      const data = { location_id: 'loc2' }
      const updated = { id: 'p1', creature_id: 'c1', ...data }
      mockApiFetch.mockResolvedValue(mockJson({ data: updated }))

      const result = await updateCreaturePlacement('game-1', 'p1', data)

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

  describe('deleteCreaturePlacement', () => {
    it('calls DELETE and returns void', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteCreaturePlacement('game-1', 'p1')

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
