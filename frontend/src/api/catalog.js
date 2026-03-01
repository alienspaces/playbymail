import { baseUrl, apiFetch, handleApiError } from './baseUrl';

export async function listCatalogGames() {
  const res = await apiFetch(`${baseUrl}/api/v1/catalog/games`, {
    headers: { 'Content-Type': 'application/json' },
  });
  await handleApiError(res, 'Failed to fetch game catalog');
  return await res.json();
}
