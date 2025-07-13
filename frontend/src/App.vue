<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div id="app">
    <nav class="main-nav">
      <router-link to="/" class="nav-logo">PlayByMail Studio</router-link>
      <router-link to="/studio" class="nav-link">Studio</router-link>
    </nav>
    <router-view />
  </div>
</template>

<style scoped>
.main-nav {
  display: flex;
  align-items: center;
  background: #2c3e50;
  padding: 1rem 2rem;
  gap: 2rem;
}
.nav-logo {
  color: #fff;
  font-weight: bold;
  font-size: 1.2rem;
  text-decoration: none;
  margin-right: auto;
}
.nav-link {
  color: #fff;
  text-decoration: none;
  font-weight: 500;
  padding: 0.5rem 1.2rem;
  border-radius: 4px;
  transition: background 0.2s;
}
.nav-link.router-link-exact-active,
.nav-link:hover {
  background: #1a232b;
}
</style>
