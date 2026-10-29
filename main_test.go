package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmbeddedTemplatesAndStatic(t *testing.T) {
	// 1. Verifikasi bahwa template HTML ter-embed dengan sempurna
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		t.Fatalf("Gagal mem-parsing embedded templates: %v", err)
	}

	expectedTemplates := []string{
		"dashboard.html",
		"login.html",
		"register.html",
		"create_todo.html",
		"detail_todo.html",
		"edit_todo.html",
		"delete_todo.html",
		"landing.html",
		"settings.html",
	}

	for _, name := range expectedTemplates {
		if tmpl.Lookup(name) == nil {
			t.Errorf("Template embedded '%s' tidak ditemukan di templateFS", name)
		}
	}

	// 2. Verifikasi bahwa static assets (css, js) ter-embed dengan benar
	staticSub, err := fs.Sub(staticFS, "static")
	if err != nil {
		t.Fatalf("Gagal mengambil sub filesystem static: %v", err)
	}

	staticServer := http.FileServer(http.FS(staticSub))

	testFiles := []string{"/app.css", "/app.js", "/landing.css"}
	for _, file := range testFiles {
		req := httptest.NewRequest(http.MethodGet, file, nil)
		rec := httptest.NewRecorder()
		staticServer.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("File statis embedded '%s' gagal diakses, status code: %d", file, rec.Code)
		}
	}
}
