package middleware

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(frontendUrl string) gin.HandlerFunc {

	frontendOrigin := frontendUrl

	if frontendOrigin == "" {
		log.Println("FRONTEND_URL is not set. Allowing all origins.")
		frontendOrigin = "*"
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{frontendOrigin}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}

	return cors.New(config)
}
