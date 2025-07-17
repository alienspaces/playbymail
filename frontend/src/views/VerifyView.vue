<template>
  <div class="verify-container">
    <h2>Enter your verification code</h2>
    <form @submit.prevent="onSubmit">
      <label for="code">Verification code</label>
      <input v-model="code" id="code" type="text" required autofocus />
      <button type="submit" :disabled="loading">Verify</button>
    </form>
    <p v-if="message" class="message">{{ message }}</p>
  </div>
</template>

<script>
import { verifyAuth } from '../api/auth';
import { useAuthStore } from '../stores/auth';
// import { mapActions } from 'pinia'; // Removed unused import

export default {
  name: 'VerifyView',
  data() {
    return {
      code: '',
      loading: false,
      message: '',
      email: '',
    };
  },
  created() {
    this.email = this.$route.query.email || '';
    if (!this.email) {
      this.$router.replace('/login');
    }
    this.authStore = useAuthStore();
  },
  methods: {
    async onSubmit() {
      this.loading = true;
      this.message = '';
      try {
        const sessionToken = await verifyAuth(this.email, this.code);
        this.authStore.setSessionToken(sessionToken);
        this.$router.push('/');
      } catch {
        this.message = 'Invalid code or verification failed.';
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>

<style scoped>
.verify-container {
  max-width: 400px;
  margin: 80px auto;
  padding: var(--space-lg);
  background: var(--color-bg);
  border-radius: var(--radius-md);
  box-shadow: 0 2px 8px rgba(0,0,0,0.07);
  display: flex;
  flex-direction: column;
  align-items: stretch;
}
h2 {
  margin-bottom: var(--space-lg);
  text-align: center;
}
label {
  margin-bottom: var(--space-sm);
  font-weight: 500;
}
input {
  padding: var(--space-sm);
  margin-bottom: var(--space-md);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-base);
}
button {
  background: #11181c; /* Keep specific dark color for verify */
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
.message {
  color: var(--color-error);
  margin-top: var(--space-md);
  text-align: center;
}
</style> 