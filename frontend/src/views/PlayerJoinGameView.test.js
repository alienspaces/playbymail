import { describe, it, expect, vi, beforeEach } from 'vitest'
import { nextTick } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import PlayerJoinGameView from './PlayerJoinGameView.vue'

const mockGetJoinGameInfo = vi.fn()
const mockSubmitJoinGame = vi.fn()

vi.mock('../api/joinGame', () => ({
  getJoinGameInfo: (...args) => mockGetJoinGameInfo(...args),
  submitJoinGame: (...args) => mockSubmitJoinGame(...args),
}))

vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({
    params: { game_subscription_id: 'sub-abc-123' },
  })),
  useRouter: vi.fn(() => ({
    push: vi.fn(),
  })),
}))

const mockGameInfo = {
  game_subscription_id: 'sub-abc-123',
  game_name: 'The Lost Kingdom',
  game_description: 'An exciting adventure game',
  game_type: 'adventure',
  turn_duration_hours: 168,
  total_capacity: 4,
  total_players: 1,
  delivery_email: true,
  delivery_physical_local: false,
  delivery_physical_post: false,
}

describe('PlayerJoinGameView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loading state while fetching game info', async () => {
    mockGetJoinGameInfo.mockReturnValue(new Promise(() => {}))

    const wrapper = mount(PlayerJoinGameView)
    await nextTick()

    expect(wrapper.find('[data-testid="join-loading"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="step-info"]').exists()).toBe(false)
  })

  it('shows game info step after successful load', async () => {
    mockGetJoinGameInfo.mockResolvedValue({ data: mockGameInfo })

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    expect(wrapper.find('[data-testid="join-loading"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="step-info"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('The Lost Kingdom')
    expect(wrapper.text()).toContain('An exciting adventure game')
    expect(wrapper.text()).toContain('Adventure')
    expect(wrapper.find('[data-testid="btn-join"]').exists()).toBe(true)
  })

  it('shows load error when game info fetch fails', async () => {
    mockGetJoinGameInfo.mockRejectedValue(new Error('Game not found'))

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    expect(wrapper.find('[data-testid="join-load-error"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('Game not found')
    expect(wrapper.find('[data-testid="step-info"]').exists()).toBe(false)
  })

  it('advances to contact step when Join this Game is clicked', async () => {
    mockGetJoinGameInfo.mockResolvedValue({ data: mockGameInfo })

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    await wrapper.find('[data-testid="btn-join"]').trigger('click')
    await nextTick()

    expect(wrapper.find('[data-testid="step-contact"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="step-info"]').exists()).toBe(false)
  })

  it('returns to info step when Back is clicked on contact step', async () => {
    mockGetJoinGameInfo.mockResolvedValue({ data: mockGameInfo })

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    await wrapper.find('[data-testid="btn-join"]').trigger('click')
    await nextTick()

    await wrapper.find('[data-testid="btn-back"]').trigger('click')
    await nextTick()

    expect(wrapper.find('[data-testid="step-info"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="step-contact"]').exists()).toBe(false)
  })

  it('hides delivery selection when only one delivery method is available', async () => {
    mockGetJoinGameInfo.mockResolvedValue({ data: mockGameInfo })

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    await wrapper.find('[data-testid="btn-join"]').trigger('click')
    await nextTick()

    expect(wrapper.find('[data-testid="delivery-selection"]').exists()).toBe(false)
  })

  it('shows delivery selection when multiple delivery methods are available', async () => {
    const multiDeliveryInfo = {
      ...mockGameInfo,
      delivery_email: true,
      delivery_physical_post: true,
    }
    mockGetJoinGameInfo.mockResolvedValue({ data: multiDeliveryInfo })

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    await wrapper.find('[data-testid="btn-join"]').trigger('click')
    await nextTick()

    expect(wrapper.find('[data-testid="delivery-selection"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="delivery-email"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="delivery-post"]').exists()).toBe(true)
  })

  it('shows all required contact fields on contact step', async () => {
    mockGetJoinGameInfo.mockResolvedValue({ data: mockGameInfo })

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    await wrapper.find('[data-testid="btn-join"]').trigger('click')
    await nextTick()

    expect(wrapper.find('[data-testid="input-email"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="input-name"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="input-address-line1"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="input-address-line2"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="input-state"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="input-postal-code"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="input-country"]').exists()).toBe(true)
  })

  it('submits form and shows success step on success', async () => {
    mockGetJoinGameInfo.mockResolvedValue({ data: mockGameInfo })
    mockSubmitJoinGame.mockResolvedValue({
      data: { game_subscription_id: 'sub-1', game_instance_id: 'inst-1', game_id: 'g1' },
    })

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    await wrapper.find('[data-testid="btn-join"]').trigger('click')
    await nextTick()

    await wrapper.find('[data-testid="input-email"]').setValue('player@example.com')
    await wrapper.find('[data-testid="input-name"]').setValue('Test Player')
    await wrapper.find('[data-testid="input-address-line1"]').setValue('123 Main St')
    await wrapper.find('[data-testid="input-state"]').setValue('VIC')
    await wrapper.find('[data-testid="input-postal-code"]').setValue('3000')
    await wrapper.find('[data-testid="input-country"]').setValue('Australia')

    await wrapper.find('.join-form').trigger('submit')
    await flushPromises()

    expect(wrapper.find('[data-testid="step-success"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('The Lost Kingdom')
    expect(wrapper.find('[data-testid="link-browse-more"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="link-browse-more"]').attributes('href')).toBe('/games')
  })

  it('shows error message when submission fails', async () => {
    mockGetJoinGameInfo.mockResolvedValue({ data: mockGameInfo })
    mockSubmitJoinGame.mockRejectedValue(new Error('Instance full'))

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    await wrapper.find('[data-testid="btn-join"]').trigger('click')
    await nextTick()

    await wrapper.find('.join-form').trigger('submit')
    await flushPromises()

    expect(wrapper.find('[data-testid="submit-error"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('Instance full')
    expect(wrapper.find('[data-testid="step-contact"]').exists()).toBe(true)
  })

  it('calls submitJoinGame with correct payload including delivery flags', async () => {
    mockGetJoinGameInfo.mockResolvedValue({ data: mockGameInfo })
    mockSubmitJoinGame.mockResolvedValue({
      data: { game_subscription_id: 'sub-1', game_instance_id: 'inst-1', game_id: 'g1' },
    })

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    await wrapper.find('[data-testid="btn-join"]').trigger('click')
    await nextTick()

    await wrapper.find('[data-testid="input-email"]').setValue('player@example.com')
    await wrapper.find('[data-testid="input-name"]').setValue('Test Player')
    await wrapper.find('[data-testid="input-address-line1"]').setValue('123 Main St')
    await wrapper.find('[data-testid="input-state"]').setValue('VIC')
    await wrapper.find('[data-testid="input-postal-code"]').setValue('3000')
    await wrapper.find('[data-testid="input-country"]').setValue('Australia')

    await wrapper.find('.join-form').trigger('submit')
    await flushPromises()

    expect(mockSubmitJoinGame).toHaveBeenCalledWith(
      'sub-abc-123',
      expect.objectContaining({
        email: 'player@example.com',
        name: 'Test Player',
        postal_address_line1: '123 Main St',
        state_province: 'VIC',
        postal_code: '3000',
        country: 'Australia',
        delivery_email: true,
        delivery_physical_local: false,
        delivery_physical_post: false,
      })
    )
  })
})
