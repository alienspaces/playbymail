import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockApiFetch = vi.fn()
const mockHandleApiError = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  getAuthHeaders: () => ({ Authorization: 'Bearer test-token' }),
  apiFetch: (...args) => mockApiFetch(...args),
  handleApiError: (...args) => mockHandleApiError(...args),
}))

import { listGames, createGame, updateGame, deleteGame, publishGame } from './games'

describe('games API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('listGames', () => {
    it('calls GET /api/v1/games without params', async () => {
      const games = [{ id: 'g1', name: 'Game 1' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: games }),
      })

      const result = await listGames()

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games',
        expect.any(Object)
      )
      expect(result.data).toEqual(games)
    })

    it('appends query params for subscriptionType and status', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
      })

      await listGames({ subscriptionType: 'basic', status: 'active' })

      const calledUrl = mockApiFetch.mock.calls[0][0]
      expect(calledUrl).toContain('subscription_type=basic')
      expect(calledUrl).toContain('status=active')
    })
  })

  describe('createGame', () => {
    it('sends POST /api/v1/games with game data', async () => {
      const gameData = { name: 'New Game', game_type: 'adventure', turn_duration_hours: 24, description: 'A game' }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: { id: 'g-new', ...gameData } }),
      })

      await createGame(gameData)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(gameData),
        })
      )
    })
  })

  describe('updateGame', () => {
    it('sends PUT /api/v1/games/:id with game data', async () => {
      const gameData = { name: 'Updated', game_type: 'adventure', turn_duration_hours: 48, description: 'Updated' }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: { id: 'g1', ...gameData } }),
      })

      await updateGame('g1', gameData)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(gameData),
        })
      )
    })
  })

  describe('deleteGame', () => {
    it('sends DELETE /api/v1/games/:id', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({}),
      })

      await deleteGame('g1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })

  describe('publishGame', () => {
    it('sends POST /api/v1/games/:id/publish', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: { id: 'g1', status: 'published' } }),
      })

      await publishGame('g1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/publish',
        expect.objectContaining({ method: 'POST' })
      )
    })
  })
})
