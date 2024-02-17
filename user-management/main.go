package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	listenAddr := flag.String("HTTP listenAddr", ":3000", "the listen address of HTTP server")
	flag.Parse()
	fmt.Println("server is listening at", ":3000")
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
