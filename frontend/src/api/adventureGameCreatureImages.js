import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Upload a portrait image for a creature
 * @param {string} gameId - The game ID
 * @param {string} creatureId - The creature ID
 * @param {File} imageFile - The image file to upload
 * @returns {Promise<Object>} - The created/updated image record
 */
export async function uploadAdventureGameCreatureImage(gameId, creatureId, imageFile) {
    const formData = new FormData();
    formData.append('image', imageFile);

    const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${gameId}/creatures/${creatureId}/image`, {
        method: 'POST',
        headers: { ...getAuthHeaders() },
        body: formData,
    });

    await handleApiError(res, 'Failed to upload creature image');

    return await res.json();
}

/**
 * Get the portrait image for a creature
 * @param {string} gameId - The game ID
 * @param {string} creatureId - The creature ID
 * @returns {Promise<Object>} - The creature image data
 */
export async function getAdventureGameCreatureImage(gameId, creatureId) {
    const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${gameId}/creatures/${creatureId}/image`, {
        headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    });

    await handleApiError(res, 'Failed to fetch creature image');

    return await res.json();
}

/**
 * Delete the portrait image for a creature
 * @param {string} gameId - The game ID
 * @param {string} creatureId - The creature ID
 * @returns {Promise<void>}
 */
export async function deleteAdventureGameCreatureImage(gameId, creatureId) {
    const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${gameId}/creatures/${creatureId}/image`, {
        method: 'DELETE',
        headers: { ...getAuthHeaders() },
    });

    if (res.status !== 204) {
        await handleApiError(res, 'Failed to delete creature image');
    }
}
