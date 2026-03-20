import { baseUrl, handleApiError } from './baseUrl';

export async function approveSubscription(gameSubscriptionId, email) {
  const url = `${baseUrl}/api/v1/game-subscriptions/${gameSubscriptionId}/approve?email=${encodeURIComponent(email)}`;
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
  });
  await handleApiError(res, 'Failed to confirm subscription');
  return await res.json();
}
