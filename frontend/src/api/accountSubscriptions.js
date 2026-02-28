import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

export async function getMyAccountSubscriptions() {
  const res = await apiFetch(`${baseUrl}/api/v1/account/subscriptions`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch account subscriptions');
  return await res.json();
}
