package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetClientIP(ctx *gin.Context) string {
	forwarded := ctx.GetHeader("X-Forwarded-For")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}

	realIP := ctx.GetHeader("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	return ctx.ClientIP()
}

func DetectDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") {
		return "mobile"
	}
	if strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad") {
		return "tablet"
	}
	return "desktop"
}

func GetCountryFromIP(ip string) string {

	return "VN"
}
