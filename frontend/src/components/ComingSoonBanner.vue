<template>
  <div class="coming-soon-banner">
    <div class="envelope">
      <!-- Stamp in top right -->
      <div class="stamp">
        <span class="stamp-text">2025</span>
      </div>
      
      <!-- Address-style text -->
      <div class="address">
        <div class="address-line">Coming Soon</div>
        <div class="address-line">Late 2025</div>
      </div>
      
      <!-- Envelope fold lines -->
      <div class="fold-line fold-line-1"></div>
      <div class="fold-line fold-line-2"></div>
    </div>
  </div>
</template>

<script setup>
// Simple static banner
</script>

<style scoped>
.coming-soon-banner {
  /* Footer banner styling */
  display: flex;
  align-items: center;
  pointer-events: none;
}

/* Alternative layout option: Top-right corner positioning
   Uncomment this and comment out the above .coming-soon-banner styles
   if you prefer the banner in the top-right corner instead of integrated in navbar

.coming-soon-banner {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
  pointer-events: none;
}
*/

.envelope {
  position: relative;
  width: 120px; /* narrower for header integration */
  height: 70px; /* shorter for header height */
  background: linear-gradient(135deg, #fef3c7, #fde68a); /* subdued yellow envelope color */
  border: 1px solid #e8c547; /* very subtle border, almost matching envelope */
  border-radius: 6px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  padding: 8px;
  box-sizing: border-box;
}

.stamp {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 20px;
  height: 24px;
  background: linear-gradient(135deg, #dc6545, #c8533f);
  border-radius: 3px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.stamp-text {
  color: white;
  font-size: 7px;
  font-weight: 700;
  text-align: center;
  line-height: 1;
}

.address {
  position: absolute;
  top: 16px;
  left: 8px;
  text-align: left;
}

.address-line {
  color: #92400e; /* darker brown text that complements the yellow */
  font-family: 'Courier New', monospace;
  font-size: 11px;
  font-weight: 600;
  line-height: 1.3;
  margin-bottom: 1px;
}

.fold-line {
  position: absolute;
  background: rgba(146, 64, 14, 0.2); /* subtle brown fold lines */
  border-radius: 1px;
}

.fold-line-1 {
  top: 35px;
  left: 6px;
  right: 6px;
  height: 1px;
}

.fold-line-2 {
  top: 6px;
  bottom: 6px;
  left: 35px;
  width: 1px;
}

/* Mobile responsive adjustments */
@media (max-width: 768px) {
  .envelope {
    width: 120px;
    height: 70px;
    padding: 4px;
  }
  
  .stamp {
    width: 20px;
    height: 24px;
    top: 2px;
    right: 2px;
  }
  
  .stamp-text {
    font-size: 7px;
  }
  
  .address {
    top: 16px;
    left: 8px;
  }
  
  .address-line {
    font-size: 11px;
  }
  
  .fold-line-1 {
    top: 35px;
    left: 6px;
    right: 6px;
  }
  
  .fold-line-2 {
    top: 6px;
    bottom: 6px;
    left: 35px;
  }
}

/* Extra small mobile adjustments */
@media (max-width: 480px) {
  .envelope {
    width: 120px;
    height: 70px;
    padding: 4px;
  }
  
  .stamp {
    width: 20px;
    height: 24px;
    top: 2px;
    right: 2px;
  }
  
  .stamp-text {
    font-size: 7px;
  }
  
  .address {
    top: 16px;
    left: 8px;
  }
  
  .address-line {
    font-size: 11px;
  }
  
  .fold-line-1 {
    top: 35px;
    left: 6px;
    right: 6px;
  }
  
  .fold-line-2 {
    top: 6px;
    bottom: 6px;
    left: 35px;
  }
}

/* Hide banner on very small screens to prevent crowding */
@media (max-width: 360px) {
  .coming-soon-banner {
    display: none;
  }
}

/* Ensure banner doesn't interfere with touch targets on mobile */
@media (max-width: 768px) {
  .coming-soon-banner {
    pointer-events: none;
  }
  
  .envelope {
    /* Ensure minimum touch target size */
    min-width: 44px;
    min-height: 44px;
  }
}
</style>
