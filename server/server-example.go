package main

import (
	"log"
	"net/http"
)

func serveHTTP(wr http.ResponseWriter, req *http.Request) {
	http.Error(wr, "Client Error: Headers", http.StatusBadRequest)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHTTP)
	log.Fatal("ListenAndServe:", http.ListenAndServe(":80", mux))
}
