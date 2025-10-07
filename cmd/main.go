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

	"fmt"
	"log"
	"os"
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

	addr := fmt.Sprintf(":%d", srv_conf.GetInt("server_port"))

	if srv_conf.UseTLS() {
		go printsrvinfo("https", addr)
		r.RunTLS(addr, srv_sec.CertFile, srv_sec.KeyFile)
	} else {
		go printsrvinfo("http", addr)
		r.Run(addr)
	}

}

func getWD() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return wd
}

func printsrvinfo(proto, adr string) {

	Hostip, _ = global.GetIPv4Addresses()

	for _, ip := range Hostip {
		if proto == "https" {
			log.Printf("Webserver %s://%s%s (TLS %d)", proto, ip, adr, srv_conf.TLSKeySize())
		} else {
			log.Printf("Webserver %s://%s%s", proto, ip, adr)
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
