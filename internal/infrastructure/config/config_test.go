package config

import (
	"os"
	"testing"
)

func TestLoadConfigSupportsEnvironmentOverridesWithoutFile(t *testing.T) {
	originalWorkingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error: %v", err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("Chdir() error: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalWorkingDirectory); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})
	t.Setenv("QZONE_HISTORY_DATABASE_DBNAME", "./environment.db")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}
	if cfg.Database.DBName != "./environment.db" {
		t.Fatalf("Database.DBName = %q, want %q", cfg.Database.DBName, "./environment.db")
	}
	if cfg.Database.Type != "sqlite" {
		t.Fatalf("Database.Type = %q, want default %q", cfg.Database.Type, "sqlite")
	}
}
