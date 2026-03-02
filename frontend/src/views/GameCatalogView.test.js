import { describe, it, expect, vi, beforeEach } from 'vitest'
import { nextTick } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import GameCatalogView from './GameCatalogView.vue'

const mockListCatalogGames = vi.fn()

vi.mock('../api/catalog', () => ({
  listCatalogGames: (...args) => mockListCatalogGames(...args),
}))

const mockCatalogData = [
  {
    game_subscription_id: 'sub-1',
    game_name: 'The Lost Kingdom',
    game_description: 'An adventure game',
    game_type: 'adventure',
    turn_duration_hours: 168,
    total_capacity: 4,
    total_players: 1,
    delivery_email: true,
    delivery_physical_local: false,
    delivery_physical_post: false,
  },
]

describe('GameCatalogView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loading state while fetching', async () => {
    mockListCatalogGames.mockReturnValue(new Promise(() => {}))

    const wrapper = mount(GameCatalogView)
    await nextTick()

    expect(wrapper.find('[data-testid="catalog-loading"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="catalog-games"]').exists()).toBe(false)
  })

  it('renders subscription cards after successful fetch', async () => {
    mockListCatalogGames.mockResolvedValue({ data: mockCatalogData })

    const wrapper = mount(GameCatalogView)
    await flushPromises()

    expect(wrapper.find('[data-testid="catalog-loading"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="catalog-games"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="sub-card-sub-1"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('The Lost Kingdom')
    expect(wrapper.text()).toContain('An adventure game')
    expect(wrapper.text()).toContain('Adventure')
  })

  it('renders join button with subscription ID', async () => {
    mockListCatalogGames.mockResolvedValue({ data: mockCatalogData })

    const wrapper = mount(GameCatalogView)
    await flushPromises()

    const joinButton = wrapper.find('[data-testid="join-button-sub-1"]')
    expect(joinButton.exists()).toBe(true)
    expect(joinButton.attributes('href')).toBe('/player/join-game/sub-1')
    expect(wrapper.text()).toContain('1 / 4 players')
    expect(wrapper.text()).toContain('Email')
  })

  it('shows empty state when catalog has no subscriptions', async () => {
    mockListCatalogGames.mockResolvedValue({ data: [] })

    const wrapper = mount(GameCatalogView)
    await flushPromises()

    expect(wrapper.find('[data-testid="catalog-empty"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="catalog-games"]').exists()).toBe(false)
    expect(wrapper.text()).toContain('No games are currently available for enrollment')
  })

  it('shows error state when fetch fails', async () => {
    mockListCatalogGames.mockRejectedValue(new Error('Network error'))

    const wrapper = mount(GameCatalogView)
    await flushPromises()

    expect(wrapper.find('[data-testid="catalog-error"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="catalog-games"]').exists()).toBe(false)
    expect(wrapper.text()).toContain('Network error')
  })

  it('retries fetch when try again button is clicked', async () => {
    mockListCatalogGames.mockRejectedValueOnce(new Error('Network error'))
    mockListCatalogGames.mockResolvedValueOnce({ data: mockCatalogData })

    const wrapper = mount(GameCatalogView)
    await flushPromises()

    expect(wrapper.find('[data-testid="catalog-error"]').exists()).toBe(true)

    await wrapper.find('.retry-button').trigger('click')
    await flushPromises()

    expect(wrapper.find('[data-testid="catalog-games"]').exists()).toBe(true)
    expect(mockListCatalogGames).toHaveBeenCalledTimes(2)
  })
})
