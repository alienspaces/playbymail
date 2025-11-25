import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import GameContext from './GameContext.vue'

describe('GameContext', () => {
  it('renders game name when provided', () => {
    const wrapper = mount(GameContext, {
      props: {
        gameName: 'Test Game'
      }
    })

    const label = wrapper.find('.game-context-label')
    const name = wrapper.find('.game-context-name')
    expect(label.text()).toBe('Game:')
    expect(name.text()).toBe('Test Game')
  })

  it('does not render when gameName is empty', () => {
    const wrapper = mount(GameContext, {
      props: {
        gameName: ''
      }
    })

    expect(wrapper.find('p').exists()).toBe(false)
  })

  it('does not render when gameName is not provided', () => {
    const wrapper = mount(GameContext)

    expect(wrapper.find('p').exists()).toBe(false)
  })
}) 