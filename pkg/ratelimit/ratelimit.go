package ratelimit

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor	struct {
	firstSeen	time.Time
	banTime		time.Time
	numOfCon	int
	restrict	bool
}

type errorInfo struct {
	NumReq		int
	TimeBan		string
}

// Change the the map to hold values of the type visitor.
var visitors = make(map[string]*visitor)
var mu sync.Mutex

// Run a background goroutine to remove old entries from the visitors map.
func init() {
	go cleanupVisitors()
}

func getIp(req *http.Request) (string, error)  {
	ips := req.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}
	return "", fmt.Errorf("No valid ip found")
}

func IpSubNet(ip string, prefix int) string {
	ipv4Addr := net.ParseIP(ip)
	ipv4Mask := net.CIDRMask(prefix, 32)
	IpSub := ipv4Addr.Mask(ipv4Mask).String()
	return IpSub
}

func getVisitor(ip string, params* settings) *visitor {
	mu.Lock()
	defer mu.Unlock()

	rateLim := params.LimitTime
	banTime := params.BanTime
	prefix	:= params.Prefix
	numCon	:= params.NumCon
	IpSub := IpSubNet(ip, prefix)
	v, exists := visitors[IpSub]
	if !exists {
		// Include the current time when creating a new visitor.
		visitors[IpSub] = &visitor{time.Now(), time.Now(), 1, false}
		return visitors[IpSub]
	}
	v.numOfCon++
	if v.restrict == false && (v.numOfCon >= numCon && time.Since(v.firstSeen) <= rateLim) {
		fmt.Println(v.numOfCon)
		v.restrict = true
		v.banTime = time.Now()
	} else if v.restrict == true && time.Since(v.banTime) > banTime {
		delete(visitors, IpSub)
	}
	return v
}

// Every minute check the map for visitors that haven't been seen for
// more than 2 minutes and delete the entries.
func cleanupVisitors() {
	params := initSettings()
	for {
		time.Sleep(params.DeleteTime)

		mu.Lock()
		banTime := params.BanTime
		for ip, v := range visitors {
			if time.Since(v.banTime) > banTime {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func limit(next http.Handler, params* settings) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, err := getIp(r)
		if err != nil {
			ip, _, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		limiter := getVisitor(ip, params)
		if limiter.restrict == true {
			tmpl, _ := template.ParseFiles("429.html")
			w.WriteHeader(429)
			desc := errorInfo{}
			desc.NumReq = params.NumCon
			desc.TimeBan = params.BanTime.String()
			tmpl.Execute(w, desc)
			return
		}
		next.ServeHTTP(w, r)
	})
}