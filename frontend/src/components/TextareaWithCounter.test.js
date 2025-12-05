import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TextareaWithCounter from './TextareaWithCounter.vue'

describe('TextareaWithCounter', () => {
    it('renders textarea with correct attributes', () => {
        const wrapper = mount(TextareaWithCounter, {
            props: {
                id: 'test-textarea',
                modelValue: '',
                maxLength: 512,
                required: true,
                placeholder: 'Enter description',
                rows: 4
            }
        })

        const textarea = wrapper.find('textarea')
        expect(textarea.exists()).toBe(true)
        expect(textarea.attributes('id')).toBe('test-textarea')
        expect(textarea.attributes('maxlength')).toBe('512')
        expect(textarea.attributes('required')).toBeDefined()
        expect(textarea.attributes('placeholder')).toBe('Enter description')
        expect(textarea.attributes('rows')).toBe('4')
    })

    it('displays character count correctly', () => {
        const wrapper = mount(TextareaWithCounter, {
            props: {
                id: 'test-textarea',
                modelValue: 'Hello world',
                maxLength: 512
            }
        })

        expect(wrapper.text()).toContain('11/512 characters')
    })

    it('updates character count when value changes', async () => {
        const wrapper = mount(TextareaWithCounter, {
            props: {
                id: 'test-textarea',
                modelValue: 'Test',
                maxLength: 512
            }
        })

        expect(wrapper.text()).toContain('4/512 characters')

        await wrapper.setProps({ modelValue: 'Longer text here' })
        expect(wrapper.text()).toContain('16/512 characters')
    })

    it('emits update:modelValue when textarea value changes', async () => {
        const wrapper = mount(TextareaWithCounter, {
            props: {
                id: 'test-textarea',
                modelValue: '',
                maxLength: 512
            }
        })

        const textarea = wrapper.find('textarea')
        await textarea.setValue('New value')

        expect(wrapper.emitted('update:modelValue')).toBeTruthy()
        expect(wrapper.emitted('update:modelValue')[0]).toEqual(['New value'])
    })

    it('handles empty value correctly', () => {
        const wrapper = mount(TextareaWithCounter, {
            props: {
                id: 'test-textarea',
                modelValue: null,
                maxLength: 512
            }
        })

        expect(wrapper.text()).toContain('0/512 characters')
    })
})

