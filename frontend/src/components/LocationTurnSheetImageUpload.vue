<template>
    <div class="turn-sheet-image-upload">
        <h3 class="upload-title">Location Turn Sheet Background Image</h3>
        <p class="upload-description">
            Upload a background image for this location's turn sheets.
            Recommended: 2480 × 3508 pixels (A4 @ 300 DPI) for best print quality.
        </p>

        <div class="image-section">
            <div class="image-preview-container">
                <div v-if="image" class="image-preview">
                    <div class="image-info">
                        <span class="image-dimensions">
                            {{ image.width }} × {{ image.height }}px
                        </span>
                        <span class="image-size">
                            {{ formatFileSize(image.file_size) }}
                        </span>
                    </div>
                    <button class="remove-btn" @click="removeImage" :disabled="loading">
                        ✕
                    </button>
                </div>
                <div v-else class="image-placeholder">
                    <span class="placeholder-text">No background image uploaded</span>
                </div>
            </div>

            <div class="upload-controls">
                <input :id="fileInputId" type="file" accept="image/webp,image/png,image/jpeg" class="file-input"
                    @change="handleFileSelect" :disabled="loading" />
                <label :for="fileInputId" class="upload-btn">
                    {{ loading ? 'Uploading...' : (image ? 'Replace' : 'Upload') }}
                </label>
            </div>

            <div v-if="error" class="error-message">
                {{ error }}
            </div>
            <div v-if="warning" class="warning-message">
                {{ warning }}
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import {
    uploadLocationTurnSheetImage,
    getLocationTurnSheetImage,
    deleteLocationTurnSheetImage
} from '../api/locationImages';

defineOptions({
    name: 'LocationTurnSheetImageUpload'
});

const props = defineProps({
    gameId: {
        type: String,
        required: true
    },
    locationId: {
        type: String,
        required: true
    }
});

const emit = defineEmits(['imagesUpdated', 'loadingChanged']);

const image = ref(null);
const loading = ref(false);
const error = ref(null);
const warning = ref(null);

// Generate unique file input ID to avoid conflicts when multiple instances exist
const fileInputId = computed(() => `file-location-${props.locationId}`);

// Computed property to check if upload is in progress
const isUploading = computed(() => loading.value);

// Watch for loading changes and emit to parent
watch(isUploading, (newValue) => {
    emit('loadingChanged', newValue);
});

// Expose isUploading for parent components
defineExpose({
    isUploading
});

async function loadImage() {
    if (!props.gameId || !props.locationId) return;

    try {
        const response = await getLocationTurnSheetImage(props.gameId, props.locationId);
        if (response.data && response.data.background) {
            image.value = response.data.background;
        } else {
            image.value = null;
        }
    } catch (err) {
        console.error('Failed to load location turn sheet image:', err);
        image.value = null;
    }
}

async function handleFileSelect(event) {
    const file = event.target.files?.[0];
    if (!file) return;

    // Reset errors and warnings
    error.value = null;
    warning.value = null;

    // Validate file type
    const validTypes = ['image/webp', 'image/png', 'image/jpeg'];
    if (!validTypes.includes(file.type)) {
        error.value = 'Invalid file type. Please use WebP, PNG, or JPEG.';
        event.target.value = '';
        return;
    }

    // Validate file size (1MB max)
    if (file.size > 1048576) {
        error.value = 'File too large. Maximum size is 1MB.';
        event.target.value = '';
        return;
    }

    loading.value = true;
    emit('loadingChanged', true);
    console.log(`[LocationTurnSheetImageUpload] Starting upload for locationId: ${props.locationId}`);

    try {
        const response = await uploadLocationTurnSheetImage(props.gameId, props.locationId, file);
        console.log(`[LocationTurnSheetImageUpload] Upload response:`, response);
        if (response.data) {
            image.value = response.data;
            if (response.data.warning) {
                warning.value = response.data.warning;
            }
            emit('imagesUpdated');
        }
    } catch (err) {
        console.error(`[LocationTurnSheetImageUpload] Upload failed:`, err);
        error.value = err.message || 'Failed to upload image';
    } finally {
        loading.value = false;
        console.log(`[LocationTurnSheetImageUpload] Upload complete`);
        emit('loadingChanged', false);
        event.target.value = '';
    }
}

async function removeImage() {
    if (!image.value) return;

    loading.value = true;
    error.value = null;
    warning.value = null;
    emit('loadingChanged', true);

    try {
        await deleteLocationTurnSheetImage(props.gameId, props.locationId);
        image.value = null;
        emit('imagesUpdated');
    } catch (err) {
        error.value = err.message || 'Failed to delete image';
    } finally {
        loading.value = false;
        emit('loadingChanged', false);
    }
}

function formatFileSize(bytes) {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
}

onMounted(() => {
    loadImage();
});

watch(() => [props.gameId, props.locationId], () => {
    loadImage();
});
</script>

<style scoped>
.turn-sheet-image-upload {
    padding: var(--space-md);
    background: var(--color-bg-light);
    border-radius: var(--radius-md);
    border: 1px solid var(--color-border);
}

.upload-title {
    margin: 0 0 var(--space-xs) 0;
    font-size: var(--font-size-md);
    color: var(--color-text);
}

.upload-description {
    margin: 0 0 var(--space-md) 0;
    font-size: var(--font-size-sm);
    color: var(--color-text-muted);
}

.image-section {
    padding: var(--space-sm);
    background: var(--color-bg);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border-light);
}

.image-preview-container {
    margin-bottom: var(--space-sm);
}

.image-preview {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: var(--space-sm);
    background: var(--color-success-light, #e8f5e9);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-success, #4caf50);
}

.image-info {
    display: flex;
    gap: var(--space-md);
    font-size: var(--font-size-sm);
    color: var(--color-text);
}

.image-dimensions {
    font-weight: var(--font-weight-bold);
}

.image-size {
    color: var(--color-text-muted);
}

.remove-btn {
    background: var(--color-danger);
    color: white;
    border: none;
    border-radius: var(--radius-sm);
    width: 24px;
    height: 24px;
    cursor: pointer;
    font-size: var(--font-size-sm);
    display: flex;
    align-items: center;
    justify-content: center;
}

.remove-btn:hover:not(:disabled) {
    background: var(--color-danger-dark);
}

.remove-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.image-placeholder {
    padding: var(--space-sm);
    background: var(--color-bg-light);
    border-radius: var(--radius-sm);
    border: 1px dashed var(--color-border);
    text-align: center;
}

.placeholder-text {
    font-size: var(--font-size-sm);
    color: var(--color-text-muted);
}

.upload-controls {
    display: flex;
    gap: var(--space-sm);
}

.file-input {
    display: none;
}

.upload-btn {
    display: inline-block;
    padding: var(--space-xs) var(--space-sm);
    background: var(--color-button);
    color: var(--color-text-light);
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: var(--font-size-sm);
    font-weight: var(--font-weight-bold);
    transition: background 0.2s ease;
}

.upload-btn:hover {
    background: var(--color-button-hover);
}

.error-message {
    margin-top: var(--space-xs);
    padding: var(--space-xs) var(--space-sm);
    background: var(--color-danger-light, #ffebee);
    color: var(--color-danger);
    border-radius: var(--radius-sm);
    font-size: var(--font-size-sm);
}

.warning-message {
    margin-top: var(--space-xs);
    padding: var(--space-xs) var(--space-sm);
    background: var(--color-warning-light, #fff3e0);
    color: var(--color-warning-dark, #e65100);
    border-radius: var(--radius-sm);
    font-size: var(--font-size-sm);
}
</style>
