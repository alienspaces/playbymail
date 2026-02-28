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
  fetchItemPlacements,
  createItemPlacement,
  updateItemPlacement,
  deleteItemPlacement,
} from './itemPlacements'

describe('itemPlacements API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('fetchItemPlacements', () => {
    it('calls GET /api/v1/adventure-games/:gameId/item-placements and returns data', async () => {
      const placements = [{ id: 'p1', item_id: 'item1', location_id: 'loc1' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: placements }),
      })

      const result = await fetchItemPlacements('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/item-placements',
        expect.objectContaining({
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual(placements)
    })

    it('returns empty array when data is null', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: null }),
      })

      const result = await fetchItemPlacements('game-1')
      expect(result).toEqual([])
    })
  })

  describe('createItemPlacement', () => {
    it('sends POST with data and returns json.data', async () => {
      const data = { item_id: 'item1', location_id: 'loc1' }
      const created = { id: 'p-new', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: created }),
      })

      const result = await createItemPlacement('game-1', data)

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

  describe('updateItemPlacement', () => {
    it('sends PUT with data and returns json.data', async () => {
      const data = { quantity: 2 }
      const updated = { id: 'p1', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: updated }),
      })

      const result = await updateItemPlacement('game-1', 'p1', data)

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

  describe('deleteItemPlacement', () => {
    it('sends DELETE and returns void', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteItemPlacement('game-1', 'p1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/item-placements/p1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })
})
