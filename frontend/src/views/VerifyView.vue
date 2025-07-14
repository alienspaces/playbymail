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
  },
  methods: {
    async onSubmit() {
      this.loading = true;
      this.message = '';
      try {
        const sessionToken = await verifyAuth(this.email, this.code);
        localStorage.setItem('session_token', sessionToken);
        this.$router.push('/');
      } catch (e) {
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