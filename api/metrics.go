package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MetricsResponse struct {
	TotalURLs   int64    `json:"total_urls"`
	TotalClicks int64    `json:"total_clicks"`
	URLsToday   int64    `json:"urls_created_today"`
	ClicksToday int64    `json:"clicks_today"`
	TopURLs     []TopURL `json:"top_urls"`
}

type TopURL struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	Clicks      int64  `json:"clicks"`
	TinyUrl     string `json:"tiny_url"`
}

func (s *Server) GetMetrics(ctx *gin.Context) {
	totalURLs, _ := s.store.CountURLs(ctx)
	totalClicks, _ := s.store.CountAllClicks(ctx)
	urlsToday, _ := s.store.CountURLsToday(ctx)
	clicksToday, _ := s.store.CountClicksToday(ctx)
	topURLs, _ := s.store.GetTopURLs(ctx, 10)

	topURLResponses := make([]TopURL, len(topURLs))
	for i, u := range topURLs {
		topURLResponses[i] = TopURL{
			ShortCode:   u.ShortCode,
			OriginalURL: u.OriginalUrl,
			Clicks:      u.ClickCount.Int64,
			TinyUrl:     s.config.BaseURL + "/" + u.ShortCode,
		}
	}

	ctx.JSON(http.StatusOK, MetricsResponse{
		TotalURLs:   totalURLs,
		TotalClicks: totalClicks,
		URLsToday:   urlsToday,
		ClicksToday: clicksToday,
		TopURLs:     topURLResponses,
	})
}
