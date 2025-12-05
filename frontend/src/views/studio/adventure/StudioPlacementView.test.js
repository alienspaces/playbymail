import { describe, it, expect, beforeEach } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import StudioPlacementView from './StudioPlacementView.vue'
import { setupGamesStore } from '../../../test-utils/studio-resource-helpers'

describe('StudioPlacementView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders heading and placeholder when a game is selected', async () => {
    await setupGamesStore({ id: 'game1', name: 'Test Game' })
    const wrapper = shallowMount(StudioPlacementView)
    expect(wrapper.text()).toContain('Placement')
    expect(wrapper.text()).toContain('Item Placements')
    expect(wrapper.text()).toContain('Creature Placements')
  })

  it('shows prompt if no game is selected', async () => {
    await setupGamesStore(null)
    const wrapper = shallowMount(StudioPlacementView)
    expect(wrapper.text()).toContain('Please select or create a game to manage placement.')
  })
}) 