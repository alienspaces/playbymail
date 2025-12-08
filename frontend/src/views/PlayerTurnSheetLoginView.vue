<template>
  <div class="login-container card">
    <h1 class="hand-drawn-title">
      <HandDrawnIcon type="shield" color="blue" class="title-icon" />
      Access Turn Sheet
    </h1>
    <form @submit.prevent="onSubmit" class="login-form">
      <div class="form-group">
        <label for="email">Email address</label>
        <input v-model="email" id="email" type="email" required autofocus autocomplete="email" />
      </div>
      <div class="form-group">
        <label for="token">Turn Sheet Token</label>
        <input v-model="turnSheetToken" id="token" type="text" readonly class="readonly-input" />
      </div>
      <div class="form-actions">
        <button type="submit" :disabled="loading">Verify & Access</button>
      </div>
    </form>
    <p v-if="message" class="error">{{ message }}</p>
  </div>
</template>

<script>
import { verifyGameSubscriptionToken } from '../api/player';
import { useAuthStore } from '../stores/auth';
import HandDrawnIcon from '../components/HandDrawnIcon.vue';

export default {
  name: 'PlayerTurnSheetLoginView',
  components: {
    HandDrawnIcon
  },
  data() {
    return {
      email: '',
      turnSheetToken: '',
      loading: false,
      message: '',
    };
  },
  created() {
    // Extract route parameters and auto-fill token
    const token = this.$route.params.turn_sheet_token;
    if (token) {
      this.turnSheetToken = token;
    } else {
      this.message = 'Invalid link. Missing turn sheet token.';
    }
  },
  methods: {
    async onSubmit() {
      const gameSubscriptionID = this.$route.params.game_subscription_id;
      const gameInstanceID = this.$route.params.game_instance_id;
      const token = this.$route.params.turn_sheet_token;

      if (!gameSubscriptionID || !gameInstanceID || !token) {
        this.message = 'Invalid link. Missing required parameters.';
        return;
      }

      this.loading = true;
      this.message = '';
      try {
        const sessionToken = await verifyGameSubscriptionToken(
          gameSubscriptionID,
          gameInstanceID,
          this.email,
          token
        );

        // Store session token
        const authStore = useAuthStore();
        authStore.setSessionToken(sessionToken);

        // Redirect to turn sheet viewer (or home for now)
        this.$router.push({ path: '/' });
      } catch (error) {
        this.message = error.message || 'Verification failed. Please check your email and try again.';
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
    0 1px 2px rgba(0, 0, 0, 0.05),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
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

.readonly-input {
  background-color: var(--color-border-light);
  cursor: not-allowed;
  opacity: 0.7;
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

