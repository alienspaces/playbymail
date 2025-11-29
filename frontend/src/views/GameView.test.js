import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import GameView from './GameView.vue'

// Mock fetch
const mockGames = [
  { id: '1', name: 'Test Game', game_type: 'adventure', created_at: '2024-07-10T12:00:00Z' }
]

window.fetch = vi.fn(() =>
  Promise.resolve({
    ok: true,
    json: () => Promise.resolve({ data: mockGames })
  })
)

describe('GameView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    fetch.mockClear()
  })

  it('renders and fetches games', async () => {
    const wrapper = mount(GameView)
    // Wait for fetchGames to complete
    await new Promise(r => setTimeout(r, 0))
    // Studio filters by Designer subscription type
    expect(fetch).toHaveBeenCalledWith('http://localhost:8080/api/v1/games?subscription_type=Designer', expect.any(Object))
    expect(wrapper.text()).toContain('Games')
    expect(wrapper.text()).toContain('Test Game')
    expect(wrapper.find('th').text()).toBe('Name')
  })

  it('shows no games found if list is empty', async () => {
    fetch.mockImplementationOnce(() =>
      Promise.resolve({ ok: true, json: () => Promise.resolve({ data: [] }) })
    )
    const wrapper = mount(GameView)
    await new Promise(r => setTimeout(r, 0))
    expect(wrapper.text()).toContain('No games found.')
  })
}) 