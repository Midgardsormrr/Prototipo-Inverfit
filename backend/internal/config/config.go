package config

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Config almacena la configuración de la aplicación obtenida mediante variables de entorno.
type Config struct {
	AppEnv             string   // development | test | production
	HTTPPort           string   // Puerto HTTP (ej. "8080")
	MongoDBURI         string   // URI de conexión a MongoDB
	MongoDBDatabase    string   // Nombre de la base de datos
	CORSAllowedOrigins []string // Orígenes CORS autorizados
	LogLevel           string   // debug | info | warn | error
}

// Load carga la configuración desde el entorno y opcionalmente desde un archivo .env si existe.
// Valida los campos obligatorios y la coherencia de los valores.
func Load() (*Config, error) {
	// Carga opcional de archivo .env en desarrollo (sin sobreescribir variables ya fijadas en el sistema)
	_ = loadEnvFile(".env")
	_ = loadEnvFile("backend/.env")

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = os.Getenv("PORT")
	}
	if httpPort == "" {
		httpPort = "8080"
	}

	cfg := &Config{
		AppEnv:          getEnvOrDefault("APP_ENV", "development"),
		HTTPPort:        httpPort,
		MongoDBURI:      os.Getenv("MONGODB_URI"),
		MongoDBDatabase: getEnvOrDefault("MONGODB_DATABASE", "inverfit_dev"),
		LogLevel:        strings.ToLower(getEnvOrDefault("LOG_LEVEL", "info")),
	}

	// Parsear CORS_ALLOWED_ORIGINS
	corsRaw := getEnvOrDefault("CORS_ALLOWED_ORIGINS", "http://localhost:5173")
	cfg.CORSAllowedOrigins = parseCommaSeparated(corsRaw)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate comprueba las variables obligatorias y los formatos permitidos.
func (c *Config) Validate() error {
	// Validación de MONGODB_URI obligatoria
	if strings.TrimSpace(c.MongoDBURI) == "" {
		return fmt.Errorf("variable de entorno requerida no configurada: MONGODB_URI")
	}

	// Validación de puerto HTTP
	portNum, err := strconv.Atoi(c.HTTPPort)
	if err != nil || portNum < 1 || portNum > 65535 {
		return fmt.Errorf("HTTP_PORT inválido: debe ser un número entero entre 1 y 65535")
	}

	// Validación de APP_ENV
	switch c.AppEnv {
	case "development", "test", "production":
		// válido
	default:
		return fmt.Errorf("APP_ENV inválido (%s): debe ser development, test o production", c.AppEnv)
	}

	// Validación de LOG_LEVEL
	switch c.LogLevel {
	case "debug", "info", "warn", "error":
		// válido
	default:
		return fmt.Errorf("LOG_LEVEL inválido (%s): debe ser debug, info, warn o error", c.LogLevel)
	}

	// Validación de base de datos
	if strings.TrimSpace(c.MongoDBDatabase) == "" {
		return fmt.Errorf("MONGODB_DATABASE no puede estar vacío")
	}

	// Validación de orígenes CORS
	if len(c.CORSAllowedOrigins) == 0 {
		return fmt.Errorf("CORS_ALLOWED_ORIGINS no puede estar vacío")
	}

	return nil
}

// RedactedMongoURI retorna la URI de MongoDB ocultando contraseñas o credenciales para logging seguro.
func (c *Config) RedactedMongoURI() string {
	u, err := url.Parse(c.MongoDBURI)
	if err != nil {
		// En caso de que no sea URL estándar parseable (ej. formato custom), ofuscamos tokens con @
		parts := strings.Split(c.MongoDBURI, "@")
		if len(parts) > 1 {
			return "mongodb://[REDACTED]@" + parts[len(parts)-1]
		}
		return "[REDACTED_URI]"
	}
	return u.Redacted()
}

// getEnvOrDefault retorna el valor de una variable de entorno o su valor por defecto si está vacía.
func getEnvOrDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if strings.TrimSpace(val) == "" {
		return defaultValue
	}
	return strings.TrimSpace(val)
}

// parseCommaSeparated divide una cadena separada por comas y descarta espacios en blanco.
func parseCommaSeparated(raw string) []string {
	parts := strings.Split(raw, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// loadEnvFile lee un archivo de texto con pares CLAVE=VALOR y los exporta a os.Setenv si no están presentes.
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			// Quitar comillas externas si las tiene
			v = strings.Trim(v, `"'`)
			// Solo setear si no está fijada en el entorno previamente
			if os.Getenv(k) == "" {
				_ = os.Setenv(k, v)
			}
		}
	}
	return scanner.Err()
}
