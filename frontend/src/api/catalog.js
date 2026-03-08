import { baseUrl, apiFetch, handleApiError } from './baseUrl';

export async function listCatalogGames() {
  const res = await apiFetch(`${baseUrl}/api/v1/catalog/game-subscriptions`, {
    headers: { 'Content-Type': 'application/json' },
  });
  await handleApiError(res, 'Failed to fetch game catalog');
  return await res.json();
}

export async function listCatalogGameInstances() {
  const res = await apiFetch(`${baseUrl}/api/v1/catalog/game-instances`, {
    headers: { 'Content-Type': 'application/json' },
  });
  await handleApiError(res, 'Failed to fetch game catalog');
  return await res.json();
}
