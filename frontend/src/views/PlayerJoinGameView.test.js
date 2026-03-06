import { describe, it, expect, vi, beforeEach } from 'vitest'
import { nextTick } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import PlayerJoinGameView from './PlayerJoinGameView.vue'

const mockGetJoinSheet = vi.fn()
const mockSubmitJoinGame = vi.fn()

vi.mock('../api/joinGame', () => ({
  getJoinSheet: (...args) => mockGetJoinSheet(...args),
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

const mockSheetHtml = '<html><body><form><input name="email" value="test@example.com"></form></body></html>'

describe('PlayerJoinGameView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loading state while fetching turn sheet', async () => {
    mockGetJoinSheet.mockReturnValue(new Promise(() => {}))

    const wrapper = mount(PlayerJoinGameView)
    await nextTick()

    expect(wrapper.find('[data-testid="join-loading"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="join-sheet"]').exists()).toBe(false)
  })

  it('shows turn sheet iframe after successful load', async () => {
    mockGetJoinSheet.mockResolvedValue(mockSheetHtml)

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    expect(wrapper.find('[data-testid="join-loading"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="join-sheet"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="join-sheet-iframe"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="btn-submit"]').exists()).toBe(true)
  })

  it('shows load error when turn sheet fetch fails', async () => {
    mockGetJoinSheet.mockRejectedValue(new Error('Unauthorized'))

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    expect(wrapper.find('[data-testid="join-load-error"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('Unauthorized')
    expect(wrapper.find('[data-testid="join-sheet"]').exists()).toBe(false)
  })

  it('has a back to catalog link', async () => {
    mockGetJoinSheet.mockResolvedValue(mockSheetHtml)

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    const backLink = wrapper.find('[data-testid="btn-back"]')
    expect(backLink.exists()).toBe(true)
    expect(backLink.attributes('href')).toBe('/games')
  })

  it('shows success step after successful submission', async () => {
    mockGetJoinSheet.mockResolvedValue(mockSheetHtml)
    mockSubmitJoinGame.mockResolvedValue({
      data: { game_subscription_id: 'sub-1', game_instance_id: 'inst-1', game_id: 'g1' },
    })

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    // Simulate a submit that bypasses iframe extraction (iframe DOM not accessible in unit tests)
    // by directly testing the success state transition
    wrapper.vm.step = 'success'
    await nextTick()

    expect(wrapper.find('[data-testid="step-success"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="link-browse-more"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="link-browse-more"]').attributes('href')).toBe('/games')
  })

  it('shows submit error when form data cannot be extracted', async () => {
    mockGetJoinSheet.mockResolvedValue(mockSheetHtml)

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    // sheetFrame ref won't have accessible contentDocument in JSDOM
    await wrapper.find('[data-testid="btn-submit"]').trigger('click')
    await flushPromises()

    expect(wrapper.find('[data-testid="submit-error"]').exists()).toBe(true)
  })

  it('calls getJoinSheet with the game subscription id', async () => {
    mockGetJoinSheet.mockResolvedValue(mockSheetHtml)

    mount(PlayerJoinGameView)
    await flushPromises()

    expect(mockGetJoinSheet).toHaveBeenCalledWith('sub-abc-123')
  })

  it('sets iframe srcdoc to the fetched HTML', async () => {
    mockGetJoinSheet.mockResolvedValue(mockSheetHtml)

    const wrapper = mount(PlayerJoinGameView)
    await flushPromises()

    const iframe = wrapper.find('[data-testid="join-sheet-iframe"]')
    expect(iframe.attributes('srcdoc')).toBe(mockSheetHtml)
  })
})
