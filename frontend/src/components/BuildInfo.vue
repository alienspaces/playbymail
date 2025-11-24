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
    <span class="build-commit">{{ buildInfo.commitRef }}</span>
    <span class="build-date">{{ formatBuildDate(buildInfo.buildDate) }}</span>
    <a href="https://github.com/alienspaces/playbymail/blob/main/RELEASE_NOTES.md" class="release-notes-link"
      target="_blank" rel="noopener noreferrer">Release Notes</a>
  </div>
</template>

<style scoped>
.build-info-panel {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: center;
  gap: 2px;
  font-size: 0.75rem;
  color: var(--color-text-light);
  opacity: 0.9;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  text-align: right;
}

.build-commit {
  color: var(--color-logo-beige-light);
  font-weight: 500;
}

.build-date {
  white-space: nowrap;
}

.release-notes-link {
  color: var(--color-text-light);
  text-decoration: underline;
  opacity: 0.9;
  transition: opacity 0.2s;
}

.release-notes-link:hover {
  opacity: 1;
  color: var(--color-logo-beige-light);
}

@media (max-width: 768px) {
  .build-info-panel {
    align-items: center;
    text-align: center;
  }
}
</style>