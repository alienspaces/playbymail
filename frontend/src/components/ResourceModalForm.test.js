import { describe, it, expect } from 'vitest'
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

  it('renders modal when visible is true', () => {
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

    expect(wrapper.find('.modal-overlay').exists()).toBe(true)
    expect(wrapper.find('.modal').exists()).toBe(true)
  })

  it('does not render when visible is false', () => {
    const wrapper = mount(ResourceModalForm, {
      props: {
        visible: false,
        mode: 'create',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: null
      }
    })

    expect(wrapper.find('.modal-overlay').exists()).toBe(false)
  })

  it('displays correct title for create mode', () => {
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

    expect(wrapper.find('h2').text()).toBe('Create Item')
  })

  it('displays correct title for edit mode', () => {
    const wrapper = mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'edit',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: null
      }
    })

    expect(wrapper.find('h2').text()).toBe('Edit Item')
  })

  it('renders form fields correctly', () => {
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

    const formGroups = wrapper.findAll('.form-group')
    expect(formGroups).toHaveLength(2)

    // Check labels
    const labels = wrapper.findAll('label')
    expect(labels[0].text()).toBe('Name *')
    expect(labels[1].text()).toBe('Description')

    // Check inputs and textareas
    const inputs = wrapper.findAll('input')
    const textareas = wrapper.findAll('textarea')
    expect(inputs).toHaveLength(1)
    expect(textareas).toHaveLength(1)
    
    expect(inputs[0].attributes('id')).toBe('name')
    expect(inputs[0].attributes('required')).toBeDefined()
    expect(inputs[0].attributes('maxlength')).toBe('100')
    
    expect(textareas[0].attributes('id')).toBe('description')
    expect(textareas[0].attributes('maxlength')).toBe('500')
  })

  it('renders select fields correctly', () => {
    const wrapper = mount(ResourceModalForm, {
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

    const selects = wrapper.findAll('select')
    const inputs = wrapper.findAll('input')
    expect(selects).toHaveLength(1)
    expect(inputs).toHaveLength(1)

    expect(selects[0].attributes('id')).toBe('category')
    expect(selects[0].attributes('required')).toBeDefined()
    
    const options = selects[0].findAll('option')
    expect(options).toHaveLength(3) // placeholder + 2 options
    expect(options[0].text()).toBe('Select a category...')
    expect(options[1].text()).toBe('Option 1')
    expect(options[2].text()).toBe('Option 2')
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

    const nameInput = wrapper.find('#name')
    const descriptionInput = wrapper.find('#description')

    expect(nameInput.element.value).toBe('Test Item')
    expect(descriptionInput.element.value).toBe('Test Description')
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

    const form = wrapper.find('form')
    await form.trigger('submit')

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

    const cancelButton = wrapper.find('button[type="button"]')
    await cancelButton.trigger('click')

    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('displays error message when error prop is provided', () => {
    const errorMessage = 'Validation failed'
    const wrapper = mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'create',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: errorMessage
      }
    })

    expect(wrapper.find('.error').text()).toBe(errorMessage)
  })

  it('does not display error when error prop is null', () => {
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

    expect(wrapper.find('.error').exists()).toBe(false)
  })

  it('shows correct button text for create mode', () => {
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

    const submitButton = wrapper.find('button[type="submit"]')
    expect(submitButton.text()).toBe('Create')
  })

  it('shows correct button text for edit mode', () => {
    const wrapper = mount(ResourceModalForm, {
      props: {
        visible: true,
        mode: 'edit',
        title: 'Item',
        fields: mockFields,
        modelValue: {},
        error: null
      }
    })

    const submitButton = wrapper.find('button[type="submit"]')
    expect(submitButton.text()).toBe('Save')
  })

  it('handles custom field slots', () => {
    const wrapper = mount(ResourceModalForm, {
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

    const selects = wrapper.findAll('select')
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
    expect(wrapper.find('#name').element.value).toBe('Initial Name')

    await wrapper.setProps({
      modelValue: { name: 'Updated Name' }
    })

    await wrapper.vm.$nextTick()
    expect(wrapper.find('#name').element.value).toBe('Updated Name')
  })
}) 