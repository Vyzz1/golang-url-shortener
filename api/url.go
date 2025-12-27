package api

import (
	"context"
	"errors"
	"net/http"
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

func (s *Server) createUrl(ctx *gin.Context) {
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
		clickData := extractClickData(ctx)
		s.store.InsertClick(context.Background(), db.InsertClickParams{
			UrlID:      pgtype.Int8{Int64: urlRecord.ID, Valid: true},
			IpAddress:  clickData.IpAddress,
			UserAgent:  clickData.UserAgent,
			Referer:    clickData.Referer,
			Country:    clickData.Country,
			DeviceType: clickData.DeviceType,
		})

		s.store.IncrementClickCount(context.Background(), shortCode)
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
	return db.InsertClickParams{
		IpAddress:  pgtype.Text{String: utils.GetClientIP(ctx), Valid: true},
		UserAgent:  pgtype.Text{String: ctx.Request.UserAgent(), Valid: true},
		Referer:    pgtype.Text{String: ctx.Request.Referer(), Valid: true},
		Country:    pgtype.Text{String: utils.GetCountryFromIP(utils.GetClientIP(ctx)), Valid: true},
		DeviceType: pgtype.Text{String: utils.DetectDeviceType(ctx.Request.UserAgent()), Valid: true},
	}
}
