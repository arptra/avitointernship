package ratelimit

import (
	"flag"
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
var RstrctdLst = make(map[string]*visitor)
var Mu sync.Mutex
var Conf *Settings

// Run a background goroutine to remove old entries from the visitors map.
func init() {
	go cleanupVisitors()
}

func InitSettings() {
	prefix := flag.Int("p", 24, "an int")
	NumCon := flag.Int("nc", 100, "an int")
	LimitTime := flag.Int("lt", 1, "an int")
	BanTime := flag.Int("bt", 2, "an int")
	DeleteTime := flag.Int("dt", 1, "an int")
	flag.Parse()
	Conf = &Settings{
		*prefix,
		*NumCon,
		time.Duration(*LimitTime) * time.Minute,
		time.Duration(*BanTime) * time.Minute,
		time.Duration(*DeleteTime) * time.Minute,
	}
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

func getVisitor(ip string) *visitor {
	Mu.Lock()
	defer Mu.Unlock()

	rateLim := Conf.LimitTime
	banTime := Conf.BanTime
	prefix := Conf.Prefix
	numCon := Conf.NumCon
	IpSub := IpSubNet(ip, prefix)
	v, exists := RstrctdLst[IpSub]
	if !exists {
		// Include the current time when creating a new visitor.
		RstrctdLst[IpSub] = &visitor{time.Now(), time.Now(), 1, false}
		return RstrctdLst[IpSub]
	}
	v.numOfCon++
	if v.Restrict == false && (v.numOfCon >= numCon && time.Since(v.firstSeen) <= rateLim) {
		v.Restrict = true
		v.BanTime = time.Now()
	} else if v.Restrict == true && time.Since(v.BanTime) > banTime {
		delete(RstrctdLst, IpSub)
	}
	return v
}

// Every minute check the map for visitors that haven't been seen for
// more than 2 minutes and delete the entries.
func cleanupVisitors() {
	for {
		time.Sleep(Conf.DeleteTime)
		Mu.Lock()
		banTime := Conf.BanTime
		for ip, v := range RstrctdLst {
			if time.Since(v.BanTime) > banTime {
				delete(RstrctdLst, ip)
			}
		}
		Mu.Unlock()
	}
}

func Limit(next http.Handler) http.Handler {
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
		if ip != "::1" {
			limiter := getVisitor(ip)
			if limiter.Restrict == true {
				tmpl, _ := template.ParseFiles("./static/429.html")
				w.WriteHeader(429)
				desc := errorInfo{Conf.NumCon, Conf.BanTime.String()}
				tmpl.Execute(w, desc)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
