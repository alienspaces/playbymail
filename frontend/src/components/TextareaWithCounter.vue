<template>
    <div class="textarea-with-counter">
        <textarea :id="id" :value="modelValue" :required="required" :maxlength="maxLength" :placeholder="placeholder"
            :rows="rows" @input="handleInput" />
        <div class="character-count">
            {{ characterCount }}/{{ maxLength }} characters
        </div>
    </div>
</template>

<script>
export default {
    name: 'TextareaWithCounter',
    props: {
        modelValue: {
            type: String,
            default: ''
        },
        id: {
            type: String,
            required: true
        },
        maxLength: {
            type: Number,
            required: true
        },
        required: {
            type: Boolean,
            default: false
        },
        placeholder: {
            type: String,
            default: ''
        },
        rows: {
            type: Number,
            default: 4
        }
    },
    emits: ['update:modelValue'],
    computed: {
        characterCount() {
            return (this.modelValue || '').length;
        }
    },
    methods: {
        handleInput(event) {
            this.$emit('update:modelValue', event.target.value);
        }
    }
};
</script>

<style scoped>
.textarea-with-counter {
    display: flex;
    flex-direction: column;
}

.textarea-with-counter textarea {
    width: 100%;
    padding: var(--space-sm);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    font-family: inherit;
    font-size: inherit;
    resize: vertical;
}

.character-count {
    font-size: 0.875rem;
    color: var(--color-text-secondary);
    text-align: right;
    margin-top: 0.25rem;
}
</style>
