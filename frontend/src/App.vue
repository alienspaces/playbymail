<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { useAuthStore } from './stores/auth';
import { storeToRefs } from 'pinia';
import { ref } from 'vue';

const authStore = useAuthStore();
const { sessionToken } = storeToRefs(authStore);

function logout() {
  authStore.logout();
  window.location.href = '/login';
}

const mobileMenuOpen = ref(false);
function toggleMobileMenu() {
  mobileMenuOpen.value = !mobileMenuOpen.value;
}
function closeMobileMenu() {
  mobileMenuOpen.value = false;
}
</script>

<template>
  <div id="app">
    <nav class="navbar">
      <div class="nav-links">
        <router-link to="/" class="logo">PlayByMail</router-link>
        <router-link to="/faq" exact-active-class="active" @click="closeMobileMenu">F.A.Q.</router-link>
        <router-link to="/studio" exact-active-class="active" @click="closeMobileMenu">Game Designer Studio</router-link>
        <router-link to="/admin" exact-active-class="active" @click="closeMobileMenu">Game Management & Admin</router-link>
      </div>
      <div class="nav-actions">
        <template v-if="sessionToken">
          <button @click="logout">Logout</button>
        </template>
        <template v-else>
          <router-link to="/login" @click="closeMobileMenu">Login</router-link>
        </template>
      </div>
      <div class="mobile-logo">PlayByMail</div>
      <button class="burger" @click="toggleMobileMenu" aria-label="Open navigation menu">
        <span :class="{ 'open': mobileMenuOpen }"></span>
        <span :class="{ 'open': mobileMenuOpen }"></span>
        <span :class="{ 'open': mobileMenuOpen }"></span>
      </button>
      <div class="mobile-menu" v-if="mobileMenuOpen">
        <router-link to="/faq" exact-active-class="active" @click="closeMobileMenu">F.A.Q.</router-link>
        <router-link to="/studio" exact-active-class="active" @click="closeMobileMenu">Game Designer Studio</router-link>
        <router-link to="/admin" exact-active-class="active" @click="closeMobileMenu">Game Management & Admin</router-link>
        <div class="mobile-actions">
          <template v-if="sessionToken">
            <button @click="logout">Logout</button>
          </template>
          <template v-else>
            <router-link to="/login" @click="closeMobileMenu">Login</router-link>
          </template>
        </div>
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
  position: relative;
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
  font-size: 2.2rem;
  line-height: 1.1;
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
.burger {
  display: none;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 40px;
  height: 40px;
  background: none;
  border: none;
  cursor: pointer;
  margin-left: 1rem;
  z-index: 1002;
}
.burger span {
  display: block;
  width: 26px;
  height: 3px;
  margin: 4px 0;
  background: #fff;
  border-radius: 2px;
  transition: 0.3s;
}
.burger span.open:nth-child(1) {
  transform: translateY(7px) rotate(45deg);
}
.burger span.open:nth-child(2) {
  opacity: 0;
}
.burger span.open:nth-child(3) {
  transform: translateY(-7px) rotate(-45deg);
}
.mobile-menu {
  display: flex;
  flex-direction: column;
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #11181c;
  padding: 1.5rem 2rem 2rem 2rem;
  z-index: 3000;
  box-shadow: 0 4px 16px rgba(0,0,0,0.15);
  animation: fadeIn 0.2s;
}
.mobile-menu a {
  color: #fff;
  text-decoration: none;
  font-weight: 500;
  padding: 0.75rem 0;
  border-radius: 3px;
  font-size: 1.1rem;
}
.mobile-menu .logo {
  font-size: 1.2rem;
  font-weight: 700;
  margin-bottom: 1rem;
}
.mobile-actions {
  margin-top: 1.5rem;
}
.mobile-actions button {
  width: 100%;
  background: #fff;
  color: #11181c;
  border: none;
  padding: 0.75rem 0;
  border-radius: 4px;
  font-weight: 600;
  cursor: pointer;
  font-size: 1.1rem;
}
.mobile-logo {
  display: none;
}
@media (max-width: 768px) {
  .nav-links,
  .nav-actions {
    display: none;
  }
  .burger {
    display: flex;
  }
  .mobile-logo {
    display: block;
    color: #fff;
    font-size: 2.2rem;
    line-height: 1.1;
    margin-right: 1rem;
    margin-left: 0.5rem;
    user-select: none;
    align-self: center;
  }
}
@media (min-width: 769px) {
  .mobile-menu {
    display: none !important;
  }
  .mobile-logo {
    display: none;
  }
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
