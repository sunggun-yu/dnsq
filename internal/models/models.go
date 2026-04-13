package models

// DNSRecord represents a single DNS record
type DNSRecord struct {
	Host string `json:"host"`
	Type string `json:"type"`
	Data string `json:"data"`
}

// NameserverResult holds all DNS results from a single nameserver
type NameserverResult struct {
	Nameserver string                  `json:"nameserver"`
	Error      string                  `json:"error,omitempty"`
	Results    map[string][]DNSRecord  `json:"results"`
}

// LookupResponse is the top-level API response
type LookupResponse struct {
	Results []NameserverResult `json:"results"`
}
