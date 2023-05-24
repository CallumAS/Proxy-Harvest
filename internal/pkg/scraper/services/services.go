package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ProxyScrape() []string {
	resp, err := http.Get("https://api.proxyscrape.com/v2/?request=getproxies&protocol=all&timeout=10000&country=all&ssl=all&anonymity=all")
	if err != nil {
		fmt.Println("Error fetching proxies:", err)
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	// Split the body by new lines to get individual proxies
	proxies := strings.Split(string(body), "\r\n")

	return proxies
}

type geonodeData struct {
	Data []struct {
		ID                 string      `json:"_id"`
		IP                 string      `json:"ip"`
		AnonymityLevel     string      `json:"anonymityLevel"`
		Asn                string      `json:"asn"`
		City               string      `json:"city"`
		Country            string      `json:"country"`
		CreatedAt          time.Time   `json:"created_at"`
		Google             bool        `json:"google"`
		Isp                string      `json:"isp"`
		LastChecked        int         `json:"lastChecked"`
		Latency            float64     `json:"latency"`
		Org                string      `json:"org"`
		Port               string      `json:"port"`
		Protocols          []string    `json:"protocols"`
		Region             interface{} `json:"region"`
		ResponseTime       int         `json:"responseTime"`
		Speed              int         `json:"speed"`
		UpdatedAt          time.Time   `json:"updated_at"`
		WorkingPercent     interface{} `json:"workingPercent"`
		UpTime             float64     `json:"upTime"`
		UpTimeSuccessCount int         `json:"upTimeSuccessCount"`
		UpTimeTryCount     int         `json:"upTimeTryCount"`
	} `json:"data"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func Geonode() []string {
	var result = []string{}
	for page := 0; page < 50; page++ {

		resp, err := http.Get("https://proxylist.geonode.com/api/proxy-list?limit=500&page=" + strconv.Itoa(page) + "&sort_by=lastChecked&sort_type=desc")
		if err != nil {
			break
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		var dat geonodeData

		if err := json.Unmarshal(body, &dat); err != nil {
			break
		}
		if len(dat.Data) <= 0 {
			break
		}
		for _, datum := range dat.Data {
			result = append(result, datum.IP+":"+datum.Port)
		}
	}
	return result
}

func FreeProxyList() []string {
	resp, err := http.Get("https://free-proxy-list.net/")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	pattern := `(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):(\d+)`

	// Compile the regex pattern
	regex := regexp.MustCompile(pattern)

	// Find the first match
	match := regex.FindAllString(string(body), -1)
	return match

}

// cloudflare going get this one
func HideMyName() []string {
	results := []string{}
	for page := 0; page < 50; page++ {
		client := cycletls.Init()
		response, err := client.Do("https://hidemy.name/en/proxy-list/?start="+strconv.Itoa(page*64)+"#list", cycletls.Options{
			Body:      "",
			Ja3:       "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-51-57-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0",
			UserAgent: "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0",
		}, "GET")
		if err != nil {
			log.Print("Request Failed: " + err.Error())
		}

		pattern := `<tr><td>(\d+\.\d+\.\d+\.\d+)</td><td>(\d+)</td>`

		// Compile the regex pattern
		regex := regexp.MustCompile(pattern)

		// Find the first match
		match := regex.FindAllString(string(response.Body), -1)
		if len(match) <= 0 {
			break
		}
		for _, s := range match {
			results = append(results, strings.Replace(strings.Replace(strings.Replace(s, "</td><td>", ":", 1), "<tr><td>", "", -1), "</td>", "", -1))
		}
	}
	return results
}

// fail2ban probably nginx related
func FreeProxy() []string {
	results := []string{}
	for page := 0; page < 100; page++ {
		client := cycletls.Init()
		response, err := client.Do("http://free-proxy.cz/en/proxylist/main/"+strconv.Itoa(page), cycletls.Options{
			Body:      "",
			Ja3:       "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-51-57-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0",
			UserAgent: "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0",
		}, "GET")
		if err != nil {
			break
		}

		pattern := `Base64\.decode\("([A-Za-z0-9+/=]+)"\).*?<span class="fport".*?>(\d+)<\/span><\/td>`

		// Compile the regex pattern
		regex := regexp.MustCompile(pattern)

		matches := regex.FindAllStringSubmatch(response.Body, -1)
		if matches == nil || len(matches) <= 0 {
			break
		}
		// Iterate over the matches
		for _, match := range matches {
			if len(match) >= 3 {
				encodedString := match[1]
				port := match[2]

				// Decode the Base64-encoded string
				decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
				if err != nil {
					fmt.Println("Error decoding Base64 string:", err)
					//					return
				}
				decodedString := string(decodedBytes)
				//println(decodedString + ":" + port)
				results = append(results, decodedString+":"+port)
			}
		}

		time.Sleep(3000)
	}
	return results
}

func TheSpeedX() []string {
	results := []string{}

	for _, url := range []string{"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt", "https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/socks4.txt", "https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/socks5.txt"} {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching proxies:", err)
			return nil
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return nil
		}

		// Split the body by new lines to get individual proxies
		proxies := strings.Split(string(body), "\n")
		for _, proxy := range proxies {
			results = append(results, proxy)
		}
	}
	return results
}
