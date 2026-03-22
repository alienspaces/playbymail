import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { ref } from 'vue'
import StudioLocationObjectEffectsView from './StudioLocationObjectEffectsView.vue'
import { findInBody, setupModalTestCleanup } from '../../../test-utils/studio-resource-helpers'

vi.mock('../../../stores/locations', () => ({
  useLocationsStore: vi.fn(() => ({
    locations: [],
    loading: false,
    error: null,
    fetchLocations: vi.fn(),
  })),
}))

vi.mock('../../../stores/locationObjects', () => ({
  useLocationObjectsStore: vi.fn(() => ({
    locationObjects: [],
    loading: false,
    error: null,
    fetchLocationObjects: vi.fn(),
  })),
}))

vi.mock('../../../stores/locationObjectEffects', () => ({
  useLocationObjectEffectsStore: vi.fn(() => ({
    locationObjectEffects: [],
    loading: false,
    error: null,
    pageNumber: 1,
    hasMore: false,
    fetchLocationObjectEffects: vi.fn(),
    createLocationObjectEffect: vi.fn(),
    updateLocationObjectEffect: vi.fn(),
    deleteLocationObjectEffect: vi.fn(),
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

vi.mock('../../../stores/creatures', () => ({
  useCreaturesStore: vi.fn(() => ({
    creatures: [],
    loading: false,
    error: null,
    fetchCreatures: vi.fn(),
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

vi.mock('../../../stores/games', () => ({
  useGamesStore: vi.fn(() => ({
    selectedGame: ref(null),
  })),
}))

describe('StudioLocationObjectEffectsView', () => {
  let pinia
  const modalCleanup = setupModalTestCleanup()

  const setupStoreMocks = async (selectedGame = null) => {
    const { useGamesStore } = await import('../../../stores/games')
    const { useLocationsStore } = await import('../../../stores/locations')
    const { useLocationObjectsStore } = await import('../../../stores/locationObjects')
    const { useLocationObjectEffectsStore } = await import('../../../stores/locationObjectEffects')
    const { useItemsStore } = await import('../../../stores/items')
    const { useCreaturesStore } = await import('../../../stores/creatures')
    const { useLocationLinksStore } = await import('../../../stores/locationLinks')

    useGamesStore.mockReturnValue({ selectedGame: ref(selectedGame) })
    useLocationsStore.mockReturnValue({ locations: [], loading: false, error: null, fetchLocations: vi.fn() })
    useLocationObjectsStore.mockReturnValue({ locationObjects: [], loading: false, error: null, fetchLocationObjects: vi.fn() })
    useLocationObjectEffectsStore.mockReturnValue({
      locationObjectEffects: [],
      loading: false,
      error: null,
      pageNumber: 1,
      hasMore: false,
      fetchLocationObjectEffects: vi.fn(),
      createLocationObjectEffect: vi.fn(),
      updateLocationObjectEffect: vi.fn(),
      deleteLocationObjectEffect: vi.fn(),
    })
    useItemsStore.mockReturnValue({ items: [], loading: false, error: null, fetchItems: vi.fn() })
    useCreaturesStore.mockReturnValue({ creatures: [], loading: false, error: null, fetchCreatures: vi.fn() })
    useLocationLinksStore.mockReturnValue({ locationLinks: [], loading: false, error: null, fetchLocationLinks: vi.fn() })
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
    const wrapper = mount(StudioLocationObjectEffectsView)

    expect(wrapper.text()).toContain('Please select or create a game to manage location object effects.')
    expect(wrapper.find('.game-table-section').exists()).toBe(false)
  })

  it('renders object effects table when game is selected', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectEffectsView)

    expect(wrapper.find('.game-table-section').exists()).toBe(true)
    expect(wrapper.find('h2').text()).toBe('Object Effects')
  })

  it('renders create button when game is selected', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectEffectsView)

    const createButton = wrapper.find('button')
    expect(createButton.text()).toBe('Create Object Effect')
  })

  it('renders ResourceTable with correct columns', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectEffectsView)

    const resourceTable = wrapper.findComponent({ name: 'ResourceTable' })
    expect(resourceTable.exists()).toBe(true)
    expect(resourceTable.props('columns')).toEqual([
      { key: 'object_name', label: 'Object' },
      { key: 'action_type', label: 'Action' },
      { key: 'effect_type', label: 'Effect' },
      { key: 'required_state_name', label: 'Required State' },
      { key: 'result_description', label: 'Description' },
      { key: 'is_repeatable', label: 'Repeatable' },
    ])
  })

  it('opens create modal when create button is clicked', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectEffectsView)

    const createButton = wrapper.find('button')
    await createButton.trigger('click')

    expect(wrapper.vm.showModal).toBe(true)
    expect(wrapper.vm.modalMode).toBe('create')
  })

  it('renders create modal with correct title', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectEffectsView)

    wrapper.vm.showModal = true
    wrapper.vm.modalMode = 'create'
    await wrapper.vm.$nextTick()

    expect(findInBody('.modal h2').textContent).toBe('Create Object Effect')
  })

  it('renders delete confirmation modal with correct props', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectEffectsView)

    wrapper.vm.showDeleteConfirm = true
    await wrapper.vm.$nextTick()

    const confirmationModal = wrapper.findComponent({ name: 'ConfirmationModal' })
    expect(confirmationModal.exists()).toBe(true)
    expect(confirmationModal.props('visible')).toBe(true)
    expect(confirmationModal.props('title')).toBe('Delete Object Effect')
    expect(confirmationModal.props('message')).toBe('Are you sure you want to delete this object effect?')
  })

  it('displays error message in modal', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })

    const wrapper = mount(StudioLocationObjectEffectsView)

    wrapper.vm.showModal = true
    wrapper.vm.modalError = 'Validation failed'
    await wrapper.vm.$nextTick()

    expect(findInBody('.modal .error').textContent).toBe('Validation failed')
  })

  it('effectFields only shows base fields when no effect_type is set', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioLocationObjectEffectsView)

    const keys = wrapper.vm.effectFields.map((f) => f.key)

    expect(keys).toContain('adventure_game_location_object_id')
    expect(keys).toContain('action_type')
    expect(keys).toContain('effect_type')
    expect(keys).toContain('result_description')
    expect(keys).not.toContain('result_value_min')
    expect(keys).not.toContain('result_adventure_game_item_id')
    expect(keys).not.toContain('result_adventure_game_location_id')
  })

  it('effectFields shows result_value_min and result_value_max as required for damage', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioLocationObjectEffectsView)

    wrapper.vm.modalForm.effect_type = 'damage'
    await wrapper.vm.$nextTick()

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).toContain('result_value_min')
    expect(keys).toContain('result_value_max')
    expect(fields.find((f) => f.key === 'result_value_min').required).toBe(true)
    expect(fields.find((f) => f.key === 'result_value_max').required).toBe(true)
    expect(keys).not.toContain('result_adventure_game_item_id')
  })

  it('effectFields shows result_adventure_game_location_id as required for teleport', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioLocationObjectEffectsView)

    wrapper.vm.modalForm.effect_type = 'teleport'
    await wrapper.vm.$nextTick()

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).toContain('result_adventure_game_location_id')
    expect(fields.find((f) => f.key === 'result_adventure_game_location_id').required).toBe(true)
    expect(keys).not.toContain('result_value_min')
  })

  it('effectFields shows result_adventure_game_item_id as required for give_item', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioLocationObjectEffectsView)

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
    const wrapper = mount(StudioLocationObjectEffectsView)

    wrapper.vm.modalForm.effect_type = 'open_link'
    await wrapper.vm.$nextTick()

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).toContain('result_adventure_game_location_link_id')
    expect(fields.find((f) => f.key === 'result_adventure_game_location_link_id').required).toBe(true)
  })

  it('effectFields shows result_adventure_game_creature_id as required for summon_creature', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioLocationObjectEffectsView)

    wrapper.vm.modalForm.effect_type = 'summon_creature'
    await wrapper.vm.$nextTick()

    const fields = wrapper.vm.effectFields
    const keys = fields.map((f) => f.key)

    expect(keys).toContain('result_adventure_game_creature_id')
    expect(fields.find((f) => f.key === 'result_adventure_game_creature_id').required).toBe(true)
    expect(keys).not.toContain('result_value_min')
  })

  it('effectFields shows no result fields for info and nothing effect types', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioLocationObjectEffectsView)

    for (const effectType of ['info', 'nothing']) {
      wrapper.vm.modalForm.effect_type = effectType
      await wrapper.vm.$nextTick()

      const keys = wrapper.vm.effectFields.map((f) => f.key)
      expect(keys).not.toContain('result_value_min')
      expect(keys).not.toContain('result_adventure_game_item_id')
      expect(keys).not.toContain('result_adventure_game_location_id')
    }
  })

  it('effectFields updates reactively when effect_type changes', async () => {
    await setupStoreMocks({ id: 'game-1', name: 'Test Game' })
    const wrapper = mount(StudioLocationObjectEffectsView)

    wrapper.vm.modalForm.effect_type = 'heal'
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.effectFields.map((f) => f.key)).toContain('result_value_min')

    wrapper.vm.modalForm.effect_type = 'give_item'
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.effectFields.map((f) => f.key)).toContain('result_adventure_game_item_id')
    expect(wrapper.vm.effectFields.map((f) => f.key)).not.toContain('result_value_min')
  })
})
