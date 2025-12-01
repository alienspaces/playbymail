import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
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

  // Helper to find elements in document body (where Teleport renders)
  const findInBody = (selector) => {
    return document.body.querySelector(selector)
  }

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)
    vi.clearAllMocks()
    // Clear any existing modals from previous tests
    document.body.innerHTML = ''
  })

  afterEach(() => {
    // Clean up after each test
    document.body.innerHTML = ''
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
    const contextLabel = wrapper.find('.game-context-label')
    const contextName = wrapper.find('.game-context-name')
    expect(contextLabel.text()).toBe('Game:')
    expect(contextName.text()).toBe('Test Game')
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
      { key: 'creature_name', label: 'Creature' },
      { key: 'location_name', label: 'Location' },
      { key: 'initial_count', label: 'Count' },
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

    expect(findInBody('.modal h2').textContent).toBe('Create Creature Placement')
  })

  it('renders edit modal with correct title', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)

    // Open modal in edit mode
    wrapper.vm.showCreaturePlacementModal = true
    wrapper.vm.creaturePlacementModalMode = 'edit'
    await wrapper.vm.$nextTick()

    expect(findInBody('.modal h2').textContent).toBe('Edit Creature Placement')
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

    const creatureSelect = findInBody('#adventure_game_creature_id')
    const locationSelect = findInBody('#adventure_game_location_id')

    expect(creatureSelect).toBeTruthy()
    expect(locationSelect).toBeTruthy()

    const creatureOptions = creatureSelect.querySelectorAll('option')
    const locationOptions = locationSelect.querySelectorAll('option')

    expect(creatureOptions).toHaveLength(3) // placeholder + 2 creatures
    expect(locationOptions).toHaveLength(3) // placeholder + 2 locations
  })

  it('closes modal when cancel button is clicked', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)

    // Open modal
    wrapper.vm.showCreaturePlacementModal = true
    await wrapper.vm.$nextTick()

    const cancelButton = findInBody('button[type="button"]')
    cancelButton.click()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.showCreaturePlacementModal).toBe(false)
  })

  it('displays error message in modal', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)

    // Open modal and set error
    wrapper.vm.showCreaturePlacementModal = true
    wrapper.vm.creaturePlacementModalError = 'Validation failed'
    await wrapper.vm.$nextTick()

    expect(findInBody('.modal .error').textContent).toBe('Validation failed')
  })

  it('renders delete confirmation modal', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)

    // Open delete confirmation
    wrapper.vm.showCreaturePlacementDeleteConfirm = true
    await wrapper.vm.$nextTick()

    // Check that ConfirmationModal is rendered with correct props
    const confirmationModal = wrapper.findComponent({ name: 'ConfirmationModal' })
    expect(confirmationModal.exists()).toBe(true)
    expect(confirmationModal.props('visible')).toBe(true)
    expect(confirmationModal.props('title')).toBe('Delete Creature Placement')
    expect(confirmationModal.props('message')).toBe('Are you sure you want to delete this creature placement?')
  })

  it('has correct CSS classes for styling', async () => {
    await setupStoreMocks({ id: 1, name: 'Test Game' })

    const wrapper = mount(StudioCreaturePlacementsView)

    expect(wrapper.find('.game-table-section').exists()).toBe(true)
    expect(wrapper.find('.game-context-name').exists()).toBe(true)
    expect(wrapper.findComponent({ name: 'PageHeader' }).exists()).toBe(true)
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

    const contextName = wrapper.find('.game-context-name')
    expect(contextName.text()).toBe('New Game')
  })
}) 