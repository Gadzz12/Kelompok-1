package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"todo-app-test/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// IsProduction mengecek apakah aplikasi berjalan di environment production atau serverless (seperti Vercel)
func IsProduction() bool {
	appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if appEnv == "production" || appEnv == "prod" {
		return true
	}
	if os.Getenv("VERCEL") != "" || os.Getenv("VERCEL_ENV") != "" {
		return true
	}
	if strings.ToLower(strings.TrimSpace(os.Getenv("ENV"))) == "production" {
		return true
	}
	return false
}

// Open membuka koneksi PostgreSQL memakai GORM dengan validasi DSN dan TLS
func Open(dsn string) (*gorm.DB, error) {
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return nil, errors.New("DATABASE_URL kosong")
	}

	u, err := url.Parse(dsn)
	if err != nil || (u.Scheme != "postgres" && u.Scheme != "postgresql") || u.Hostname() == "" || u.User == nil || len(u.Path) < 2 {
		return nil, errors.New("DATABASE_URL wajib berupa connection string PostgreSQL Neon yang valid")
	}

	query := u.Query()
	switch query.Get("sslmode") {
	case "require", "verify-ca", "verify-full":
	default:
		return nil, errors.New("DATABASE_URL wajib memakai TLS: sslmode=require atau verify-full")
	}

	inProd := IsProduction()

	if query.Get("connect_timeout") == "" {
		if inProd {
			query.Set("connect_timeout", "10")
		} else {
			query.Set("connect_timeout", "25")
		}
	}
	u.RawQuery = query.Encode()

	// Level log GORM: silent saat di cloud/production supaya tidak menghasilkan overhead log query yang lambat
	gormLogLevel := logger.Warn
	if inProd {
		gormLogLevel = logger.Silent
	}
	if envLevel := strings.ToLower(strings.TrimSpace(os.Getenv("GORM_LOG_LEVEL"))); envLevel != "" {
		switch envLevel {
		case "silent":
			gormLogLevel = logger.Silent
		case "error":
			gormLogLevel = logger.Error
		case "warn":
			gormLogLevel = logger.Warn
		case "info":
			gormLogLevel = logger.Info
		}
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  u.String(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		TranslateError:       true,
		DisableAutomaticPing: true,
		Logger:               logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("konfigurasi koneksi PostgreSQL tidak valid: %w", err)
	}

	conn, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Konfigurasi koneksi pool database yang optimal untuk serverless (Vercel / Cloud Functions)
	if inProd {
		conn.SetMaxOpenConns(2)
		conn.SetMaxIdleConns(1)
		conn.SetConnMaxLifetime(10 * time.Minute)
		conn.SetConnMaxIdleTime(1 * time.Minute)
	} else {
		conn.SetMaxOpenConns(10)
		conn.SetMaxIdleConns(5)
		conn.SetConnMaxLifetime(30 * time.Minute)
		conn.SetConnMaxIdleTime(5 * time.Minute)
	}

	pingTimeout := 10 * time.Second
	if !inProd {
		pingTimeout = 25 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("koneksi PostgreSQL gagal; periksa DATABASE_URL dan akses jaringan: %w", err)
	}

	return db, nil
}

// AutoMigrate membuat/memperbarui tabel database secara otomatis.
// Di environment production / serverless, migrasi dilewati secara default agar tidak
// menyebabkan timeout startup (batas 30 detik di Vercel), kecuali diaktifkan via
// environment variable AUTO_MIGRATE=true atau RUN_MIGRATIONS=true.
func AutoMigrate(db *gorm.DB) error {
	autoMigrateEnv := strings.ToLower(strings.TrimSpace(os.Getenv("AUTO_MIGRATE")))
	runMigrationsEnv := strings.ToLower(strings.TrimSpace(os.Getenv("RUN_MIGRATIONS")))

	// Jika secara eksplisit dinonaktifkan
	if autoMigrateEnv == "false" || autoMigrateEnv == "0" || runMigrationsEnv == "false" || runMigrationsEnv == "0" {
		log.Println("[Database] AutoMigrate dilewati (dinonaktifkan via environment variable).")
		return nil
	}

	// Jangan jalankan AutoMigrate secara otomatis jika aplikasi berada di environment production
	if IsProduction() && autoMigrateEnv != "true" && autoMigrateEnv != "1" && runMigrationsEnv != "true" && runMigrationsEnv != "1" {
		log.Println("[Database] AutoMigrate dilewati pada environment production untuk mempercepat startup server.")
		return nil
	}

	log.Println("[Database] Menjalankan AutoMigrate skema tabel...")
	return db.AutoMigrate(
		&model.User{},
		&model.Todo{},
		&model.SubTodo{},
		&model.Collaborator{},
		&model.Invitation{},
		&model.InboxMessage{},
		&model.Friend{},
		&model.FriendRequest{},
	)
}

// SeedDefaultData membuat user awal jika tabel users di database masih kosong
func SeedDefaultData(db *gorm.DB) error {
	// Lewati seeding otomatis di production kecuali diminta secara eksplisit
	if IsProduction() && strings.ToLower(strings.TrimSpace(os.Getenv("SEED_DATA"))) != "true" {
		return nil
	}

	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}

	// Jika tabel user sudah terisi, lewati seeding
	if count > 0 {
		return nil
	}

	log.Println("Tabel database kosong. Memasukkan data awal default ke database...")
	defaultUser := model.User{
		Nama:      "abizard rahadyan",
		Username:  "Gadzz12",
		Email:     "abizardrahadyan5@gmail.com",
		Password:  "wasdwasd",
		CreatedAt: time.Now(),
	}

	return db.Create(&defaultUser).Error
}

// SeedFromJSON menyalin data dari file JSON ke database jika tabel masih kosong (opsional)
func SeedFromJSON(db *gorm.DB, jsonPath string) error {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}

	// Jika tabel user sudah terisi, lewati seeding
	if count > 0 {
		return nil
	}

	file, err := os.ReadFile(jsonPath)
	if err != nil {
		// Jika file json tidak ditemukan, lewati
		return nil
	}

	var data model.DataStore
	if err := json.Unmarshal(file, &data); err != nil {
		return fmt.Errorf("format %s tidak valid: %w", jsonPath, err)
	}

	log.Println("Tabel database kosong. Memindahkan data dari", jsonPath, "ke database PostgreSQL...")

	return db.Transaction(func(tx *gorm.DB) error {
		for _, user := range data.Users {
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
		}

		for _, todo := range data.Todos {
			cleanTodo := todo
			cleanTodo.SubTodos = nil
			cleanTodo.Collaborators = nil
			cleanTodo.CollabRecords = nil
			if err := tx.Create(&cleanTodo).Error; err != nil {
				return err
			}
		}

		for _, sub := range data.SubTodos {
			if err := tx.Create(&sub).Error; err != nil {
				return err
			}
		}

		for _, collab := range data.Collaborators {
			if err := tx.Create(&collab).Error; err != nil {
				return err
			}
		}

		for _, inv := range data.Invitations {
			if err := tx.Create(&inv).Error; err != nil {
				return err
			}
		}

		for _, msg := range data.InboxMessages {
			if err := tx.Create(&msg).Error; err != nil {
				return err
			}
		}

		log.Println("Migrasi data dari JSON ke database berhasil dilakukan.")
		return nil
	})
}
