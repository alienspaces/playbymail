<template>
  <div>
    <div v-if="!isLoggedIn" class="studio-entry">
      <h1>Game Designer Studio</h1>
      <p>
        The Game Designer Studio is your workspace for creating, editing, and managing play-by-mail games. Here, designers can:
      </p>
      <ul>
        <li>Create and organize games</li>
        <li>Design locations, items, and creatures</li>
        <li>Link locations and define requirements</li>
        <li>Assign items and creatures to locations</li>
        <li>Test and iterate on adventure game designs</li>
      </ul>
      <p><em>Sign in to access the full studio features.</em></p>
    </div>
    <div v-else class="studio-entry">
      <div class="studio-help">
        <h2>Studio Help</h2>
        <ul>
          <li>Use the sidebar to navigate between games, locations, items, and creatures.</li>
          <li>Select or create a game to begin managing its resources.</li>
          <li>Each section allows you to add, edit, or remove resources.</li>
          <li>Placement lets you assign items and creatures to locations.</li>
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
.studio-entry {
  max-width: 700px;
  margin: var(--space-lg) auto;
  padding: var(--space-xl);
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
}
.studio-entry h1 {
  margin-bottom: var(--space-lg);
  font-size: var(--font-size-xl);
}
.studio-entry p {
  margin-bottom: var(--space-md);
  font-size: var(--font-size-md);
}
.studio-entry ul {
  margin-bottom: var(--space-lg);
  margin-left: var(--space-lg);
  padding-left: var(--space-sm);
}
.studio-entry li {
  margin-bottom: var(--space-sm);
  font-size: var(--font-size-sm);
}
.studio-entry em {
  display: block;
  margin-top: var(--space-lg);
  font-size: var(--font-size-sm);
}
.studio-main {
  display: flex;
  gap: var(--space-lg);
  align-items: flex-start;
}
.studio-help {
  min-width: 250px;
  background: var(--color-bg-alt);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-xl);
  margin: var(--space-lg) 0 var(--space-lg) var(--space-lg);
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
  font-size: var(--font-size-sm);
  text-align: left;
}
.studio-help h2 {
  font-size: var(--font-size-xl);
  margin-bottom: var(--space-lg);
  text-align: left;
}
.studio-help ul {
  margin-left: var(--space-lg);
  padding-left: var(--space-sm);
}
.studio-help li {
  margin-bottom: var(--space-sm);
}
</style> 