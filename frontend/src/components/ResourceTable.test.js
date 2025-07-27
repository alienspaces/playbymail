import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ResourceTable from './ResourceTable.vue'

describe('ResourceTable', () => {
  const mockColumns = [
    { key: 'name', label: 'Name' },
    { key: 'description', label: 'Description' }
  ]

  const mockRows = [
    { id: 1, name: 'Test Item 1', description: 'Description 1' },
    { id: 2, name: 'Test Item 2', description: 'Description 2' }
  ]

  it('renders table with data', () => {
    const wrapper = mount(ResourceTable, {
      props: {
        columns: mockColumns,
        rows: mockRows,
        loading: false,
        error: null
      }
    })

    expect(wrapper.find('table').exists()).toBe(true)
    expect(wrapper.find('thead').exists()).toBe(true)
    expect(wrapper.find('tbody').exists()).toBe(true)
    
    // Check column headers (no actions column when no actions slot)
    const headers = wrapper.findAll('th')
    expect(headers).toHaveLength(2) // Only the data columns
    expect(headers[0].text()).toBe('Name')
    expect(headers[1].text()).toBe('Description')
    
    // Check data rows
    const rows = wrapper.findAll('tbody tr')
    expect(rows).toHaveLength(2)
    
    const firstRowCells = rows[0].findAll('td')
    expect(firstRowCells[0].text()).toBe('Test Item 1')
    expect(firstRowCells[1].text()).toBe('Description 1')
  })

  it('renders loading state', () => {
    const wrapper = mount(ResourceTable, {
      props: {
        columns: mockColumns,
        rows: [],
        loading: true,
        error: null
      }
    })

    expect(wrapper.find('table').exists()).toBe(false)
    expect(wrapper.text()).toContain('Loading...')
  })

  it('renders error state', () => {
    const errorMessage = 'Failed to load data'
    const wrapper = mount(ResourceTable, {
      props: {
        columns: mockColumns,
        rows: [],
        loading: false,
        error: errorMessage
      }
    })

    expect(wrapper.find('table').exists()).toBe(false)
    expect(wrapper.find('.error').text()).toBe(errorMessage)
  })

  it('renders empty state when no data', () => {
    const wrapper = mount(ResourceTable, {
      props: {
        columns: mockColumns,
        rows: [],
        loading: false,
        error: null
      }
    })

    expect(wrapper.find('table').exists()).toBe(false)
    expect(wrapper.text()).toContain('No records found.')
  })

  it('renders actions column when actions slot is provided', () => {
    const wrapper = mount(ResourceTable, {
      props: {
        columns: mockColumns,
        rows: mockRows,
        loading: false,
        error: null
      },
      slots: {
        actions: '<button>Edit</button>'
      }
    })

    const headers = wrapper.findAll('th')
    expect(headers).toHaveLength(3)
    expect(headers[2].text()).toBe('Actions')
    
    const actionCells = wrapper.findAll('td').filter(cell => cell.text().includes('Edit'))
    expect(actionCells).toHaveLength(2)
  })

  it('does not render actions column when no actions slot', () => {
    const wrapper = mount(ResourceTable, {
      props: {
        columns: mockColumns,
        rows: mockRows,
        loading: false,
        error: null
      }
    })

    const headers = wrapper.findAll('th')
    expect(headers).toHaveLength(2) // Only the data columns
    expect(headers[0].text()).toBe('Name')
    expect(headers[1].text()).toBe('Description')
  })

  it('passes row data to actions slot', () => {
    const wrapper = mount(ResourceTable, {
      props: {
        columns: mockColumns,
        rows: mockRows,
        loading: false,
        error: null
      },
      slots: {
        actions: '<template #actions="{ row }"><span>{{ row.name }}</span></template>'
      }
    })

    // Find action cells (they should be the last td in each row)
    const rows = wrapper.findAll('tbody tr')
    const actionCells = rows.map(row => row.findAll('td').slice(-1)[0])
    expect(actionCells).toHaveLength(2)
    expect(actionCells[0].text()).toBe('Test Item 1')
    expect(actionCells[1].text()).toBe('Test Item 2')
  })

  it('handles empty columns array', () => {
    const wrapper = mount(ResourceTable, {
      props: {
        columns: [],
        rows: mockRows,
        loading: false,
        error: null
      },
      slots: {
        actions: '<button>Edit</button>'
      }
    })

    expect(wrapper.find('table').exists()).toBe(true)
    const headers = wrapper.findAll('th')
    expect(headers).toHaveLength(1) // Only actions column
    expect(headers[0].text()).toBe('Actions')
  })

  it('handles null/undefined rows gracefully', () => {
    const wrapper = mount(ResourceTable, {
      props: {
        columns: mockColumns,
        rows: null,
        loading: false,
        error: null
      }
    })

    expect(wrapper.find('table').exists()).toBe(false)
    expect(wrapper.text()).toContain('No records found.')
  })
}) 