import { baseUrl, getAuthHeaders, apiFetch } from './baseUrl';

/**
 * Upload a turn sheet background image for a location
 * @param {string} gameId - The game ID
 * @param {string} locationId - The location ID
 * @param {File} imageFile - The image file to upload
 * @returns {Promise<Object>} - The created/updated image record
 */
export async function uploadLocationTurnSheetImage(gameId, locationId, imageFile) {
    const formData = new FormData();
    formData.append('image', imageFile);

    const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${gameId}/locations/${locationId}/turn-sheet-image`, {
        method: 'POST',
        headers: { ...getAuthHeaders() },
        body: formData,
    });

    if (!res.ok) {
        const errorData = await res.json().catch(() => ({}));
        throw new Error(errorData.error?.message || 'Failed to upload image');
    }

    return await res.json();
}

/**
 * Get the turn sheet background image for a location
 * @param {string} gameId - The game ID
 * @param {string} locationId - The location ID
 * @returns {Promise<Object>} - The turn sheet image data
 */
export async function getLocationTurnSheetImage(gameId, locationId) {
    const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${gameId}/locations/${locationId}/turn-sheet-image`, {
        headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    });

    if (!res.ok) {
        const errorData = await res.json().catch(() => ({}));
        throw new Error(errorData.error?.message || 'Failed to fetch turn sheet image');
    }

    return await res.json();
}

/**
 * Delete the turn sheet background image for a location
 * @param {string} gameId - The game ID
 * @param {string} locationId - The location ID
 * @returns {Promise<void>}
 */
export async function deleteLocationTurnSheetImage(gameId, locationId) {
    const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${gameId}/locations/${locationId}/turn-sheet-image`, {
        method: 'DELETE',
        headers: { ...getAuthHeaders() },
    });

    if (!res.ok && res.status !== 204) {
        const errorData = await res.json().catch(() => ({}));
        throw new Error(errorData.error?.message || 'Failed to delete turn sheet image');
    }
}

/**
 * Get the preview URL for a location's choice turn sheet
 * @param {string} gameId - The game ID
 * @param {string} locationId - The location ID
 * @returns {string} - The preview URL
 */
export function getLocationChoiceTurnSheetPreviewUrl(gameId, locationId) {
    return `${baseUrl}/api/v1/adventure-games/${gameId}/locations/${locationId}/turn-sheets/preview`;
}

