import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { ref } from 'vue'
import StudioItemEffectsView from './StudioItemEffectsView.vue'
import { setupModalTestCleanup } from '../../../test-utils/studio-resource-helpers'

vi.mock('../../../stores/itemEffects', () => ({
  useItemEffectsStore: vi.fn(() => ({
    itemEffects: [],
    loading: false,
    error: null,
    pageNumber: 1,
    hasMore: false,
    fetchItemEffects: vi.fn(),
    createItemEffect: vi.fn(),
    updateItemEffect: vi.fn(),
    deleteItemEffect: vi.fn(),
  })),
}))

vi.mock('../../../stores/items', () => ({
  useItemsStore: vi.fn(() => ({
    items: [],
    loading: false,
    error: null,
    fetchItems: vi.fn(),
  })),
}))

vi.mock('../../../stores/locations', () => ({
  useLocationsStore: vi.fn(() => ({
    locations: [],
    loading: false,
    error: null,
    fetchLocations: vi.fn(),
  })),
}))

vi.mock('../../../stores/locationLinks', () => ({
  useLocationLinksStore: vi.fn(() => ({
    locationLinks: [],
    loading: false,
    error: null,
    fetchLocationLinks: vi.fn(),
  })),
}))

vi.mock('../../../stores/creatures', () => ({
  useCreaturesStore: vi.fn(() => ({
    creatures: [],
    loading: false,
    error: null,
    fetchCreatures: vi.fn(),
  })),
}))

vi.mock('../../../stores/games', () => ({
  useGamesStore: vi.fn(() => ({
    selectedGame: ref(null),
  })),
}))

describe('StudioItemEffectsView', () => {
  let pinia
  const modalCleanup = setupModalTestCleanup()

  const setupStoreMocks = async (selectedGame = null) => {
    const { useGamesStore } = await import('../../../stores/games')
    const { useItemEffectsStore } = await import('../../../stores/itemEffects')
    const { useItemsStore } = await import('../../../stores/items')
    const { useLocationsStore } = await import('../../../stores/locations')
    const { useLocationLinksStore } = await import('../../../stores/locationLinks')
    const { useCreaturesStore } = await import('../../../stores/creatures')

    useGamesStore.mockReturnValue({ selectedGame: ref(selectedGame) })
    useItemEffectsStore.mockReturnValue({
      itemEffects: [],
      loading: false,
      error: null,
      fetchItemEffects: vi.fn(),
      createItemEffect: vi.fn(),
      updateItemEffect: vi.fn(),
      deleteItemEffect: vi.fn(),
    })
    useItemsStore.mockReturnValue({ items: [], loading: false, error: null, fetchItems: vi.fn() })
    useLocationsStore.mockReturnValue({ locations: [], loading: false, error: null, fetchLocations: vi.fn() })
    useLocationLinksStore.mockReturnValue({ locationLinks: [], loading: false, error: null, fetchLocationLinks: vi.fn() })
    useCreaturesStore.mockReturnValue({ creatures: [], loading: false, error: null, fetchCreatures: vi.fn() })
  }

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)
    vi.clearAllMocks()
    modalCleanup.beforeEach()
  })

  afterEach(() => {
    modalCleanup.afterEach()
  })

  it('renders prompt when no game is selected', () => {
    const wrapper = mount(StudioItemEffectsView)
    expect(wrapper.text()).toContain('Please select or create a game to manage item effects.')
    expect(wrapper.find('.game-table-section').exists()).toBe(false)
  })

  it('renders item effects table when game is selected', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioItemEffectsView)
    expect(wrapper.find('.game-table-section').exists()).toBe(true)
    expect(wrapper.find('h2').text()).toBe('Item Effects')
  })

  it('effectFields only shows base fields when no effect_type is set', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioItemEffectsView)

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).toContain('adventure_game_item_id')
    expect(keys).toContain('action_type')
    expect(keys).toContain('effect_type')
    expect(keys).toContain('result_description')
    expect(keys).toContain('is_repeatable')

    // No result fields shown when effect_type is not set
    expect(keys).not.toContain('result_value_min')
    expect(keys).not.toContain('result_value_max')
    expect(keys).not.toContain('result_adventure_game_item_id')
    expect(keys).not.toContain('result_adventure_game_location_id')
  })

  it('effectFields shows result_value_min and result_value_max as required for heal_wielder', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioItemEffectsView)

    wrapper.vm.modalForm.effect_type = 'heal_wielder'
    await wrapper.vm.$nextTick()

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).toContain('result_value_min')
    expect(keys).toContain('result_value_max')

    const minField = fields.find((f) => f.key === 'result_value_min')
    const maxField = fields.find((f) => f.key === 'result_value_max')
    expect(minField.required).toBe(true)
    expect(maxField.required).toBe(true)

    // Other result fields should NOT be shown
    expect(keys).not.toContain('result_adventure_game_item_id')
    expect(keys).not.toContain('result_adventure_game_location_id')
  })

  it('effectFields shows result_value_min and result_value_max as required for weapon_damage', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioItemEffectsView)

    wrapper.vm.modalForm.effect_type = 'weapon_damage'
    await wrapper.vm.$nextTick()

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).toContain('result_value_min')
    expect(keys).toContain('result_value_max')

    const minField = fields.find((f) => f.key === 'result_value_min')
    expect(minField.required).toBe(true)
  })

  it('effectFields shows result_value_min and result_value_max as required for armor_defense', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioItemEffectsView)

    wrapper.vm.modalForm.effect_type = 'armor_defense'
    await wrapper.vm.$nextTick()

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).toContain('result_value_min')
    expect(keys).toContain('result_value_max')
    expect(keys).not.toContain('result_adventure_game_location_id')
  })

  it('effectFields shows result_adventure_game_location_id as required for teleport', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioItemEffectsView)

    wrapper.vm.modalForm.effect_type = 'teleport'
    await wrapper.vm.$nextTick()

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).toContain('result_adventure_game_location_id')
    expect(fields.find((f) => f.key === 'result_adventure_game_location_id').required).toBe(true)
    expect(keys).not.toContain('result_value_min')
    expect(keys).not.toContain('result_adventure_game_item_id')
  })

  it('effectFields shows result_adventure_game_item_id as required for give_item', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioItemEffectsView)

    wrapper.vm.modalForm.effect_type = 'give_item'
    await wrapper.vm.$nextTick()

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).toContain('result_adventure_game_item_id')
    expect(fields.find((f) => f.key === 'result_adventure_game_item_id').required).toBe(true)
    expect(keys).not.toContain('result_value_min')
  })

  it('effectFields shows result_adventure_game_location_link_id as required for open_link', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioItemEffectsView)

    wrapper.vm.modalForm.effect_type = 'open_link'
    await wrapper.vm.$nextTick()

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).toContain('result_adventure_game_location_link_id')
    expect(fields.find((f) => f.key === 'result_adventure_game_location_link_id').required).toBe(true)
    expect(keys).not.toContain('result_value_min')
  })

  it('effectFields shows result_adventure_game_creature_id as required for summon_creature', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioItemEffectsView)

    wrapper.vm.modalForm.effect_type = 'summon_creature'
    await wrapper.vm.$nextTick()

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).toContain('result_adventure_game_creature_id')
    expect(fields.find((f) => f.key === 'result_adventure_game_creature_id').required).toBe(true)
    expect(keys).not.toContain('result_value_min')
  })

  it('effectFields shows no result fields for info effect_type', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioItemEffectsView)

    wrapper.vm.modalForm.effect_type = 'info'
    await wrapper.vm.$nextTick()

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).not.toContain('result_value_min')
    expect(keys).not.toContain('result_adventure_game_item_id')
    expect(keys).not.toContain('result_adventure_game_location_id')
  })

  it('effectFields shows no result fields for nothing effect_type', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioItemEffectsView)

    wrapper.vm.modalForm.effect_type = 'nothing'
    await wrapper.vm.$nextTick()

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).not.toContain('result_value_min')
    expect(keys).not.toContain('result_adventure_game_item_id')
  })

  it('effectFields updates reactively when effect_type changes', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioItemEffectsView)

    wrapper.vm.modalForm.effect_type = 'heal_wielder'
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.effectFields.map((f) => f.key)).toContain('result_value_min')

    wrapper.vm.modalForm.effect_type = 'teleport'
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.effectFields.map((f) => f.key)).toContain('result_adventure_game_location_id')
    expect(wrapper.vm.effectFields.map((f) => f.key)).not.toContain('result_value_min')
  })
})
