// middleware/cors.go
package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(frontendUrl string) gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{frontendUrl}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}

	return cors.New(config)
}
