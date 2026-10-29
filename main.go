package main

import (
	"embed"
	"io/fs"
	"net/http"

	"todo-app-test/internal/app"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {
	staticSub, err := fs.Sub(staticFS, "static")
	var staticHttpFS http.FileSystem
	if err == nil {
		staticHttpFS = http.FS(staticSub)
	}
	app.RunWithAssets(templateFS, staticHttpFS)
}
