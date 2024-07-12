package dnslookup

import (
	"math/rand"
	"net"
	"strings"

	"github.com/sunggun-yu/dnsq/internal/models"
)

// randomHostname generates a random hostname
func randomHostname() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	length := 8
	// TODO: need to check length of hostname and see if it is over 255 long including the subdomain
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// GetDNSRecords returns DNS records for a given hostname
func GetDNSRecords(hostname string) []models.DNSRecord {
	var records []models.DNSRecord
	currentHost := hostname
	isWildcard := false

	// check if it is a wildcard
	if strings.HasPrefix(currentHost, "*.") {
		randomSubdomain := randomHostname()
		currentHost = randomSubdomain + hostname[1:]
		isWildcard = true
	}

	// CNAME lookup
	cname, err := net.LookupCNAME(currentHost)

	if err == nil && cname != currentHost+"." {
		// remove trailing dot
		cname = strings.TrimSuffix(cname, ".")
		cnameHost := currentHost

		// if it is a wildcard, set the original hostname to the records
		if isWildcard {
			cnameHost = hostname
			isWildcard = false // reset the wildcard flag
		}

		records = append(records, models.DNSRecord{Host: cnameHost, Type: "CNAME", Data: cname})
		// replace currentHost with the CNAME
		currentHost = cname
	}

	// A and AAAA lookup
	ips, err := net.LookupIP(currentHost)

	// if it is a wildcard, set the original hostname to the records
	if isWildcard {
		currentHost = hostname
	}

	if err == nil {
		for _, ip := range ips {
			if ipv4 := ip.To4(); ipv4 != nil {
				records = append(records, models.DNSRecord{Host: currentHost, Type: "A", Data: ipv4.String()})
			} else {
				records = append(records, models.DNSRecord{Host: currentHost, Type: "AAAA", Data: ip.String()})
			}
		}
	}

	return records
}
