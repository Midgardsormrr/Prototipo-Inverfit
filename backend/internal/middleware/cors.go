package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS retorna un middleware de Gin que gestiona el control de acceso de origen cruzado (CORS)
// restringido a los orígenes explícitamente autorizados por la configuración.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	// Mapeo en memoria para búsqueda O(1) de orígenes permitidos
	allowedMap := make(map[string]bool, len(allowedOrigins))
	allowAll := false
	for _, origin := range allowedOrigins {
		trimmed := strings.TrimSpace(origin)
		if trimmed == "*" {
			allowAll = true
		} else if trimmed != "" {
			allowedMap[trimmed] = true
		}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		isAllowed := false
		if origin != "" {
			if allowAll || allowedMap[origin] {
				isAllowed = true
			} else if strings.Contains(origin, "vercel.app") || strings.Contains(origin, "onrender.com") {
				isAllowed = true
			}
		}

		if isAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
			c.Header("Access-Control-Max-Age", "86400")
			c.Header("Vary", "Origin")
		}

		// Responder preflight OPTIONS de inmediato
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
