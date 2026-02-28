import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'

const mockCreateAccountContact = vi.fn()
const mockUpdateAccountContact = vi.fn()

vi.mock('@/api/account', () => ({
  createAccountContact: (...args) => mockCreateAccountContact(...args),
  updateAccountContact: (...args) => mockUpdateAccountContact(...args),
}))

import ContactModal from './ContactModal.vue'

describe('ContactModal', () => {
  const baseProps = {
    visible: true,
    contact: null,
    accountId: 'acct-1',
    accountUserId: 'user-1',
  }

  beforeEach(() => {
    vi.clearAllMocks()
    mockCreateAccountContact.mockResolvedValue({ id: 'c-new' })
    mockUpdateAccountContact.mockResolvedValue({ id: 'c1' })
  })

  it('renders create form when no contact is provided', () => {
    const wrapper = mount(ContactModal, { props: baseProps })
    expect(wrapper.find('h2').text()).toBe('Add Contact')
    expect(wrapper.find('button[type="submit"]').text()).toBe('Create')
  })

  it('renders edit form when contact is provided', () => {
    const wrapper = mount(ContactModal, {
      props: {
        ...baseProps,
        contact: { id: 'c1', name: 'Existing', postal_address_line1: '123 St', state_province: 'CA', country: 'US', postal_code: '90210' },
      },
    })
    expect(wrapper.find('h2').text()).toBe('Edit Contact')
    expect(wrapper.find('button[type="submit"]').text()).toBe('Update')
  })

  it('populates form fields from contact prop', () => {
    const contact = {
      id: 'c1',
      name: 'Test Name',
      postal_address_line1: '456 Ave',
      postal_address_line2: 'Suite B',
      state_province: 'NY',
      country: 'US',
      postal_code: '10001',
    }
    const wrapper = mount(ContactModal, { props: { ...baseProps, contact } })

    expect(wrapper.find('#name').element.value).toBe('Test Name')
    expect(wrapper.find('#postal_address_line1').element.value).toBe('456 Ave')
    expect(wrapper.find('#postal_address_line2').element.value).toBe('Suite B')
    expect(wrapper.find('#state_province').element.value).toBe('NY')
    expect(wrapper.find('#country').element.value).toBe('US')
    expect(wrapper.find('#postal_code').element.value).toBe('10001')
  })

  it('calls createAccountContact with accountId and accountUserId on submit for new contact', async () => {
    const wrapper = mount(ContactModal, { props: baseProps })

    await wrapper.find('#name').setValue('New Person')
    await wrapper.find('#postal_address_line1').setValue('789 Blvd')
    await wrapper.find('#state_province').setValue('TX')
    await wrapper.find('#country').setValue('US')
    await wrapper.find('#postal_code').setValue('75001')

    await wrapper.find('form').trigger('submit')
    await vi.dynamicImportSettled()

    expect(mockCreateAccountContact).toHaveBeenCalledWith(
      'acct-1',
      'user-1',
      expect.objectContaining({
        name: 'New Person',
        postal_address_line1: '789 Blvd',
      })
    )
  })

  it('calls updateAccountContact with accountId, accountUserId, and contactId on submit for existing contact', async () => {
    const contact = {
      id: 'c1',
      name: 'Old Name',
      postal_address_line1: '123 St',
      state_province: 'CA',
      country: 'US',
      postal_code: '90210',
    }
    const wrapper = mount(ContactModal, { props: { ...baseProps, contact } })

    await wrapper.find('#name').setValue('Updated Name')
    await wrapper.find('form').trigger('submit')
    await vi.dynamicImportSettled()

    expect(mockUpdateAccountContact).toHaveBeenCalledWith(
      'acct-1',
      'user-1',
      'c1',
      expect.objectContaining({ name: 'Updated Name' })
    )
  })

  it('shows error when accountId or accountUserId is missing', async () => {
    const wrapper = mount(ContactModal, {
      props: { ...baseProps, accountId: '', accountUserId: '' },
    })

    await wrapper.find('#name').setValue('Test')
    await wrapper.find('#postal_address_line1').setValue('123 St')
    await wrapper.find('#state_province').setValue('CA')
    await wrapper.find('#country').setValue('US')
    await wrapper.find('#postal_code').setValue('90210')
    await wrapper.find('form').trigger('submit')

    expect(wrapper.find('.error').exists()).toBe(true)
    expect(mockCreateAccountContact).not.toHaveBeenCalled()
  })

  it('emits close when cancel button is clicked', async () => {
    const wrapper = mount(ContactModal, { props: baseProps })
    await wrapper.find('button[type="button"]').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('does not render when visible is false', () => {
    const wrapper = mount(ContactModal, { props: { ...baseProps, visible: false } })
    expect(wrapper.find('.modal-overlay').exists()).toBe(false)
  })
})
