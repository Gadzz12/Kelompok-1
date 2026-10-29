package database

import (
	"os"
	"testing"
)

func TestRejectInvalidDatabaseURL(t *testing.T) {
	for _, dsn := range []string{
		"", "gin-auth.db", "%", "https://example.com/db",
		"postgresql://user:secret@localhost/db",
		"postgresql://user:secret@localhost/db?sslmode=disable",
		"postgresql://user:secret@localhost/db?sslmode=prefer",
		"postgresql://user:secret@/db?sslmode=require",
	} {
		if _, err := Open(dsn); err == nil {
			t.Fatal("URL tidak valid atau tanpa TLS harus ditolak")
		}
	}
}

func TestIsProduction(t *testing.T) {
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("VERCEL")
		os.Unsetenv("VERCEL_ENV")
		os.Unsetenv("ENV")
	}()

	os.Unsetenv("APP_ENV")
	os.Unsetenv("VERCEL")
	os.Unsetenv("VERCEL_ENV")
	os.Unsetenv("ENV")

	if IsProduction() {
		t.Fatal("Seharusnya bernilai false ketika bukan production")
	}

	os.Setenv("APP_ENV", "production")
	if !IsProduction() {
		t.Fatal("Seharusnya bernilai true saat APP_ENV=production")
	}

	os.Setenv("APP_ENV", "development")
	os.Setenv("VERCEL", "1")
	if !IsProduction() {
		t.Fatal("Seharusnya bernilai true saat VERCEL=1")
	}
}

func TestAutoMigrateSkippedInProduction(t *testing.T) {
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("AUTO_MIGRATE")
		os.Unsetenv("RUN_MIGRATIONS")
	}()

	os.Setenv("APP_ENV", "production")
	os.Unsetenv("AUTO_MIGRATE")
	os.Unsetenv("RUN_MIGRATIONS")

	// Di production, AutoMigrate dilewati sehingga nil db tidak akan menyebabkan panic
	if err := AutoMigrate(nil); err != nil {
		t.Fatalf("AutoMigrate seharusnya dilewati di production tanpa error: %v", err)
	}

	// Di non-production tapi AUTO_MIGRATE=false
	os.Setenv("APP_ENV", "development")
	os.Setenv("AUTO_MIGRATE", "false")
	if err := AutoMigrate(nil); err != nil {
		t.Fatalf("AutoMigrate seharusnya dilewati saat AUTO_MIGRATE=false tanpa error: %v", err)
	}
}

func TestSeedDefaultDataSkippedInProduction(t *testing.T) {
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("SEED_DATA")
	}()

	os.Setenv("APP_ENV", "production")
	os.Unsetenv("SEED_DATA")

	// Di production, SeedDefaultData dilewati sehingga nil db tidak akan menyebabkan panic
	if err := SeedDefaultData(nil); err != nil {
		t.Fatalf("SeedDefaultData seharusnya dilewati di production tanpa error: %v", err)
	}
}
