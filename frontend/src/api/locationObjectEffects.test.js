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
  fetchLocationObjectEffects,
  createLocationObjectEffect,
  updateLocationObjectEffect,
  deleteLocationObjectEffect,
} from './locationObjectEffects'

describe('locationObjectEffects API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('fetchLocationObjectEffects', () => {
    it('calls GET /api/v1/adventure-games/:gameId/location-object-effects and returns data with hasMore', async () => {
      const effects = [{ id: 'eff1', action_type: 'inspect', effect_type: 'info' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: effects }),
        headers: { get: (name) => name === 'X-Pagination' ? '{"has_more":false}' : null },
      })

      const result = await fetchLocationObjectEffects('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/adventure-games/game-1/location-object-effects'),
        expect.objectContaining({
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual({ data: effects, hasMore: false })
    })

    it('returns empty data when response data is null', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: null }),
        headers: { get: () => null },
      })

      const result = await fetchLocationObjectEffects('game-1')
      expect(result).toEqual({ data: [], hasMore: false })
    })
  })

  describe('createLocationObjectEffect', () => {
    it('sends POST with data and returns json.data', async () => {
      const data = {
        adventure_game_location_object_id: 'obj1',
        action_type: 'inspect',
        effect_type: 'info',
        result_description: 'You inspect the object.',
      }
      const created = { id: 'eff-new', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: created }),
      })

      const result = await createLocationObjectEffect('game-1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/location-object-effects',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ 'Content-Type': 'application/json', Authorization: 'Bearer test-token' }),
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(created)
    })
  })

  describe('updateLocationObjectEffect', () => {
    it('sends PUT with data and returns json.data', async () => {
      const data = { result_description: 'Updated description' }
      const updated = { id: 'eff1', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: updated }),
      })

      const result = await updateLocationObjectEffect('game-1', 'eff1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/location-object-effects/eff1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(updated)
    })
  })

  describe('deleteLocationObjectEffect', () => {
    it('sends DELETE and returns void', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteLocationObjectEffect('game-1', 'eff1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/location-object-effects/eff1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })
})
