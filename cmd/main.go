package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func main() {

	//var httpsSrv *http.Server
	mgr := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("wthpd.com"),
		Cache:      autocert.DirCache("cert-cache"),
	}

	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("Listening on 80...")
	log.Fatal(http.ListenAndServe(":80", mgr.HTTPHandler(nil)))
}
