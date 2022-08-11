package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func serveHTTP(wr http.ResponseWriter, req *http.Request) {
	http.Error(wr, "Client Error: Headers", http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHTTP)
	logrus.Info("Server online...")
	logrus.Fatal("ListenAndServe:", http.ListenAndServe(":80", mux))
}
