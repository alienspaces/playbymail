import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import PlayerConfirmSubscriptionView from './PlayerConfirmSubscriptionView.vue'

const mockApproveSubscription = vi.fn()

vi.mock('../api/approveSubscription', () => ({
  approveSubscription: (...args) => mockApproveSubscription(...args),
}))

vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({
    params: { game_subscription_id: 'sub-abc-123' },
    query: { email: 'player@example.com' },
  })),
}))

describe('PlayerConfirmSubscriptionView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loading state while confirming', async () => {
    mockApproveSubscription.mockReturnValue(new Promise(() => {}))

    const wrapper = mount(PlayerConfirmSubscriptionView)

    expect(wrapper.find('[data-testid="confirm-loading"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="confirm-success"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="confirm-error"]').exists()).toBe(false)
  })

  it('shows success state after successful confirmation', async () => {
    mockApproveSubscription.mockResolvedValue({ data: { status: 'active' } })

    const wrapper = mount(PlayerConfirmSubscriptionView)
    await flushPromises()

    expect(wrapper.find('[data-testid="confirm-loading"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="confirm-success"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="confirm-error"]').exists()).toBe(false)
  })

  it('shows browse more link in success state', async () => {
    mockApproveSubscription.mockResolvedValue({ data: { status: 'active' } })

    const wrapper = mount(PlayerConfirmSubscriptionView)
    await flushPromises()

    const link = wrapper.find('[data-testid="link-browse-more"]')
    expect(link.exists()).toBe(true)
    expect(link.attributes('href')).toBe('/games')
  })

  it('shows error state when confirmation fails', async () => {
    mockApproveSubscription.mockRejectedValue(new Error('Subscription already confirmed'))

    const wrapper = mount(PlayerConfirmSubscriptionView)
    await flushPromises()

    expect(wrapper.find('[data-testid="confirm-loading"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="confirm-error"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="confirm-success"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="error-message"]').text()).toContain('Subscription already confirmed')
  })

  it('calls approveSubscription with route params and query email', async () => {
    mockApproveSubscription.mockResolvedValue({ data: { status: 'active' } })

    mount(PlayerConfirmSubscriptionView)
    await flushPromises()

    expect(mockApproveSubscription).toHaveBeenCalledWith('sub-abc-123', 'player@example.com')
  })

  it('shows fallback error message when error has no message', async () => {
    mockApproveSubscription.mockRejectedValue({})

    const wrapper = mount(PlayerConfirmSubscriptionView)
    await flushPromises()

    expect(wrapper.find('[data-testid="error-message"]').text()).toContain(
      'Failed to confirm your subscription'
    )
  })
})
