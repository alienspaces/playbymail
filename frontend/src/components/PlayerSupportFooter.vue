<template>
  <footer class="player-support-footer">
    <div class="footer-content">
      <p>
        Need help? Contact us at
        <a href="mailto:support@playbymail.games">support@playbymail.games</a>
      </p>
      <p class="copyright">&copy; {{ currentYear }} PlayByMail. All rights reserved.</p>
    </div>
  </footer>
</template>

<script setup>
import { computed } from 'vue'

const currentYear = computed(() => new Date().getFullYear())
</script>

<style scoped>
.player-support-footer {
  margin-top: auto;
  padding: var(--space-lg) var(--space-md);
  background: var(--color-bg-alt);
  border-top: 1px solid var(--color-border);
  text-align: center;
}

.footer-content {
  max-width: 600px;
  margin: 0 auto;
}

.footer-content p {
  margin: var(--space-xs) 0;
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.footer-content a {
  color: var(--color-primary);
  text-decoration: none;
}

.footer-content a:hover {
  text-decoration: underline;
}

.copyright {
  font-size: var(--font-size-xs);
}
</style>
