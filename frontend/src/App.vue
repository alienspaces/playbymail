<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { useAuthStore } from './stores/auth';
import { storeToRefs } from 'pinia';

const authStore = useAuthStore();
const { sessionToken } = storeToRefs(authStore);

function logout() {
  authStore.logout();
  window.location.href = '/login';
}
</script>

<template>
  <div id="app">
    <nav class="navbar">
      <div class="nav-links">
        <router-link to="/" class="logo" exact-active-class="active">PlayByMail</router-link>
        <router-link to="/faq" exact-active-class="active">F.A.Q.</router-link>
        <router-link to="/studio" exact-active-class="active">Game Designer Studio</router-link>
        <router-link to="/admin" exact-active-class="active">Game Management & Admin</router-link>
      </div>
      <div class="nav-actions">
        <template v-if="sessionToken">
          <button @click="logout">Logout</button>
        </template>
        <template v-else>
          <router-link to="/login">Login</router-link>
        </template>
      </div>
    </nav>
    <router-view />
  </div>
</template>

<style>
body {
  background: #f6f8fa;
  margin: 0;
  font-family: 'Inter', Arial, sans-serif;
}
#app {
  min-height: 100vh;
}
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #11181c;
  color: #fff;
  padding: 1rem 2rem;
}
.nav-links {
  display: flex;
  gap: 2rem;
  align-items: center;
}
.nav-links a {
  color: #fff;
  text-decoration: none;
  font-weight: 500;
  padding: 0.25rem 0.5rem;
  border-radius: 3px;
  transition: background 0.2s;
}
.nav-links a.active {
  background: #1976d2;
  color: #fff;
}
.logo {
  font-weight: 700;
  font-size: 1.2rem;
}
.nav-actions {
  display: flex;
  gap: 1rem;
  align-items: center;
}
.nav-actions button {
  background: #fff;
  color: #11181c;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  font-weight: 600;
  cursor: pointer;
}
.nav-actions a {
  color: #fff;
  text-decoration: underline;
  font-weight: 500;
}
</style>
