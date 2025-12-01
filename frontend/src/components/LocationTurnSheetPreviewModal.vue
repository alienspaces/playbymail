<template>
    <!-- Teleport modal to body to avoid z-index stacking context issues -->
    <Teleport to="body">
        <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
            <div class="preview-modal">
                <div class="modal-header">
                    <h2>{{ title }}</h2>
                    <button class="close-btn" @click="$emit('close')" aria-label="Close">✕</button>
                </div>

                <div class="modal-content">
                    <!-- Show loading spinner while iframe is loading -->
                    <div v-if="loading && !error" class="loading-state">
                        <div class="spinner"></div>
                        <p>Generating preview...</p>
                    </div>

                    <!-- Show error state -->
                    <div v-if="error" class="error-state">
                        <p class="error-text">{{ error }}</p>
                        <button class="retry-btn" @click="retryPreview">Retry</button>
                    </div>

                    <!-- Always render iframe when visible (but hide it while loading) -->
                    <iframe v-show="!loading && !error" :key="iframeKey" :src="previewUrl" class="pdf-preview"
                        title="Turn Sheet Preview" @load="onIframeLoad" @error="onIframeError"></iframe>
                </div>

                <div class="modal-footer">
                    <button class="close-btn-secondary" @click="$emit('close')">Close</button>
                </div>
            </div>
        </div>
    </Teleport>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { getLocationChoiceTurnSheetPreviewUrl } from '../api/locationImages';
import { useAuthStore } from '../stores/auth';

defineOptions({
    name: 'LocationTurnSheetPreviewModal'
});

const props = defineProps({
    visible: {
        type: Boolean,
        default: false
    },
    gameId: {
        type: String,
        required: true
    },
    locationId: {
        type: String,
        required: true
    },
    locationName: {
        type: String,
        default: 'Location'
    },
    title: {
        type: String,
        default: 'Location Choice Turn Sheet Preview'
    }
});

defineEmits(['close']);

const loading = ref(true);
const error = ref(null);
const iframeKey = ref(0); // Used to force iframe reload

const authStore = useAuthStore();

const previewUrl = computed(() => {
    if (!props.gameId || !props.locationId || !props.visible) return '';
    const baseUrl = getLocationChoiceTurnSheetPreviewUrl(props.gameId, props.locationId);
    const token = authStore.sessionToken;
    // Append token as query param for iframe authentication
    const url = token ? `${baseUrl}?token=${encodeURIComponent(token)}` : baseUrl;
    console.log('[LocationTurnSheetPreviewModal] Preview URL:', url);
    return url;
});

function startLoading() {
    console.log('[LocationTurnSheetPreviewModal] Starting to load preview');
    loading.value = true;
    error.value = null;
}

function retryPreview() {
    console.log('[LocationTurnSheetPreviewModal] Retrying preview');
    iframeKey.value++; // Force iframe to reload
    startLoading();
}

function onIframeLoad() {
    console.log('[LocationTurnSheetPreviewModal] Iframe loaded successfully');
    loading.value = false;
}

function onIframeError() {
    console.log('[LocationTurnSheetPreviewModal] Iframe load error');
    loading.value = false;
    error.value = 'Failed to load preview. Please try again.';
}

watch(() => props.visible, (newVal) => {
    if (newVal) {
        console.log('[LocationTurnSheetPreviewModal] Modal opened, gameId:', props.gameId, 'locationId:', props.locationId);
        startLoading();
    } else {
        // Reset state when modal closes
        loading.value = true;
        error.value = null;
    }
});
</script>

<style scoped>
.modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
    /* High z-index to appear above all page elements */
    padding: var(--space-md);
}

.preview-modal {
    background: var(--color-bg);
    border-radius: var(--radius-lg);
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    width: 100%;
    max-width: 900px;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: var(--space-md) var(--space-lg);
    border-bottom: 1px solid var(--color-border);
    background: var(--color-bg-light);
}

.modal-header h2 {
    margin: 0;
    font-size: var(--font-size-lg);
    color: var(--color-text);
}

.close-btn {
    background: none;
    border: none;
    font-size: var(--font-size-lg);
    cursor: pointer;
    color: var(--color-text-muted);
    padding: var(--space-xs);
    border-radius: var(--radius-sm);
    transition: all 0.2s ease;
}

.close-btn:hover {
    background: var(--color-border);
    color: var(--color-text);
}

.modal-content {
    flex: 1;
    overflow: hidden;
    min-height: 500px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.loading-state,
.error-state {
    text-align: center;
    padding: var(--space-xl);
}

.spinner {
    width: 48px;
    height: 48px;
    border: 4px solid var(--color-border);
    border-top-color: var(--color-button);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin: 0 auto var(--space-md);
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}

.error-text {
    color: var(--color-danger);
    margin-bottom: var(--space-md);
}

.retry-btn {
    padding: var(--space-sm) var(--space-md);
    background: var(--color-button);
    color: var(--color-text-light);
    border: none;
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-weight: var(--font-weight-bold);
}

.retry-btn:hover {
    background: var(--color-button-hover);
}

.pdf-preview {
    width: 100%;
    height: 100%;
    min-height: 500px;
    border: none;
}

.modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: var(--space-sm);
    padding: var(--space-md) var(--space-lg);
    border-top: 1px solid var(--color-border);
    background: var(--color-bg-light);
}

.close-btn-secondary {
    padding: var(--space-sm) var(--space-md);
    background: var(--color-bg);
    color: var(--color-text);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-weight: var(--font-weight-bold);
}

.close-btn-secondary:hover {
    background: var(--color-border);
}

@media (max-width: 768px) {
    .preview-modal {
        max-height: 95vh;
    }

    .modal-content {
        min-height: 400px;
    }

    .modal-footer {
        flex-direction: column-reverse;
    }

    .close-btn-secondary {
        width: 100%;
        text-align: center;
        justify-content: center;
    }
}
</style>
