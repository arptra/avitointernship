package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"avitointernship/server/pkg/ratelimit"
)

type infoBan struct {
	Id       int
	Subnet   string
	TimeLeft string
}

type infoIps struct {
	Id int
	Ip string
}

type info struct {
	Ips     []infoIps
	Ban     []infoBan
	ServSet *ratelimit.Settings
}

var conections = []infoIps{}

func addIp(ip string) {
	var id int
	if len(conections) == 0 {
		id = 0
	} else {
		id = conections[len(conections)-1].Id
	}
	id++
	conections = append(conections, infoIps{id, ip})
}

func getInfo() info {
	ratelimit.Mu.Lock()
	vectorBan := []infoBan{}
	i := 1
	for subnet, v := range ratelimit.RstrctdLst {
		if v.Restrict == true {
			ban := infoBan{i, subnet, time.Now().Sub(v.BanTime).String()}
			vectorBan = append(vectorBan, ban)
			i++
		}
	}
	ratelimit.Mu.Unlock()
	IpsBan := info{conections, vectorBan, ratelimit.Conf}
	return IpsBan
}

func showIps(w http.ResponseWriter) {
	ips := getInfo()
	tmpl, _ := template.ParseFiles("./static/index.html")
	tmpl.Execute(w, ips)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	ip, _ := ratelimit.GetIp(r)
	if len(ip) > 0 {
		addIp(ip)
	}
	showIps(w)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/delete" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "./static/index.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "%v\n", r.PostForm)
		address := r.FormValue("address")
		ratelimit.Mu.Lock()
		delete(ratelimit.RstrctdLst, address)
		ratelimit.Mu.Unlock()
		fmt.Fprintf(w, "Address = %s have deleted \n", address)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func ChangeConfHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/change_settings" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "./static/index.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "%v\n", r.PostForm)
		Prefix, _ := strconv.ParseInt(r.FormValue("Prefix"), 10, 64)
		NumCon, _ := strconv.ParseInt(r.FormValue("NumCon"), 10, 64)
		LimitTime, _ := strconv.ParseInt(r.FormValue("LimitTime"), 10, 64)
		BanTime, _ := strconv.ParseInt(r.FormValue("BanTime"), 10, 64)
		DeleteTime, _ := strconv.ParseInt(r.FormValue("DeleteTime"), 10, 64)
		ratelimit.Conf = &ratelimit.Settings{
			Prefix:     int(Prefix),
			NumCon:     int(NumCon),
			LimitTime:  time.Duration(int(LimitTime)) * time.Minute,
			BanTime:    time.Duration(int(BanTime)) * time.Minute,
			DeleteTime: time.Duration(int(DeleteTime)) * time.Minute,
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	ratelimit.InitSettings()
	fmt.Println("Server is listening...")
	router := mux.NewRouter()
	router.HandleFunc("/", infoHandler)
	router.HandleFunc("/delete", DeleteHandler)
	router.HandleFunc("/change_settings", ChangeConfHandler)
	http.Handle("/", router)
	http.ListenAndServe("0.0.0.0:8181", ratelimit.Limit(router))
}
