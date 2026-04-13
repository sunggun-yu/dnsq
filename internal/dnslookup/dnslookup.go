package dnslookup

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/sunggun-yu/dnsq/internal/models"
)

const (
	maxCNAMEDepth = 10
	dnsTimeout    = 5 * time.Second
)

// randomHostname generates a random hostname for wildcard domain probing
func randomHostname() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	length := 8
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// GetDefaultNameservers reads nameservers from /etc/resolv.conf.
// Falls back to 8.8.8.8:53 if the file cannot be read.
func GetDefaultNameservers() []string {
	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil || len(config.Servers) == 0 {
		return []string{"8.8.8.8:53"}
	}
	servers := make([]string, 0, len(config.Servers))
	for _, s := range config.Servers {
		servers = append(servers, normalizeNameserver(s))
	}
	return servers
}

// normalizeNameserver ensures a nameserver address has a port (default :53).
func normalizeNameserver(ns string) string {
	ns = strings.TrimSpace(ns)
	if ns == "" {
		return ""
	}
	_, _, err := net.SplitHostPort(ns)
	if err != nil {
		// No port specified, add default
		return net.JoinHostPort(ns, "53")
	}
	return ns
}

// GetDNSRecords performs DNS lookups for the given hosts against each nameserver.
// Queries are run concurrently across nameservers.
func GetDNSRecords(hosts []string, nameservers []string, includeAAAA bool) []models.NameserverResult {
	// Normalize nameservers
	normalized := make([]string, 0, len(nameservers))
	for _, ns := range nameservers {
		if n := normalizeNameserver(ns); n != "" {
			normalized = append(normalized, n)
		}
	}
	if len(normalized) == 0 {
		normalized = GetDefaultNameservers()
	}

	results := make([]models.NameserverResult, len(normalized))
	var wg sync.WaitGroup

	for i, ns := range normalized {
		wg.Add(1)
		go func(idx int, nameserver string) {
			defer wg.Done()
			// Check nameserver reachability with a simple query
			if err := checkNameserver(nameserver); err != nil {
				results[idx] = models.NameserverResult{
					Nameserver: nameserver,
					Error:      err.Error(),
					Results:    make(map[string][]models.DNSRecord),
				}
				return
			}
			results[idx] = models.NameserverResult{
				Nameserver: nameserver,
				Results:    lookupAllHosts(hosts, nameserver, includeAAAA),
			}
		}(i, ns)
	}

	wg.Wait()
	return results
}

// checkNameserver verifies that a nameserver is reachable by sending a simple query.
func checkNameserver(nameserver string) error {
	msg := new(dns.Msg)
	msg.SetQuestion(".", dns.TypeNS)
	msg.RecursionDesired = true

	client := &dns.Client{Timeout: dnsTimeout}
	_, _, err := client.Exchange(msg, nameserver)
	if err != nil {
		return friendlyError(err)
	}
	return nil
}

// friendlyError converts low-level DNS errors into user-friendly messages.
func friendlyError(err error) error {
	s := err.Error()
	if strings.Contains(s, "i/o timeout") {
		return fmt.Errorf("connection timed out — nameserver is unreachable or not responding")
	}
	if strings.Contains(s, "connection refused") {
		return fmt.Errorf("connection refused — no DNS service running on this address")
	}
	if strings.Contains(s, "no such host") {
		return fmt.Errorf("nameserver hostname could not be resolved")
	}
	if strings.Contains(s, "network is unreachable") {
		return fmt.Errorf("network is unreachable — cannot reach this nameserver")
	}
	return err
}

// lookupAllHosts resolves all hosts against a single nameserver.
func lookupAllHosts(hosts []string, nameserver string, includeAAAA bool) map[string][]models.DNSRecord {
	result := make(map[string][]models.DNSRecord)
	for _, host := range hosts {
		host = strings.TrimSpace(host)
		if host != "" {
			result[host] = resolveHost(host, nameserver, includeAAAA)
		}
	}
	return result
}

// resolveHost resolves a single host, following CNAME chains and returning A (and optionally AAAA) records.
func resolveHost(hostname, nameserver string, includeAAAA bool) []models.DNSRecord {
	var records []models.DNSRecord
	resolveDNS(hostname, nameserver, includeAAAA, &records, 0)
	return records
}

// resolveDNS recursively resolves DNS records for a hostname.
func resolveDNS(hostname, nameserver string, includeAAAA bool, records *[]models.DNSRecord, depth int) {
	if depth >= maxCNAMEDepth {
		return
	}

	currentHost := hostname
	isWildcard := false

	// Handle wildcard domains
	if strings.HasPrefix(currentHost, "*.") {
		randomSubdomain := randomHostname()
		currentHost = randomSubdomain + hostname[1:]
		isWildcard = true
	}

	// Query for CNAME
	cname := queryCNAME(currentHost, nameserver)
	if cname != "" && cname != currentHost {
		displayHost := currentHost
		if isWildcard {
			displayHost = hostname
		}
		*records = append(*records, models.DNSRecord{Host: displayHost, Type: "CNAME", Data: cname})
		// Recursively resolve the CNAME target
		resolveDNS(cname, nameserver, includeAAAA, records, depth+1)
		return
	}

	// Query for A records
	displayHost := currentHost
	if isWildcard {
		displayHost = hostname
	}

	aRecords := queryA(currentHost, nameserver)
	for _, ip := range aRecords {
		*records = append(*records, models.DNSRecord{Host: displayHost, Type: "A", Data: ip})
	}

	// Optionally query for AAAA records
	if includeAAAA {
		aaaaRecords := queryAAAA(currentHost, nameserver)
		for _, ip := range aaaaRecords {
			*records = append(*records, models.DNSRecord{Host: displayHost, Type: "AAAA", Data: ip})
		}
	}
}

// queryCNAME queries a nameserver for CNAME records.
// Returns the CNAME target (without trailing dot) or empty string.
func queryCNAME(host, nameserver string) string {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(host), dns.TypeCNAME)
	msg.RecursionDesired = true

	client := &dns.Client{Timeout: dnsTimeout}
	resp, _, err := client.Exchange(msg, nameserver)
	if err != nil || resp == nil {
		return ""
	}

	for _, ans := range resp.Answer {
		if cname, ok := ans.(*dns.CNAME); ok {
			return strings.TrimSuffix(cname.Target, ".")
		}
	}
	return ""
}

// queryA queries a nameserver for A records.
func queryA(host, nameserver string) []string {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(host), dns.TypeA)
	msg.RecursionDesired = true

	client := &dns.Client{Timeout: dnsTimeout}
	resp, _, err := client.Exchange(msg, nameserver)
	if err != nil || resp == nil {
		return nil
	}

	var ips []string
	for _, ans := range resp.Answer {
		if a, ok := ans.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}
	return ips
}

// queryAAAA queries a nameserver for AAAA records.
func queryAAAA(host, nameserver string) []string {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(host), dns.TypeAAAA)
	msg.RecursionDesired = true

	client := &dns.Client{Timeout: dnsTimeout}
	resp, _, err := client.Exchange(msg, nameserver)
	if err != nil || resp == nil {
		return nil
	}

	var ips []string
	for _, ans := range resp.Answer {
		if aaaa, ok := ans.(*dns.AAAA); ok {
			ips = append(ips, aaaa.AAAA.String())
		}
	}
	return ips
}

// FormatNameserver returns a display-friendly nameserver string.
// Attempts reverse DNS lookup to show hostname alongside IP.
func FormatNameserver(ns string) string {
	host, _, err := net.SplitHostPort(ns)
	if err != nil {
		return ns
	}
	names, err := net.LookupAddr(host)
	if err == nil && len(names) > 0 {
		name := strings.TrimSuffix(names[0], ".")
		return fmt.Sprintf("%s (%s)", ns, name)
	}
	return ns
}
