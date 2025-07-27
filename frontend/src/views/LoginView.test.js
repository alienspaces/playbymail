import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import LoginView from './LoginView.vue'

// Mock the auth API
vi.mock('../api/auth', () => ({
  requestAuth: vi.fn()
}))

// Mock router - not used in current tests but available for future use
// const router = createRouter({
//   history: createWebHistory(),
//   routes: [
//     { path: '/login', component: LoginView },
//     { path: '/verify', component: { template: '<div>Verify</div>' } }
//   ]
// })

describe('LoginView', () => {
  const mountWithMocks = (routeQuery = {}) => mount(LoginView, {
    global: {
      mocks: {
        $route: {
          query: routeQuery
        },
        $router: {
          push: vi.fn()
        }
      }
    }
  })

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders the login form', () => {
    const wrapper = mountWithMocks()
    
    expect(wrapper.find('h2').text()).toBe('Sign in with Email')
    expect(wrapper.find('form').exists()).toBe(true)
    expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('has correct form structure', () => {
    const wrapper = mountWithMocks()
    
    const emailInput = wrapper.find('#email')
    expect(emailInput.exists()).toBe(true)
    expect(emailInput.attributes('type')).toBe('email')
    expect(emailInput.attributes('required')).toBeDefined()
    expect(emailInput.attributes('autofocus')).toBeDefined()
    
    const label = wrapper.find('label[for="email"]')
    expect(label.text()).toBe('Email address')
  })

  it('updates email value when input changes', async () => {
    const wrapper = mountWithMocks()
    
    const emailInput = wrapper.find('#email')
    await emailInput.setValue('test@example.com')
    
    expect(wrapper.vm.email).toBe('test@example.com')
  })

  it('shows loading state when form is submitted', async () => {
    const { requestAuth } = await import('../api/auth')
    requestAuth.mockImplementation(() => new Promise(resolve => setTimeout(() => resolve(true), 100)))
    
    const wrapper = mountWithMocks()
    
    const emailInput = wrapper.find('#email')
    await emailInput.setValue('test@example.com')
    
    const form = wrapper.find('form')
    await form.trigger('submit')
    
    // Check loading state immediately after submission
    expect(wrapper.vm.loading).toBe(true)
    const submitButton = wrapper.find('button[type="submit"]')
    expect(submitButton.attributes('disabled')).toBeDefined()
  })

  it('displays error message on failed submission', async () => {
    const { requestAuth } = await import('../api/auth')
    requestAuth.mockResolvedValue(false)
    
    const wrapper = mountWithMocks()
    
    const emailInput = wrapper.find('#email')
    await emailInput.setValue('test@example.com')
    
    const form = wrapper.find('form')
    await form.trigger('submit')
    
    await wrapper.vm.$nextTick()
    
    expect(wrapper.find('.message').text()).toBe('Failed to send verification email.')
  })

  it('displays error message on API exception', async () => {
    const { requestAuth } = await import('../api/auth')
    requestAuth.mockRejectedValue(new Error('Network error'))
    
    const wrapper = mountWithMocks()
    
    const emailInput = wrapper.find('#email')
    await emailInput.setValue('test@example.com')
    
    const form = wrapper.find('form')
    await form.trigger('submit')
    
    await wrapper.vm.$nextTick()
    
    expect(wrapper.find('.message').text()).toBe('Failed to send verification email.')
  })

  it('navigates to verify page on successful submission', async () => {
    const { requestAuth } = await import('../api/auth')
    requestAuth.mockResolvedValue(true)
    
    const mockPush = vi.fn()
    const wrapper = mount(LoginView, {
      global: {
        mocks: {
          $route: { query: {} },
          $router: {
            push: mockPush
          }
        }
      }
    })
    
    const emailInput = wrapper.find('#email')
    await emailInput.setValue('test@example.com')
    
    const form = wrapper.find('form')
    await form.trigger('submit')
    
    expect(mockPush).toHaveBeenCalledWith({
      path: '/verify',
      query: { email: 'test@example.com' }
    })
  })

  it('clears message when form is submitted', async () => {
    const { requestAuth } = await import('../api/auth')
    requestAuth.mockResolvedValue(true)
    
    const wrapper = mountWithMocks()
    
    // Set initial message
    wrapper.vm.message = 'Previous error message'
    await wrapper.vm.$nextTick()
    
    const emailInput = wrapper.find('#email')
    await emailInput.setValue('test@example.com')
    
    const form = wrapper.find('form')
    await form.trigger('submit')
    
    expect(wrapper.vm.message).toBe('')
  })

  it('displays session expired message from query parameter', () => {
    const wrapper = mountWithMocks({ code: 'session_expired' })
    
    expect(wrapper.find('.message').text()).toBe('Session expired. Please log in again.')
  })

  it('does not display message when no query parameter', () => {
    const wrapper = mountWithMocks()
    
    expect(wrapper.find('.message').exists()).toBe(false)
  })

  it('has correct CSS classes for styling', () => {
    const wrapper = mountWithMocks()
    
    expect(wrapper.find('.login-container').exists()).toBe(true)
    expect(wrapper.find('.card').exists()).toBe(true)
    expect(wrapper.find('.login-form').exists()).toBe(true)
    expect(wrapper.find('.form-group').exists()).toBe(true)
    expect(wrapper.find('.form-actions').exists()).toBe(true)
  })

  it('button text is correct', () => {
    const wrapper = mountWithMocks()
    expect(wrapper.find('button[type="submit"]').text()).toBe('Send Code')
  })

  it('form prevents default submission', async () => {
    const wrapper = mountWithMocks()
    
    const form = wrapper.find('form')
    
    // Mock the onSubmit method to check if it's called
    const onSubmitSpy = vi.spyOn(wrapper.vm, 'onSubmit')
    
    await form.trigger('submit')
    
    expect(onSubmitSpy).toHaveBeenCalled()
  })
}) 