package manager

import (
	"proxy-harvester/internal/pkg/checker"
	"proxy-harvester/internal/pkg/scanner"
	"proxy-harvester/internal/pkg/scraper"
	"time"
)

var minutes = time.Duration(60)

var Scan = false
var Scrape = false

func process() {
	if Scan {
		scanner.Start()
	}
	if Scrape {
		scraper.ScrapeAll()
	}
}

func Start() {
	//start proxy checking results
	go checker.Start()

	process()

	//start scraping &/or checking every x amount of minutes
	ticker := time.NewTicker(minutes * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				process()
			}
		}
	}()
}
