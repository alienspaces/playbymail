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
  listGameInstanceParameters,
  getGameInstanceParameter,
  createGameInstanceParameter,
  updateGameInstanceParameter,
  deleteGameInstanceParameter,
} from './gameInstanceParameters'

describe('gameInstanceParameters API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  const mockJson = (data) => ({
    ok: true,
    json: () => Promise.resolve(data),
  })

  describe('listGameInstanceParameters', () => {
    it('calls GET /api/v1/games/:gameId/instances/:gameInstanceId/parameters', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: [] }))
      await listGameInstanceParameters('g1', 'i1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1/parameters',
        expect.any(Object)
      )
    })

    it('appends config_key query param when configKey provided', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: [] }))
      await listGameInstanceParameters('g1', 'i1', { configKey: 'foo' })
      const calledUrl = mockApiFetch.mock.calls[0][0]
      expect(calledUrl).toContain('config_key=foo')
    })
  })

  describe('getGameInstanceParameter', () => {
    it('calls GET .../parameters/:parameterId', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: { id: 'p1' } }))
      await getGameInstanceParameter('g1', 'i1', 'p1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1/parameters/p1',
        expect.any(Object)
      )
    })
  })

  describe('createGameInstanceParameter', () => {
    it('calls POST .../parameters with body', async () => {
      const data = { config_key: 'foo', value: 'bar' }
      mockApiFetch.mockResolvedValue(mockJson({ data: { id: 'p1', ...data } }))
      await createGameInstanceParameter('g1', 'i1', data)
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1/parameters',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(data),
        })
      )
    })
  })

  describe('updateGameInstanceParameter', () => {
    it('calls PUT .../parameters/:parameterId with body', async () => {
      const data = { value: 'updated' }
      mockApiFetch.mockResolvedValue(mockJson({ data: { id: 'p1', ...data } }))
      await updateGameInstanceParameter('g1', 'i1', 'p1', data)
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1/parameters/p1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(data),
        })
      )
    })
  })

  describe('deleteGameInstanceParameter', () => {
    it('calls DELETE .../parameters/:parameterId', async () => {
      mockApiFetch.mockResolvedValue(mockJson({}))
      await deleteGameInstanceParameter('g1', 'i1', 'p1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/g1/instances/i1/parameters/p1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })
})
