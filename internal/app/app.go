package app

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"todo-app-test/internal/config"
	"todo-app-test/internal/database"
	"todo-app-test/internal/handler"
	"todo-app-test/internal/repository"
)

// Run menjalankan seluruh siklus hidup aplikasi Todo menggunakan aset dari disk
func Run() {
	RunWithAssets(nil, nil)
}

// RunWithAssets menjalankan seluruh siklus hidup aplikasi Todo dengan dukungan embed.FS untuk hosting mandiri / cloud
func RunWithAssets(templateFS fs.FS, staticFS http.FileSystem) {
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

	var tmpl *template.Template
	if templateFS != nil {
		tmpl, err = template.ParseFS(templateFS, "templates/*.html")
		if err != nil {
			// Coba pola *.html jika fsys sudah berupa sub-filesystem dari templates/
			if tmpl2, err2 := template.ParseFS(templateFS, "*.html"); err2 == nil {
				tmpl, err = tmpl2, nil
			} else {
				log.Printf("Peringatan saat membaca template dari embed: %v. Mencoba fallback ke disk...", err)
			}
		}
	}

	if tmpl == nil {
		tmpl, err = template.ParseGlob("templates/*.html")
		if err != nil {
			// Fallback: coba dari direktori parent (jika dijalankan dari subdirektori)
			if tmpl2, err2 := template.ParseGlob("../todo-list/templates/*.html"); err2 == nil {
				tmpl, err = tmpl2, nil
			} else {
				log.Fatal("Gagal membaca template HTML:", err)
			}
		}
	}

	appHandler := handler.New(repo, tmpl, staticFS)
	mux := http.NewServeMux()
	appHandler.RegisterRoutes(mux)

	fmt.Printf("Server Todo List aktif pada port %s\n", cfg.Port)
	fmt.Printf("-> Akses Lokal     : http://localhost%s\n", cfg.Port)
	fmt.Printf("-> Akses Cloudflare: Siap di-tunnel melalui Cloudflare atau Cloud Hosting\n")
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
