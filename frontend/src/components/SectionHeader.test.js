import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import SectionHeader from './SectionHeader.vue'

describe('SectionHeader', () => {
  it('renders title correctly', () => {
    const wrapper = mount(SectionHeader, {
      props: {
        title: 'Locations',
        resourceName: 'Location'
      }
    })
    
    expect(wrapper.find('h2').text()).toBe('Locations')
  })

  it('renders create button with correct text', () => {
    const wrapper = mount(SectionHeader, {
      props: {
        title: 'Locations',
        resourceName: 'Location'
      }
    })
    
    const button = wrapper.find('button')
    expect(button.exists()).toBe(true)
    expect(button.text()).toBe('Create New Location')
  })

  it('emits create event when button is clicked', async () => {
    const wrapper = mount(SectionHeader, {
      props: {
        title: 'Locations',
        resourceName: 'Location'
      }
    })
    
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('create')).toBeTruthy()
  })
}) 