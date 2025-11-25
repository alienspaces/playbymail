package generator

import (
	"os"
	"os/exec"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// findChromePath finds the Chrome executable path by checking environment variables,
// PATH, and common file system locations.
func findChromePath(l logger.Logger) string {
	// Check if Chrome is available
	chromePath := os.Getenv("GOOGLE_CHROME_SHIM")
	if chromePath != "" {
		l.Info("using Chrome from GOOGLE_CHROME_SHIM chrome_path=%s", chromePath)
		return chromePath
	}

	l.Info("GOOGLE_CHROME_SHIM not set, searching for Chrome")

	// First, try to find Chrome in PATH (buildpacks often add it there)
	if path, err := exec.LookPath("chrome"); err == nil {
		l.Info("found Chrome in PATH chrome_path=%s", path)
		return path
	}
	if path, err := exec.LookPath("google-chrome"); err == nil {
		l.Info("found Chrome in PATH chrome_path=%s", path)
		return path
	}
	if path, err := exec.LookPath("chromium"); err == nil {
		l.Info("found Chrome in PATH chrome_path=%s", path)
		return path
	}

	// Try to find Chrome in common locations
	commonPaths := []string{
		"/usr/bin/google-chrome",
		"/usr/bin/chromium-browser",
		"/usr/bin/chromium",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			l.Info("found Chrome at path chrome_path=%s", path)
			return path
		}
	}

	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
