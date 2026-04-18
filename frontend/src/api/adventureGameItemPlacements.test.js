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
  fetchAdventureGameItemPlacements,
  createAdventureGameItemPlacement,
  updateAdventureGameItemPlacement,
  deleteAdventureGameItemPlacement,
} from './adventureGameItemPlacements'

describe('itemPlacements API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('fetchAdventureGameItemPlacements', () => {
    it('calls GET /api/v1/adventure-games/:gameId/item-placements and returns data with hasMore', async () => {
      const placements = [{ id: 'p1', item_id: 'item1', location_id: 'loc1' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: placements }),
        headers: { get: (name) => name === 'X-Pagination' ? '{"has_more":false}' : null },
      })

      const result = await fetchAdventureGameItemPlacements('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/adventure-games/game-1/item-placements'),
        expect.objectContaining({
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual({ data: placements, hasMore: false })
    })

    it('returns empty data when response data is null', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: null }),
        headers: { get: () => null },
      })

      const result = await fetchAdventureGameItemPlacements('game-1')
      expect(result).toEqual({ data: [], hasMore: false })
    })
  })

  describe('createAdventureGameItemPlacement', () => {
    it('sends POST with data and returns json.data', async () => {
      const data = { item_id: 'item1', location_id: 'loc1' }
      const created = { id: 'p-new', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: created }),
      })

      const result = await createAdventureGameItemPlacement('game-1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/item-placements',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ 'Content-Type': 'application/json', Authorization: 'Bearer test-token' }),
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(created)
    })
  })

  describe('updateAdventureGameItemPlacement', () => {
    it('sends PUT with data and returns json.data', async () => {
      const data = { quantity: 2 }
      const updated = { id: 'p1', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: updated }),
      })

      const result = await updateAdventureGameItemPlacement('game-1', 'p1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/item-placements/p1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(updated)
    })
  })

  describe('deleteAdventureGameItemPlacement', () => {
    it('sends DELETE and returns void', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteAdventureGameItemPlacement('game-1', 'p1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/item-placements/p1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })
})
