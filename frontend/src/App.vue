<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { useAuthStore } from './stores/auth';
import { storeToRefs } from 'pinia';
import { ref } from 'vue';
import BuildInfo from './components/BuildInfo.vue';
import ComingSoonBanner from './components/ComingSoonBanner.vue';
import titleImage from './assets/title-logo-v4.png';

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
    <div class="navbar-container">
      <nav class="navbar">
        <div class="nav-links">
          <router-link to="/" class="logo">
            <img :src="titleImage" alt="Play By Mail" class="logo-image" />
          </router-link>
          <router-link to="/studio" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">Game
            Designer Studio</router-link>
          <router-link to="/admin" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">Game
            Management</router-link>
          <router-link to="/faq" class="navbar-link" exact-active-class="active"
            @click="closeMobileMenu">F.A.Q.</router-link>
        </div>
        <div class="nav-actions">
          <template v-if="sessionToken">
            <router-link to="/account" class="navbar-link" exact-active-class="active"
              @click="closeMobileMenu">Account</router-link>
            <button @click="logout">Logout</button>
          </template>
          <template v-else>
            <router-link to="/login" class="navbar-link" exact-active-class="active"
              @click="closeMobileMenu">Login</router-link>
          </template>
        </div>
        <router-link to="/" class="mobile-logo">
          <img :src="titleImage" alt="Play By Mail" class="logo-image" />
        </router-link>
        <button class="burger icon-btn" @click="toggleMobileMenu" aria-label="Open navigation menu">
          <span :class="{ 'open': mobileMenuOpen }"></span>
          <span :class="{ 'open': mobileMenuOpen }"></span>
          <span :class="{ 'open': mobileMenuOpen }"></span>
        </button>
        <div class="mobile-menu" v-if="mobileMenuOpen">
          <router-link to="/faq" class="navbar-link" exact-active-class="active"
            @click="closeMobileMenu">F.A.Q.</router-link>
          <router-link to="/studio" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">Game
            Designer
            Studio</router-link>
          <router-link to="/admin" class="navbar-link" exact-active-class="active" @click="closeMobileMenu">Game
            Management</router-link>
          <div class="mobile-actions">
            <template v-if="sessionToken">
              <router-link to="/account" class="navbar-link" exact-active-class="active"
                @click="closeMobileMenu">Account</router-link>
              <button @click="logout">Logout</button>
            </template>
            <template v-else>
              <router-link to="/login" class="navbar-link" exact-active-class="active"
                @click="closeMobileMenu">Login</router-link>
            </template>
          </div>
        </div>
      </nav>
    </div>
    <main class="app-main">
      <router-view />
    </main>
    <footer>
      <ComingSoonBanner class="footer-banner-floating" />
      <div class="footer-content">
        <div class="footer-center">
          <p>&copy; 2025 PlayByMail. All rights reserved.</p>
          <p>
            <a href="mailto:support@playbymail.games" class="footer-link">support@playbymail.games</a>
          </p>
        </div>
        <div class="footer-meta">
          <span class="footer-meta-label">Latest release</span>
          <BuildInfo />
        </div>
      </div>
    </footer>
  </div>
</template>

<style>
html,
body {
  height: 100%;
  min-height: 100%;
  overflow-x: hidden;
  overflow-y: auto;
}

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: visible;
  position: relative;
}

/* Ensure app fills viewport height on larger screens */
@media (min-width: 769px) {
  #app {
    height: 100vh;
    min-height: 100vh;
  }
}

/* Make main content area grow to fill space and push footer down */
.app-main {
  flex: 1 0 auto;
  display: flex;
  flex-direction: column;
  min-height: 0;
  margin-bottom: 0;
}

main.app-main>.studio-layout,
main.app-main>.account-layout,
main.app-main>.management-layout {
  flex: 1 1 auto;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

/* Prevent standalone views from stretching */
main.app-main> :not(.studio-layout):not(.account-layout):not(.management-layout) {
  flex: 0 0 auto;
}

.navbar-container {
  position: relative;
  overflow: visible;
  margin-top: var(--space-md);
  min-height: 140px;
  /* Accommodate 140px image with bleed */
  display: flex;
  align-items: center;
  /* No z-index to avoid creating stacking context that traps modals */
}

.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: transparent;
  color: var(--color-text-light);
  padding: 0 var(--space-lg);
  position: relative;
  width: 100%;
  min-height: 67.5px;
  /* Original narrower green band */
}

.navbar::before {
  content: "";
  position: absolute;
  left: 0;
  right: 0;
  top: 50%;
  height: 67.5px;
  background: var(--color-logo-teal);
  transform: translateY(-50%);
  z-index: 0;
}

.navbar>* {
  position: relative;
  z-index: 1;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.nav-links .logo {
  margin-right: var(--space-md);
  position: relative;
  z-index: 10;
}

.nav-links .coming-soon-banner {
  margin: 0 var(--space-md);
}

.navbar-link {
  color: var(--color-text-light);
  text-decoration: none;
  font-weight: 500;
  padding: var(--space-xs) var(--space-sm);
  border-radius: var(--radius-sm);
  transition: background-color 0.2s ease, color 0.2s ease;
  white-space: nowrap;
}

.nav-actions .navbar-link {
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-sm);
  padding: 6px 12px;
}

.logo {
  display: flex;
  align-items: center;
  text-decoration: none !important;
  padding: 0;
  margin: 0;
}

.logo:hover,
.logo:focus,
.logo:active,
.logo:visited {
  text-decoration: none !important;
}

.logo-image {
  height: 140px;
  width: auto;
  display: block;
  flex-shrink: 0;
  margin: -10px 0;
  padding: 0;
  border-radius: var(--radius-md);
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.nav-actions button {
  background: var(--color-accent);
  color: var(--color-text-light);
  border: none;
  padding: 6px 12px;
  border-radius: var(--radius-sm);
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-sm);
  cursor: pointer;
  transition: background-color 0.2s ease;
  white-space: nowrap;
}

.nav-actions button:hover {
  background: var(--color-accent-dark);
}

.nav-actions a {
  color: var(--color-text-light);
  text-decoration: none;
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-sm);
  padding: 6px 12px;
  border-radius: var(--radius-sm);
  transition: background-color 0.2s ease, color 0.2s ease;
  white-space: nowrap;
}

.burger {
  display: none;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 40px;
  height: 40px;
  margin-left: var(--space-md);
  z-index: 1;
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
  background: var(--color-logo-teal-light);
  /* Lighter teal matching footer */
  padding: var(--space-lg) var(--space-lg) var(--space-lg) var(--space-lg);
  z-index: 500;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
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
  background: var(--color-accent);
  color: var(--color-text-light);
  border: none;
  padding: var(--space-md) 0;
  border-radius: var(--radius-sm);
  font-weight: var(--font-weight-bold);
  cursor: pointer;
  font-size: var(--font-size-md);
  transition: background 0.2s;
}

.mobile-actions button:hover {
  background: var(--color-accent-dark);
}

.mobile-logo {
  display: none;
}

@media (max-width: 768px) {
  .nav-actions {
    display: none;
  }

  .nav-links {
    display: none;
  }

  .burger {
    display: flex;
  }

  .mobile-logo {
    display: block;
    margin: 0;
    padding: 0;
    user-select: none;
    align-self: center;
    flex-shrink: 1;
    text-decoration: none !important;
  }

  .mobile-logo:hover,
  .mobile-logo:focus,
  .mobile-logo:active,
  .mobile-logo:visited {
    text-decoration: none !important;
  }

  .mobile-logo .logo-image {
    height: 100px;
    margin: -10px 0;
    padding: 0;
    border-radius: var(--radius-md);
  }





  /* Ensure navbar container and navbar have proper spacing on mobile */
  .navbar-container {
    min-height: 100px;
    /* Accommodate 100px mobile logo with bleed */
  }

  .navbar {
    padding: var(--space-xs) var(--space-sm);
    /* Reduced padding on mobile for larger logo */
    min-height: 45px;
    /* Original narrower green band on mobile */
  }

  .navbar::before {
    height: 45px;
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
  from {
    opacity: 0;
    transform: translateY(-10px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}



/* Footer styling */
footer {
  margin-top: var(--space-xl);
  padding: var(--space-xs) 0;
  min-height: 70px;
  background: transparent;
  position: relative;
  overflow: visible;
  margin-bottom: var(--space-lg);
  z-index: 1;
  width: 100%;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  --side-bleed: max(0px, (100vw - 100%)/2);
}

footer::before {
  content: "";
  position: absolute;
  left: calc(0px - var(--side-bleed));
  right: calc(0px - var(--side-bleed));
  top: 50%;
  height: 70px;
  background: var(--color-logo-teal);
  transform: translateY(-50%);
  z-index: 0;
}

.footer-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  gap: var(--space-sm);
  flex-wrap: nowrap;
  position: relative;
  z-index: 1;
  padding: 0 var(--space-lg);
  padding-left: calc(var(--space-lg) + 170px);
}

.footer-banner-floating {
  position: absolute;
  top: 50%;
  left: var(--space-lg);
  transform: translateY(-50%);
  pointer-events: none;
  display: flex;
  align-items: center;
  gap: var(--space-xs);
  z-index: 2;
}

.footer-banner-floating :deep(.envelope-icon) {
  height: 140px;
}

.footer-banner-floating :deep(.coming-soon-title) {
  font-size: var(--font-size-xs);
}

.footer-banner-floating :deep(.coming-soon-date) {
  font-size: var(--font-size-xxs);
}

.footer-center {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-xs);
  text-align: center;
  min-width: 0;
}

footer p {
  margin: 0;
  color: var(--color-text-light);
  font-size: var(--font-size-sm);
}

.footer-link,
.footer-meta-label {
  color: var(--color-text-light);
  text-decoration: underline;
}

.footer-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: var(--space-xxs);
  text-align: right;
  min-width: 0;
}

.footer-meta-label {
  font-size: var(--font-size-xs);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  text-decoration: none;
  opacity: 0.8;
}

.footer-meta :deep(.build-info-panel) {
  margin: 0;
}

.footer-meta :deep(.build-info-content) {
  justify-content: flex-end;
}

.footer-link:hover {
  color: var(--color-logo-beige-light);
}

/* Mobile footer adjustments */
@media (max-width: 768px) {
  .app-main {
    margin-bottom: 0;
  }

  footer {
    margin-top: var(--space-lg);
    margin-bottom: 0;
    padding: 0;
    min-height: auto;
    overflow: visible;
  }

  footer {
    --side-bleed: 0;
  }

  footer::before {
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    height: 100%;
    transform: none;
    background: var(--color-logo-teal);
  }

  .footer-content {
    flex-direction: column;
    gap: var(--space-md);
    flex-wrap: nowrap;
    padding: var(--space-xl) var(--space-md) var(--space-md);
    padding-top: calc(var(--space-xl) + 60px);
    position: relative;
    z-index: 1;
    width: 100%;
    align-items: center;
  }

  .footer-banner-floating {
    position: absolute;
    top: -60px;
    left: 50%;
    transform: translateX(-50%);
    width: auto;
    justify-content: center;
    margin-bottom: 0;
    z-index: 2;
  }

  .footer-banner-floating :deep(.envelope-icon) {
    height: 120px;
  }

  .footer-banner-floating :deep(.coming-soon-title) {
    font-size: var(--font-size-sm);
  }

  .footer-banner-floating :deep(.coming-soon-date) {
    font-size: var(--font-size-xs);
  }

  .footer-center {
    padding: 0;
    width: 100%;
    gap: var(--space-xs);
  }

  footer p {
    font-size: var(--font-size-base);
    color: var(--color-text-light);
    opacity: 1;
    margin: 0;
  }

  .footer-link {
    color: var(--color-text-light);
    opacity: 1;
  }

  .footer-meta {
    width: 100%;
    align-items: center;
    text-align: center;
    gap: var(--space-xs);
  }

  .footer-meta-label {
    font-size: var(--font-size-sm);
    opacity: 1;
  }

  .footer-meta :deep(.build-info-panel) {
    align-items: center;
    text-align: center;
    font-size: var(--font-size-sm);
    opacity: 1;
  }

  .footer-meta :deep(.build-info-content) {
    justify-content: center;
  }
}
</style>
