<template>
  <div class="login-container card">
    <h1 class="hand-drawn-title">
      <HandDrawnIcon type="shield" color="blue" class="title-icon" />
      Sign in with Email
    </h1>
    <form @submit.prevent="onSubmit" class="login-form">
      <div class="form-group">
        <label for="email">Email address</label>
        <input v-model="email" id="email" type="email" required autofocus autocomplete="off" />
      </div>
      <div class="form-actions">
        <button type="submit" :disabled="loading">Send Code</button>
      </div>
    </form>
    <p v-if="message" class="error">{{ message }}</p>
  </div>
</template>

<script>
import { requestAuth } from '../api/auth';
import HandDrawnIcon from '../components/HandDrawnIcon.vue';

const codeToMessage = {
  session_expired: 'Session expired. Please log in again.',
  // Add more codes as needed
};

export default {
  name: 'LoginView',
  data() {
    return {
      email: '',
      loading: false,
      message: '',
    };
  },
  created() {
    const code = this.$route.query.code;
    if (code && codeToMessage[code]) {
      this.message = codeToMessage[code];
    }
  },
  methods: {
    async onSubmit() {
      this.loading = true;
      this.message = '';
      try {
        const ok = await requestAuth(this.email);
        if (ok) {
          this.$router.push({ path: '/verify', query: { email: this.email } });
        } else {
          this.message = 'Failed to send verification email.';
        }
      } catch {
        this.message = 'Failed to send verification email.';
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>

<style scoped>
.login-container {
  max-width: 400px;
  margin: 80px auto;
  display: flex;
  flex-direction: column;
  align-items: stretch;
}
h1 {
  margin-bottom: var(--space-lg);
  font-size: var(--font-size-xl);
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  position: relative;
  color: var(--color-text);
}

.hand-drawn-title {
  position: relative;
}

.title-icon {
  font-size: 0.8em;
  margin-right: 0.2em;
}
button {
  background: transparent;
  color: var(--color-button);
  border: 2px solid var(--color-button);
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-base);
  cursor: pointer;
  font-weight: var(--font-weight-bold);
  transition: all 0.2s;
  box-shadow: 
    0 1px 2px rgba(0,0,0,0.05),
    inset 0 1px 0 rgba(255,255,255,0.1);
}
button:hover,
button:focus {
  background: var(--color-button);
  color: var(--color-text-light);
}
button:disabled {
  background: var(--color-border);
  color: #aaa;
  cursor: not-allowed;
  opacity: 0.6;
}
.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.error {
  color: var(--color-warning-dark);
  background: var(--color-warning-light);
  padding: var(--space-sm) var(--space-md);
  margin-top: var(--space-md);
  text-align: center;
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-warning);
}
</style> 