package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"todo-app-test/internal/config"
	"todo-app-test/internal/database"
	"todo-app-test/internal/handler"
	"todo-app-test/internal/repository"
)

// Run menjalankan seluruh siklus hidup aplikasi Todo
func Run() {
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL belum diatur di file .env. Aplikasi ini memerlukan database PostgreSQL (Neon / Supabase). Silakan isi DATABASE_URL.")
	}

	log.Println("Mencoba menghubungkan ke database PostgreSQL...")
	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database PostgreSQL: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Gagal melakukan AutoMigrate database: %v", err)
	}

	if err := database.SeedDefaultData(db); err != nil {
		log.Printf("Peringatan saat seeding data awal database: %v", err)
	}

	repo := repository.NewGORMRepository(db)
	log.Println("STATUS: Berhasil terhubung ke database PostgreSQL.")

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		// Fallback: coba dari direktori parent (jika dijalankan dari root proyek)
		if tmpl2, err2 := template.ParseGlob("../todo-list/templates/*.html"); err2 == nil {
			tmpl, err = tmpl2, nil
		} else {
			log.Fatal("Gagal membaca template HTML:", err)
		}
	}

	appHandler := handler.New(repo, tmpl)
	mux := http.NewServeMux()
	appHandler.RegisterRoutes(mux)

	fmt.Printf("Server berjalan di http://localhost%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
