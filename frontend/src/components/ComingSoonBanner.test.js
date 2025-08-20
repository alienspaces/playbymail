import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ComingSoonBanner from './ComingSoonBanner.vue'

describe('ComingSoonBanner', () => {
  it('renders correctly', () => {
    const wrapper = mount(ComingSoonBanner)
    
    expect(wrapper.find('.coming-soon-banner').exists()).toBe(true)
    expect(wrapper.find('.envelope').exists()).toBe(true)
    expect(wrapper.find('.stamp').exists()).toBe(true)
    expect(wrapper.find('.address').exists()).toBe(true)
  })

  it('has correct CSS classes', () => {
    const wrapper = mount(ComingSoonBanner)
    
    expect(wrapper.classes()).toContain('coming-soon-banner')
    expect(wrapper.find('.envelope').classes()).toContain('envelope')
    expect(wrapper.find('.stamp').classes()).toContain('stamp')
    expect(wrapper.find('.address').classes()).toContain('address')
  })

  it('displays envelope content correctly', () => {
    const wrapper = mount(ComingSoonBanner)
    
    const addressLines = wrapper.findAll('.address-line')
    expect(addressLines).toHaveLength(2)
    expect(addressLines[0].text()).toBe('Coming Soon')
    expect(addressLines[1].text()).toBe('Late 2025')
    
    const stampText = wrapper.find('.stamp-text')
    expect(stampText.text()).toBe('2025')
  })
})
