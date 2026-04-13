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
		c.JSON(http.StatusBadRequest, gin.H{"error": "No hosts provided"})
		return
	}

	hosts := strings.Split(hostsParam, ",")

	// Parse nameservers (optional)
	var nameservers []string
	nameserversParam := c.Query("nameservers")
	if nameserversParam != "" {
		nameservers = strings.Split(nameserversParam, ",")
	}

	// Parse includeDefault (default: true)
	includeDefault := c.DefaultQuery("includeDefault", "true") == "true"

	// Parse includeAAAA (default: false)
	includeAAAA := c.DefaultQuery("includeAAAA", "false") == "true"

	// Build final nameserver list
	finalNameservers := buildNameserverList(nameservers, includeDefault)

	results := dnslookup.GetDNSRecords(hosts, finalNameservers, includeAAAA)

	c.JSON(http.StatusOK, models.LookupResponse{Results: results})
}

// buildNameserverList constructs the final list of nameservers based on user input and includeDefault flag.
func buildNameserverList(nameservers []string, includeDefault bool) []string {
	defaults := dnslookup.GetDefaultNameservers()

	if len(nameservers) == 0 {
		// No custom nameservers provided — always use default
		return defaults
	}

	if includeDefault {
		// Prepend default nameservers to custom ones
		return append(defaults, nameservers...)
	}

	return nameservers
}

// InfoHandler returns application version and repo information
func InfoHandler(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": version,
			"repo":    "https://github.com/sunggun-yu/dnsq",
		})
	}
}
