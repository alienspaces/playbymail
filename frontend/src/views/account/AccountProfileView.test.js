import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

const mockGetMe = vi.fn()
const mockGetAccount = vi.fn()
const mockUpdateAccount = vi.fn()
const mockDeleteAccountUser = vi.fn()
const mockSetAccountTimezone = vi.fn()

vi.mock('@/api/account', () => ({
  getMe: (...args) => mockGetMe(...args),
  getAccount: (...args) => mockGetAccount(...args),
  updateAccount: (...args) => mockUpdateAccount(...args),
  deleteAccountUser: (...args) => mockDeleteAccountUser(...args),
}))

vi.mock('@/stores/auth', () => ({
  useAuthStore: vi.fn(() => ({
    accountTimezone: null,
    setAccountTimezone: mockSetAccountTimezone,
    logout: vi.fn(),
  })),
}))

vi.mock('@/utils/dateFormat', () => ({
  formatDateTime: vi.fn((str) => (str ? 'Formatted: ' + str : 'N/A')),
}))

vi.mock('@/components/ConfirmationModal.vue', () => ({
  default: {
    name: 'ConfirmationModal',
    props: [
      'visible',
      'title',
      'message',
      'warning',
      'confirmText',
      'loading',
      'requireConfirmation',
      'confirmationText',
    ],
    emits: ['confirm', 'cancel'],
    template: '<div class="mock-confirmation-modal" />',
  },
}))

vi.mock('@/components/PageHeader.vue', () => ({
  default: {
    name: 'PageHeader',
    props: ['title', 'showIcon', 'titleLevel', 'subtitle'],
    template: '<div class="mock-page-header"><slot /></div>',
  },
}))

vi.mock('@/components/DataCard.vue', () => ({
  default: {
    name: 'DataCard',
    props: ['title', 'variant'],
    template: '<div class="mock-data-card"><slot /><slot name="primary" /></div>',
  },
}))

vi.mock('@/components/DataItem.vue', () => ({
  default: {
    name: 'DataItem',
    props: ['label', 'value'],
    template: '<div class="mock-data-item">{{ label }}: {{ value }}</div>',
  },
}))

vi.mock('@/components/Button.vue', () => ({
  default: {
    name: 'AppButton',
    props: ['variant', 'size', 'disabled'],
    emits: ['click'],
    template: '<button :disabled="disabled" @click="$emit(\'click\')"><slot /></button>',
  },
}))

import AccountProfileView from './AccountProfileView.vue'

describe('AccountProfileView', () => {
  const meData = {
    id: 'user-1',
    account_id: 'acct-1',
    email: 'player@example.com',
    created_at: '2026-01-01T00:00:00Z',
  }
  const accountData = { id: 'acct-1', name: 'Test Player', timezone: null }

  beforeEach(() => {
    vi.clearAllMocks()
    mockGetMe.mockResolvedValue(meData)
    mockGetAccount.mockResolvedValue(accountData)
    mockUpdateAccount.mockResolvedValue({ ...accountData, name: 'Updated' })
    mockDeleteAccountUser.mockResolvedValue(true)
  })

  it('loads account data on mount', async () => {
    mount(AccountProfileView)
    await flushPromises()

    expect(mockGetMe).toHaveBeenCalled()
    expect(mockGetAccount).toHaveBeenCalledWith('acct-1')
  })

  it('propagates account timezone to auth store on load', async () => {
    mockGetAccount.mockResolvedValue({ ...accountData, timezone: 'Australia/Sydney' })
    mount(AccountProfileView)
    await flushPromises()

    expect(mockSetAccountTimezone).toHaveBeenCalledWith('Australia/Sydney')
  })

  it('calls setAccountTimezone with null when account has no timezone', async () => {
    mount(AccountProfileView)
    await flushPromises()

    expect(mockSetAccountTimezone).toHaveBeenCalledWith(null)
  })

  it('displays email from account user', async () => {
    const wrapper = mount(AccountProfileView)
    await flushPromises()

    expect(wrapper.text()).toContain('player@example.com')
  })

  it('displays formatted created_at date', async () => {
    const wrapper = mount(AccountProfileView)
    await flushPromises()

    expect(wrapper.text()).toContain('Formatted: 2026-01-01T00:00:00Z')
  })

  it('shows loading state before data arrives', () => {
    mockGetMe.mockReturnValue(new Promise(() => {}))
    const wrapper = mount(AccountProfileView)

    expect(wrapper.text()).toContain('Loading account information')
  })

  it('shows error when getMe fails', async () => {
    mockGetMe.mockRejectedValue(new Error('Auth failed'))
    const wrapper = mount(AccountProfileView)
    await flushPromises()

    expect(wrapper.text()).toContain('Auth failed')
  })

  it('shows error when getAccount fails', async () => {
    mockGetAccount.mockRejectedValue(new Error('Account not found'))
    const wrapper = mount(AccountProfileView)
    await flushPromises()

    expect(wrapper.text()).toContain('Account not found')
  })

  it('saves updated account name via updateAccount', async () => {
    const wrapper = mount(AccountProfileView)
    await flushPromises()

    // Click the Edit button for account name
    const editBtn = wrapper.findAll('button').find((b) => b.text().trim() === 'Edit')
    await editBtn.trigger('click')

    // Fill in name input
    const input = wrapper.find('input.name-input')
    await input.setValue('New Name')

    // Click Save
    const saveBtn = wrapper.findAll('button').find((b) => b.text().trim() === 'Save')
    await saveBtn.trigger('click')
    await flushPromises()

    expect(mockUpdateAccount).toHaveBeenCalledWith('acct-1', { name: 'New Name' })
  })

  it('saves timezone via updateAccount when timezone is changed', async () => {
    mockUpdateAccount.mockResolvedValue({ ...accountData, timezone: 'America/New_York' })
    const wrapper = mount(AccountProfileView)
    await flushPromises()

    // Find the timezone Edit button (second Edit button)
    const editBtns = wrapper.findAll('button').filter((b) => b.text().trim() === 'Edit')
    await editBtns[editBtns.length - 1].trigger('click')

    // Set timezone select
    const select = wrapper.find('select.timezone-select')
    await select.setValue('America/New_York')

    // Click Save
    const saveBtn = wrapper.findAll('button').find((b) => b.text().trim() === 'Save')
    await saveBtn.trigger('click')
    await flushPromises()

    expect(mockUpdateAccount).toHaveBeenCalledWith('acct-1', { timezone: 'America/New_York' })
    expect(mockSetAccountTimezone).toHaveBeenCalledWith('America/New_York')
  })
})
