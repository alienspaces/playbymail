import { describe, it, expect, beforeEach } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { useGamesStore } from '../../../stores/games'
import StudioPlacementView from './StudioPlacementView.vue'

describe('StudioPlacementView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders heading and placeholder when a game is selected', () => {
    const gamesStore = useGamesStore()
    gamesStore.selectedGame = { id: 'game1', name: 'Test Game' }
    const wrapper = shallowMount(StudioPlacementView)
    expect(wrapper.text()).toContain('Placement')
    expect(wrapper.text()).toContain('Assign items and creatures to locations here.')
  })

  it('shows prompt if no game is selected', () => {
    const gamesStore = useGamesStore()
    gamesStore.selectedGame = null
    const wrapper = shallowMount(StudioPlacementView)
    expect(wrapper.text()).toContain('Please select or create a game to manage placement.')
  })
}) 