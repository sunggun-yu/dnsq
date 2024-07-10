package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunggun-yu/dnsq/internal/dnslookup"
	"github.com/sunggun-yu/dnsq/internal/models"
)

// DNSLookupHandler handles DNS lookup requests
func DNSLookupHandler(c *gin.Context) {
	hostsParam := c.Query("hosts")
	if hostsParam == "" {
		// return 400 Bad Request if no hosts provided
		c.JSON(http.StatusBadRequest, gin.H{"error": "No hosts provided"})
		return
	}

	hosts := strings.Split(hostsParam, ",")
	results := make(map[string][]models.DNSRecord)
	for _, host := range hosts {
		host = strings.TrimSpace(host)
		if host != "" {
			results[host] = dnslookup.GetDNSRecords(host)
		}
	}

	// return 200 OK with the DNS records
	c.JSON(http.StatusOK, results)
} 
