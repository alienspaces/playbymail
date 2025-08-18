import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ComingSoonBanner from './ComingSoonBanner.vue'

describe('ComingSoonBanner', () => {
  it('renders correctly', () => {
    const wrapper = mount(ComingSoonBanner)
    
    expect(wrapper.find('.coming-soon-banner').exists()).toBe(true)
    expect(wrapper.find('.banner-content').exists()).toBe(true)
    expect(wrapper.find('.banner-text').text()).toBe('Coming Soon')
    expect(wrapper.find('.banner-date').text()).toBe('Late 2025')
  })

  it('has correct CSS classes', () => {
    const wrapper = mount(ComingSoonBanner)
    
    expect(wrapper.classes()).toContain('coming-soon-banner')
    expect(wrapper.find('.banner-content').classes()).toContain('banner-content')
    expect(wrapper.find('.banner-text').classes()).toContain('banner-text')
    expect(wrapper.find('.banner-date').classes()).toContain('banner-date')
  })

  it('displays banner text correctly', () => {
    const wrapper = mount(ComingSoonBanner)
    
    const bannerText = wrapper.find('.banner-text')
    const bannerDate = wrapper.find('.banner-date')
    
    expect(bannerText.text()).toBe('Coming Soon')
    expect(bannerDate.text()).toBe('Late 2025')
  })
})
