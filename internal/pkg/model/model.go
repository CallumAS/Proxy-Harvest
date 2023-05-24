package model

import "time"

type Proxy struct {
	IP          string
	Port        int
	LastChecked time.Time
	ProxyType   int
}

type ScannerSettings struct {
	Ranges []string `json:"ranges"`
	Ports  []int    `json:"ports"`
	Active bool     `json:"active"`
	Tasks  int      `json:"tasks"`
}
type ScraperSettings struct {
	Active bool `json:"active"`
}
type CheckerSettings struct {
	Tasks   int    `json:"tasks"`
	Timeout string `json:"timeout"`
}
type Settings struct {
	Scanner ScannerSettings `json:"scanner"`
	Scraper ScraperSettings `json:"scraper"`
	Checker CheckerSettings `json:"checker"`
}
