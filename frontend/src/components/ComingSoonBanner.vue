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
  position: absolute;
  top: 72px; /* desktop offset */
  left: 20px;
  z-index: 9999;
  pointer-events: none;
}

.envelope {
  position: relative;
  width: 140px; /* narrower, more envelope-like */
  height: 90px; /* taller for envelope proportions */
  background: linear-gradient(135deg, #fef3c7, #fde68a); /* subdued yellow envelope color */
  border: 2px solid #d97706; /* warm brown-orange border */
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  padding: 12px;
  box-sizing: border-box;
}

.stamp {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 24px;
  height: 28px;
  background: linear-gradient(135deg, #dc3545, #c82333);
  border-radius: 3px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.stamp-text {
  color: white;
  font-size: 8px;
  font-weight: 700;
  text-align: center;
  line-height: 1;
}

.address {
  position: absolute;
  top: 20px;
  left: 12px;
  text-align: left;
}

.address-line {
  color: #92400e; /* darker brown text that complements the yellow */
  font-family: 'Courier New', monospace;
  font-size: 11px;
  font-weight: 600;
  line-height: 1.3;
  margin-bottom: 2px;
}

.fold-line {
  position: absolute;
  background: rgba(146, 64, 14, 0.2); /* subtle brown fold lines */
  border-radius: 1px;
}

.fold-line-1 {
  top: 45px;
  left: 8px;
  right: 8px;
  height: 1px;
}

.fold-line-2 {
  top: 8px;
  bottom: 8px;
  left: 45px;
  width: 1px;
}

@media (max-width: 768px) {
  .coming-soon-banner {
    top: 73px; /* mobile offset */
    left: 15px;
  }
  
  .envelope {
    width: 120px;
    height: 80px;
    padding: 10px;
  }
  
  .stamp {
    width: 20px;
    height: 24px;
    top: 6px;
    right: 6px;
  }
  
  .stamp-text {
    font-size: 7px;
  }
  
  .address {
    top: 18px;
    left: 10px;
  }
  
  .address-line {
    font-size: 10px;
  }
  
  .fold-line-1 {
    top: 40px;
    left: 6px;
    right: 6px;
  }
  
  .fold-line-2 {
    top: 6px;
    bottom: 6px;
    left: 40px;
  }
}
</style>
