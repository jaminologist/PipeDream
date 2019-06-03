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
	http.HandleFunc("/connectToServer", server.HandleNewConnection)
	go server.Run()
	//Test git webhook
	log.Fatal(http.ListenAndServe(":5080", nil))
}
