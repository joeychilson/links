package main

import (
	"log"

	"github.com/joeychilson/starter-templ/server"
)

func main() {
	server := server.New()

	log.Println("Serving application @ http://localhost:8080")
	if err := server.ListenAndServe(":8080"); err != nil {
		log.Fatal(err)
	}
}
