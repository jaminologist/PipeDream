package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	/*mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("/certs"),
		HostPolicy: autocert.HostWhitelist("www.wthpd.com"),
	}
	server := &http.Server{
		Addr:    ":443",
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}*/

	fmt.Println("Listening on 17000...")

	//go func() {

	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":17000", nil))
	//}()

	//fmt.Println("Hold up?")

	//server.ListenAndServeTLS("", "")
}
