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

func main() {

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
