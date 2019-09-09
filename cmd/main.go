package main

import (
	"fmt"
	"log"
	"net/http"

	"bryjamin.com/multiplayer"
)

func main() {
	fmt.Println("Listening on 5080...")
	server := multiplayer.NewServer()
	go server.Run()
	http.HandleFunc("/singlePlayerBlitzGame", server.CreateSinglePlayerSession)
	http.HandleFunc("/versusBlitzGame", server.FindTwoPlayerSession)
	http.HandleFunc("/aiBlitzGame", server.FindAISession)

	log.Fatal(http.ListenAndServe(":5080", nil))
}
