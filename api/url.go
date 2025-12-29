package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	db "url-shortener/db/sqlc"
	"url-shortener/utils"

	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateUrlRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
}

type CreateUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

func isValidURL(raw string) bool {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	if u.Host == "" {
		return false
	}

	return true
}

func (s *Server) CreateUrl(ctx *gin.Context) {
	var req CreateUrlRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	const (
		maxRetries = 5
		codeLen    = 7
	)

	if !isValidURL(req.LongUrl) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		return
	}

	var shortCode string

	for attempt := 0; attempt < maxRetries; attempt++ {
		newCode, genErr := utils.GenerateShortCode(codeLen)
		if genErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short code"})
			return
		}

		_, err := s.store.CreateURL(ctx, db.CreateURLParams{
			OriginalUrl: req.LongUrl,
			ShortCode:   newCode,
		})
		if err != nil {
			if isDuplicateKeyError(err) {
				continue
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
			return
		}

		shortCode = newCode
		break
	}

	if shortCode == "" {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Could not generate a unique short URL, please retry"})
		return
	}

	shortUrl := s.config.BaseURL + "/" + shortCode
	ctx.JSON(http.StatusOK, CreateUrlResponse{ShortUrl: shortUrl})
}

func (s *Server) RedirectToLongUrl(ctx *gin.Context) {
	shortCode := ctx.Param("short_code")

	if utils.ValidateShortCode(shortCode) == false {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid short code format"})
		return
	}

	urlRecord, err := s.store.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		return
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		clickData := extractClickData(ctx)

		_, err := s.store.InsertClick(bgCtx, db.InsertClickParams{
			UrlID:      pgtype.Int8{Int64: urlRecord.ID, Valid: true},
			IpAddress:  clickData.IpAddress,
			UserAgent:  clickData.UserAgent,
			Referer:    clickData.Referer,
			Country:    clickData.Country,
			ClickedAt:  pgtype.Timestamp{Time: time.Now(), Valid: true},
			DeviceType: clickData.DeviceType,
		})

		if err != nil {
			fmt.Println("Failed to log click data:", err)
			return
		}

		err = s.store.IncrementClickCount(bgCtx, shortCode)
		if err != nil {
			fmt.Println("Failed to increment click count:", err)
			return
		}
	}()

	ctx.Redirect(http.StatusFound, urlRecord.OriginalUrl)
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func extractClickData(ctx *gin.Context) db.InsertClickParams {
	ip := utils.GetClientIP(ctx)

	return db.InsertClickParams{
		IpAddress:  pgtype.Text{String: ip, Valid: true},
		UserAgent:  pgtype.Text{String: ctx.Request.UserAgent(), Valid: true},
		Referer:    pgtype.Text{String: ctx.Request.Referer(), Valid: true},
		Country:    pgtype.Text{String: utils.GetCountryFromIP(ip), Valid: true},
		DeviceType: pgtype.Text{String: utils.DetectDeviceType(ctx.Request.UserAgent()), Valid: true},
	}
}

type GetListUrlsRequest struct {
	Limit int32 `json:"limit" form:"limit,default=10"`
	Page  int32 `json:"page" form:"page,default=0"`
}

type UrlResponse struct {
	Id          int64            `json:"id"`
	OriginalUrl string           `json:"original_url"`
	ShortCode   string           `json:"short_code"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	ClickCount  int64            `json:"click_count"`
	TinyUrl     string           `json:"tiny_url"`
}

type GetListUrlsResponse struct {
	Content     []UrlResponse `json:"content"`
	IsLast      bool          `json:"is_last"`
	IsFirst     bool          `json:"is_first"`
	IsPrevious  bool          `json:"is_previous"`
	IsNext      bool          `json:"is_next"`
	CurrentPage int32         `json:"current_page"`
	TotalCount  int64         `json:"total_count"`
}

func (s *Server) GetListUrls(ctx *gin.Context) {
	var req GetListUrlsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	offset := req.Page * req.Limit

	urls, err := s.store.ListURLs(ctx, db.ListURLsParams{
		Limit:  req.Limit,
		Offset: offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URLs"})
		return
	}

	total, err := s.store.CountURLs(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL count"})
		return
	}

	content := make([]UrlResponse, len(urls))
	for i, u := range urls {
		content[i] = UrlResponse{
			Id:          u.ID,
			OriginalUrl: u.OriginalUrl,
			ShortCode:   u.ShortCode,
			CreatedAt:   u.CreatedAt,
			ClickCount:  u.ClickCount.Int64,
			TinyUrl:     s.config.BaseURL + "/" + u.ShortCode,
		}
	}

	isLast := offset+int32(len(urls)) >= int32(total)

	ctx.JSON(http.StatusOK, GetListUrlsResponse{
		Content:     content,
		IsFirst:     req.Page == 0,
		IsLast:      isLast,
		IsPrevious:  req.Page > 0,
		IsNext:      !isLast,
		CurrentPage: req.Page,
		TotalCount:  total,
	})

}

type GetUrlStatsRequest struct {
	Limit int32 `json:"limit" form:"limit,default=10"`
	Page  int32 `json:"page" form:"page,default=0"`
}

type GetUrlStatsResponse struct {
	Content     []UrlStats `json:"content"`
	IsLast      bool       `json:"is_last"`
	IsFirst     bool       `json:"is_first"`
	IsPrevious  bool       `json:"is_previous"`
	IsNext      bool       `json:"is_next"`
	CurrentPage int32      `json:"current_page"`
	TotalCount  int64      `json:"total_count"`
}

type UrlStats struct {
	Id         int64            `json:"id"`
	ClickedAt  pgtype.Timestamp `json:"clicked_at"`
	IpAddress  pgtype.Text      `json:"ip_address"`
	UserAgent  pgtype.Text      `json:"user_agent"`
	Referer    pgtype.Text      `json:"referer"`
	DeviceType pgtype.Text      `json:"device_type"`
	Country    pgtype.Text      `json:"country"`
}

func (s *Server) GetUrlStats(ctx *gin.Context) {
	urlID, err := strconv.ParseInt(ctx.Param("url_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}

	var req GetUrlStatsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	urlIDPg := pgtype.Int8{Int64: urlID, Valid: true}
	offset := req.Page * req.Limit

	stats, err := s.store.GetClicksByURLID(ctx, db.GetClicksByURLIDParams{
		UrlID:  urlIDPg,
		Limit:  req.Limit,
		Offset: offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL stats"})
		return
	}

	content := make([]UrlStats, len(stats))
	for i, s := range stats {
		content[i] = UrlStats{
			Id:         s.ID,
			ClickedAt:  s.ClickedAt,
			IpAddress:  s.IpAddress,
			UserAgent:  s.UserAgent,
			Referer:    s.Referer,
			DeviceType: s.DeviceType,
			Country:    s.Country,
		}
	}

	total, err := s.store.CountClicksByURLID(ctx, urlIDPg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve click count"})
		return
	}

	isLast := offset+int32(len(stats)) >= int32(total)

	ctx.JSON(http.StatusOK, GetUrlStatsResponse{
		Content:     content,
		IsFirst:     req.Page == 0,
		IsLast:      isLast,
		IsPrevious:  req.Page > 0,
		IsNext:      !isLast,
		CurrentPage: req.Page,
		TotalCount:  total,
	})
}

type GetUrlClickCountResponse struct {
	ClickCount int64 `json:"click_count"`
}

func (s *Server) GetUrlClickCount(ctx *gin.Context) {
	urlId := ctx.Param("url_id")
	urlIdInt, err := strconv.ParseInt(urlId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL ID"})
		return
	}
	clickCount, err := s.store.CountClicksByURLID(ctx, pgtype.Int8{Int64: urlIdInt, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve click count"})
		return
	}
	ctx.JSON(http.StatusOK, GetUrlClickCountResponse{ClickCount: clickCount})
}
