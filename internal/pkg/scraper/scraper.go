package scraper

import (
	"proxy-harvester/internal/pkg/checker"
	"proxy-harvester/internal/pkg/scraper/services"
)

// interface that all scrapers have to follow
type proxySource func() []string

// list of all active proxy scrapers
var sources = []proxySource{
	services.ProxyScrape,
	services.Geonode,
	services.FreeProxyList,
	services.HideMyName,
	services.FreeProxy,
	services.TheSpeedX,
}

func ScrapeAll() {
	for _, source := range sources {
		for _, s := range source() { //scrape
			checker.Queue.Enqueue(s) //queue results
		}
	}
}
