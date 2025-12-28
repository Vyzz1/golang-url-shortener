package api

import (
	middleware "url-shortener/middlewares"
	"url-shortener/utils"

	db "url-shortener/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	config *utils.Config
	store  db.Store
}

func NewServer(config *utils.Config, store db.Store) *Server {
	server := &Server{
		config: config,
		store:  store,
	}
	server.setupRouter()
	return server
}

func (s *Server) setupRouter() {
	s.router = gin.Default()

	s.router.Use(middleware.CORS(s.config.FrontendURL))
	s.router.Use(middleware.RateLimit())

	s.router.GET("/:short_code", s.RedirectToLongUrl)

	s.router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "UP",
		})
	})
	apiRoutes := s.router.Group("/api")

	apiRoutes.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "URL Shortener API is running!",
		})
	})

	apiRoutes.GET("/url", s.GetListUrls)

	apiRoutes.POST("/url/shorten", s.CreateUrl)

	apiRoutes.GET("/url/:url_id/stats", s.GetUrlStats)
	apiRoutes.GET("/url/:url_id/stats/count", s.GetUrlClickCount)

	apiRoutes.GET("/metrics", s.GetMetrics)

}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
