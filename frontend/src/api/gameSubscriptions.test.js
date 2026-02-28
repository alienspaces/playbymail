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
  getMyGameSubscriptions,
  createGameSubscription,
  cancelGameSubscription,
  linkGameInstanceToSubscription,
  unlinkGameInstanceFromSubscription,
  getSubscriptionInstances,
} from './gameSubscriptions'

describe('gameSubscriptions API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  const mockJson = (data) => ({
    ok: true,
    json: () => Promise.resolve(data),
  })

  describe('getMyGameSubscriptions', () => {
    it('calls GET /api/v1/game-subscriptions', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: [] }))
      await getMyGameSubscriptions()
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/game-subscriptions',
        expect.any(Object)
      )
    })
  })

  describe('createGameSubscription', () => {
    it('calls POST /api/v1/game-subscriptions with body { game_id, subscription_type, instance_limit? }', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: { id: 's1' } }))
      await createGameSubscription('g1', 'basic', 5)
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/game-subscriptions',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            game_id: 'g1',
            subscription_type: 'basic',
            instance_limit: 5,
          }),
        })
      )
    })

    it('omits instance_limit when null', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: { id: 's1' } }))
      await createGameSubscription('g1', 'basic')
      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          body: JSON.stringify({
            game_id: 'g1',
            subscription_type: 'basic',
          }),
        })
      )
    })
  })

  describe('cancelGameSubscription', () => {
    it('calls DELETE /api/v1/game-subscriptions/:subscriptionId', async () => {
      mockApiFetch.mockResolvedValue(mockJson({}))
      await cancelGameSubscription('s1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/game-subscriptions/s1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })

    it('returns null on 204 status', async () => {
      mockApiFetch.mockResolvedValue({ ok: true, status: 204 })
      const result = await cancelGameSubscription('s1')
      expect(result).toBeNull()
    })
  })

  describe('linkGameInstanceToSubscription', () => {
    it('calls POST /api/v1/game-subscriptions/:subscriptionId/instances with body', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: {} }))
      await linkGameInstanceToSubscription('s1', 'i1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/game-subscriptions/s1/instances',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            game_subscription_id: 's1',
            game_instance_id: 'i1',
          }),
        })
      )
    })
  })

  describe('unlinkGameInstanceFromSubscription', () => {
    it('calls DELETE /api/v1/game-subscriptions/:subscriptionId/instances/:instanceId', async () => {
      mockApiFetch.mockResolvedValue(mockJson({}))
      await unlinkGameInstanceFromSubscription('s1', 'i1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/game-subscriptions/s1/instances/i1',
        expect.objectContaining({ method: 'DELETE' })
      )
    })

    it('returns null on 204 status', async () => {
      mockApiFetch.mockResolvedValue({ ok: true, status: 204 })
      const result = await unlinkGameInstanceFromSubscription('s1', 'i1')
      expect(result).toBeNull()
    })
  })

  describe('getSubscriptionInstances', () => {
    it('calls GET /api/v1/game-subscriptions/:subscriptionId', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: { game_instance_ids: ['i1'] } }))
      const result = await getSubscriptionInstances('s1')
      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/game-subscriptions/s1',
        expect.any(Object)
      )
      expect(result).toEqual(['i1'])
    })

    it('returns data.data?.game_instance_ids or empty array', async () => {
      mockApiFetch.mockResolvedValue(mockJson({ data: {} }))
      const result = await getSubscriptionInstances('s1')
      expect(result).toEqual([])
    })
  })
})
