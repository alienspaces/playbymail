<template>
  <div>
    <div v-if="!isLoggedIn" class="admin-entry">
      <h1>Game Management & Administration</h1>
      <p>
        This area is for administrators to configure and monitor the automated systems that handle scanning, printing, and mailing for play-by-mail games. Most processes run automatically, but you can adjust operational settings, monitor system status, and intervene if issues arise.
      </p>
      <ul>
        <li>Configure scanning, printing, and mailing parameters</li>
        <li>Monitor job queues, system health, and logs</li>
        <li>Handle exceptions and review flagged items</li>
        <li>Manage running games and player accounts</li>
      </ul>
      <p><em>Sign in to access administration tools.</em></p>
    </div>
    <div v-else class="admin-main">
      <div class="admin-help">
        <h2>Admin Help</h2>
        <ul>
          <li>Configure and monitor automated scanning, printing, and mailing systems.</li>
          <li>View job queues, logs, and system health dashboards.</li>
          <li>Handle exceptions, review flagged items, and re-run failed jobs if needed.</li>
          <li>Manage running games and player accounts.</li>
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
  max-width: 700px;
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
  border-radius: 6px;
  padding: 1.5rem 1rem;
  margin: 2rem 0 2rem 2rem;
  font-size: 0.98em;
}
</style> 