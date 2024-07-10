package models

type DNSRecord struct {
	Host string `json:"host"`
	Type string `json:"type"`
	Data string `json:"data"`
}
