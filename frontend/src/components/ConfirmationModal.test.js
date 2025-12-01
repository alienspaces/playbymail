import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import ConfirmationModal from './ConfirmationModal.vue'

describe('ConfirmationModal', () => {
  const defaultProps = {
    visible: true,
    title: 'Delete Item',
    message: 'Are you sure?',
    confirmText: 'Delete'
  }

  // Helper to find elements in document body (where Teleport renders)
  const findInBody = (selector) => {
    return document.body.querySelector(selector)
  }

  const findAllInBody = (selector) => {
    return document.body.querySelectorAll(selector)
  }

  beforeEach(() => {
    // Clear any existing modals from previous tests
    document.body.innerHTML = ''
  })

  afterEach(() => {
    // Clean up after each test
    document.body.innerHTML = ''
  })

  it('renders when visible is true', () => {
    mount(ConfirmationModal, {
      props: defaultProps
    })

    expect(findInBody('.modal-overlay')).toBeTruthy()
    expect(findInBody('h2').textContent).toBe('Delete Item')
    expect(findInBody('p').textContent).toBe('Are you sure?')
  })

  it('does not render when visible is false', () => {
    mount(ConfirmationModal, {
      props: { ...defaultProps, visible: false }
    })

    expect(findInBody('.modal-overlay')).toBeNull()
  })

  it('renders warning text when provided', () => {
    mount(ConfirmationModal, {
      props: {
        ...defaultProps,
        warning: 'This action cannot be undone'
      }
    })

    expect(findInBody('.warning-text').textContent).toBe('This action cannot be undone')
  })

  it('renders confirmation input when requireConfirmation is true', () => {
    mount(ConfirmationModal, {
      props: {
        ...defaultProps,
        requireConfirmation: true,
        confirmationText: 'DELETE'
      }
    })

    expect(findInBody('.confirmation-input')).toBeTruthy()
    expect(findInBody('input').getAttribute('placeholder')).toBe('DELETE')
  })

  it('emits cancel when overlay is clicked', async () => {
    const wrapper = mount(ConfirmationModal, {
      props: defaultProps
    })

    const overlay = findInBody('.modal-overlay')
    overlay.click()
    await wrapper.vm.$nextTick()

    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('emits confirm when confirm button is clicked', async () => {
    const wrapper = mount(ConfirmationModal, {
      props: defaultProps
    })

    const confirmBtn = findInBody('.danger-btn')
    confirmBtn.click()
    await wrapper.vm.$nextTick()

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

    const input = findInBody('input')
    input.value = 'WRONG'
    input.dispatchEvent(new Event('input'))
    await wrapper.vm.$nextTick()

    const confirmBtn = findInBody('.danger-btn')
    expect(confirmBtn.hasAttribute('disabled')).toBe(true)
  })

  it('enables confirm button when confirmation text matches', async () => {
    const wrapper = mount(ConfirmationModal, {
      props: {
        ...defaultProps,
        requireConfirmation: true,
        confirmationText: 'DELETE'
      }
    })

    const input = findInBody('input')
    input.value = 'DELETE'
    input.dispatchEvent(new Event('input'))
    await wrapper.vm.$nextTick()

    const confirmBtn = findInBody('.danger-btn')
    expect(confirmBtn.hasAttribute('disabled')).toBe(false)
  })

  it('shows loading text when loading is true', () => {
    mount(ConfirmationModal, {
      props: {
        ...defaultProps,
        loading: true,
        loadingText: 'Deleting...'
      }
    })

    expect(findInBody('.danger-btn').textContent).toBe('Deleting...')
  })

  it('renders custom cancel text', () => {
    mount(ConfirmationModal, {
      props: {
        ...defaultProps,
        cancelText: 'Never mind'
      }
    })

    const buttons = findAllInBody('button')
    expect(buttons[0].textContent).toBe('Never mind')
  })

  it('renders error text when provided', () => {
    mount(ConfirmationModal, {
      props: {
        ...defaultProps,
        error: 'Something went wrong'
      }
    })

    expect(findInBody('.error-text').textContent).toBe('Something went wrong')
  })
}) 