import { describe, it, expect, vi, beforeEach } from 'vitest'
import { nextTick } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import GameCatalogView from './GameCatalogView.vue'

const mockListCatalogGameInstances = vi.fn()

vi.mock('../api/catalog', () => ({
  listCatalogGameInstances: (...args) => mockListCatalogGameInstances(...args),
}))

const mockCatalogData = [
  {
    game_instance_id: 'inst-1',
    game_id: 'g1',
    game_subscription_id: 'sub-1',
    game_name: 'The Lost Kingdom',
    game_description: 'An adventure game',
    game_type: 'adventure',
    account_name: 'Test Host',
    turn_duration_hours: 168,
    required_player_count: 4,
    player_count: 0,
    remaining_capacity: 4,
    delivery_email: true,
    delivery_physical_local: false,
    delivery_physical_post: false,
    is_closed_testing: false,
    created_at: '2026-01-01T00:00:00Z',
  },
]

describe('GameCatalogView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loading state while fetching', async () => {
    mockListCatalogGameInstances.mockReturnValue(new Promise(() => {}))

    const wrapper = mount(GameCatalogView)
    await nextTick()

    expect(wrapper.find('[data-testid="catalog-loading"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="catalog-games"]').exists()).toBe(false)
  })

  it('renders instance cards after successful fetch', async () => {
    mockListCatalogGameInstances.mockResolvedValue({ data: mockCatalogData })

    const wrapper = mount(GameCatalogView)
    await flushPromises()

    expect(wrapper.find('[data-testid="catalog-loading"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="catalog-games"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="instance-card-inst-1"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('The Lost Kingdom')
    expect(wrapper.text()).toContain('An adventure game')
    expect(wrapper.text()).toContain('Hosted by Test Host')
    expect(wrapper.text()).toContain('Adventure')
  })

  it('renders join button with subscription ID', async () => {
    mockListCatalogGameInstances.mockResolvedValue({ data: mockCatalogData })

    const wrapper = mount(GameCatalogView)
    await flushPromises()

    const joinButton = wrapper.find('[data-testid="join-button-sub-1"]')
    expect(joinButton.exists()).toBe(true)
    expect(joinButton.attributes('href')).toBe('/player/join-game/sub-1')
    expect(wrapper.text()).toContain('4 players needed')
    expect(wrapper.text()).toContain('Email')
  })

  it('shows empty state when catalog has no instances', async () => {
    mockListCatalogGameInstances.mockResolvedValue({ data: [] })

    const wrapper = mount(GameCatalogView)
    await flushPromises()

    expect(wrapper.find('[data-testid="catalog-empty"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="catalog-games"]').exists()).toBe(false)
    expect(wrapper.text()).toContain('No games are currently available for enrollment')
  })

  it('shows error state when fetch fails', async () => {
    mockListCatalogGameInstances.mockRejectedValue(new Error('Network error'))

    const wrapper = mount(GameCatalogView)
    await flushPromises()

    expect(wrapper.find('[data-testid="catalog-error"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="catalog-games"]').exists()).toBe(false)
    expect(wrapper.text()).toContain('Network error')
  })

  it('retries fetch when try again button is clicked', async () => {
    mockListCatalogGameInstances.mockRejectedValueOnce(new Error('Network error'))
    mockListCatalogGameInstances.mockResolvedValueOnce({ data: mockCatalogData })

    const wrapper = mount(GameCatalogView)
    await flushPromises()

    expect(wrapper.find('[data-testid="catalog-error"]').exists()).toBe(true)

    await wrapper.find('.retry-button').trigger('click')
    await flushPromises()

    expect(wrapper.find('[data-testid="catalog-games"]').exists()).toBe(true)
    expect(mockListCatalogGameInstances).toHaveBeenCalledTimes(2)
  })
})
