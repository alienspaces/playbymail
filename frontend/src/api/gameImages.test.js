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
  uploadGameTurnSheetImage,
  getGameTurnSheetImages,
  getGameTurnSheetImage,
  deleteGameTurnSheetImage,
  getGameTurnSheetPreviewUrl,
  getGameTurnSheetImageLegacy,
  getGameJoinTurnSheetPreviewUrl,
} from './gameImages'

describe('gameImages API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('uploadGameTurnSheetImage', () => {
    it('calls POST with FormData and turn_sheet_type query param', async () => {
      const imageFile = { name: 'test.png' }
      const responseData = { data: { id: 'img1' } }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(responseData),
      })

      const result = await uploadGameTurnSheetImage('game-1', imageFile, 'adventure_game_inventory_management')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/games/game-1/turn-sheet-images?turn_sheet_type=adventure_game_inventory_management'),
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      const callOptions = mockApiFetch.mock.calls[0][1]
      expect(callOptions.body).toBeInstanceOf(FormData)
      expect(callOptions.body.has('image')).toBe(true)
      expect(result).toEqual(responseData)
    })

    it('defaults turn_sheet_type to adventure_game_join_game', async () => {
      mockApiFetch.mockResolvedValue({ ok: true, json: () => Promise.resolve({}) })

      await uploadGameTurnSheetImage('game-1', { name: 'test.png' })

      const url = mockApiFetch.mock.calls[0][0]
      expect(url).toContain('turn_sheet_type=adventure_game_join_game')
    })
  })

  describe('getGameTurnSheetImages', () => {
    it('calls GET with optional turnSheetType query param', async () => {
      const responseData = { data: [] }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(responseData),
      })

      const result = await getGameTurnSheetImages('game-1', 'adventure_game_join_game')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/games/game-1/turn-sheet-images?turn_sheet_type=adventure_game_join_game'),
        expect.any(Object)
      )
      expect(result).toEqual(responseData)
    })

    it('calls GET without turn_sheet_type when not provided', async () => {
      mockApiFetch.mockResolvedValue({ ok: true, json: () => Promise.resolve({ data: [] }) })

      await getGameTurnSheetImages('game-1')

      const url = mockApiFetch.mock.calls[0][0]
      expect(url).not.toContain('turn_sheet_type')
    })
  })

  describe('getGameTurnSheetImage', () => {
    it('calls GET /api/v1/games/:gameId/turn-sheet-images/:imageId', async () => {
      const responseData = { data: { id: 'img1' } }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(responseData),
      })

      const result = await getGameTurnSheetImage('game-1', 'img1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/game-1/turn-sheet-images/img1',
        expect.objectContaining({
          headers: expect.objectContaining({ 'Content-Type': 'application/json', Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual(responseData)
    })
  })

  describe('deleteGameTurnSheetImage', () => {
    it('calls DELETE and does not call handleApiError when status is 204', async () => {
      mockApiFetch.mockResolvedValue({ ok: true, status: 204 })

      await deleteGameTurnSheetImage('game-1', 'img1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/games/game-1/turn-sheet-images/img1',
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

      await expect(deleteGameTurnSheetImage('game-1', 'img1')).rejects.toThrow()
      expect(mockHandleApiError).toHaveBeenCalled()
    })
  })

  describe('getGameTurnSheetPreviewUrl', () => {
    it('returns synchronous URL string with turn_sheet_type param', () => {
      const result = getGameTurnSheetPreviewUrl('game-1', 'adventure_game_join_game')

      expect(result).toContain('http://localhost:8080/api/v1/games/game-1/turn-sheets/preview')
      expect(result).toContain('turn_sheet_type=adventure_game_join_game')
      expect(mockApiFetch).not.toHaveBeenCalled()
    })
  })

  describe('getGameTurnSheetImageLegacy', () => {
    it('calls getGameTurnSheetImages with adventure_game_join_game and returns transformed format', async () => {
      const apiResponse = { data: [{ id: 'img1', url: 'http://example.com/img1' }] }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(apiResponse),
      })

      const result = await getGameTurnSheetImageLegacy('game-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        expect.stringContaining('turn_sheet_type=adventure_game_join_game'),
        expect.any(Object)
      )
      expect(result).toEqual({
        data: {
          background: apiResponse.data[0],
        },
      })
    })

    it('returns background null when no images', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: [] }),
      })

      const result = await getGameTurnSheetImageLegacy('game-1')

      expect(result).toEqual({ data: { background: null } })
    })
  })

  describe('getGameJoinTurnSheetPreviewUrl', () => {
    it('returns URL from getGameTurnSheetPreviewUrl with adventure_game_join_game', () => {
      const result = getGameJoinTurnSheetPreviewUrl('game-1')

      expect(result).toContain('http://localhost:8080/api/v1/games/game-1/turn-sheets/preview')
      expect(result).toContain('turn_sheet_type=adventure_game_join_game')
      expect(mockApiFetch).not.toHaveBeenCalled()
    })
  })
})
