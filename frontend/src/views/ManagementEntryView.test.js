import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { useAuthStore } from '../stores/auth'
import ManagementEntryView from './ManagementEntryView.vue'

describe('ManagementEntryView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders admin entry content when user is not authenticated', () => {
    const authStore = useAuthStore()
    authStore.sessionToken = null
    
    const routerLinkStub = {
      template: '<a><slot /></a>'
    }
    const wrapper = mount(ManagementEntryView, {
      global: {
        stubs: {
          'router-link': routerLinkStub,
          'router-view': true
        }
      }
    })
    
    // Should show admin entry content
    expect(wrapper.text()).toContain('Game Management')
    expect(wrapper.text()).toContain('Sign in to access game management tools')
    expect(wrapper.text()).toContain('Sign In to Game Management')
    
    // Should not show admin help section
    expect(wrapper.find('.entry-help').exists()).toBe(false)
  })

  it('renders admin help and router-view when user is authenticated', () => {
    const authStore = useAuthStore()
    authStore.sessionToken = 'test-token'
    
    const wrapper = mount(ManagementEntryView, {
      global: {
        stubs: {
          'router-link': true,
          'router-view': true
        }
      }
    })
    
    // Should show admin help section
    expect(wrapper.find('.entry-help').exists()).toBe(true)
    expect(wrapper.text()).toContain('Game Management Help')
    
    // Should show router-view for admin content
    expect(wrapper.find('router-view-stub').exists()).toBe(true)
    
    // Should not show sign-in button
    expect(wrapper.text()).not.toContain('Sign In to Game Management')
  })
}) 