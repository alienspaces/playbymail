import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockApiFetch = vi.fn()
const mockHandleApiError = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  apiFetch: (...args) => mockApiFetch(...args),
  handleApiError: (...args) => mockHandleApiError(...args),
}))

import { listCatalogGames } from './catalog'

describe('catalog API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation(() => {})
  })

  describe('listCatalogGames', () => {
    it('calls GET /api/v1/catalog/game-subscriptions with no auth headers', async () => {
      const catalogData = [{ game_subscription_id: 'sub-1', game_name: 'Test Game', game_type: 'adventure', turn_duration_hours: 168, game_description: '', total_capacity: 4, total_players: 1, delivery_email: true, delivery_physical_local: false, delivery_physical_post: false }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: catalogData }),
      })

      const result = await listCatalogGames()

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/catalog/game-subscriptions',
        expect.objectContaining({ headers: { 'Content-Type': 'application/json' } })
      )
      expect(result.data).toEqual(catalogData)
    })

    it('returns empty data array when catalog is empty', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
      })

      const result = await listCatalogGames()
      expect(result.data).toEqual([])
    })

    it('calls handleApiError on request failure', async () => {
      const errorRes = { ok: false, status: 500 }
      mockApiFetch.mockResolvedValue(errorRes)
      mockHandleApiError.mockRejectedValue(new Error('Failed to fetch game catalog'))

      await expect(listCatalogGames()).rejects.toThrow('Failed to fetch game catalog')
      expect(mockHandleApiError).toHaveBeenCalledWith(errorRes, 'Failed to fetch game catalog')
    })
  })
})
