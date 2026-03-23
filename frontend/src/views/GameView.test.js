import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import GameView from './GameView.vue'

const mockGames = [
  { id: '1', name: 'Test Game', game_type: 'adventure', created_at: '2024-07-10T12:00:00Z' }
]

const mockSubscriptions = [
  { subscription_type: 'basic_game_designer', status: 'active' }
]

const mockFetch = (url) => {
  if (url && url.includes('/account/subscriptions')) {
    return Promise.resolve({ ok: true, json: () => Promise.resolve({ data: mockSubscriptions }) })
  }
  return Promise.resolve({ ok: true, json: () => Promise.resolve({ data: mockGames }) })
}

// Default fetch mock: returns games for game API calls, subscriptions for subscription calls
window.fetch = vi.fn(mockFetch)

describe('GameView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    fetch.mockClear()
    fetch.mockImplementation(mockFetch)
  })

  it('renders and fetches games', async () => {
    const wrapper = mount(GameView)
    await new Promise(r => setTimeout(r, 0))
    expect(fetch).toHaveBeenCalledWith('http://localhost:8080/api/v1/games?is_designer=true', expect.any(Object))
    expect(wrapper.text()).toContain('Games')
    expect(wrapper.text()).toContain('Test Game')
    expect(wrapper.find('th').text()).toBe('Name')
  })

  it('shows no games found if list is empty', async () => {
    fetch.mockImplementation((url) => {
      if (url && url.includes('/account/subscriptions')) {
        return Promise.resolve({ ok: true, json: () => Promise.resolve({ data: mockSubscriptions }) })
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve({ data: [] }) })
    })
    const wrapper = mount(GameView)
    await new Promise(r => setTimeout(r, 0))
    expect(wrapper.text()).toContain('No games found.')
  })

  it('opens create modal with no game type pre-selected', async () => {
    const wrapper = mount(GameView, {
      global: { mocks: { $router: { push: vi.fn() } } }
    })
    await new Promise(r => setTimeout(r, 0))

    await wrapper.vm.openCreate()
    await new Promise(r => setTimeout(r, 0))

    expect(wrapper.vm.showModal).toBe(true)
    expect(wrapper.vm.modalMode).toBe('create')
    expect(wrapper.vm.modalForm.game_type).toBe('')
  })

  it('game type field includes adventure and mecha options', () => {
    const wrapper = mount(GameView, {
      global: { mocks: { $router: { push: vi.fn() } } }
    })

    const gameTypeField = wrapper.vm.gameFields.find(f => f.key === 'game_type')
    expect(gameTypeField).toBeDefined()

    const optionValues = gameTypeField.options.map(o => o.value)
    expect(optionValues).toContain('adventure')
    expect(optionValues).toContain('mecha')
  })

  it('creates a mecha game and closes the modal', async () => {
    const createdGame = { id: '99', name: 'Battle for the Planet', game_type: 'mecha' }
    fetch.mockImplementation((url, opts) => {
      if (url && url.includes('/account/subscriptions')) {
        return Promise.resolve({ ok: true, json: () => Promise.resolve({ data: mockSubscriptions }) })
      }
      if (opts && opts.method === 'POST') {
        return Promise.resolve({ ok: true, json: () => Promise.resolve({ data: createdGame }) })
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve({ data: mockGames }) })
    })

    const mockPush = vi.fn()
    const wrapper = mount(GameView, {
      global: { mocks: { $router: { push: mockPush } } }
    })
    await new Promise(r => setTimeout(r, 0))

    await wrapper.vm.openCreate()
    await new Promise(r => setTimeout(r, 0))

    await wrapper.vm.handleSubmit({
      name: 'Battle for the Planet',
      game_type: 'mecha',
      turn_duration_hours: 168,
      description: 'Mech warriors battling for a planet',
    })
    await new Promise(r => setTimeout(r, 0))

    expect(wrapper.vm.showModal).toBe(false)
    expect(mockPush).toHaveBeenCalledWith('/studio/99/turn-sheet-backgrounds')
  })
}) 