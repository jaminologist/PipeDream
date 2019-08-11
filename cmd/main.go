package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("Listening on 80...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
