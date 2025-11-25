import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ComingSoonBanner from './ComingSoonBanner.vue'

describe('ComingSoonBanner', () => {
  it('renders correctly', () => {
    const wrapper = mount(ComingSoonBanner)

    expect(wrapper.find('.coming-soon-banner').exists()).toBe(true)
    expect(wrapper.find('.envelope-icon').exists()).toBe(true)
    expect(wrapper.find('.coming-soon-title').exists()).toBe(true)
    expect(wrapper.find('.coming-soon-date').exists()).toBe(true)
  })

  it('has correct CSS classes', () => {
    const wrapper = mount(ComingSoonBanner)

    expect(wrapper.classes()).toContain('coming-soon-banner')
    expect(wrapper.find('.envelope-icon').classes()).toContain('envelope-icon')
    expect(wrapper.find('.coming-soon-title').classes()).toContain('coming-soon-title')
    expect(wrapper.find('.coming-soon-date').classes()).toContain('coming-soon-date')
  })

  it('displays envelope content correctly', () => {
    const wrapper = mount(ComingSoonBanner)

    expect(wrapper.find('.coming-soon-title').text()).toBe('Coming Soon')
    expect(wrapper.find('.coming-soon-date').text()).toBe('Early 2026')
    expect(wrapper.find('img.envelope-icon').attributes('src')).toBeTruthy()
  })
})
