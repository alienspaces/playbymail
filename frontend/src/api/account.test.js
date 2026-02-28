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
  getMe,
  deleteAccountUser,
  getAccountContacts,
  getAccountContact,
  createAccountContact,
  updateAccountContact,
  deleteAccountContact,
} from './account'

describe('account API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockHandleApiError.mockImplementation((res) => res)
  })

  describe('getMe', () => {
    it('calls GET /api/v1/me and returns data', async () => {
      const meData = { id: 'user-1', account_id: 'acct-1', email: 'test@example.com' }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: meData }),
      })

      const result = await getMe()

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/me',
        expect.objectContaining({
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual(meData)
    })
  })

  describe('deleteAccountUser', () => {
    it('calls DELETE /api/v1/accounts/:accountId/users/:accountUserId', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      const result = await deleteAccountUser('acct-1', 'user-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/accounts/acct-1/users/user-1',
        expect.objectContaining({ method: 'DELETE' })
      )
      expect(result).toBe(true)
    })
  })

  describe('getAccountContacts', () => {
    it('builds correct nested path with accountId and accountUserId', async () => {
      const contacts = [{ id: 'c1', name: 'Contact 1' }]
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: contacts }),
      })

      const result = await getAccountContacts('acct-1', 'user-1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/accounts/acct-1/users/user-1/contacts',
        expect.objectContaining({
          headers: expect.objectContaining({ Authorization: 'Bearer test-token' }),
        })
      )
      expect(result).toEqual(contacts)
    })

    it('returns empty array when data is null', async () => {
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: null }),
      })

      const result = await getAccountContacts('acct-1', 'user-1')
      expect(result).toEqual([])
    })
  })

  describe('getAccountContact', () => {
    it('builds correct path with all three IDs', async () => {
      const contact = { id: 'c1', name: 'Contact 1' }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: contact }),
      })

      const result = await getAccountContact('acct-1', 'user-1', 'c1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/accounts/acct-1/users/user-1/contacts/c1',
        expect.any(Object)
      )
      expect(result).toEqual(contact)
    })
  })

  describe('createAccountContact', () => {
    it('sends POST with contact data to the correct nested path', async () => {
      const contactData = { name: 'New Contact', postal_address_line1: '123 St' }
      const created = { id: 'c-new', ...contactData }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: created }),
      })

      const result = await createAccountContact('acct-1', 'user-1', contactData)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/accounts/acct-1/users/user-1/contacts',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(contactData),
        })
      )
      expect(result).toEqual(created)
    })
  })

  describe('updateAccountContact', () => {
    it('sends PUT with contact data to the correct nested path', async () => {
      const contactData = { name: 'Updated Contact' }
      const updated = { id: 'c1', ...contactData }
      mockApiFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ data: updated }),
      })

      const result = await updateAccountContact('acct-1', 'user-1', 'c1', contactData)

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/accounts/acct-1/users/user-1/contacts/c1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(contactData),
        })
      )
      expect(result).toEqual(updated)
    })
  })

  describe('deleteAccountContact', () => {
    it('sends DELETE to the correct nested path', async () => {
      mockApiFetch.mockResolvedValue({ ok: true })

      const result = await deleteAccountContact('acct-1', 'user-1', 'c1')

      expect(mockApiFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/accounts/acct-1/users/user-1/contacts/c1',
        expect.objectContaining({ method: 'DELETE' })
      )
      expect(result).toBe(true)
    })
  })
})
