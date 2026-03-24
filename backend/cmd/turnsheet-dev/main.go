// turnsheet-dev is a development HTTP server for iterating on turn sheet HTML templates.
//
// It renders all turn sheet types using the exact same DocumentRenderer code path as
// production, serves the gallery, watches for template changes, and auto-reloads the
// browser via Server-Sent Events (SSE) — no DB, Chrome, or test harness required.
//
// The sample fixture data is sourced from turnsheet.DevFixtures(), the same data used
// by the automated rendering tests (TestRenderAllSheets).
//
// Usage:
//
//	go run ./cmd/turnsheet-dev/main.go [--port 8090] [--templates backend/templates] [--testdata backend/internal/turnsheet/testdata]
package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"

	coreconfig "gitlab.com/alienspaces/playbymail/core/config"
	corelog "gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/internal/generator"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

func main() {
	port := flag.String("port", "8090", "HTTP port to listen on")
	templatesRoot := flag.String("templates", "templates", "Root templates directory (contains turnsheet/ subdirectory)")
	testdataDir := flag.String("testdata", "internal/turnsheet/testdata", "Output directory served as static files (also contains background images)")
	flag.Parse()

	l, err := corelog.NewLogger(coreconfig.Config{})
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	renderer, err := generator.NewDocumentRenderer(l)
	if err != nil {
		log.Fatalf("failed to create document renderer: %v", err)
	}
	renderer.SetTemplatePath(*templatesRoot)

	srv := &devServer{
		renderer:      renderer,
		templatesRoot: *templatesRoot,
		testdataDir:   *testdataDir,
		clients:       make(map[chan string]struct{}),
	}

	if err := srv.renderAll(); err != nil {
		log.Fatalf("initial render failed: %v", err)
	}

	go srv.watch()

	mux := http.NewServeMux()
	mux.HandleFunc("/events", srv.handleSSE)
	mux.HandleFunc("/", srv.handleStatic)

	addr := ":" + *port
	fmt.Printf("turnsheet-dev listening on http://localhost%s\n", addr)
	fmt.Printf("Open http://localhost%s/gallery.html\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// devServer coordinates rendering, file watching, and SSE broadcasting.
type devServer struct {
	renderer      *generator.DocumentRenderer
	templatesRoot string
	testdataDir   string

	mu      sync.Mutex
	clients map[chan string]struct{}
}

// loadBackground reads a background PNG from testdata and returns a base64 data URI.
func (s *devServer) loadBackground(filename string) string {
	path := filepath.Join(s.testdataDir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("warning: could not load background image %s: %v", path, err)
		return ""
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	return "data:image/png;base64," + encoded
}

// renderAll renders all sheet fixtures to HTML files in testdataDir.
func (s *devServer) renderAll() error {
	ctx := context.Background()
	// devTurnSheetCode is a sample code used for the gallery only; it is not scannable.
	const devTurnSheetCode = "ZGV2LXR1cm4tc2hlZXQtY29kZQ=="

	for _, f := range turnsheet.DevFixtures() {
		bg := s.loadBackground(f.BackgroundFile)
		data := f.MakeData(bg, devTurnSheetCode)

		html, err := s.renderer.GenerateHTML(ctx, f.TemplatePath, data)
		if err != nil {
			return fmt.Errorf("render %s: %w", f.TemplatePath, err)
		}

		out := filepath.Join(s.testdataDir, f.OutputBaseName+".html")
		if err := os.WriteFile(out, []byte(html), 0644); err != nil {
			return fmt.Errorf("write %s: %w", out, err)
		}
		log.Printf("rendered %s.html", f.OutputBaseName)
	}
	return nil
}

// watch watches the turnsheet templates directory and re-renders on changes.
func (s *devServer) watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("could not create watcher: %v", err)
		return
	}
	defer watcher.Close()

	watchDir := filepath.Join(s.templatesRoot, "turnsheet")
	if err := watcher.Add(watchDir); err != nil {
		log.Printf("could not watch %s: %v", watchDir, err)
		return
	}
	log.Printf("watching %s for changes", watchDir)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				log.Printf("template changed: %s — re-rendering all sheets", filepath.Base(event.Name))
				if err := s.renderAll(); err != nil {
					log.Printf("render error: %v", err)
				} else {
					s.broadcast("reload")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("watcher error: %v", err)
		}
	}
}

// broadcast sends a named SSE event to all connected clients.
func (s *devServer) broadcast(event string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ch := range s.clients {
		select {
		case ch <- event:
		default:
		}
	}
}

// sseScript is injected into gallery.html responses to enable auto-reload.
const sseScript = `<script>
(function() {
  var es = new EventSource('/events');
  es.addEventListener('reload', function() {
    document.querySelectorAll('iframe').forEach(function(f) { f.src = f.src; });
  });
})();
</script>`

// handleSSE serves the SSE endpoint for browser auto-reload.
func (s *devServer) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan string, 4)
	s.mu.Lock()
	s.clients[ch] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, ch)
		s.mu.Unlock()
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	// Send an initial comment to establish the connection.
	fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-ch:
			fmt.Fprintf(w, "event: %s\ndata: {}\n\n", event)
			flusher.Flush()
		}
	}
}

// handleStatic serves files from testdataDir, injecting the SSE script into gallery.html.
func (s *devServer) handleStatic(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	if urlPath == "/" {
		urlPath = "/gallery.html"
	}
	// Sanitise path to prevent directory traversal.
	clean := filepath.Clean(strings.TrimPrefix(urlPath, "/"))
	if strings.HasPrefix(clean, "..") {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	filePath := filepath.Join(s.testdataDir, clean)

	if strings.HasSuffix(urlPath, "gallery.html") {
		s.serveGallery(w, filePath)
		return
	}

	http.ServeFile(w, r, filePath)
}

// serveGallery reads gallery.html, injects the SSE reload script before </body>, and writes it.
func (s *devServer) serveGallery(w http.ResponseWriter, filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "gallery.html not found", http.StatusNotFound)
		return
	}

	html := string(data)
	if idx := strings.LastIndex(html, "</body>"); idx >= 0 {
		html = html[:idx] + sseScript + html[idx:]
	} else {
		html += sseScript
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}
