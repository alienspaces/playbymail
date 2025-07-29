import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BuildInfo from './BuildInfo.vue'

describe('BuildInfo', () => {
  it('renders commit ref and build date', () => {
    const wrapper = mount(BuildInfo)
    
    // Should contain commit info (either dev or actual commit)
    const text = wrapper.text()
    expect(text).toMatch(/(dev|[a-f0-9]{7})/) // dev or 7-char commit hash
    // Should contain a date (any month)
    expect(wrapper.text()).toMatch(/\d{4}/) // year
  })

  it('displays build information in panel', () => {
    const wrapper = mount(BuildInfo)
    
    // Check that the component renders the expected structure
    expect(wrapper.find('.build-info-panel').exists()).toBe(true)
    expect(wrapper.find('.build-info-content').exists()).toBe(true)
    expect(wrapper.find('.build-commit').exists()).toBe(true)
    expect(wrapper.find('.build-date').exists()).toBe(true)
  })

  it('handles date formatting gracefully', () => {
    const wrapper = mount(BuildInfo)
    
    // Should contain a formatted date string
    const buildDateText = wrapper.find('.build-date').text()
    // Should contain a year
    expect(buildDateText).toMatch(/\d{4}/)
    
    // Should not contain "Unknown" unless there's an error
    expect(buildDateText).not.toContain('Unknown')
  })
})