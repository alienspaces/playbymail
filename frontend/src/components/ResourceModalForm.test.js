import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import ResourceModalForm from './ResourceModalForm.vue'

describe('ResourceModalForm', () => {
  const mockFields = [
    { key: 'name', label: 'Name', required: true, maxlength: 100 },
    { key: 'description', label: 'Description', type: 'textarea', maxlength: 500 }
  ]

  const mockModelValue = {
    name: 'Test Item',
    description: 'Test Description'
  }

  const mockSelectFields = [
    { key: 'category', label: 'Category', type: 'select', required: true, placeholder: 'Select a category...' },
    { key: 'count', label: 'Count', type: 'number', required: true, min: 1 }
  ]

  const mockOptions = {
    category: [
      { value: 'option1', label: 'Option 1' },
      { value: 'option2', label: 'Option 2' }
    ]
  }

  // Helper to find elements in document body (where Teleport renders)
  const findInBody = (selector) => {
    return document.body.querySelector(selector)
  }

  const findAllInBody = (selector) => {
    return document.body.querySelectorAll(selector)
  }

  beforeEach(() => {
    // Clear any existing modals from previous tests
    document.body.innerHTML = ''
  })

  afterEach(() => {
    // Clean up after each test
    document.body.innerHTML = ''
  })

  it('renders modal when visible is true', () => {
    mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'create',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: null
      }
    })

    expect(findInBody('.modal-overlay')).toBeTruthy()
    expect(findInBody('.modal')).toBeTruthy()
  })

  it('does not render when visible is false', () => {
    mount(ResourceModalForm, {
      props: {
        visible: false,
        mode: 'create',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: null
      }
    })

    expect(findInBody('.modal-overlay')).toBeNull()
  })

  it('displays correct title for create mode', () => {
    mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'create',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: null
      }
    })

    expect(findInBody('h2').textContent).toBe('Create Item')
  })

  it('displays correct title for edit mode', () => {
    mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'edit',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: null
      }
    })

    expect(findInBody('h2').textContent).toBe('Edit Item')
  })

  it('renders form fields correctly', () => {
    mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'create',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: null
      }
    })

    const formGroups = findAllInBody('.form-group')
    expect(formGroups).toHaveLength(2)

    // Check labels
    const labels = findAllInBody('label')
    expect(labels[0].textContent.trim()).toBe('Name *')
    expect(labels[1].textContent.trim()).toBe('Description')

    // Check inputs and textareas
    const inputs = findAllInBody('input')
    const textareas = findAllInBody('textarea')
    expect(inputs).toHaveLength(1)
    expect(textareas).toHaveLength(1)
    
    expect(inputs[0].getAttribute('id')).toBe('name')
    expect(inputs[0].hasAttribute('required')).toBe(true)
    expect(inputs[0].getAttribute('maxlength')).toBe('100')
    
    expect(textareas[0].getAttribute('id')).toBe('description')
    expect(textareas[0].getAttribute('maxlength')).toBe('500')
  })

  it('renders select fields correctly', () => {
    mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'create',
        title: 'Item',
        fields: mockSelectFields,
        modelValue: {},
        error: null,
        options: mockOptions
      }
    })

    const selects = findAllInBody('select')
    const inputs = findAllInBody('input')
    expect(selects).toHaveLength(1)
    expect(inputs).toHaveLength(1)

    expect(selects[0].getAttribute('id')).toBe('category')
    expect(selects[0].hasAttribute('required')).toBe(true)
    
    const options = selects[0].querySelectorAll('option')
    expect(options).toHaveLength(3) // placeholder + 2 options
    expect(options[0].textContent).toBe('Select a category...')
    expect(options[1].textContent).toBe('Option 1')
    expect(options[2].textContent).toBe('Option 2')
  })

  it('populates form with modelValue', async () => {
    const wrapper = mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'edit',
        title: 'Item',
        fields: mockFields,
        modelValue: mockModelValue,
        error: null
      }
    })

    await wrapper.vm.$nextTick()

    const nameInput = findInBody('#name')
    const descriptionInput = findInBody('#description')

    expect(nameInput.value).toBe('Test Item')
    expect(descriptionInput.value).toBe('Test Description')
  })

  it('emits submit event with form data', async () => {
    const wrapper = mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'create',
        title: 'Item',
        fields: mockFields,
        modelValue: mockModelValue,
        error: null
      }
    })

    await wrapper.vm.$nextTick()

    const form = findInBody('form')
    form.dispatchEvent(new Event('submit', { cancelable: true }))
    await wrapper.vm.$nextTick()

    expect(wrapper.emitted('submit')).toBeTruthy()
    expect(wrapper.emitted('submit')[0][0]).toEqual(mockModelValue)
  })

  it('emits cancel event when cancel button is clicked', async () => {
    const wrapper = mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'create',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: null
      }
    })

    const cancelButton = findInBody('button[type="button"]')
    cancelButton.click()
    await wrapper.vm.$nextTick()

    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('displays error message when error prop is provided', () => {
    const errorMessage = 'Validation failed'
    mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'create',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: errorMessage
      }
    })

    expect(findInBody('.error').textContent).toBe(errorMessage)
  })

  it('does not display error when error prop is null', () => {
    mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'create',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: null
      }
    })

    expect(findInBody('.error')).toBeNull()
  })

  it('shows correct button text for create mode', () => {
    mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'create',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: null
      }
    })

    const submitButton = findInBody('button[type="submit"]')
    expect(submitButton.textContent).toBe('Create')
  })

  it('shows correct button text for edit mode', () => {
    mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'edit',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: null
      }
    })

    const submitButton = findInBody('button[type="submit"]')
    expect(submitButton.textContent).toBe('Save')
  })

  it('handles custom field slots', () => {
    mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'create',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: null
      },
      slots: {
        field: '<template #field="{ field, value, update }"><select :value="value" @change="update($event.target.value)"><option value="">Select...</option></select></template>'
      }
    })

    const selects = findAllInBody('select')
    expect(selects).toHaveLength(2)
  })

  it('updates form data when modelValue changes', async () => {
    const wrapper = mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'edit',
        title: 'Item',
        fields: mockFields,
        modelValue: { name: 'Initial Name' },
        error: null
      }
    })

    await wrapper.vm.$nextTick()
    expect(findInBody('#name').value).toBe('Initial Name')

    await wrapper.setProps({
      modelValue: { name: 'Updated Name' }
    })

    await wrapper.vm.$nextTick()
    expect(findInBody('#name').value).toBe('Updated Name')
  })
}) 