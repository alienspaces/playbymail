<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { useAuthStore } from './stores/auth';
import { storeToRefs } from 'pinia';
import { ref } from 'vue';
import BuildInfo from './components/BuildInfo.vue';
import ComingSoonBanner from './components/ComingSoonBanner.vue';

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
        <router-link to="/" class="logo">
          <svg class="logo-icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
          </svg>
          <span class="logo-text">
            <span class="logo-capital">P</span>lay<span class="logo-capital">B</span>y<span class="logo-capital">M</span>ail
          </span>
        </router-link>
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
      <router-link to="/" class="mobile-logo">
        <span class="logo-text">
          <span class="logo-capital">P</span>lay<span class="logo-capital">B</span>y<span class="logo-capital">M</span>ail
        </span>
      </router-link>
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
          <footer>
        <ComingSoonBanner class="footer-banner" />
        <div class="footer-center">
          <BuildInfo />
          <p>&copy; 2025 PlayByMail. All rights reserved.</p>
          <p>
            <a href="mailto:support@playbymail.games" class="footer-link">support@playbymail.games</a>
          </p>
        </div>
      </footer>
  </div>
</template>

<style>
#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #1F1B3D; /* Deep purple-black for navbar */
  color: var(--color-text-light);
  padding: var(--space-md) var(--space-lg);
  position: relative;
}
.nav-links {
  display: flex;
  align-items: center;
  gap: var(--space-lg);
}

.nav-links .coming-soon-banner {
  margin: 0 var(--space-md);
}

.navbar-link {
  color: var(--color-text-light);
  text-decoration: underline;
  font-weight: 500;
}
.logo {
  font-family: var(--font-family-heading);
  font-weight: var(--font-weight-bold);
  font-size: var(--font-size-xl);
  line-height: 1.1;
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  color: var(--color-accent);
  letter-spacing: 0.5px;
  text-decoration: none !important;
}
.logo:hover,
.logo:focus,
.logo:active,
.logo:visited {
  text-decoration: none !important;
  color: var(--color-accent);
}

.logo-icon {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
}

.logo-text {
  color: var(--color-accent-light);
}

.logo-capital {
  color: var(--color-accent);
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
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  font-weight: var(--font-weight-bold);
  cursor: pointer;
  transition: background 0.2s;
}
.nav-actions button:hover {
  background: var(--color-accent-dark);
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
  background: #1F1B3D; /* Deep purple-black */
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
    color: var(--color-text-light);
    font-size: var(--font-size-lg); /* Reduced from xl to lg */
    line-height: 1.1;
    margin-right: var(--space-sm); /* Reduced from md to sm */
    margin-left: var(--space-sm);
    user-select: none;
    align-self: center;
    flex-shrink: 1; /* Allow logo to shrink if needed */
    text-decoration: none !important;
  }
  .mobile-logo:hover,
  .mobile-logo:focus,
  .mobile-logo:active,
  .mobile-logo:visited {
    text-decoration: none !important;
    color: var(--color-text-light);
  }
  

  


  /* Ensure navbar has proper spacing */
  .navbar {
    padding: var(--space-md) var(--space-sm); /* Reduced horizontal padding on mobile */
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



/* Footer styling */
footer {
  margin-top: 0;
  padding: var(--space-sm) var(--space-lg);
  background: var(--color-background-soft);
  border-top: 1px solid var(--color-border);
  position: relative;
}

.footer-banner {
  position: absolute;
  left: var(--space-lg);
  top: 50%;
  transform: translateY(-50%);
}

.footer-center {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-sm);
  text-align: center;
  padding: var(--space-md) 0;
}

footer p {
  margin: 0;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.footer-link {
  color: var(--color-text-muted);
  text-decoration: underline;
}

.footer-link:hover {
  color: var(--color-text);
}

/* Mobile footer adjustments */
@media (max-width: 768px) {
  footer {
    padding: var(--space-xs) var(--space-md);
  }
  
  .footer-banner {
    position: relative;
    left: auto;
    top: auto;
    transform: none;
    margin-bottom: var(--space-sm);
  }
  
  .footer-center {
    padding: 0;
  }
  
  footer {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
}
</style>
