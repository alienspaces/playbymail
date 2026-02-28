import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

const mockGetMe = vi.fn()
const mockGetAccountContacts = vi.fn()
const mockDeleteAccountContact = vi.fn()

vi.mock('@/api/account', () => ({
  getMe: (...args) => mockGetMe(...args),
  getAccountContacts: (...args) => mockGetAccountContacts(...args),
  deleteAccountContact: (...args) => mockDeleteAccountContact(...args),
}))

vi.mock('@/components/ContactModal.vue', () => ({
  default: {
    name: 'ContactModal',
    props: ['visible', 'contact', 'accountId', 'accountUserId'],
    template: '<div class="mock-contact-modal" />',
  },
}))

vi.mock('@/components/AppButton.vue', () => ({
  default: {
    name: 'AppButton',
    template: '<button><slot /></button>',
  },
}))

vi.mock('@/components/ConfirmationModal.vue', () => ({
  default: {
    name: 'ConfirmationModal',
    props: ['visible', 'title', 'message', 'confirmText', 'loading'],
    template: '<div class="mock-confirmation-modal" />',
  },
}))

vi.mock('@/components/PageHeader.vue', () => ({
  default: {
    name: 'PageHeader',
    props: ['title', 'actionText', 'showIcon', 'titleLevel', 'subtitle'],
    template: '<div class="mock-page-header"><slot /></div>',
  },
}))

import AccountContactsView from './AccountContactsView.vue'

describe('AccountContactsView', () => {
  const meData = { id: 'user-1', account_id: 'acct-1', email: 'test@example.com' }
  const contactsData = [
    { id: 'c1', name: 'Home', postal_address_line1: '123 Main St', state_province: 'CA', country: 'US', postal_code: '90210' },
    { id: 'c2', name: 'Office', postal_address_line1: '456 Work Ave', state_province: 'NY', country: 'US', postal_code: '10001' },
  ]

  beforeEach(() => {
    vi.clearAllMocks()
    mockGetMe.mockResolvedValue(meData)
    mockGetAccountContacts.mockResolvedValue(contactsData)
    mockDeleteAccountContact.mockResolvedValue(true)
  })

  it('fetches me and contacts with correct IDs on mount', async () => {
    mount(AccountContactsView)
    await flushPromises()

    expect(mockGetMe).toHaveBeenCalled()
    expect(mockGetAccountContacts).toHaveBeenCalledWith('acct-1', 'user-1')
  })

  it('displays contacts after loading', async () => {
    const wrapper = mount(AccountContactsView)
    await flushPromises()

    expect(wrapper.text()).toContain('Home')
    expect(wrapper.text()).toContain('Office')
  })

  it('shows loading state initially', () => {
    mockGetMe.mockReturnValue(new Promise(() => {}))
    const wrapper = mount(AccountContactsView)
    expect(wrapper.text()).toContain('Loading contacts...')
  })

  it('shows error when getMe fails', async () => {
    mockGetMe.mockRejectedValue(new Error('Network error'))
    const wrapper = mount(AccountContactsView)
    await flushPromises()
    expect(wrapper.text()).toContain('Network error')
  })
})
