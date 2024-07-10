package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sunggun-yu/dnsq/internal/server/handlers"
)

func Run(port int) {
	// Create a new Gin router
	r := gin.Default()

	// Serve static files, such as index.html
	r.Static("/static", "./static")

	// API endpoint for DNS lookup
	r.GET("/api/lookup", handlers.DNSLookupHandler)

	// Serve the index.html at the root
	r.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Run the server on the given port
	r.Run(fmt.Sprintf(":%d", port))
}
