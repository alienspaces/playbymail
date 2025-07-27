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
        <router-link to="/faq" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">F.A.Q.</router-link>
        <router-link to="/studio" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">Game Designer Studio</router-link>
        <router-link to="/admin" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">Game Management</router-link>
      </div>
      <div class="nav-actions">
        <template v-if="sessionToken">
          <router-link to="/account" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">Account</router-link>
          <button @click="logout">Logout</button>
        </template>
        <template v-else>
          <router-link to="/login" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">Login</router-link>
        </template>
      </div>
      <div class="mobile-logo">PlayByMail</div>
      <button class="burger icon-btn" @click="toggleMobileMenu" aria-label="Open navigation menu">
        <span :class="{ 'open': mobileMenuOpen }"></span>
        <span :class="{ 'open': mobileMenuOpen }"></span>
        <span :class="{ 'open': mobileMenuOpen }"></span>
      </button>
      <div class="mobile-menu" v-if="mobileMenuOpen">
        <router-link to="/faq" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">F.A.Q.</router-link>
        <router-link to="/studio" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">Game Designer Studio</router-link>
        <router-link to="/admin" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">Game Management</router-link>
        <div class="mobile-actions">
          <template v-if="sessionToken">
            <router-link to="/account" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">Account</router-link>
            <button @click="logout">Logout</button>
          </template>
          <template v-else>
            <router-link to="/login" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">Login</router-link>
          </template>
        </div>
      </div>
    </nav>
    <router-view />
  </div>
</template>

<style>
#app {
  min-height: 100vh;
}
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #11181c; /* Keep this specific dark color for navbar */
  color: var(--color-text-light);
  padding: var(--space-md) var(--space-lg);
  position: relative;
}
.nav-links {
  display: flex;
  gap: var(--space-lg);
  align-items: center;
}
.logo {
  font-weight: var(--font-weight-bold);
  font-size: var(--font-size-xl);
  line-height: 1.1;
}
.nav-actions {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}
.nav-actions button {
  background: var(--color-text-light);
  color: #11181c; /* Keep this specific dark color */
  border: none;
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  font-weight: var(--font-weight-bold);
  cursor: pointer;
}
.nav-actions a {
  color: var(--color-text-light);
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
  margin-left: var(--space-md);
  z-index: 1002;
}
.burger span {
  display: block;
  width: 26px;
  height: 3px;
  margin: 4px 0;
  background: var(--color-text-light);
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
  background: #11181c; /* Keep this specific dark color */
  padding: var(--space-lg) var(--space-lg) var(--space-lg) var(--space-lg);
  z-index: 3000;
  box-shadow: 0 4px 16px rgba(0,0,0,0.15);
  animation: fadeIn 0.2s;
}
.mobile-menu a {
  color: var(--color-text-light);
  text-decoration: none;
  font-weight: 500;
  padding: var(--space-md) 0;
  border-radius: 3px;
  font-size: var(--font-size-md);
}
.mobile-menu .logo {
  font-size: var(--font-size-md);
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--space-md);
}
.mobile-actions {
  margin-top: var(--space-lg);
}
.mobile-actions button {
  width: 100%;
  background: var(--color-text-light);
  color: #11181c; /* Keep this specific dark color */
  border: none;
  padding: var(--space-md) 0;
  border-radius: var(--radius-sm);
  font-weight: var(--font-weight-bold);
  cursor: pointer;
  font-size: var(--font-size-md);
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
    color: var(--color-text-light);
    font-size: var(--font-size-xl);
    line-height: 1.1;
    margin-right: var(--space-md);
    margin-left: var(--space-sm);
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
