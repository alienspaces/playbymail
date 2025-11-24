<template>
  <div class="coming-soon-banner">
    <img :src="envelopeImage" alt="Envelope" class="envelope-icon" />
    <div class="coming-soon-text">
      <div class="coming-soon-title">Coming Soon</div>
      <div class="coming-soon-date">Early 2026</div>
    </div>
  </div>
</template>

<script setup>
import envelopeImage from '../assets/breakout-envelope-v2.png';
</script>

<style scoped>
.coming-soon-banner {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  pointer-events: none;
}

.envelope-icon {
  height: 140px;
  width: auto;
  display: block;
  flex-shrink: 0;
  align-self: center;
  border-radius: var(--radius-md);
}

.coming-soon-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  justify-content: center;
  align-self: center;
}

.coming-soon-title {
  color: var(--color-text-light);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-normal);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  line-height: 1.2;
  white-space: nowrap;
  opacity: 0.8;
}

.coming-soon-date {
  color: var(--color-text-light);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  line-height: 1.2;
  white-space: nowrap;
  opacity: 0.95;
}

/* Mobile responsive adjustments */
@media (max-width: 768px) {
  .coming-soon-banner {
    pointer-events: none;
    flex-direction: column;
    gap: var(--space-sm);
  }
  
  .envelope-icon {
    height: 100px;
  }

  .coming-soon-text {
    text-align: center;
    gap: 2px;
  }

  .coming-soon-title {
    font-size: var(--font-size-sm);
  }

  .coming-soon-date {
    font-size: var(--font-size-xs);
  }
}

/* Hide banner on very small screens to prevent crowding */
@media (max-width: 360px) {
  .coming-soon-banner {
    display: none;
  }
}
</style>
