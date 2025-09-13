package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/orhantugrul/medium-rsjs/src/feed"
)

type Config struct {
	Port           string
	TrustedProxies []string
}

func NewConfig() Config {
	config := Config{
		Port:           ":8080",
		TrustedProxies: []string{"127.0.0.1"},
	}

	if os.Getenv("GIN_MODE") == "release" {
		if port := os.Getenv("PORT"); port != "" {
			config.Port = ":" + port
		}

		if proxies := os.Getenv("TRUSTED_PROXIES"); proxies != "" {
			config.TrustedProxies = strings.Split(proxies, ",")
		}
	}

	return config
}

func main() {
	config := NewConfig()
	router := gin.New()

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %s %s %s %d %s %v\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.MethodColor(),
			param.Method,
			param.ResetColor(),
			param.Path,
			param.StatusCodeColor(),
			param.StatusCode,
			param.ResetColor(),
			param.Latency,
		)
	}))
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		feed.BindRoutes(api)

		api.GET("/health", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"status":  "✨ healthy",
				"service": "Medium RSJS API",
				"version": "1.0.0",
				"uptime":  "running",
			})
		})
	}

	router.Run(config.Port)
}
