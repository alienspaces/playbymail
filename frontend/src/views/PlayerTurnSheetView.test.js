import { describe, it, expect, vi, beforeEach } from 'vitest'
import { nextTick } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import PlayerTurnSheetView from './PlayerTurnSheetView.vue'

const mockGetGSITurnSheets = vi.fn()
const mockSubmitGSITurnSheets = vi.fn()
const mockDownloadGSITurnSheetPDF = vi.fn()
const mockUploadGSITurnSheetScan = vi.fn()
const mockVerifyGameSubscriptionToken = vi.fn()
const mockRequestNewTurnSheetToken = vi.fn()

vi.mock('../api/player', () => ({
  verifyGameSubscriptionToken: (...args) => mockVerifyGameSubscriptionToken(...args),
  requestNewTurnSheetToken: (...args) => mockRequestNewTurnSheetToken(...args),
  getGSITurnSheets: (...args) => mockGetGSITurnSheets(...args),
  submitGSITurnSheets: (...args) => mockSubmitGSITurnSheets(...args),
  downloadGSITurnSheetPDF: (...args) => mockDownloadGSITurnSheetPDF(...args),
  uploadGSITurnSheetScan: (...args) => mockUploadGSITurnSheetScan(...args),
}))

const mockSetSessionToken = vi.fn()

vi.mock('../stores/auth', () => ({
  useAuthStore: vi.fn(() => ({
    sessionToken: '',
    setSessionToken: mockSetSessionToken,
  })),
}))

let routeParams = { game_subscription_instance_id: 'gsi-abc-123' }

vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({
    params: routeParams,
  })),
  useRouter: vi.fn(() => ({ push: vi.fn() })),
}))

const mockSheets = [
  { id: 'ts-1', sheet_type: 'location_choice', turn_number: 1, is_completed: false },
  { id: 'ts-2', sheet_type: 'inventory_management', turn_number: 1, is_completed: false },
]

describe('PlayerTurnSheetView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    routeParams = { game_subscription_instance_id: 'gsi-abc-123' }
  })

  describe('without turn_sheet_token in route', () => {
    it('shows loading state while fetching', async () => {
      mockGetGSITurnSheets.mockReturnValue(new Promise(() => {}))

      const wrapper = mount(PlayerTurnSheetView)
      await nextTick()

      expect(wrapper.find('[data-testid="ts-loading"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="ts-list"]').exists()).toBe(false)
    })

    it('shows turn sheet list after successful load', async () => {
      mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      expect(wrapper.find('[data-testid="ts-loading"]').exists()).toBe(false)
      expect(wrapper.find('[data-testid="ts-list"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="ts-title"]').text()).toBe('Your Turn Sheets')
      expect(wrapper.find('[data-testid="ts-card-ts-1"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="ts-card-ts-2"]').exists()).toBe(true)
    })

    it('shows empty state when no turn sheets', async () => {
      mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: [] })

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      expect(wrapper.find('[data-testid="ts-empty"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="btn-submit"]').exists()).toBe(false)
    })

    it('shows load error when fetch fails', async () => {
      mockGetGSITurnSheets.mockRejectedValue(new Error('Server error'))

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      expect(wrapper.find('[data-testid="ts-load-error"]').exists()).toBe(true)
      expect(wrapper.text()).toContain('Server error')
    })
  })

  describe('with turn_sheet_token in route (auto-verify)', () => {
    beforeEach(() => {
      routeParams = {
        game_subscription_instance_id: 'gsi-abc-123',
        turn_sheet_token: 'token-xyz',
      }
    })

    it('auto-verifies token and loads turn sheets on success', async () => {
      mockVerifyGameSubscriptionToken.mockResolvedValue('session-abc')
      mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      expect(mockVerifyGameSubscriptionToken).toHaveBeenCalledWith('gsi-abc-123', 'token-xyz')
      expect(mockSetSessionToken).toHaveBeenCalledWith('session-abc')
      expect(wrapper.find('[data-testid="ts-list"]').exists()).toBe(true)
    })

    it('shows expired-token UI when verify fails', async () => {
      mockVerifyGameSubscriptionToken.mockRejectedValue(new Error('Token expired'))

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      expect(wrapper.find('[data-testid="ts-token-expired"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="expired-email-input"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="expired-request-btn"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="ts-list"]').exists()).toBe(false)
    })

    it('requests new link when email is submitted on expired-token UI', async () => {
      mockVerifyGameSubscriptionToken.mockRejectedValue(new Error('Token expired'))
      mockRequestNewTurnSheetToken.mockResolvedValue()

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      await wrapper.find('[data-testid="expired-email-input"]').setValue('player@example.com')
      await wrapper.find('[data-testid="expired-request-btn"]').trigger('submit')
      await flushPromises()

      expect(mockRequestNewTurnSheetToken).toHaveBeenCalledWith('gsi-abc-123', 'player@example.com')
      expect(wrapper.find('[data-testid="expired-success"]').exists()).toBe(true)
    })
  })

  describe('turn sheet interactions', () => {
    it('shows pending and completed status for each sheet', async () => {
      const sheets = [
        { id: 'ts-1', sheet_type: 'location_choice', turn_number: 1, is_completed: false },
        { id: 'ts-2', sheet_type: 'inventory_management', turn_number: 1, is_completed: true },
      ]
      mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: sheets })

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      expect(wrapper.find('[data-testid="ts-status-ts-1"]').text()).toBe('Pending')
      expect(wrapper.find('[data-testid="ts-status-ts-2"]').text()).toBe('Completed')
    })

    it('submit button is disabled when all sheets are already completed', async () => {
      const completedSheets = mockSheets.map((s) => ({ ...s, is_completed: true }))
      mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: completedSheets })

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      expect(wrapper.find('[data-testid="btn-submit"]').attributes('disabled')).toBeDefined()
      expect(wrapper.find('[data-testid="btn-submit"]').text()).toBe('Already submitted')
    })

    it('submit button is enabled when there are pending sheets', async () => {
      mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      expect(wrapper.find('[data-testid="btn-submit"]').attributes('disabled')).toBeUndefined()
    })

    it('shows success state after submission', async () => {
      mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
      mockSubmitGSITurnSheets.mockResolvedValue({ submitted_count: 2, total_count: 2 })

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      await wrapper.find('[data-testid="btn-submit"]').trigger('click')
      await flushPromises()

      expect(wrapper.find('[data-testid="ts-success"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="link-browse-more"]').attributes('href')).toBe('/games')
      expect(wrapper.find('[data-testid="ts-list"]').exists()).toBe(false)
    })

    it('shows submit error when submission fails', async () => {
      mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
      mockSubmitGSITurnSheets.mockRejectedValue(new Error('Server error'))

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      await wrapper.find('[data-testid="btn-submit"]').trigger('click')
      await flushPromises()

      expect(wrapper.find('[data-testid="submit-error"]').exists()).toBe(true)
      expect(wrapper.text()).toContain('Server error')
    })

    it('renders a Download PDF button for each turn sheet', async () => {
      mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      expect(wrapper.find('[data-testid="btn-download-ts-1"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="btn-download-ts-2"]').exists()).toBe(true)
    })

    it('calls downloadGSITurnSheetPDF when Download PDF is clicked', async () => {
      mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })

      const mockBlob = new Blob(['%PDF'], { type: 'application/pdf' })
      mockDownloadGSITurnSheetPDF.mockResolvedValue({
        blob: () => Promise.resolve(mockBlob),
      })

      URL.createObjectURL = vi.fn(() => 'blob:test-url')
      URL.revokeObjectURL = vi.fn()

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      await wrapper.find('[data-testid="btn-download-ts-1"]').trigger('click')
      await flushPromises()

      expect(mockDownloadGSITurnSheetPDF).toHaveBeenCalledWith('gsi-abc-123', 'ts-1')
    })

    it('renders an Upload Scan button for each pending turn sheet', async () => {
      mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      expect(wrapper.find('[data-testid="btn-upload-ts-1"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="btn-upload-ts-2"]').exists()).toBe(true)
    })

    it('calls uploadGSITurnSheetScan when a file is selected', async () => {
      mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
      mockUploadGSITurnSheetScan.mockResolvedValue({ data: {} })

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      const file = new File(['img'], 'scan.png', { type: 'image/png' })
      const input = wrapper.find('[data-testid="input-upload-ts-1"]')
      Object.defineProperty(input.element, 'files', {
        value: [file],
        configurable: true,
      })
      await input.trigger('change')
      await flushPromises()

      expect(mockUploadGSITurnSheetScan).toHaveBeenCalledWith('gsi-abc-123', 'ts-1', file)
    })

    it('calls submitGSITurnSheets with the correct gsi id', async () => {
      mockGetGSITurnSheets.mockResolvedValue({ turn_sheets: mockSheets })
      mockSubmitGSITurnSheets.mockResolvedValue({ submitted_count: 2, total_count: 2 })

      const wrapper = mount(PlayerTurnSheetView)
      await flushPromises()

      await wrapper.find('[data-testid="btn-submit"]').trigger('click')
      await flushPromises()

      expect(mockSubmitGSITurnSheets).toHaveBeenCalledWith('gsi-abc-123')
    })
  })
})
