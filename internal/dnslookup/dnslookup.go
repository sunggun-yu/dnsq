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
func GetDNSRecords(hosts []string) map[string][]models.DNSRecord {
	result := make(map[string][]models.DNSRecord)
	for _, host := range hosts {
		host = strings.TrimSpace(host)
		if host != "" {
			result[host] = dnsRecords(host)
		}
	}
	return result
}

// GetDNSRecords returns DNS records for a given hostname
func dnsRecords(hostname string) []models.DNSRecord {
	var records []models.DNSRecord
	resolveDNS(hostname, &records)
	return records
}

// resolveDNS resolves DNS records for a given hostname. It is a recursive function that resolves CNAME and A/AAAA records.
func resolveDNS(hostname string, records *[]models.DNSRecord) {
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

		*records = append(*records, models.DNSRecord{Host: cnameHost, Type: "CNAME", Data: cname})
		// replace currentHost with the CNAME
		currentHost = cname

		// recursively resolve the CNAME
		resolveDNS(currentHost, records)
		// no need to continue if it is a CNAME
		// so that, final CNAME will resolve the A/AAAA records
		return
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
				*records = append(*records, models.DNSRecord{Host: currentHost, Type: "A", Data: ipv4.String()})
			} else {
				*records = append(*records, models.DNSRecord{Host: currentHost, Type: "AAAA", Data: ip.String()})
			}
		}
	}
}
