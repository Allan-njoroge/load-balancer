package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Name string
	Port int
}

func main() {
	name := flag.String("name", "server", "server name")
	port := flag.Int("port", 8081, "server port")

	flag.Parse()

	server := Server{
		Name: *name,
		Port: *port,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from %s\n", server.Name)
	})

	address := fmt.Sprintf(":%d", server.Port)

	log.Printf("%s is listening on %s", server.Name, address)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
