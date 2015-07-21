package main

import (
	"log"
	"net/http"
)

func main() {

	emulators.init()
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
