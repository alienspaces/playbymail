import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { useGamesStore } from '../stores/games'
import StudioLayout from './StudioLayout.vue'
import * as vueRouter from 'vue-router'

// Mock useRoute
vi.mock('vue-router', async () => {
  const actual = await vi.importActual('vue-router')
  return {
    ...actual,
    useRoute: vi.fn()
  }
})

describe('StudioLayout', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders basic sidebar navigation links', () => {
    vueRouter.useRoute.mockReturnValue({ path: '/studio' })
    
    const gamesStore = useGamesStore()
    gamesStore.selectedGame = null
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
    expect(linkTexts).toContain('Games')
    expect(linkTexts).not.toContain('Locations')
    expect(linkTexts).not.toContain('Items')
    expect(linkTexts).not.toContain('Creatures')
    expect(linkTexts).not.toContain('Placement')
  })

  it('renders adventure menu when on adventure game route', () => {
    vueRouter.useRoute.mockReturnValue({ path: '/studio/game1/locations' })
    
    const gamesStore = useGamesStore()
    gamesStore.selectedGame = { id: 'game1', name: 'Test Game', game_type: 'adventure' }
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
    expect(linkTexts).toContain('Games')
    expect(linkTexts).toContain('Locations')
    expect(linkTexts).toContain('Items')
    expect(linkTexts).toContain('Creatures')
    expect(linkTexts).toContain('Placement')
  })
}) 