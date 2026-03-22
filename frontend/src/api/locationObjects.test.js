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
  fetchLocationObjects,
  createLocationObject,
  updateLocationObject,
  deleteLocationObject,
} from './locationObjects'

describe('locationObjects API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('fetchLocationObjects', () => {
    it('calls GET /api/v1/adventure-games/:gameId/location-objects and returns data with hasMore', async () => {
      const objects = [{ id: 'obj1', name: 'Ancient Shrine', initial_state: 'intact' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: objects }),
        headers: { get: (name) => name === 'X-Pagination' ? '{"has_more":false}' : null },
      })

      const result = await fetchLocationObjects('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/adventure-games/game-1/location-objects'),
        expect.objectContaining({
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual({ data: objects, hasMore: false })
    })

    it('returns empty data when response data is null', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: null }),
        headers: { get: () => null },
      })

      const result = await fetchLocationObjects('game-1')
      expect(result).toEqual({ data: [], hasMore: false })
    })
  })

  describe('createLocationObject', () => {
    it('sends POST with data and returns json.data', async () => {
      const data = { name: 'New Object', description: 'A new object', adventure_game_location_id: 'loc1' }
      const created = { id: 'obj-new', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: created }),
      })

      const result = await createLocationObject('game-1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/location-objects',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ 'Content-Type': 'application/json', Authorization: 'Bearer test-token' }),
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(created)
    })
  })

  describe('updateLocationObject', () => {
    it('sends PUT with data and returns json.data', async () => {
      const data = { name: 'Updated Object', description: 'Updated description' }
      const updated = { id: 'obj1', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: updated }),
      })

      const result = await updateLocationObject('game-1', 'obj1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/location-objects/obj1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(updated)
    })
  })

  describe('deleteLocationObject', () => {
    it('sends DELETE and returns void', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteLocationObject('game-1', 'obj1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/location-objects/obj1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })
})
