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
type visitor struct {
	firstSeen time.Time
	BanTime   time.Time
	numOfCon  int
	Restrict  bool
}

type errorInfo struct {
	NumReq  int
	TimeBan string
}

type Settings struct {
	Prefix     int
	NumCon     int
	LimitTime  time.Duration
	BanTime    time.Duration
	DeleteTime time.Duration
}

// Change the the map to hold values of the type visitor.
var Visitors = make(map[string]*visitor)
var Mu sync.Mutex

// Run a background goroutine to remove old entries from the visitors map.
func init() {
	go cleanupVisitors()
}

func InitSettings() *Settings {
	defVal := &Settings{
		24,
		100,
		1 * time.Minute,
		2 * time.Minute,
		1 * time.Minute,
	}
	return defVal
}

func GetIp(req *http.Request) (string, error) {
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

func getVisitor(ip string, params *Settings) *visitor {
	Mu.Lock()
	defer Mu.Unlock()

	rateLim := params.LimitTime
	banTime := params.BanTime
	prefix := params.Prefix
	numCon := params.NumCon
	IpSub := IpSubNet(ip, prefix)
	v, exists := Visitors[IpSub]
	if !exists {
		// Include the current time when creating a new visitor.
		Visitors[IpSub] = &visitor{time.Now(), time.Now(), 1, false}
		return Visitors[IpSub]
	}
	v.numOfCon++
	if v.Restrict == false && (v.numOfCon >= numCon && time.Since(v.firstSeen) <= rateLim) {
		fmt.Println(v.numOfCon)
		v.Restrict = true
		v.BanTime = time.Now()
	} else if v.Restrict == true && time.Since(v.BanTime) > banTime {
		delete(Visitors, IpSub)
	}
	return v
}

// Every minute check the map for visitors that haven't been seen for
// more than 2 minutes and delete the entries.
func cleanupVisitors() {
	params := InitSettings()
	for {
		time.Sleep(params.DeleteTime)
		Mu.Lock()
		banTime := params.BanTime
		for ip, v := range Visitors {
			if time.Since(v.BanTime) > banTime {
				delete(Visitors, ip)
			}
		}
		Mu.Unlock()
	}
}

func Limit(next http.Handler, params *Settings) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, err := GetIp(r)
		if err != nil {
			ip, _, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		limiter := getVisitor(ip, params)
		if limiter.Restrict == true {
			tmpl, _ := template.ParseFiles("./static/429.html")
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
