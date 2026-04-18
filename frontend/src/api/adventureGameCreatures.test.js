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
  fetchAdventureGameCreatures,
  createAdventureGameCreature,
  updateAdventureGameCreature,
  deleteAdventureGameCreature,
} from './adventureGameCreatures'

describe('creatures API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('fetchAdventureGameCreatures', () => {
    it('calls GET /api/v1/adventure-games/:gameId/creatures and returns data with hasMore', async () => {
      const creatures = [{ id: 'creature1', name: 'Creature 1' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: creatures }),
        headers: { get: (name) => name === 'X-Pagination' ? '{"has_more":false}' : null },
      })

      const result = await fetchAdventureGameCreatures('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/adventure-games/game-1/creatures'),
        expect.objectContaining({
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual({ data: creatures, hasMore: false })
    })

    it('returns empty data when response data is null', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: null }),
        headers: { get: () => null },
      })

      const result = await fetchAdventureGameCreatures('game-1')
      expect(result).toEqual({ data: [], hasMore: false })
    })
  })

  describe('createAdventureGameCreature', () => {
    it('sends POST with data and returns json.data', async () => {
      const data = { name: 'New Creature', description: 'A creature', max_health: 50, attack_damage: 10, defense: 0 }
      const created = { id: 'creature-new', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: created }),
      })

      const result = await createAdventureGameCreature('game-1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/creatures',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ 'Content-Type': 'application/json', Authorization: 'Bearer test-token' }),
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(created)
    })
  })

  describe('updateAdventureGameCreature', () => {
    it('sends PUT with data and returns json.data', async () => {
      const data = { name: 'Updated Creature' }
      const updated = { id: 'creature1', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: updated }),
      })

      const result = await updateAdventureGameCreature('game-1', 'creature1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/creatures/creature1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(updated)
    })
  })

  describe('deleteAdventureGameCreature', () => {
    it('sends DELETE and returns void', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteAdventureGameCreature('game-1', 'creature1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/creatures/creature1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })
})
