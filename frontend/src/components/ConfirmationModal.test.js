import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ConfirmationModal from './ConfirmationModal.vue'

describe('ConfirmationModal', () => {
  const defaultProps = {
    visible: true,
    title: 'Delete Item',
    message: 'Are you sure?',
    confirmText: 'Delete'
  }

  it('renders when visible is true', () => {
    const wrapper = mount(ConfirmationModal, {
      props: defaultProps
    })
    
    expect(wrapper.find('.modal-overlay').exists()).toBe(true)
    expect(wrapper.find('h2').text()).toBe('Delete Item')
    expect(wrapper.find('p').text()).toBe('Are you sure?')
  })

  it('does not render when visible is false', () => {
    const wrapper = mount(ConfirmationModal, {
      props: { ...defaultProps, visible: false }
    })
    
    expect(wrapper.find('.modal-overlay').exists()).toBe(false)
  })

  it('renders warning text when provided', () => {
    const wrapper = mount(ConfirmationModal, {
      props: {
        ...defaultProps,
        warning: 'This action cannot be undone'
      }
    })
    
    expect(wrapper.find('.warning-text').text()).toBe('This action cannot be undone')
  })

  it('renders confirmation input when requireConfirmation is true', () => {
    const wrapper = mount(ConfirmationModal, {
      props: {
        ...defaultProps,
        requireConfirmation: true,
        confirmationText: 'DELETE'
      }
    })
    
    expect(wrapper.find('.confirmation-input').exists()).toBe(true)
    expect(wrapper.find('input').attributes('placeholder')).toBe('DELETE')
  })

  it('emits cancel when overlay is clicked', async () => {
    const wrapper = mount(ConfirmationModal, {
      props: defaultProps
    })
    
    await wrapper.find('.modal-overlay').trigger('click')
    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('emits confirm when confirm button is clicked', async () => {
    const wrapper = mount(ConfirmationModal, {
      props: defaultProps
    })
    
    await wrapper.find('.danger-btn').trigger('click')
    expect(wrapper.emitted('confirm')).toBeTruthy()
  })

  it('disables confirm button when confirmation text does not match', async () => {
    const wrapper = mount(ConfirmationModal, {
      props: {
        ...defaultProps,
        requireConfirmation: true,
        confirmationText: 'DELETE'
      }
    })
    
    const input = wrapper.find('input')
    await input.setValue('WRONG')
    
    const confirmBtn = wrapper.find('.danger-btn')
    expect(confirmBtn.attributes('disabled')).toBeDefined()
  })

  it('enables confirm button when confirmation text matches', async () => {
    const wrapper = mount(ConfirmationModal, {
      props: {
        ...defaultProps,
        requireConfirmation: true,
        confirmationText: 'DELETE'
      }
    })
    
    const input = wrapper.find('input')
    await input.setValue('DELETE')
    
    const confirmBtn = wrapper.find('.danger-btn')
    expect(confirmBtn.attributes('disabled')).toBeUndefined()
  })

  it('shows loading text when loading is true', () => {
    const wrapper = mount(ConfirmationModal, {
      props: {
        ...defaultProps,
        loading: true,
        loadingText: 'Deleting...'
      }
    })
    
    expect(wrapper.find('.danger-btn').text()).toBe('Deleting...')
  })
}) 