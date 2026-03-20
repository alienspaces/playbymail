import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockFetch = vi.fn()
globalThis.fetch = mockFetch

const mockHandleApiError = vi.fn()

vi.mock('./baseUrl', () => ({
  baseUrl: 'http://localhost:8080',
  handleApiError: (...args) => mockHandleApiError(...args),
}))

import { approveSubscription } from './approveSubscription'

const SUB_ID = 'sub-abc-123'
const EMAIL = 'player@example.com'
const BASE = 'http://localhost:8080/api/v1/game-subscriptions'

describe('approveSubscription API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation(() => {})
  })

  it('calls POST /api/v1/game-subscriptions/:id/approve with email query param', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: { status: 'active' } }),
    })

    await approveSubscription(SUB_ID, EMAIL)

    expect(mockFetch).toHaveBeenCalledWith(
      `${BASE}/${SUB_ID}/approve?email=${encodeURIComponent(EMAIL)}`,
      expect.objectContaining({
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      })
    )
  })

  it('returns parsed JSON response on success', async () => {
    const responseData = { data: { status: 'active' } }
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(responseData),
    })

    const result = await approveSubscription(SUB_ID, EMAIL)

    expect(result).toEqual(responseData)
  })

  it('calls handleApiError on failure', async () => {
    const errorRes = { ok: false, status: 400 }
    mockFetch.mockResolvedValue(errorRes)
    mockHandleApiError.mockRejectedValue(new Error('Failed to confirm subscription'))

    await expect(approveSubscription(SUB_ID, EMAIL)).rejects.toThrow('Failed to confirm subscription')
    expect(mockHandleApiError).toHaveBeenCalledWith(errorRes, 'Failed to confirm subscription')
  })

  it('URL-encodes the email address', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: { status: 'active' } }),
    })

    const emailWithPlus = 'player+tag@example.com'
    await approveSubscription(SUB_ID, emailWithPlus)

    const calledUrl = mockFetch.mock.calls[0][0]
    expect(calledUrl).toContain(encodeURIComponent(emailWithPlus))
    expect(calledUrl).not.toContain('+')
  })
})
