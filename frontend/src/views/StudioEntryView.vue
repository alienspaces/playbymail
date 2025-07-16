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
  margin: 2rem auto;
  padding: 2.5rem 2.5rem 2.5rem 2.5rem;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
}
.studio-entry h1 {
  margin-bottom: 1.5rem;
  font-size: 2.2rem;
}
.studio-entry p {
  margin-bottom: 1.25rem;
  font-size: 1.1rem;
}
.studio-entry ul {
  margin-bottom: 1.5rem;
  margin-left: 2rem;
  padding-left: 0.5rem;
}
.studio-entry li {
  margin-bottom: 0.75rem;
  font-size: 1.05rem;
}
.studio-entry em {
  display: block;
  margin-top: 1.5rem;
  font-size: 1.05rem;
}
.studio-main {
  display: flex;
  gap: 2rem;
  align-items: flex-start;
}
.studio-help {
  min-width: 250px;
  background: #f7f7f7;
  border: 1px solid #eee;
  border-radius: 6px;
  padding: 1.5rem 1rem;
  margin: 2rem 0 2rem 2rem;
  font-size: 0.98em;
}
</style> 