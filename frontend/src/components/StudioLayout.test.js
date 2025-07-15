import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { useGamesStore } from '../stores/games'
import StudioLayout from './StudioLayout.vue'

describe('StudioLayout', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders sidebar navigation links', () => {
    const gamesStore = useGamesStore()
    gamesStore.selectedGame = { id: 'game1', name: 'Test Game' }
    const routerLinkStub = {
      template: '<a><slot /></a>'
    }
    const wrapper = mount(StudioLayout, {
      global: {
        stubs: {
          'router-link': routerLinkStub,
          'router-view': true
        }
      }
    })
    const navLinks = wrapper.findAll('.sidebar a')
    const linkTexts = navLinks.map(link => link.text())
    expect(linkTexts).toContain('Locations')
    expect(linkTexts).toContain('Items')
    expect(linkTexts).toContain('Creatures')
    expect(linkTexts).toContain('Placement')
  })
}) 