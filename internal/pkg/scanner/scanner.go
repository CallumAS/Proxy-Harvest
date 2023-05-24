package scanner

import (
	"fmt"
	"net"
	"proxy-harvester/internal/pkg/checker"
	"sync"
	"time"
)

var Ports = []int{
	8080,
	4145,
	//80,
	3128,
	999,
	9091,
	5678,
	8089,
	4153,
	9002,
	1080,
	8585,
	3629,
	//443,
	8081,
	8888,
	32650,
	32767,
	9090,
	1981,
	9999,
	10080,
	8181,
	83,
	1976,
	57775,
	8085,
	10800,
	7777,
	8083,
	30001,
	39593,
	50001,
	55443,
	8090,
	8899,
	8999,
	1088,
	33080,
	53281,
	8082,
	8118,
	8192,
	20000,
	25566,
	45787,
	5001,
	5566,
	58208,
	59166,
	7890,
	808,
	8123,
	84,
	9000,
	9300,
	9991,
	10001,
	10801,
	1081,
	111,
	1313,
	14282,
	3125,
	3127,
	3129,
	3333,
	4006,
	41317,
	41476,
	41485,
	41890,
	45005,
	5443,
	6969,
	8001,
	81,
	8104,
	82,
	8443,
	87,
	8989,
	9443,
	9992,
	10,
	10004,
	10808,
	10919,
	1100,
	11166,
	12370,
	12391,
	13267,
	1337,
	14287,
	14455,
	14921,
	15280,
	15291,
	15294,
	15303,
	15864,
	16003,
	16894,
	18080,
	18081,
	18301,
	18762,
	18765,
	19132,
	19404,
	25100,
	27360,
	27391,
	30000,
	31034,
	31433,
	31653,
	31679,
	32030,
	35904,
	36181,
	37726,
	41460,
	41679,
	4444,
	4480,
	46164,
	4673,
	49547,
	5002,
	54321,
	60671,
	63389,
	64149,
	64581,
	64943,
	6667,
	7657,
	7788,
	8000,
	8002,
	8009,
	8010,
	8084,
	8088,
	8282,
	8998,
	9191,
	9980,
	9981,
	9982,
	9994,
	9995,
	9996,
	9998,
}
var Ranges = []string{}
var RangeTasks = 100

func Start() {
	//for each cidr specified
	for _, s := range Ranges {
		//check for an open port and send to checker
		grab(s)
	}
}

// process an ip cidr
func grab(ipRange string) {

	wg := sync.WaitGroup{}                       // WaitGroup to ensure all goroutines finish
	semaphore := make(chan struct{}, RangeTasks) // Semaphore to limit the number of concurrent goroutines

	// Start goroutines for scraping IP ranges concurrently
	for _, ip := range hosts(ipRange) {
		semaphore <- struct{}{} // Acquire a semaphore before starting a goroutine
		wg.Add(1)
		go func(ip string) {
			scrapeIP(ip, &wg)
			<-semaphore // Release the semaphore when the goroutine finishes
		}(ip)
	}
	wg.Wait()
	fmt.Println("Scanner complete!")
}

// get all ips in the cidr
// https://gist.github.com/kotakanbe/d3059af990252ba89a82
func hosts(cidr string) []string {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1]
}

// https://gist.github.com/kotakanbe/d3059af990252ba89a82
// http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// checks if a specific IP has a port is open
func scrapeIP(ip string, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, port := range Ports {
		address := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err == nil {
			_ = conn.Close() // Close the connection if successful

			//add to proxy checker
			checker.Queue.Enqueue(address)
		}
	}
}
