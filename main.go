package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)



type settings	struct {
	Prefix     int
	NumCon     int
	LimitTime  time.Duration
	BanTime    time.Duration
	DeleteTime time.Duration
}

type infoBan	struct {
	Id			int
	Subnet		string
	TimeLeft	string
}

type infoIps	struct {
	Id			int
	Ip      	string
}

type info		struct {
	Ips			[]infoIps
	Ban			[]infoBan
	ServSet		*settings
}

func writeIpToFile(ip string, db_name string){
	file, err := os.OpenFile(db_name, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		file, err := os.Create(db_name)
		if err != nil{
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		file.WriteString(ip)
		file.WriteString("\n")
		defer file.Close()
	}
	file.WriteString(ip)
	file.WriteString("\n")
}

func readDB(path string) info {
	vectorIps := []infoIps{}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	i := 1
	for scanner.Scan() {
		ip := infoIps{}
		ip.Ip = scanner.Text()
		ip.Id = i
		vectorIps = append(vectorIps, ip)
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	mu.Lock()
	vectorBan := []infoBan{}
	i = 1
	for subnet, v := range visitors {
		if v.restrict == true {
			ban := infoBan{}
			ban.Subnet = subnet
			ban.Id = i
			ban.TimeLeft = time.Now().Sub(v.banTime).String()
			vectorBan = append(vectorBan, ban)
			i++
		}
	}
	mu.Unlock()
	IpsBan := info{}
	IpsBan.Ips = vectorIps
	IpsBan.Ban = vectorBan
	IpsBan.ServSet = initSettings()
	return IpsBan
}

func showIps(w http.ResponseWriter) {

	ips := readDB("db_ip.txt")
	tmpl, _ := template.ParseFiles("index.html")
	tmpl.Execute(w, ips)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {

	ip, _ := getIp(r)
	if len(ip) > 0 {
	writeIpToFile(ip, "db_ip.txt")
	}
	showIps(w)
}

func initSettings() *settings{
	defVal := &settings{
					24,
					100,
					1*time.Minute,
					2*time.Minute,
					1*time.Minute,
	}
	return defVal
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/delete" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "index.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "%v\n", r.PostForm)
		address := r.FormValue("address")
		mu.Lock()
		delete(visitors, address)
		mu.Unlock()
		fmt.Fprintf(w, "Address = %s have deleted \n", address)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	params := initSettings()
	fmt.Println("Server is listening...")
	router := mux.NewRouter()
	router.HandleFunc("/", infoHandler)
	router.HandleFunc("/delete", DeleteHandler)
	http.Handle("/",router)
	http.ListenAndServe("localhost:8181", limit(router, params))
}