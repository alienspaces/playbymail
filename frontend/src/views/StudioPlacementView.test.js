import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StudioPlacementView from './StudioPlacementView.vue'

describe('StudioPlacementView', () => {
  it('renders heading and placeholder', () => {
    const wrapper = mount(StudioPlacementView)
    expect(wrapper.text()).toContain('Placement')
    expect(wrapper.text()).toContain('Assign items and creatures to locations here.')
  })
}) 