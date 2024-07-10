package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sunggun-yu/dnsq/server/handlers"
)

func main() {
	r := gin.Default()
	r.GET("/api/lookup", handlers.DNSLookupHandler)
	r.Run(":8080")
}
