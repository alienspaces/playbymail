import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ManagementEntryView from './ManagementEntryView.vue'

describe('ManagementEntryView', () => {

  it('renders admin entry content', () => {
    const routerLinkStub = {
      template: '<a><slot /></a>'
    }
    const wrapper = mount(ManagementEntryView, {
      global: {
        stubs: {
          'router-link': routerLinkStub
        }
      }
    })

    // Should show admin entry content
    expect(wrapper.text()).toContain('Game Management')
    expect(wrapper.text()).toContain('Sign in to access game management tools')
    expect(wrapper.text()).toContain('Sign In to Game Management')
  })
}) 