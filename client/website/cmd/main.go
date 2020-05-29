package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Listening on 17000...")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":17000", nil))
}
