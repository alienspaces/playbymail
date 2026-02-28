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
  uploadLocationTurnSheetImage,
  getLocationTurnSheetImage,
  deleteLocationTurnSheetImage,
  getLocationChoiceTurnSheetPreviewUrl,
} from './locationImages'

describe('locationImages API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('uploadLocationTurnSheetImage', () => {
    it('calls POST with FormData body and auth headers (no Content-Type)', async () => {
      const imageFile = { name: 'test.png' }
      const responseData = { data: { id: 'img1' } }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(responseData),
      })

      const result = await uploadLocationTurnSheetImage('game-1', 'loc-1', imageFile)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/locations/loc-1/turn-sheet-image',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      const callOptions = mockApiFetch.mock.calls[0][1]
      expect(callOptions.body).toBeInstanceOf(FormData)
      expect(callOptions.body.has('image')).toBe(true)
      expect(callOptions.headers['Content-Type']).toBeUndefined()
      expect(result).toEqual(responseData)
    })
  })

  describe('getLocationTurnSheetImage', () => {
    it('calls GET and returns json', async () => {
      const responseData = { data: { id: 'img1' } }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(responseData),
      })

      const result = await getLocationTurnSheetImage('game-1', 'loc-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/locations/loc-1/turn-sheet-image',
        expect.objectContaining({
          headers: expect.objectContaining({ 'Content-Type': 'application/json', Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual(responseData)
    })
  })

  describe('deleteLocationTurnSheetImage', () => {
    it('calls DELETE and does not call handleApiError when status is 204', async () => {
      mockApiFetch.mockResolvedValue({ ok: true, status: 204 })

      await deleteLocationTurnSheetImage('game-1', 'loc-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/adventure-games/game-1/locations/loc-1/turn-sheet-image',
        expect.objectContaining({
          method: 'DELETE',
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(mockHandleApiError).not.toHaveBeenCalled()
    })

    it('calls handleApiError when status is not 204', async () => {
      mockApiFetch.mockResolvedValue({ ok: false, status: 500 })
      mockHandleApiError.mockImplementation(() => {
        throw new Error('Failed to delete turn sheet image')
      })

      await expect(deleteLocationTurnSheetImage('game-1', 'loc-1')).rejects.toThrow()
      expect(mockHandleApiError).toHaveBeenCalled()
    })
  })

  describe('getLocationChoiceTurnSheetPreviewUrl', () => {
    it('returns synchronous URL string without fetching', () => {
      const result = getLocationChoiceTurnSheetPreviewUrl('game-1', 'loc-1')

      expect(result).toBe(
        'http://localhost:8080/api/v1/adventure-games/game-1/locations/loc-1/turn-sheets/preview'
      )
      expect(mockApiFetch).not.toHaveBeenCalled()
    })
  })
})
