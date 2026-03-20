<template>
  <div class="confirmation-card card">
    <h1>{{ title }}</h1>
    <p class="confirmation-message">{{ message }}</p>
    <a :href="linkHref" class="confirmation-link" :data-testid="linkTestId">{{ linkText }}</a>
  </div>
</template>

<script setup>
defineOptions({ name: 'ConfirmationCard' })

defineProps({
  title: {
    type: String,
    required: true,
  },
  message: {
    type: String,
    required: true,
  },
  linkText: {
    type: String,
    default: 'Browse more games',
  },
  linkHref: {
    type: String,
    default: '/games',
  },
  linkTestId: {
    type: String,
    default: 'link-browse-more',
  },
})
</script>

<style scoped>
.confirmation-card {
  max-width: 600px;
  margin: 0 auto;
  padding: var(--space-xl, 2rem);
  text-align: center;
}

.confirmation-card h1 {
  font-size: var(--font-size-xl, 1.5rem);
  margin-bottom: var(--space-md, 1rem);
}

.confirmation-message {
  margin-bottom: var(--space-xl, 2rem);
  color: var(--color-text-muted, #444);
}

.confirmation-link {
  display: inline-block;
  padding: var(--space-sm, 0.5rem) var(--space-xl, 2rem);
  background: transparent;
  color: var(--color-button, #006ecd);
  border: 2px solid var(--color-button, #006ecd);
  border-radius: var(--radius-sm, 4px);
  text-decoration: none;
  font-weight: var(--font-weight-bold, 600);
  font-size: var(--font-size-sm, 0.875rem);
  transition: all 0.2s;
}

.confirmation-link:hover {
  background: var(--color-button, #006ecd);
  color: var(--color-text-light, #fff);
}
</style>
