package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery retorna un middleware que captura panics inesperados y responde con un JSON uniforme de error 500,
// protegiendo al cliente de filtraciones de trazas internas o detalles de infraestructura.
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Registrar el panic en el log del servidor
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				logger.Error("Panic recuperado en handler HTTP",
					slog.String("path", c.Request.URL.Path),
					slog.String("method", c.Request.Method),
					slog.String("error", err.Error()),
				)

				// Respuesta de error uniforme sin exponer detalles técnicos
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{
						"code":    "INTERNAL_SERVER_ERROR",
						"message": "Ha ocurrido un error interno en el servidor",
					},
				})
			}
		}()

		c.Next()
	}
}
