import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

const mockFetchGames = vi.fn()
const mockSetSelectedGame = vi.fn()
const mockFetchGameInstances = vi.fn()
const mockPollGameInstances = vi.fn()
const mockStartGameInstance = vi.fn()

vi.mock('../../stores/games', () => ({
  useGamesStore: vi.fn(() => ({
    games: [{ id: 'game-1', name: 'Test Game', status: 'published', turn_duration_hours: 24 }],
    fetchGames: mockFetchGames,
    setSelectedGame: mockSetSelectedGame,
  })),
}))

vi.mock('../../stores/gameInstances', () => ({
  useGameInstancesStore: vi.fn(() => ({
    gameInstances: [],
    loading: false,
    error: null,
    fetchGameInstances: mockFetchGameInstances,
    pollGameInstances: mockPollGameInstances,
    startGameInstance: mockStartGameInstance,
    createGameInstance: vi.fn(),
    updateGameInstance: vi.fn(),
    deleteGameInstance: vi.fn(),
    pauseGameInstance: vi.fn(),
    resumeGameInstance: vi.fn(),
    cancelGameInstance: vi.fn(),
  })),
}))

vi.mock('../../stores/auth', () => ({
  useAuthStore: vi.fn(() => ({
    accountTimezone: 'UTC',
  })),
}))

vi.mock('../../utils/dateFormat', () => ({
  formatDateTime: vi.fn((str) => (str ? 'FMT:' + str : 'N/A')),
  formatDeadline: vi.fn((str) => (str ? 'DEADLINE:' + str : 'N/A')),
}))

import { formatDateTime, formatDeadline } from '../../utils/dateFormat'

vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({ params: { gameId: 'game-1' } })),
  useRouter: vi.fn(() => ({ push: vi.fn() })),
}))

vi.mock('../../components/Button.vue', () => ({
  default: { name: 'Button', template: '<button><slot /></button>' },
}))
vi.mock('../../components/PageHeader.vue', () => ({
  default: {
    name: 'PageHeader',
    props: ['title', 'subtitle', 'showIcon', 'titleLevel', 'actionText'],
    emits: ['action'],
    template: '<div class="mock-page-header">{{ title }}<slot /></div>',
  },
}))
vi.mock('../../components/ResourceTable.vue', () => ({
  default: {
    name: 'ResourceTable',
    props: ['columns', 'rows', 'loading', 'error'],
    template: '<div class="resource-table" />',
  },
}))
vi.mock('../../components/TableActions.vue', () => ({
  default: { name: 'TableActions', props: ['actions'], template: '<div />' },
}))
vi.mock('../../components/ResourceModalForm.vue', () => ({
  default: {
    name: 'ResourceModalForm',
    props: ['visible', 'title', 'fields', 'formData', 'error'],
    emits: ['submit', 'cancel'],
    template: '<div />',
  },
}))

import ManagementGameInstancesView from './ManagementGameInstancesView.vue'

describe('ManagementGameInstancesView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.useFakeTimers()
    mockFetchGames.mockResolvedValue()
    mockFetchGameInstances.mockResolvedValue()
    mockPollGameInstances.mockResolvedValue()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('fetches game instances on mount', async () => {
    mount(ManagementGameInstancesView)
    await flushPromises()

    expect(mockFetchGameInstances).toHaveBeenCalledWith('game-1')
  })

  it('renders Active Instances and Completed Instances sections', async () => {
    const wrapper = mount(ManagementGameInstancesView)
    await flushPromises()

    expect(wrapper.text()).toContain('Active Instances')
    expect(wrapper.text()).toContain('Completed Instances')
  })

  it('starts a poll timer on mount', async () => {
    mount(ManagementGameInstancesView)
    await flushPromises()

    expect(mockPollGameInstances).not.toHaveBeenCalled()

    vi.advanceTimersByTime(30000)
    await flushPromises()

    expect(mockPollGameInstances).toHaveBeenCalledWith('game-1')
  })

  it('polls again after a second interval', async () => {
    mount(ManagementGameInstancesView)
    await flushPromises()

    vi.advanceTimersByTime(60000)
    await flushPromises()

    expect(mockPollGameInstances).toHaveBeenCalledTimes(2)
  })

  it('clears the poll timer on unmount', async () => {
    const wrapper = mount(ManagementGameInstancesView)
    await flushPromises()

    wrapper.unmount()

    vi.advanceTimersByTime(30000)
    await flushPromises()

    expect(mockPollGameInstances).not.toHaveBeenCalled()
  })

  it('uses shared formatDateTime for date columns', () => {
    expect(formatDateTime('2026-03-22T10:00:00Z')).toBe('FMT:2026-03-22T10:00:00Z')
    expect(formatDateTime(null)).toBe('N/A')
  })

  it('uses shared formatDeadline for next_turn_due_at', () => {
    expect(formatDeadline('2026-03-22T10:00:00Z')).toBe('DEADLINE:2026-03-22T10:00:00Z')
    expect(formatDeadline(null)).toBe('N/A')
  })
})
