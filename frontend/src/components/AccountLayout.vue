<!--
  AccountLayout.vue
  Layout component for the Account management interface.
  Uses the shared SidebarLayout component.
-->
<template>
  <SidebarLayout title="Account" icon-type="person" icon-color="blue">
    <!-- Sidebar navigation -->
    <template #sidebar>
      <ul>
        <li>
          <router-link to="/account" active-class="active" exact>
            <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
              <path
                d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z" />
            </svg>
            Profile
          </router-link>
        </li>
        <li>
          <router-link to="/account/contacts" active-class="active">
            <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
              <path
                d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z" />
            </svg>
            Contacts
          </router-link>
        </li>
        <li>
          <router-link to="/account/subscriptions" active-class="active">
            <svg class="nav-icon" viewBox="0 0 24 24" fill="currentColor">
              <path
                d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm-1 14H5c-.55 0-1-.45-1-1V7c0-.55.45-1 1-1h14c.55 0 1 .45 1 1v10c0 .55-.45 1-1 1z" />
            </svg>
            Subscriptions
          </router-link>
        </li>
      </ul>
    </template>

    <!-- Help content -->
    <template #help>
      <h2>Account Help</h2>
      <p>This is context-sensitive help for the account section. (Stub)</p>
    </template>

    <!-- Main content -->
    <router-view />
  </SidebarLayout>
</template>

<script setup>
import SidebarLayout from './SidebarLayout.vue';
</script>
