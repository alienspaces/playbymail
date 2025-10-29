<script setup>
// Build info injected by Vite during build
const buildInfo = {
  commitRef: import.meta.env.VITE_COMMIT_REF || 'dev',
  buildDate: import.meta.env.VITE_BUILD_DATE || new Date().toISOString(),
  buildTime: import.meta.env.VITE_BUILD_TIME || new Date().toISOString()
}

// Format the build date for display
const formatBuildDate = (dateString) => {
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      timeZoneName: 'short'
    })
  } catch {
    return 'Unknown'
  }
}
</script>

<template>
  <div class="build-info-panel">
    <div class="build-info-content">
      <span class="build-commit">{{ buildInfo.commitRef }}</span>
      <span class="build-separator">•</span>
      <span class="build-date">{{ formatBuildDate(buildInfo.buildDate) }}</span>
      <span class="build-separator">•</span>
      <a href="https://github.com/alienspaces/playbymail/blob/main/RELEASE_NOTES.md" class="release-notes-link" target="_blank" rel="noopener noreferrer">Release Notes</a>
    </div>
  </div>
</template>

<style scoped>
.build-info-panel {
  display: inline-block;
  background: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: var(--space-xs) var(--space-sm);
  margin: 0 auto;
}

.build-info-content {
  display: flex;
  align-items: center;
  gap: var(--space-xs);
  font-size: 0.75rem;
  color: var(--color-text-muted);
  opacity: 0.6;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.build-commit {
  color: var(--color-primary);
  font-weight: 500;
}

.build-separator {
  opacity: 0.4;
}

.build-date {
  white-space: nowrap;
}

.release-notes-link {
  color: var(--color-text-muted);
  text-decoration: underline;
  opacity: 0.8;
  transition: opacity 0.2s;
}

.release-notes-link:hover {
  opacity: 1;
  color: var(--color-primary);
}

@media (max-width: 768px) {
  .build-info-content {
    flex-direction: column;
    gap: 2px;
    text-align: center;
  }
  
  .build-separator {
    display: none;
  }
}
</style>