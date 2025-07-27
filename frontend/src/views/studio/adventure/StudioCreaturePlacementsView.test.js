import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { ref } from 'vue'
import StudioCreaturePlacementsView from './StudioCreaturePlacementsView.vue'

// Mock the stores
vi.mock('../../../stores/creatures', () => ({
  useCreaturesStore: vi.fn(() => ({
    creatures: [],
    loading: false,
    error: null,
    fetchCreatures: vi.fn()
  }))
}))

vi.mock('../../../stores/locations', () => ({
  useLocationsStore: vi.fn(() => ({
    locations: [],
    loading: false,
    error: null,
    fetchLocations: vi.fn()
  }))
}))

vi.mock('../../../stores/creaturePlacements', () => ({
  useCreaturePlacementsStore: vi.fn(() => ({
    creaturePlacements: [],
    loading: false,
    error: null,
    createCreaturePlacement: vi.fn(),
    updateCreaturePlacement: vi.fn(),
    deleteCreaturePlacement: vi.fn(),
    fetchCreaturePlacements: vi.fn()
  }))
}))

vi.mock('../../../stores/games', () => ({
  useGamesStore: vi.fn(() => ({
    selectedGame: ref(null)
  }))
}))

describe('StudioCreaturePlacementsView', () => {
  let pinia

  const setupStoreMocks = async (selectedGame = null) => {
    const { useGamesStore } = await import('../../../stores/games')
    const { useCreaturesStore } = await import('../../../stores/creatures')
    const { useLocationsStore } = await import('../../../stores/locations')
    const { useCreaturePlacementsStore } = await import('../../../stores/creaturePlacements')
    
    useGamesStore.mockReturnValue({
      selectedGame: ref(selectedGame)
    })
    useCreaturesStore.mockReturnValue({
      creatures: [],
      loading: false,
      error: null,
      fetchCreatures: vi.fn()
    })
    useLocationsStore.mockReturnValue({
      locations: [],
      loading: false,
      error: null,
      fetchLocations: vi.fn()
    })
    useCreaturePlacementsStore.mockReturnValue({
      creaturePlacements: [],
      loading: false,
      error: null,
      createCreaturePlacement: vi.fn(),
      updateCreaturePlacement: vi.fn(),
      deleteCreaturePlacement: vi.fn(),
      fetchCreaturePlacements: vi.fn()
    })
  }

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)
    vi.clearAllMocks()
  })

  it('renders prompt when no game is selected', () => {
    const wrapper = mount(StudioCreaturePlacementsView)
    
    expect(wrapper.text()).toContain('Please select or create a game to manage creature placements.')
    expect(wrapper.find('.game-table-section').exists()).toBe(false)
  })

  it('renders creature placements table when game is selected', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)
    
    expect(wrapper.find('.game-table-section').exists()).toBe(true)
    expect(wrapper.find('.game-context-name').text()).toBe('Game: Test Game')
    expect(wrapper.find('h2').text()).toBe('Creature Placements')
  })

  it('renders create button when game is selected', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)
    
    const createButton = wrapper.find('button')
    expect(createButton.text()).toBe('Create Creature Placement')
  })

  it('renders ResourceTable with correct props', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)
    
    const resourceTable = wrapper.findComponent({ name: 'ResourceTable' })
    expect(resourceTable.exists()).toBe(true)
    expect(resourceTable.props('columns')).toEqual([
      { key: 'adventure_game_creature_id', label: 'Creature ID' },
      { key: 'adventure_game_location_id', label: 'Location ID' },
      { key: 'created_at', label: 'Created' }
    ])
  })

  it('opens create modal when create button is clicked', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)
    
    const createButton = wrapper.find('button')
    await createButton.trigger('click')
    
    expect(wrapper.vm.showCreaturePlacementModal).toBe(true)
    expect(wrapper.vm.creaturePlacementModalMode).toBe('create')
  })

  it('renders create modal with correct title', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)
    
    // Open modal
    wrapper.vm.showCreaturePlacementModal = true
    wrapper.vm.creaturePlacementModalMode = 'create'
    await wrapper.vm.$nextTick()
    
    expect(wrapper.find('.modal h2').text()).toBe('Create Creature Placement')
  })

  it('renders edit modal with correct title', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)
    
    // Open modal in edit mode
    wrapper.vm.showCreaturePlacementModal = true
    wrapper.vm.creaturePlacementModalMode = 'edit'
    await wrapper.vm.$nextTick()
    
    expect(wrapper.find('.modal h2').text()).toBe('Edit Creature Placement')
  })

  it('renders creature and location select options', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })
    
    // Override the creatures and locations for this specific test
    const { useCreaturesStore } = await import('../../../stores/creatures')
    const { useLocationsStore } = await import('../../../stores/locations')
    
    useCreaturesStore.mockReturnValue({
      creatures: [
        { id: 1, name: 'Dragon' },
        { id: 2, name: 'Goblin' }
      ],
      loading: false,
      error: null,
      fetchCreatures: vi.fn()
    })
    useLocationsStore.mockReturnValue({
      locations: [
        { id: 1, name: 'Cave' },
        { id: 2, name: 'Forest' }
      ],
      loading: false,
      error: null,
      fetchLocations: vi.fn()
    })

    const wrapper = mount(StudioCreaturePlacementsView)
    
    // Open modal
    wrapper.vm.showCreaturePlacementModal = true
    await wrapper.vm.$nextTick()
    
    const creatureSelect = wrapper.find('#adventure_game_creature_id')
    const locationSelect = wrapper.find('#adventure_game_location_id')
    
    expect(creatureSelect.exists()).toBe(true)
    expect(locationSelect.exists()).toBe(true)
    
    const creatureOptions = creatureSelect.findAll('option')
    const locationOptions = locationSelect.findAll('option')
    
    expect(creatureOptions).toHaveLength(3) // placeholder + 2 creatures
    expect(locationOptions).toHaveLength(3) // placeholder + 2 locations
  })

  it('closes modal when cancel button is clicked', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)
    
    // Open modal
    wrapper.vm.showCreaturePlacementModal = true
    await wrapper.vm.$nextTick()
    
    const cancelButton = wrapper.find('button[type="button"]')
    await cancelButton.trigger('click')
    
    expect(wrapper.vm.showCreaturePlacementModal).toBe(false)
  })

  it('displays error message in modal', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)
    
    // Open modal and set error
    wrapper.vm.showCreaturePlacementModal = true
    wrapper.vm.creaturePlacementModalError = 'Validation failed'
    await wrapper.vm.$nextTick()
    
    expect(wrapper.find('.modal .error').text()).toBe('Validation failed')
  })

  it('renders delete confirmation modal', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)
    
    // Open delete confirmation
    wrapper.vm.showCreaturePlacementDeleteConfirm = true
    await wrapper.vm.$nextTick()
    
    expect(wrapper.find('.modal h2').text()).toBe('Delete Creature Placement')
    expect(wrapper.text()).toContain('Are you sure you want to delete this creature placement?')
  })

  it('has correct CSS classes for styling', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)
    
    expect(wrapper.find('.game-table-section').exists()).toBe(true)
    expect(wrapper.find('.game-context-name').exists()).toBe(true)
    expect(wrapper.find('.section-header').exists()).toBe(true)
  })

  it('watches for selectedGame changes', async () => {
    const selectedGameRef = ref(null)
    const { useGamesStore } = await import('../../../stores/games')
    const { useCreaturesStore } = await import('../../../stores/creatures')
    const { useLocationsStore } = await import('../../../stores/locations')
    const { useCreaturePlacementsStore } = await import('../../../stores/creaturePlacements')
    
    useGamesStore.mockReturnValue({
      selectedGame: selectedGameRef
    })
    useCreaturesStore.mockReturnValue({
      creatures: [],
      loading: false,
      error: null,
      fetchCreatures: vi.fn()
    })
    useLocationsStore.mockReturnValue({
      locations: [],
      loading: false,
      error: null,
      fetchLocations: vi.fn()
    })
    useCreaturePlacementsStore.mockReturnValue({
      creaturePlacements: [],
      loading: false,
      error: null,
      createCreaturePlacement: vi.fn(),
      updateCreaturePlacement: vi.fn(),
      deleteCreaturePlacement: vi.fn(),
      fetchCreaturePlacements: vi.fn()
    })

    const wrapper = mount(StudioCreaturePlacementsView)
    
    // Initially no game selected
    expect(wrapper.text()).toContain('Please select or create a game')
    
    // Change selected game
    selectedGameRef.value = { id: 1, name: 'New Game' }
    await wrapper.vm.$nextTick()
    
    expect(wrapper.find('.game-context-name').text()).toBe('Game: New Game')
  })
}) 