package server

import (
	"log"
	"srv/filefunc"
	"srv/srv_sec"
)

func CheckTLS(app_path string, keysize int) {

	srv_sec.CertFile = app_path + "/app.crt"
	srv_sec.KeyFile = app_path + "/app.key"
	if !filefunc.IsExists(srv_sec.CertFile) || !filefunc.IsExists(srv_sec.KeyFile) {
		log.Println("No RSA files found. Creating key pair ...")

		err := srv_sec.GenerateTLS(keysize)
		if err != nil {
			log.Fatal(err)
		}
	}
}
