package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sunggun-yu/dnsq/internal/server/handlers"
)

func Run(port int) {
	r := gin.Default()
	r.GET("/api/lookup", handlers.DNSLookupHandler)
	r.Run(fmt.Sprintf(":%d", port))
}
