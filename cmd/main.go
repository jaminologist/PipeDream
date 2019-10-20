package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"bryjamin.com/multiplayer"
	"golang.org/x/crypto/acme/autocert"
)

const (
	productionEnvironment string = "production"
)

func main() {

	var (
		environment = flag.String("env", "test", "environment")
	)

	flag.Parse()

	if *environment == productionEnvironment {
		fmt.Println("Production Environment on 5080...")

		mux := http.NewServeMux()
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
		}

		go func() {
			log.Fatal(http.ListenAndServe(":5080", certManager.HTTPHandler(nil)))
		}()

		fmt.Println("Listening on 5080...with TLS")

		server.ListenAndServeTLS("", "")

	} else {
		fmt.Println("Listening on 5080...")
		server := multiplayer.NewServer()
		go server.Run()
		http.HandleFunc("/singlePlayerBlitzGame", server.CreateSinglePlayerSession)
		http.HandleFunc("/versusBlitzGame", server.FindTwoPlayerSession)
		http.HandleFunc("/aiBlitzGame", server.FindAISession)
		http.HandleFunc("/versusAiBlitzGame", server.FindVersusAISession)
		log.Fatal(http.ListenAndServe(":5080", nil))
	}
}
