package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logging retorna un middleware que registra cada petición HTTP de forma estructurada.
// No registra cuerpos de solicitudes, cadenas de consulta (query strings), contraseñas, tokens ni credenciales.
func Logging(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		// Construir campos de registro estructurado (excluyendo query parameters por seguridad)
		fields := []any{
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("latency", latency),
			slog.String("client_ip", clientIP),
		}

		// Registrar errores de Gin si ocurrieron durante el ciclo de vida de la petición
		if len(c.Errors) > 0 {
			fields = append(fields, slog.String("gin_errors", c.Errors.String()))
		}

		// Determinar nivel según código HTTP
		switch {
		case status >= 500:
			logger.Error("Petición HTTP completada con error del servidor", fields...)
		case status >= 400:
			logger.Warn("Petición HTTP completada con error del cliente", fields...)
		default:
			logger.Info("Petición HTTP completada", fields...)
		}
	}
}
