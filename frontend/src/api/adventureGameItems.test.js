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
  fetchAdventureGameItems,
  createAdventureGameItem,
  updateAdventureGameItem,
  deleteAdventureGameItem,
} from './adventureGameItems'

describe('items API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('fetchAdventureGameItems', () => {
    it('calls GET /api/v1/adventure-games/:gameId/items and returns data with hasMore', async () => {
      const items = [{ id: 'item1', name: 'Item 1' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: items }),
        headers: { get: (name) => name === 'X-Pagination' ? '{"has_more":false}' : null },
      })

      const result = await fetchAdventureGameItems('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/adventure-games/game-1/items'),
        expect.objectContaining({
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual({ data: items, hasMore: false })
    })

    it('returns empty data when response data is null', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: null }),
        headers: { get: () => null },
      })

      const result = await fetchAdventureGameItems('game-1')
      expect(result).toEqual({ data: [], hasMore: false })
    })
  })

  describe('createAdventureGameItem', () => {
    it('sends POST with data and returns json.data', async () => {
      const data = { name: 'New Item', description: 'An item' }
      const created = { id: 'item-new', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: created }),
      })

      const result = await createAdventureGameItem('game-1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/items',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ 'Content-Type': 'application/json', Authorization: 'Bearer test-token' }),
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(created)
    })
  })

  describe('updateAdventureGameItem', () => {
    it('sends PUT with data and returns json.data', async () => {
      const data = { name: 'Updated Item' }
      const updated = { id: 'item1', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: updated }),
      })

      const result = await updateAdventureGameItem('game-1', 'item1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/items/item1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(updated)
    })
  })

  describe('deleteAdventureGameItem', () => {
    it('sends DELETE and returns void', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteAdventureGameItem('game-1', 'item1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/items/item1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })
})
