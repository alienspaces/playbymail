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
  fetchLocationLinks,
  createLocationLink,
  updateLocationLink,
  deleteLocationLink,
} from './locationLinks'

describe('locationLinks API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('fetchLocationLinks', () => {
    it('calls GET /api/v1/adventure-games/:gameId/location-links and returns data', async () => {
      const links = [{ id: 'link1', from_location_id: 'loc1', to_location_id: 'loc2' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: links }),
      })

      const result = await fetchLocationLinks('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/location-links',
        expect.objectContaining({
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual(links)
    })

    it('returns empty array when data is null', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: null }),
      })

      const result = await fetchLocationLinks('game-1')
      expect(result).toEqual([])
    })
  })

  describe('createLocationLink', () => {
    it('sends POST with data and returns json.data', async () => {
      const data = { from_location_id: 'loc1', to_location_id: 'loc2', label: 'North' }
      const created = { id: 'link-new', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: created }),
      })

      const result = await createLocationLink('game-1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/location-links',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ 'Content-Type': 'application/json', Authorization: 'Bearer test-token' }),
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(created)
    })
  })

  describe('updateLocationLink', () => {
    it('sends PUT with data and returns json.data', async () => {
      const data = { label: 'South' }
      const updated = { id: 'link1', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: updated }),
      })

      const result = await updateLocationLink('game-1', 'link1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/location-links/link1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(updated)
    })
  })

  describe('deleteLocationLink', () => {
    it('sends DELETE and returns void', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteLocationLink('game-1', 'link1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/location-links/link1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })
})
