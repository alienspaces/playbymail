<!--
  ManagementTurnSheetsView.vue
  Turn sheet management: download join game PDFs and upload scanned turn sheets.
-->
<template>
  <div class="turn-sheets-view">
    <PageHeader 
      title="Turn Sheets" 
      titleLevel="h2" 
      :showIcon="false"
      subtitle="Download join game forms and upload scanned turn sheets"
    />

    <!-- Loading state -->
    <div v-if="loading" class="loading-state">
      <p>Loading game information...</p>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
      <button @click="loadGame">Retry</button>
    </div>

    <!-- Main content -->
    <div v-else-if="game" class="turn-sheets-content">
      <!-- Download Section -->
      <DataCard title="Download Join Game Turn Sheet">
        <div class="section-content">
          <p class="section-description">
            Generate and download a join game turn sheet for <strong>{{ game.name }}</strong>. 
            This can be printed and distributed to new players who want to join the game.
          </p>
        </div>

        <template #primary>
          <Button 
            @click="downloadJoinGameSheet" 
            variant="primary" 
            size="small"
            :disabled="downloading"
          >
            <svg class="btn-icon" viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z"/>
            </svg>
            {{ downloading ? 'Downloading...' : 'Download PDF' }}
          </Button>
        </template>
      </DataCard>

      <!-- Upload Section -->
      <DataCard title="Upload Scanned Turn Sheets">
        <div class="section-content">
          <p class="section-description">
            Upload scanned or photographed turn sheets for processing. 
            Supported formats: JPEG, PNG. Maximum file size: 10MB.
          </p>

          <div class="upload-area" 
            :class="{ 'drag-over': isDragOver }"
            @dragover.prevent="isDragOver = true"
            @dragleave.prevent="isDragOver = false"
            @drop.prevent="handleDrop"
          >
            <input 
              ref="fileInput"
              type="file" 
              accept="image/jpeg,image/png"
              @change="handleFileSelect"
              class="file-input"
            />
            <div class="upload-content">
              <svg class="upload-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M9 16h6v-6h4l-7-7-7 7h4zm-4 2h14v2H5z"/>
              </svg>
              <p class="upload-text">
                Drag and drop an image here, or <button type="button" class="browse-link" @click="triggerFileInput">browse</button>
              </p>
              <p class="upload-hint">JPEG or PNG, max 10MB</p>
            </div>
          </div>

          <!-- Selected file preview -->
          <div v-if="selectedFile" class="selected-file">
            <div class="file-info">
              <svg class="file-icon" viewBox="0 0 24 24" fill="currentColor">
                <path d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.5 13.5l2.5 3.01L14.5 12l4.5 6H5l3.5-4.5z"/>
              </svg>
              <span class="file-name">{{ selectedFile.name }}</span>
              <span class="file-size">{{ formatFileSize(selectedFile.size) }}</span>
            </div>
            <button type="button" class="remove-file" @click="clearSelectedFile">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
              </svg>
            </button>
          </div>
        </div>

        <template #primary>
          <Button 
            @click="uploadTurnSheet" 
            variant="primary" 
            size="small"
            :disabled="!selectedFile || uploading"
          >
            <svg class="btn-icon" viewBox="0 0 24 24" fill="currentColor">
              <path d="M9 16h6v-6h4l-7-7-7 7h4zm-4 2h14v2H5z"/>
            </svg>
            {{ uploading ? 'Processing...' : 'Upload & Process' }}
          </Button>
        </template>
      </DataCard>

      <!-- Upload Result -->
      <DataCard v-if="uploadResult" title="Processing Result">
        <div class="result-content">
          <div class="result-status" :class="uploadResult.success ? 'success' : 'error'">
            <svg v-if="uploadResult.success" class="status-icon" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
            </svg>
            <svg v-else class="status-icon" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2C6.47 2 2 6.47 2 12s4.47 10 10 10 10-4.47 10-10S17.53 2 12 2zm5 13.59L15.59 17 12 13.41 8.41 17 7 15.59 10.59 12 7 8.41 8.41 7 12 10.59 15.59 7 17 8.41 13.41 12 17 15.59z"/>
            </svg>
            <span>{{ uploadResult.message }}</span>
          </div>

          <div v-if="uploadResult.data" class="result-details">
            <DataItem label="Turn Sheet ID" :value="uploadResult.data.turn_sheet_id" />
            <DataItem label="Sheet Type" :value="uploadResult.data.sheet_type" />
            <DataItem label="Status" :value="uploadResult.data.processing_status" />
            
            <div v-if="uploadResult.data.scanned_data" class="scanned-data">
              <h4>Extracted Data</h4>
              <div class="data-grid">
                <DataItem 
                  v-for="(value, key) in uploadResult.data.scanned_data" 
                  :key="key"
                  :label="formatLabel(key)" 
                  :value="value || '(empty)'" 
                />
              </div>
            </div>
          </div>
        </div>

        <template #primary>
          <Button @click="clearUploadResult" variant="secondary" size="small">
            Clear Result
          </Button>
        </template>
      </DataCard>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useGamesStore } from '../../stores/games';
import { useAuthStore } from '../../stores/auth';
import Button from '../../components/Button.vue';
import DataCard from '../../components/DataCard.vue';
import DataItem from '../../components/DataItem.vue';
import PageHeader from '../../components/PageHeader.vue';

const route = useRoute();
const gamesStore = useGamesStore();
const authStore = useAuthStore();

const gameId = computed(() => route.params.gameId);
const game = computed(() => gamesStore.games.find(g => g.id === gameId.value));

const loading = ref(false);
const error = ref(null);

// Download state
const downloading = ref(false);

// Upload state
const fileInput = ref(null);
const selectedFile = ref(null);
const isDragOver = ref(false);
const uploading = ref(false);
const uploadResult = ref(null);

onMounted(async () => {
  await loadGame();
});

const loadGame = async () => {
  loading.value = true;
  error.value = null;
  try {
    if (!game.value) {
      await gamesStore.fetchGames();
    }
    if (game.value) {
      gamesStore.setSelectedGame(game.value);
    }
  } catch (err) {
    error.value = err.message || 'Failed to load game information';
  } finally {
    loading.value = false;
  }
};

const downloadJoinGameSheet = async () => {
  downloading.value = true;
  error.value = null;
  
  try {
    const token = authStore.sessionToken;
    const url = `/api/v1/games/${gameId.value}/turn-sheets`;
    
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    if (!response.ok) {
      throw new Error(`Download failed: ${response.statusText}`);
    }

    // Get the PDF blob and trigger download
    const blob = await response.blob();
    const downloadUrl = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = `${game.value.name.replace(/\s+/g, '_')}_join_game.pdf`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(downloadUrl);
  } catch (err) {
    error.value = err.message || 'Failed to download turn sheet';
  } finally {
    downloading.value = false;
  }
};

const triggerFileInput = () => {
  fileInput.value?.click();
};

const handleFileSelect = (event) => {
  const file = event.target.files?.[0];
  if (file) {
    validateAndSetFile(file);
  }
};

const handleDrop = (event) => {
  isDragOver.value = false;
  const file = event.dataTransfer?.files?.[0];
  if (file) {
    validateAndSetFile(file);
  }
};

const validateAndSetFile = (file) => {
  // Validate file type
  if (!['image/jpeg', 'image/png'].includes(file.type)) {
    error.value = 'Invalid file type. Please upload a JPEG or PNG image.';
    return;
  }
  
  // Validate file size (10MB max)
  if (file.size > 10 * 1024 * 1024) {
    error.value = 'File too large. Maximum size is 10MB.';
    return;
  }
  
  selectedFile.value = file;
  error.value = null;
};

const clearSelectedFile = () => {
  selectedFile.value = null;
  if (fileInput.value) {
    fileInput.value.value = '';
  }
};

const uploadTurnSheet = async () => {
  if (!selectedFile.value) return;
  
  uploading.value = true;
  uploadResult.value = null;
  error.value = null;
  
  try {
    const token = authStore.sessionToken;
    
    // Read file as array buffer
    const arrayBuffer = await selectedFile.value.arrayBuffer();
    
    const response = await fetch('/api/v1/turn-sheets', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': selectedFile.value.type
      },
      body: arrayBuffer
    });

    const data = await response.json();

    if (!response.ok) {
      // Backend returns an array of errors like [{code, message}]
      const errorMessage = Array.isArray(data) && data.length > 0 
        ? data[0].message 
        : (data.message || data.error || 'Upload failed');
      uploadResult.value = {
        success: false,
        message: errorMessage,
        data: null
      };
      return;
    }

    uploadResult.value = {
      success: true,
      message: 'Turn sheet processed successfully!',
      data: data
    };
    
    // Clear the selected file after successful upload
    clearSelectedFile();
  } catch (err) {
    uploadResult.value = {
      success: false,
      message: err.message || 'Failed to upload turn sheet',
      data: null
    };
  } finally {
    uploading.value = false;
  }
};

const clearUploadResult = () => {
  uploadResult.value = null;
};

const formatFileSize = (bytes) => {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
};

const formatLabel = (key) => {
  return key
    .replace(/_/g, ' ')
    .replace(/\b\w/g, c => c.toUpperCase());
};
</script>

<style scoped>
.turn-sheets-view {
  max-width: 900px;
  margin: 0 auto;
}

.loading-state,
.error-state {
  text-align: center;
  padding: var(--space-xl);
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.error-state button {
  margin-top: var(--space-md);
  padding: var(--space-sm) var(--space-md);
  background: var(--color-primary);
  color: var(--color-text-light);
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
}

.turn-sheets-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-xl);
}

.section-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.section-description {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  line-height: 1.5;
  margin: 0;
}

.btn-icon {
  width: 16px;
  height: 16px;
  margin-right: var(--space-xs);
}

/* Upload area */
.upload-area {
  border: 2px dashed var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-xl);
  text-align: center;
  transition: all 0.2s ease;
  position: relative;
  cursor: pointer;
}

.upload-area:hover {
  border-color: var(--color-primary);
  background: var(--color-bg-light);
}

.upload-area.drag-over {
  border-color: var(--color-primary);
  background: var(--color-primary-light);
}

.file-input {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}

.upload-content {
  pointer-events: none;
}

.upload-icon {
  width: 48px;
  height: 48px;
  color: var(--color-text-muted);
  margin-bottom: var(--space-md);
}

.upload-text {
  font-size: var(--font-size-md);
  color: var(--color-text);
  margin: 0 0 var(--space-xs) 0;
}

.browse-link {
  background: none;
  border: none;
  color: var(--color-primary);
  cursor: pointer;
  font-size: inherit;
  text-decoration: underline;
  pointer-events: auto;
}

.upload-hint {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
  margin: 0;
}

/* Selected file */
.selected-file {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-md);
  background: var(--color-bg-light);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
}

.file-info {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.file-icon {
  width: 24px;
  height: 24px;
  color: var(--color-primary);
}

.file-name {
  font-weight: var(--font-weight-medium);
  color: var(--color-text);
}

.file-size {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}

.remove-file {
  background: none;
  border: none;
  cursor: pointer;
  padding: var(--space-xs);
  border-radius: var(--radius-sm);
  transition: background 0.2s;
}

.remove-file:hover {
  background: var(--color-danger-light);
}

.remove-file svg {
  width: 20px;
  height: 20px;
  color: var(--color-danger);
}

/* Upload result */
.result-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.result-status {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-md);
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-medium);
}

.result-status.success {
  background: var(--color-success-light);
  color: var(--color-success);
}

.result-status.error {
  background: var(--color-danger-light);
  color: var(--color-danger);
}

.status-icon {
  width: 24px;
  height: 24px;
}

.result-details {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.scanned-data {
  margin-top: var(--space-md);
  padding-top: var(--space-md);
  border-top: 1px solid var(--color-border);
}

.scanned-data h4 {
  margin: 0 0 var(--space-sm) 0;
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.data-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: var(--space-sm);
}
</style>

