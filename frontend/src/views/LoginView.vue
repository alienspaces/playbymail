<template>
  <div class="login-container">
    <h2>Sign in with Email</h2>
    <form @submit.prevent="onSubmit">
      <label for="email">Email address</label>
      <input v-model="email" id="email" type="email" required autofocus />
      <button type="submit" :disabled="loading">Send Code</button>
    </form>
    <p v-if="message" class="message">{{ message }}</p>
  </div>
</template>

<script>
import { requestAuth } from '../api/auth';

export default {
  name: 'LoginView',
  data() {
    return {
      email: '',
      loading: false,
      message: '',
    };
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
  padding: 2rem;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.07);
  display: flex;
  flex-direction: column;
  align-items: stretch;
}
h2 {
  margin-bottom: 1.5rem;
  text-align: center;
}
label {
  margin-bottom: 0.5rem;
  font-weight: 500;
}
input {
  padding: 0.5rem;
  margin-bottom: 1rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 1rem;
}
button {
  background: #11181c;
  color: #fff;
  border: none;
  padding: 0.75rem;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  font-weight: 600;
}
button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.message {
  color: #d33;
  margin-top: 1rem;
  text-align: center;
}
</style> 