import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockApiFetch = vi.fn()
const mockHandleApiError = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  getAuthHeaders: () => ({ Authorization: 'Bearer test-token' }),
  apiFetch: (...args) => mockApiFetch(...args),
  handleApiError: (...args) => mockHandleApiError(...args),
}))

import { fetchMechaGameChassis, createMechaGameChassis, updateMechaGameChassis, deleteMechaGameChassis } from './mechaGameChassis'

describe('mechaGameChassis API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('fetchMechaGameChassis', () => {
    it('calls GET /api/v1/mecha-games/:gameId/chassis and returns data with hasMore', async () => {
      const chassis = [{ id: 'ch1', name: 'Atlas' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: chassis }),
        headers: { get: (name) => name === 'X-Pagination' ? '{"has_more":false}' : null },
      })

      const result = await fetchMechaGameChassis('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/mecha-games/game-1/chassis'),
        expect.objectContaining({ headers: expect.objectContaining({ Authorization: 'Bearer test-token' }) })
      )
      expect(result).toEqual({ data: chassis, hasMore: false })
    })

    it('returns hasMore true when header indicates more pages', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: [{ id: 'ch1' }] }),
        headers: { get: (name) => name === 'X-Pagination' ? '{"has_more":true}' : null },
      })

      const result = await fetchMechaGameChassis('game-1', { page_number: 1 })
      expect(result.hasMore).toBe(true)
    })

    it('returns empty array when response data is null', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: null }),
        headers: { get: () => null },
      })

      const result = await fetchMechaGameChassis('game-1')
      expect(result).toEqual({ data: [], hasMore: false })
    })

    it('encodes gameId in URL', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
        headers: { get: () => null },
      })

      await fetchMechaGameChassis('game/id')

      expect(mockApiFetch.mock.calls[0][0]).toContain('game%2Fid')
    })
  })

  describe('createMechaGameChassis', () => {
    it('sends POST with data and returns json.data', async () => {
      const data = { name: 'Atlas', tonnage: 100 }
      const created = { id: 'ch-new', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: created }),
      })

      const result = await createMechaGameChassis('game-1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/mecha-games/game-1/chassis',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ 'Content-Type': 'application/json', Authorization: 'Bearer test-token' }),
          body: JSON.stringify(data),
        })
      )
      expect(result).toEqual(created)
    })
  })

  describe('updateMechaGameChassis', () => {
    it('sends PUT with data and returns json.data', async () => {
      const data = { name: 'Atlas MkII' }
      const updated = { id: 'ch1', ...data }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: updated }),
      })

      const result = await updateMechaGameChassis('game-1', 'ch1', data)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/mecha-games/game-1/chassis/ch1',
        expect.objectContaining({ method: 'PUT', body: JSON.stringify(data) })
      )
      expect(result).toEqual(updated)
    })
  })

  describe('deleteMechaGameChassis', () => {
    it('sends DELETE request', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      await deleteMechaGameChassis('game-1', 'ch1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/mecha-games/game-1/chassis/ch1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })
  })
})
