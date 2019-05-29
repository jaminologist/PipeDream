package main

import (
	"fmt"
	"log"
	"net/http"

	"bryjamin.com/multiplayer"
)

func main() {
	fmt.Println("Listening on 8080...")
	server := multiplayer.NewServer()
	http.HandleFunc("/connectToServer", server.HandleNewConnection)
	go server.Run()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
