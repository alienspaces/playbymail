<template>
  <div class="coming-soon-banner">
    <div class="banner-content">
      <span class="banner-text">Coming Soon</span>
      <span class="banner-date">Late 2025</span>
    </div>
  </div>
</template>

<script setup>
// Simple static banner
</script>

<style scoped>
.coming-soon-banner {
  position: absolute;
  top: 72px; /* desktop offset */
  left: 20px;
  z-index: 9999;
  pointer-events: none;
}

.banner-content {
  background: linear-gradient(135deg, #ffb3b3, #ffcc99); /* softer pastel colors */
  color: #8b4513; /* darker brown text for better contrast */
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 6px 12px 8px 12px;
  border-radius: 6px;
  border: 2px solid rgba(255, 255, 255, 0.4);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.6);
  white-space: nowrap;
}

.banner-text {
  font-size: 13px;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-right: 8px;
}

.banner-date {
  font-size: 12px;
  font-weight: 700;
  opacity: 0.9;
}

@media (max-width: 768px) {
  .coming-soon-banner {
    top: 73px; /* mobile offset */
    left: 15px;
  }
  
  .banner-content {
    padding: 6px 12px;
  }
  
  .banner-text { 
    font-size: 12px; 
    margin-right: 6px;
  }
  
  .banner-date { 
    font-size: 11px; 
  }
}
</style>
