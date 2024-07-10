package dnslookup

import (
	"net"
	"strings"

	"github.com/sunggun-yu/dnsq/internal/models"
)

// GetDNSRecords returns DNS records for a given hostname
func GetDNSRecords(hostname string) []models.DNSRecord {
	var records []models.DNSRecord
	currentHost := hostname

	// CNAME lookup
	cname, err := net.LookupCNAME(currentHost)
	if err == nil && cname != currentHost+"." {
		// remove trailing dot
		cname = strings.TrimSuffix(cname, ".")
		records = append(records, models.DNSRecord{Host: currentHost, Type: "CNAME", Data: cname})
		// replace currentHost with the CNAME
		currentHost = cname
	}

	// A and AAAA lookup
	ips, err := net.LookupIP(currentHost)
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
