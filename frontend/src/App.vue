<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div id="app">
    <nav class="navbar">
      <router-link to="/" class="logo">PlayByMail Studio</router-link>
      <div class="nav-actions">
        <template v-if="isAuthenticated">
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

<script>
export default {
  name: 'App',
  computed: {
    isAuthenticated() {
      return !!localStorage.getItem('session_token');
    },
  },
  methods: {
    logout() {
      localStorage.removeItem('session_token');
      this.$router.push('/login');
    },
  },
};
</script>

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
.logo {
  color: #fff;
  text-decoration: none;
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
