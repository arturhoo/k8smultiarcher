package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cache Cache

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	configureCache()

	r := gin.Default()
	r.POST("/mutate", mutateHandler)
	r.GET("/healthz", healthzHandler)
	r.GET("/livez", livezHandler)
	startServer(r)
}

func mutateHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Msgf("failed to read request body: %s", err)
	}

	review, err := ProcessAdmissionReview(cache, body)
	if err != nil {
		log.Printf("failed process pod admission review: %s", err)
		return
	}
	c.JSON(200, review)
}

func healthzHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func livezHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func configureCache() {
	var cacheChoice, cacheSizeStr string
	var cacheSize int
	var err error

	if cacheSizeStr = os.Getenv("CACHE_SIZE"); cacheSizeStr == "" {
		cacheSize = cacheSizeDefault
	} else {
		cacheSize, err = strconv.Atoi(cacheSizeStr)
		if err != nil {
			log.Fatal().Msgf("invalid cache size: %s", cacheSizeStr)
		}
	}

	if cacheChoice = os.Getenv("CACHE"); cacheChoice == "" || cacheChoice == "inmemory" {
		log.Print("using in-memory cache")
		cache = NewInMemoryCache(cacheSize)
	} else if cacheChoice == "redis" {
		var redisAddr string
		if redisAddr = os.Getenv("REDIS_ADDR"); redisAddr == "" {
			redisAddr = redisAddrDefault
		}
		log.Print("using redis cache")
		cache = NewRedisCache(redisAddr)
	} else {
		log.Fatal().Msgf("invalid cache choice. must be either 'redis' or 'inmemory': %s.", cacheChoice)
	}
}

func startServer(r *gin.Engine) {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	tlsEnabled := os.Getenv("TLS_ENABLED")
	if port == "" {
		if tlsEnabled == "true" {
			port = "8443"
		} else {
			port = "8080"
		}
	}
	addr := fmt.Sprintf("%s:%s", host, port)
	if tlsEnabled == "true" {
		var certPath, keyPath string
		if certPath = os.Getenv("CERT_PATH"); certPath == "" {
			certPath = "./certs/tls.crt"
		}
		if keyPath = os.Getenv("KEY_PATH"); keyPath == "" {
			keyPath = "./certs/tls.key"
		}
		r.RunTLS(addr, certPath, keyPath)
	}
	r.Run(addr)
}
