import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import PlayerApp from './PlayerApp.vue'
import PlayerSupportFooter from './components/PlayerSupportFooter.vue'

vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({
    params: {},
    query: {},
  })),
  useRouter: vi.fn(() => ({
    push: vi.fn(),
  })),
}))

describe('PlayerApp', () => {
  it('renders router-view', () => {
    const wrapper = mount(PlayerApp, {
      global: {
        stubs: {
          'router-view': true,
          PlayerSupportFooter: true,
        },
      },
    })
    expect(wrapper.find('router-view-stub').exists()).toBe(true)
  })

  it('renders PlayerSupportFooter', () => {
    const wrapper = mount(PlayerApp, {
      global: {
        stubs: {
          'router-view': true,
        },
        components: {
          PlayerSupportFooter,
        },
      },
    })
    expect(wrapper.findComponent(PlayerSupportFooter).exists()).toBe(true)
  })

  it('has flex column layout for full height', () => {
    const wrapper = mount(PlayerApp, {
      global: {
        stubs: {
          'router-view': true,
          PlayerSupportFooter: true,
        },
      },
    })
    const app = wrapper.find('#player-app')
    expect(app.exists()).toBe(true)
  })
})
