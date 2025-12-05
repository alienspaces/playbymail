import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Upload a turn sheet background image for a game
 * @param {string} gameId - The game ID
 * @param {File} imageFile - The image file to upload
 * @param {string} turnSheetType - The turn sheet type (e.g., 'adventure_game_join_game', 'adventure_game_inventory_management')
 * @returns {Promise<Object>} - The created/updated image record
 */
export async function uploadGameTurnSheetImage(gameId, imageFile, turnSheetType = 'adventure_game_join_game') {
  const formData = new FormData();
  formData.append('image', imageFile);

  const url = new URL(`${baseUrl}/api/v1/games/${gameId}/turn-sheet-images`);
  url.searchParams.set('turn_sheet_type', turnSheetType);

  const res = await apiFetch(url.toString(), {
    method: 'POST',
    headers: { ...getAuthHeaders() },
    body: formData,
  });

  await handleApiError(res, 'Failed to upload image');

  return await res.json();
}

/**
 * Get turn sheet background images for a game
 * @param {string} gameId - The game ID
 * @param {string} [turnSheetType] - Optional turn sheet type filter
 * @returns {Promise<Object>} - The turn sheet images data
 */
export async function getGameTurnSheetImages(gameId, turnSheetType = null) {
  const url = new URL(`${baseUrl}/api/v1/games/${gameId}/turn-sheet-images`);
  if (turnSheetType) {
    url.searchParams.set('turn_sheet_type', turnSheetType);
  }

  const res = await apiFetch(url.toString(), {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });

  await handleApiError(res, 'Failed to fetch turn sheet images');

  return await res.json();
}

/**
 * Get a specific turn sheet background image by ID
 * @param {string} gameId - The game ID
 * @param {string} imageId - The image ID
 * @returns {Promise<Object>} - The turn sheet image data
 */
export async function getGameTurnSheetImage(gameId, imageId) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/turn-sheet-images/${imageId}`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });

  await handleApiError(res, 'Failed to fetch turn sheet image');

  return await res.json();
}

/**
 * Delete a turn sheet background image by ID
 * @param {string} gameId - The game ID
 * @param {string} imageId - The image ID
 * @returns {Promise<void>}
 */
export async function deleteGameTurnSheetImage(gameId, imageId) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${gameId}/turn-sheet-images/${imageId}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });

  if (res.status !== 204) {
    await handleApiError(res, 'Failed to delete turn sheet image');
  }
}

/**
 * Get the preview URL for a game's turn sheet
 * @param {string} gameId - The game ID
 * @param {string} turnSheetType - The turn sheet type (e.g., 'adventure_game_join_game', 'adventure_game_inventory_management')
 * @returns {string} - The preview URL
 */
export function getGameTurnSheetPreviewUrl(gameId, turnSheetType) {
  const url = new URL(`${baseUrl}/api/v1/games/${gameId}/turn-sheets/preview`);
  url.searchParams.set('turn_sheet_type', turnSheetType);
  return url.toString();
}

/**
 * Legacy function for backward compatibility
 * @deprecated Use getGameTurnSheetImages instead
 */
export async function getGameTurnSheetImageLegacy(gameId) {
  const response = await getGameTurnSheetImages(gameId, 'adventure_game_join_game');
  // Return in old format for backward compatibility
  return {
    data: {
      background: response.data?.[0] || null
    }
  };
}

/**
 * Legacy function for backward compatibility
 * @deprecated Use getGameJoinTurnSheetPreviewUrl instead
 */
export function getGameJoinTurnSheetPreviewUrl(gameId) {
  return getGameTurnSheetPreviewUrl(gameId, 'adventure_game_join_game');
}
