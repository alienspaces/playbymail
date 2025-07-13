import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StudioLayout from './StudioLayout.vue'

describe('StudioLayout', () => {
  it('renders sidebar navigation links', () => {
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
    const navLinks = wrapper.findAll('li')
    const linkTexts = navLinks.map(li => li.text())
    expect(linkTexts).toContain('Games')
    expect(linkTexts).toContain('Locations')
    expect(linkTexts).toContain('Items')
    expect(linkTexts).toContain('Creatures')
    expect(linkTexts).toContain('Placement')
  })
}) 