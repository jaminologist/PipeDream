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
	http.HandleFunc("/connectToServer", server.HandleNewConnection)
	http.HandleFunc("/singlePlayerBlitzGame", server.CreateSinglePlayerSession)
	http.HandleFunc("/versusBlitzGame", server.FindTwoPlayerSession)

	//Test git webhook
	log.Fatal(http.ListenAndServe(":5080", nil))
}
