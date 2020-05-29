package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/bryjammin/pipedream/lobby/multiplayer"
)

const (
	productionEnvironment string = "production"
)

func main() {

	var (
		environment = flag.String("env", "test", "environment")
	)

	flag.Parse()

	switch *environment {
	case productionEnvironment:
		fmt.Println("Production Environment on 5080...")
		fmt.Println("Listening on 5080...with TLS")
		startup()
		log.Fatal(http.ListenAndServeTLS(":5080", "", "", nil))
	default:
		fmt.Println("Listening on 5080...")
		startup()
		log.Fatal(http.ListenAndServe(":5080", nil))
	}
}

func startup() {
	server := multiplayer.NewServer()
	go server.Run()
	addRoutes(server)
}

func addRoutes(s *multiplayer.Server) {
	http.HandleFunc("/singlePlayerBlitzGame", s.CreateSinglePlayerSession)
	http.HandleFunc("/versusBlitzGame", s.FindTwoPlayerSession)
	http.HandleFunc("/aiBlitzGame", s.FindAISession)
	http.HandleFunc("/versusAiBlitzGame", s.FindVersusAISession)
	http.HandleFunc("/tutorialGame", s.FindTutorialSession)
}
