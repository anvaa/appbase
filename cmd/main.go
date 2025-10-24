package main

import (
	"server"
	"server/global"
	"server/srv_conf"
	"server/srv_sec"

	"user"

	"app"
	"app/app_conf"
	"app/app_db"

	"os/signal"
	"syscall"
	"time"
	"strings"

	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	app_dir string = getWD()
	Hostip  []string

	CloseChan = make(chan os.Signal, 1)
)

func init() {

	setupCloseHandler()

	// Check/make srv.yaml
	err := server.ServerInit(app_dir)
	if err != nil {
		log.Fatal(err)
	}

	// Check/make userdb, usr.yaml & user http.FileSystem
	err = user.UserInit(app_dir)
	if err != nil {
		log.Fatal(err)
	}

	// Check/make appdb, app.yaml & app http.FileSystem
	err = app.AppInit(app_dir)
	if err != nil {
		log.Fatal(err)
	}

	app_conf.StartTime = time.Now().Unix()

	err = server.SaveHostInfo()
	if err != nil {
		log.Fatal(err)
	}

}

func readCmdLineArgs() {
	
	if len(os.Args) > 1 {
		for i, arg := range os.Args {
			switch arg {
			case "-tls":
				if i+1 < len(os.Args) {
					if os.Args[i+1] == "true" {
						srv_conf.SetVal("use_tls", true)
					} else {
						srv_conf.SetVal("use_tls", false)
					}
				}
			case "-port":
				if i+1 < len(os.Args) {
					port_int := global.StringToInt(os.Args[i+1])
					if port_int > 0 && port_int < 65536 {
						srv_conf.SetVal("server_port", port_int)
					} else {
						log.Fatalf("Invalid port number: %s", os.Args[i+1])
						log.Println("Using default port from config.")
					}
				}
			case "-debug":
				if i+1 < len(os.Args) {
					if os.Args[i+1] == "true" {
						srv_conf.SetVal("gin_mode", "debug")
					} else {
						srv_conf.SetVal("gin_mode", "release")
					}
				}
			}
		}	
	}

}

func main() {

	// Read command line args
	readCmdLineArgs()

	r := server.InitWebServer()

	startwebserver(r, srv_conf.UseTLS())

}

func startwebserver(r *gin.Engine, useTLS bool) {

	_hostinfo := srv_conf.GetHostInfo().(map[string]string)
	_port := fmt.Sprintf(":%s", _hostinfo["port"])

	if useTLS {
		go printsrvinfo("https", _hostinfo)
		err := r.RunTLS(_port, srv_sec.CertFilePath(), srv_sec.KeyFilePath())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		go printsrvinfo("http", _hostinfo)
		err := r.Run(_port)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func printsrvinfo(proto string, hostinfo map[string]string) {

	_port := fmt.Sprintf(":%s", hostinfo["port"])
	_ips := strings.SplitSeq(hostinfo["ip"], ",")
	_name := hostinfo["name"]

	log.Printf("HostName: %s", _name)
	
	for ip := range _ips {
		if proto == "https" {
			log.Printf("Webserver %s://%s%s (TLS %d)\n", proto, ip, _port, srv_conf.TLSKeySize())
		} else {

			log.Printf("Webserver %s://%s%s", proto, ip, _port)
		}
	}
}

func setupCloseHandler() {
	log.Println("SetupCloseHandler ...")
	// closeChan := make(chan os.Signal, 1) // Make the channel buffered with a capacity of 1
	signal.Notify(CloseChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-CloseChan
		log.Printf("\nClosing app ...\n")

		app_db.CloseAppDB()

		os.Exit(0)
	}()
}

func getWD() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return wd
}
