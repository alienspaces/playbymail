<!--
  StudioTurnSheetBackgroundsView.vue
  View for managing game-level turn sheet background images (join game, inventory management).
  Location-specific backgrounds are managed from the location edit form.
-->
<template>
    <div>
        <div v-if="!selectedGame">
            <p>Select a game to manage turn sheet backgrounds.</p>
        </div>
        <div v-else class="game-table-section">
            <GameContext :gameName="selectedGame.name" />
            <PageHeader title="Turn Sheet Backgrounds" :showIcon="false" titleLevel="h2" />

            <div class="turn-sheet-backgrounds">
                <p class="description">
                    Upload background images for game-level turn sheets. Location-specific backgrounds are managed from
                    the
                    <router-link :to="`/studio/${selectedGame.id}/locations`">Locations</router-link> page.
                </p>

                <!-- Tabs for different turn sheet types -->
                <div class="tabs">
                    <button v-for="sheetType in availableSheetTypes" :key="sheetType.value"
                        :class="['tab', { active: activeTab === sheetType.value }]"
                        @click="setActiveTab(sheetType.value)">
                        {{ sheetType.label }}
                    </button>
                </div>

                <!-- Content for active turn sheet type -->
                <div v-if="activeSheetType" class="tab-content">
                    <div class="sheet-type-section">
                        <h3>{{ activeSheetType.label }} Turn Sheet Background</h3>
                        <p class="sheet-description">{{ activeSheetType.description }}</p>

                        <GameTurnSheetImageUpload :gameId="String(selectedGame.id)" :turnSheetType="activeSheetType.value"
                            :image="getImageForType(activeSheetType.value)" @imagesUpdated="handleImagesUpdated"
                            @loadingChanged="handleLoadingChanged" />

                        <div class="preview-section">
                            <button class="preview-btn" @click="openPreview(activeSheetType.value)" :disabled="loading">
                                Preview Turn Sheet
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Preview Modal -->
        <GameTurnSheetPreviewModal :visible="showPreviewModal" :gameId="selectedGame?.id ? String(selectedGame.id) : ''"
            :gameName="selectedGame?.name || 'Game'" :turnSheetType="previewTurnSheetType"
            :title="`${getSheetTypeLabel(previewTurnSheetType)} Turn Sheet Preview`" @close="closePreviewModal" />
    </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { useGamesStore } from '../../../stores/games';
import { getGameTurnSheetImages } from '../../../api/gameImages';
import PageHeader from '../../../components/PageHeader.vue';
import GameContext from '../../../components/GameContext.vue';
import GameTurnSheetImageUpload from '../../../components/GameTurnSheetImageUpload.vue';
import GameTurnSheetPreviewModal from '../../../components/GameTurnSheetPreviewModal.vue';

const gamesStore = useGamesStore();
const { selectedGame } = storeToRefs(gamesStore);

// Turn sheet types available for game-level backgrounds
const availableSheetTypes = [
    {
        value: 'adventure_game_join_game',
        label: 'Join Game',
        description: 'Background image for the join game turn sheet that new players receive when joining.'
    },
    {
        value: 'adventure_game_inventory_management',
        label: 'Inventory Management',
        description: 'Background image for the inventory management turn sheet used by players to manage their items.'
    }
];

const activeTab = ref(availableSheetTypes[0].value);
const images = ref([]);
const loading = ref(false);
const showPreviewModal = ref(false);
const previewTurnSheetType = ref('');

// Computed property to get the active sheet type object
const activeSheetType = computed(() => {
    return availableSheetTypes.find(st => st.value === activeTab.value) || null;
});

// Set active tab
function setActiveTab(tabValue) {
    console.log('[StudioTurnSheetBackgroundsView] setActiveTab called with:', tabValue);
    activeTab.value = tabValue;
    console.log('[StudioTurnSheetBackgroundsView] activeTab is now:', activeTab.value);
}

// Get image for a specific turn sheet type
function getImageForType(turnSheetType) {
    return images.value.find(img => img.turn_sheet_type === turnSheetType) || null;
}

// Get label for a turn sheet type
function getSheetTypeLabel(turnSheetType) {
    const sheetType = availableSheetTypes.find(st => st.value === turnSheetType);
    return sheetType ? sheetType.label : turnSheetType;
}

// Load images for the selected game
async function loadImages() {
    if (!selectedGame.value) return;

    loading.value = true;
    try {
        const response = await getGameTurnSheetImages(selectedGame.value.id);
        images.value = response.data || [];
    } catch (error) {
        console.error('Failed to load turn sheet images:', error);
        images.value = [];
    } finally {
        loading.value = false;
    }
}

// Handle images updated event
function handleImagesUpdated() {
    loadImages();
}

// Handle loading changed event
function handleLoadingChanged(isLoading) {
    loading.value = isLoading;
}

// Open preview modal
function openPreview(turnSheetType) {
    console.log('[StudioTurnSheetBackgroundsView] openPreview called with turnSheetType:', turnSheetType);
    previewTurnSheetType.value = turnSheetType;
    console.log('[StudioTurnSheetBackgroundsView] previewTurnSheetType set to:', previewTurnSheetType.value);
    showPreviewModal.value = true;
}

// Close preview modal
function closePreviewModal() {
    showPreviewModal.value = false;
    previewTurnSheetType.value = '';
}

// Watch for game selection changes
watch(
    () => selectedGame.value,
    (newGame) => {
        if (newGame) {
            loadImages();
        }
    },
    { immediate: true }
);
</script>

<style scoped>
.turn-sheet-backgrounds {
    margin-top: var(--space-md);
}

.description {
    margin-bottom: var(--space-md);
    color: var(--color-text-muted);
    font-size: var(--font-size-sm);
}

.description a {
    color: var(--color-primary);
    text-decoration: none;
}

.description a:hover {
    text-decoration: underline;
}

.tabs {
    display: flex;
    gap: var(--space-xs);
    border-bottom: 2px solid var(--color-border);
    margin-bottom: var(--space-md);
}

.tab {
    padding: var(--space-sm) var(--space-md);
    background: none;
    border: none;
    border-bottom: 2px solid transparent;
    cursor: pointer;
    font-size: var(--font-size-base);
    font-weight: var(--font-weight-semibold);
    color: var(--color-text-muted);
    transition: all 0.2s ease;
    margin-bottom: -2px;
}

.tab:hover {
    color: var(--color-text);
    background: var(--color-bg-light);
}

.tab.active {
    color: var(--color-primary);
    border-bottom-color: var(--color-primary);
}

.tab-content {
    margin-top: var(--space-md);
}

.sheet-type-section {
    padding: var(--space-md);
    background: var(--color-bg-light);
    border-radius: var(--radius-md);
    border: 1px solid var(--color-border);
}

.sheet-type-section h3 {
    margin: 0 0 var(--space-xs) 0;
    font-size: var(--font-size-lg);
    color: var(--color-text);
}

.sheet-description {
    margin: 0 0 var(--space-md) 0;
    color: var(--color-text-muted);
    font-size: var(--font-size-sm);
}

.preview-section {
    margin-top: var(--space-md);
    padding-top: var(--space-md);
    border-top: 1px solid var(--color-border);
}

.preview-btn {
    padding: var(--space-sm) var(--space-md);
    background: var(--color-button);
    color: var(--color-text-light);
    border: none;
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: var(--font-size-base);
    font-weight: var(--font-weight-bold);
    transition: background 0.2s ease;
}

.preview-btn:hover:not(:disabled) {
    background: var(--color-button-hover);
}

.preview-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}
</style>
