package config

import (
	"os"
	"testing"
)

func clearEnv() {
	_ = os.Unsetenv("APP_ENV")
	_ = os.Unsetenv("HTTP_PORT")
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("MONGODB_URI")
	_ = os.Unsetenv("MONGODB_DATABASE")
	_ = os.Unsetenv("CORS_ALLOWED_ORIGINS")
	_ = os.Unsetenv("LOG_LEVEL")
}

func TestConfigLoadDefaults(t *testing.T) {
	clearEnv()
	_ = os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	defer clearEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("esperado sin error, obtenido: %v", err)
	}

	if cfg.AppEnv != "development" {
		t.Errorf("AppEnv esperado 'development', obtenido '%s'", cfg.AppEnv)
	}
	if cfg.HTTPPort != "8080" {
		t.Errorf("HTTPPort esperado '8080', obtenido '%s'", cfg.HTTPPort)
	}
	if cfg.MongoDBDatabase != "inverfit_dev" {
		t.Errorf("MongoDBDatabase esperado 'inverfit_dev', obtenido '%s'", cfg.MongoDBDatabase)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel esperado 'info', obtenido '%s'", cfg.LogLevel)
	}
	if len(cfg.CORSAllowedOrigins) != 1 || cfg.CORSAllowedOrigins[0] != "http://localhost:5173" {
		t.Errorf("CORSAllowedOrigins esperado ['http://localhost:5173'], obtenido %v", cfg.CORSAllowedOrigins)
	}
}

func TestConfigMissingMongoDBURI(t *testing.T) {
	clearEnv()
	defer clearEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("esperaba error por MONGODB_URI ausente, pero fue nil")
	}
}

func TestConfigInvalidPort(t *testing.T) {
	clearEnv()
	_ = os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	_ = os.Setenv("HTTP_PORT", "999999")
	defer clearEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("esperaba error por HTTP_PORT fuera de rango, pero fue nil")
	}

	_ = os.Setenv("HTTP_PORT", "abc")
	_, err = Load()
	if err == nil {
		t.Fatal("esperaba error por HTTP_PORT no numérico, pero fue nil")
	}
}

func TestConfigInvalidAppEnv(t *testing.T) {
	clearEnv()
	_ = os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	_ = os.Setenv("APP_ENV", "staging_invalido")
	defer clearEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("esperaba error por APP_ENV inválido, pero fue nil")
	}
}

func TestConfigInvalidLogLevel(t *testing.T) {
	clearEnv()
	_ = os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	_ = os.Setenv("LOG_LEVEL", "superverbose")
	defer clearEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("esperaba error por LOG_LEVEL inválido, pero fue nil")
	}
}

func TestConfigCORSOrigins(t *testing.T) {
	clearEnv()
	_ = os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	_ = os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173, https://inverfit.pages.dev , https://admin.inverfit.cl")
	defer clearEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("esperaba carga exitosa, obtenido: %v", err)
	}

	if len(cfg.CORSAllowedOrigins) != 3 {
		t.Fatalf("esperaba 3 orígenes, obtenido %d", len(cfg.CORSAllowedOrigins))
	}
	if cfg.CORSAllowedOrigins[1] != "https://inverfit.pages.dev" {
		t.Errorf("origen 1 esperado 'https://inverfit.pages.dev', obtenido '%s'", cfg.CORSAllowedOrigins[1])
	}
}

func TestRedactedMongoURI(t *testing.T) {
	cfg := &Config{
		MongoDBURI: "mongodb://admin:supersecretpassword@cluster0.abcde.mongodb.net/inverfit_prod?retryWrites=true&w=majority",
	}
	redacted := cfg.RedactedMongoURI()

	if redacted == "" {
		t.Fatal("redacted URI no debe ser vacía")
	}
	if stringContains(redacted, "supersecretpassword") {
		t.Errorf("la URI redactada no debe contener la contraseña: %s", redacted)
	}
	if !stringContains(redacted, "cluster0.abcde.mongodb.net") {
		t.Errorf("la URI redactada debe conservar el host: %s", redacted)
	}
}

func stringContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && searchSubstring(s, substr)))
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
