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
  listAllGameInstances,
  listGameInstances,
  getGameInstance,
  createGameInstance,
  updateGameInstance,
  deleteGameInstance,
  startGameInstance,
  pauseGameInstance,
  resumeGameInstance,
  cancelGameInstance,
  resetGameInstance,
  getJoinGameLink,
  inviteTester,
} from './gameInstances'

describe('gameInstances API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  const mockJson = (data) => ({
    ok: true,
    json: () => Promise.resolve(data),
  })

  describe('listAllGameInstances', () => {
    it('calls GET /api/v1/game-instances with correct headers', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: [] }))
      await listAllGameInstances()
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/game-instances',
        expect.objectContaining({
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            Authorization: 'Bearer test-token',
          }),
        })
      )
    })
  })

  describe('listGameInstances', () => {
    it('calls GET /api/v1/games/:gameId/instances', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: [] }))
      await listGameInstances('g1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances',
        expect.any(Object)
      )
    })
  })

  describe('getGameInstance', () => {
    it('calls GET /api/v1/games/:gameId/instances/:instanceId', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: { id: 'i1' } }))
      await getGameInstance('g1', 'i1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1',
        expect.any(Object)
      )
    })
  })

  describe('createGameInstance', () => {
    it('calls POST /api/v1/games/:gameId/instances with body', async () => {
      const data = { name: 'Instance 1' }
      mockApiFetch.mockResolvedValue(mockJson({ data: { id: 'i1', ...data } }))
      await createGameInstance('g1', data)
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(data),
        })
      )
    })
  })

  describe('updateGameInstance', () => {
    it('calls PUT /api/v1/games/:gameId/instances/:instanceId with body', async () => {
      const data = { name: 'Updated' }
      mockApiFetch.mockResolvedValue(mockJson({ data: { id: 'i1', ...data } }))
      await updateGameInstance('g1', 'i1', data)
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(data),
        })
      )
    })
  })

  describe('deleteGameInstance', () => {
    it('calls DELETE /api/v1/games/:gameId/instances/:instanceId', async () => {
      mockApiFetch.mockResolvedValue(mockJson({}))
      await deleteGameInstance('g1', 'i1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })

  describe('startGameInstance', () => {
    it('calls POST .../instances/:instanceId/start', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: {} }))
      await startGameInstance('g1', 'i1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1/start',
        expect.objectContaining({ method: 'POST' })
      )
    })
  })

  describe('pauseGameInstance', () => {
    it('calls POST .../instances/:instanceId/pause', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: {} }))
      await pauseGameInstance('g1', 'i1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1/pause',
        expect.objectContaining({ method: 'POST' })
      )
    })
  })

  describe('resumeGameInstance', () => {
    it('calls POST .../instances/:instanceId/resume', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: {} }))
      await resumeGameInstance('g1', 'i1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1/resume',
        expect.objectContaining({ method: 'POST' })
      )
    })
  })

  describe('cancelGameInstance', () => {
    it('calls POST .../instances/:instanceId/cancel', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: {} }))
      await cancelGameInstance('g1', 'i1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1/cancel',
        expect.objectContaining({ method: 'POST' })
      )
    })
  })

  describe('resetGameInstance', () => {
    it('calls POST .../instances/:instanceId/reset', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: {} }))
      await resetGameInstance('g1', 'i1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1/reset',
        expect.objectContaining({ method: 'POST' })
      )
    })
  })

  describe('getJoinGameLink', () => {
    it('calls GET .../instances/:instanceId/join-link', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: { link: 'http://...' } }))
      await getJoinGameLink('g1', 'i1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1/join-link',
        expect.any(Object)
      )
    })
  })

  describe('inviteTester', () => {
    it('calls POST .../instances/:instanceId/invite-tester with body { email }', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: {} }))
      await inviteTester('g1', 'i1', 'tester@example.com')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1/invite-tester',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ email: 'tester@example.com' }),
        })
      )
    })
  })
})
