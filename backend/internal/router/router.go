package router

import (
	"log/slog"
	"net/http"

	"inverfit-backend/internal/config"
	"inverfit-backend/internal/database"
	"inverfit-backend/internal/handlers"
	"inverfit-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Setup configura e inicializa el router de Gin con middlewares, rutas y manejadores de error uniformes.
func Setup(cfg *config.Config, db database.Pinger, logger *slog.Logger) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Middlewares globales obligatorios
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Logging(logger))
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))

	// Manejo uniforme de rutas inexistentes (404)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, handlers.ErrorResponse{
			Error: handlers.ErrorDetail{
				Code:    "NOT_FOUND",
				Message: "Ruta no encontrada",
			},
		})
	})

	// Manejo uniforme de métodos no permitidos (405)
	r.HandleMethodNotAllowed = true
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, handlers.ErrorResponse{
			Error: handlers.ErrorDetail{
				Code:    "METHOD_NOT_ALLOWED",
				Message: "Método HTTP no permitido",
			},
		})
	})

	// Inicializar Handlers
	healthHandler := handlers.NewHealthHandler(db)

	// Grupo de endpoints de salud y diagnóstico
	healthGroup := r.Group("/health")
	{
		healthGroup.GET("/live", healthHandler.Live)
		healthGroup.GET("/ready", healthHandler.Ready)
	}

	return r
}
