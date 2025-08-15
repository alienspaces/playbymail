<template>
  <div class="login-container card">
    <h2>Sign in with Email</h2>
    <form @submit.prevent="onSubmit" class="login-form">
      <div class="form-group">
        <label for="email">Email address</label>
        <input v-model="email" id="email" type="email" required autofocus autocomplete="off" />
      </div>
      <div class="form-actions">
        <button type="submit" :disabled="loading">Send Code</button>
      </div>
    </form>
    <p v-if="message" class="message">{{ message }}</p>
  </div>
</template>

<script>
import { requestAuth } from '../api/auth';

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
h2 {
  margin-bottom: var(--space-lg);
  text-align: center;
}
button {
  background: #11181c; /* Keep specific dark color for login */
  color: var(--color-text-light);
  border: none;
  padding: var(--space-md);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-base);
  cursor: pointer;
  font-weight: var(--font-weight-bold);
}
button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}
</style> 