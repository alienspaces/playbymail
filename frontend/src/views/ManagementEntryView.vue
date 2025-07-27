<template>
  <div>
    <div v-if="!isLoggedIn" class="admin-entry">
      <h1>Game Management</h1>
      <p>
        This area is for game managers to subscribe to or purchase games, configure runtime parameters, review automated turn processing, attend to player queries, and provide assistance. The platform automatically processes player turns, but game managers provide human oversight and support.
      </p>
      <ul>
        <li>Subscribe to or purchase games to run</li>
        <li>Configure game runtime parameters and settings</li>
        <li>Review automated turn processing and outcomes</li>
        <li>Attend to player queries and provide assistance</li>
        <li>Review potential issues highlighted by the platform</li>
        <li>Monitor game states and player progress</li>
      </ul>
      <p><em>Sign in to access game management tools.</em></p>
      <div class="admin-cta">
        <router-link to="/login" class="cta-button">Sign In to Game Management</router-link>
      </div>
    </div>
    <div v-else class="admin-entry">
      <div class="admin-help">
        <h2>Game Management Help</h2>
        <ul>
          <li>Subscribe to or purchase games you want to run for players.</li>
          <li>Configure game runtime parameters like turn frequency and player limits.</li>
          <li>Review automated turn processing results and player submissions.</li>
          <li>Attend to player queries and provide personalized assistance.</li>
          <li>Review potential issues flagged by the platform for human oversight.</li>
          <li>Monitor game states, player progress, and overall game health.</li>
        </ul>
      </div>
      <router-view />
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth';
import { computed } from 'vue';
const authStore = useAuthStore();
const isLoggedIn = computed(() => !!authStore.sessionToken);
</script>

<style scoped>
.admin-entry {
  max-width: 900px;
  margin: 2rem auto;
  padding: 2.5rem 2.5rem 2.5rem 2.5rem;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
}
.admin-entry h1 {
  margin-bottom: 1.5rem;
  font-size: 2.2rem;
}
.admin-entry p {
  margin-bottom: 1.25rem;
  font-size: 1.1rem;
}
.admin-entry ul {
  margin-bottom: 1.5rem;
  margin-left: 2rem;
  padding-left: 0.5rem;
}
.admin-entry li {
  margin-bottom: 0.75rem;
  font-size: 1.05rem;
}
.admin-entry em {
  display: block;
  margin-top: 1.5rem;
  font-size: 1.05rem;
}
.admin-main {
  display: flex;
  gap: 2rem;
  align-items: flex-start;
}
.admin-help {
  min-width: 250px;
  background: #f7f7f7;
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 2.5rem;
  margin: 2rem 0 2rem 2rem;
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
  font-size: 1.05rem;
  text-align: left;
}
.admin-help h2 {
  font-size: 2.2rem;
  margin-bottom: 1.5rem;
  text-align: left;
}
.admin-help ul {
  margin-left: 2rem;
  padding-left: 0.5rem;
}
.admin-help li {
  margin-bottom: 0.75rem;
}
.admin-cta {
  text-align: center;
  margin-top: var(--space-lg);
}
.cta-button {
  display: inline-block;
  padding: var(--space-md) var(--space-xl);
  background: var(--color-primary);
  color: var(--color-text-light);
  border-radius: var(--radius-sm);
  text-decoration: none;
  font-weight: var(--font-weight-bold);
  transition: background 0.2s;
}
.cta-button:hover {
  background: var(--color-primary-dark);
}
</style> 