package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed web/templates/*.gohtml
var templateFS embed.FS

//go:embed web/static/*
var staticFS embed.FS

type Tool struct {
	ID          string
	Name        string
	Description string
	URL         string
	TorURL      string
	Icon        template.HTML
	Action      string
	TorAction   string
}

type PageData struct {
	Tools []Tool
}

var tools = []Tool{
	{
		ID:          "redlib",
		Name:        "Redlib",
		Description: "A private gateway to Reddit's realms, free from trackers and the burden of accounts.",
		URL:         "", // set clearnet URL when deployed, e.g. https://redlib.example.com
		TorURL:      "", // set .onion URL when available
		Action:      "Enter the Archive",
		TorAction:   "Tor Access",
	},
	{
		ID:          "wikiless",
		Name:        "Wikiless",
		Description: "The great library of knowledge, veiled from prying eyes. Wikipedia, liberated.",
		URL:         "",
		TorURL:      "",
		Action:      "Seek Knowledge",
		TorAction:   "Tor Access",
	},
	{
		ID:          "searxng",
		Name:        "SearXNG",
		Description: "A metasearch oracle that queries many sources while shielding your identity from all.",
		URL:         "",
		TorURL:      "", // e.g. http://searx....onion
		Action:      "Begin Search",
		TorAction:   "Tor Access",
	},
	{
		ID:          "whoogle",
		Name:        "Whoogle",
		Description: "Google's vast knowledge, stripped of its surveillance. Results without the watching.",
		URL:         "",
		TorURL:      "",
		Action:      "Search Freely",
		TorAction:   "Tor Access",
	},
}

var tmpl *template.Template

func main() {
	var err error
	tmpl, err = template.ParseFS(templateFS, "web/templates/*.gohtml")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	staticContent, err := fs.Sub(staticFS, "web/static")
	if err != nil {
		log.Fatalf("failed to create static fs: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	staticHandler := http.StripPrefix("/static/", cacheHeaders(http.FileServer(http.FS(staticContent))))
	mux.Handle("GET /static/", staticHandler)

	mux.HandleFunc("GET /{$}", handleIndex)

	mux.HandleFunc("GET /favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.ServeFileFS(w, r, staticContent, "favicon.svg")
	})

	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, staticContent, "robots.txt")
	})

	mux.HandleFunc("GET /sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		http.ServeFileFS(w, r, staticContent, "sitemap.xml")
	})

	addr := ":" + port()
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, securityHeaders(mux)); err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	toolsWithIcons := make([]Tool, len(tools))
	copy(toolsWithIcons, tools)
	for i := range toolsWithIcons {
		toolsWithIcons[i].Icon = getToolIcon(toolsWithIcons[i].ID)
	}

	data := PageData{
		Tools: toolsWithIcons,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, must-revalidate")

	if err := tmpl.ExecuteTemplate(w, "index.gohtml", data); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csp := "default-src 'self'; " +
			"style-src 'self'; " +
			"font-src 'self'; " +
			"script-src 'self' 'sha256-E2R/YbsbrAICx7FLw7KAday2OnDuTO8Sm8GKrPiZJ+g='; " +
			"img-src 'self' data:; " +
			"base-uri 'self'; " +
			"form-action 'self'"
		w.Header().Set("Content-Security-Policy", csp)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

func cacheHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := strings.ToLower(filepath.Ext(r.URL.Path))
		switch ext {
		case ".woff2", ".woff", ".ttf":
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		case ".css":
			w.Header().Set("Cache-Control", "public, max-age=86400")
		case ".svg", ".png", ".ico":
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		default:
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		next.ServeHTTP(w, r)
	})
}

func getToolIcon(id string) template.HTML {
	icons := map[string]string{
		"redlib": `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
			<circle cx="12" cy="8" r="5"/>
			<path d="M12 13v3"/>
			<path d="M8 21h8"/>
			<path d="M7 8h0M17 8h0"/>
			<path d="M9 5.5c0-1 .5-2 1.5-2.5M15 5.5c0-1-.5-2-1.5-2.5"/>
		</svg>`,
		"wikiless": `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
			<path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/>
			<path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/>
			<path d="M8 7h8"/>
			<path d="M8 11h8"/>
			<path d="M8 15h5"/>
			<circle cx="16" cy="15" r="1" fill="currentColor" stroke="none"/>
		</svg>`,
		"searxng": `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
			<circle cx="11" cy="11" r="7"/>
			<path d="m21 21-4.35-4.35"/>
			<path d="M11 8v6"/>
			<path d="M8 11h6"/>
			<circle cx="11" cy="11" r="3" stroke-dasharray="2 2"/>
		</svg>`,
		"whoogle": `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
			<circle cx="12" cy="12" r="3"/>
			<path d="M12 5v2"/>
			<path d="M12 17v2"/>
			<path d="M5 12h2"/>
			<path d="M17 12h2"/>
			<circle cx="12" cy="12" r="8"/>
			<path d="M3 3l4 4"/>
			<path d="M17 17l4 4"/>
		</svg>`,
	}
	if icon, ok := icons[id]; ok {
		return template.HTML(icon)
	}
	return template.HTML(`<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/></svg>`)
}
