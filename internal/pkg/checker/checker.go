package checker

import (
	"github.com/enriquebris/goconcurrentqueue"
	"h12.io/socks"
	"net/http"
	"net/url"
	"proxy-harvester/internal/pkg/model"
	"strconv"
	"strings"
	"time"
)

var Timeout = "5s"

const (
	Unknown int = 0
	HTTP        = 1
	HTTPS       = 2
	SOCKS4      = 3
	Socks4a     = 4
	SOCKS5      = 5
	INVALID     = 6
)

func Check(proxy string, proxyType int) model.Proxy {
	var parsed = strings.Split(proxy, ":")
	if len(parsed) < 2 {
		return model.Proxy{
			IP:          "",
			Port:        0,
			LastChecked: time.Now(),
			ProxyType:   INVALID,
		}
	}
	var ip = parsed[0]
	var port, err = strconv.Atoi(parsed[1])
	if err != nil {
		return model.Proxy{
			IP:          "",
			Port:        0,
			LastChecked: time.Now(),
			ProxyType:   INVALID,
		}

	}
	var client http.Client
	switch proxyType {
	case HTTP:
		dialSocksProxy := socks.Dial("http://" + proxy + "?timeout=" + Timeout)
		proxyUrl, _ := url.Parse("http://" + proxy + "?timeout=" + Timeout)
		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		tr := &http.Transport{Dial: dialSocksProxy}
		client = http.Client{Transport: tr}
		break
	case HTTPS:
		dialSocksProxy := socks.Dial("http://" + proxy + "?timeout=" + Timeout)
		proxyUrl, _ := url.Parse("http://" + proxy + "?timeout=" + Timeout)
		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		tr := &http.Transport{Dial: dialSocksProxy}
		client = http.Client{Transport: tr}
		break
	case SOCKS4:
		dialSocksProxy := socks.Dial("socks4://" + proxy + "?timeout=" + Timeout)
		tr := &http.Transport{Dial: dialSocksProxy}
		client = http.Client{Transport: tr}
		break
	case Socks4a:
		dialSocksProxy := socks.Dial("socks4a://" + proxy + "?timeout=" + Timeout)
		tr := &http.Transport{Dial: dialSocksProxy}
		client = http.Client{Transport: tr}
		break
	case SOCKS5:
		dialSocksProxy := socks.Dial("socks5://" + proxy + "?timeout=" + Timeout)
		tr := &http.Transport{Dial: dialSocksProxy}
		client = http.Client{Transport: tr}
		break
	case Unknown:
		for _, pType := range []int{HTTPS, SOCKS4, Socks4a, SOCKS5} {
			var res = Check(proxy, pType)
			if res.ProxyType != INVALID {
				return res
			}
		}
		return model.Proxy{
			IP:          ip,
			Port:        port,
			LastChecked: time.Now(),
			ProxyType:   INVALID,
		}
	}

	resp, err := client.Get("https://httpbin.org/get")
	if err != nil {
		return model.Proxy{
			IP:          ip,
			Port:        port,
			LastChecked: time.Now(),
			ProxyType:   INVALID,
		}
	}
	defer resp.Body.Close()
	return model.Proxy{
		IP:          ip,
		Port:        port,
		LastChecked: time.Now(),
		ProxyType:   proxyType,
	}

}

func check(proxy string) {
	//println("Checking", proxy)
	var result = Check(proxy, Unknown)
	if result.ProxyType != INVALID {
		Results[proxy] = result
	} else {
		_, ok := Results[proxy]
		if ok {
			delete(Results, proxy)

		}
	}

	//fmt.Println("Unknown Proxy Type or Invalid Proxy")
}

var minutes = time.Duration(5)
var recheckAfterMinutes = time.Duration(5)

func Recheck() {
	ticker := time.NewTicker(minutes * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				for proxy, data := range Results {
					currentTimestamp := time.Now()
					thresholdTime := currentTimestamp.Add(-time.Duration(minutes) * time.Minute)

					if data.LastChecked.Before(thresholdTime) {
						Queue.Enqueue(proxy)
					}
				}

				// do stuff
			}
		}
	}()
}

var Tasks = 50
var (
	Queue   = goconcurrentqueue.NewFIFO()
	Results = map[string]model.Proxy{}
)

func SetTasks(task int) {
	if task > Tasks {
		// Increase the number of goroutines
		for i := Tasks; i < task; i++ {
			go func(routineId int) {
				for {
					value, _ := Queue.DequeueOrWaitForNextElement()
					proxy := value.(string)
					check(proxy)
					if routineId >= Tasks {
						break
					}
				}
			}(i)
		}
	}

	Tasks = task
}

func Start() {
	// Recheck results that have not been rechecked from new scans/scrapes
	Recheck()

	// Check new results coming in from scans/scrapes
	for i := 0; i < Tasks; i++ {
		go func(routineId int) {
			for {
				value, _ := Queue.DequeueOrWaitForNextElement()
				proxy := value.(string)
				check(proxy)
				if routineId >= Tasks {
					break
				}
			}
		}(i)
	}
}
