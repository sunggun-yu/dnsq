package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sunggun-yu/dnsq/internal/server/handlers"
)

func Execute() {
	r := gin.Default()
	r.GET("/api/lookup", handlers.DNSLookupHandler)
	r.Run(":8080")
}
