package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func serveHTTP(wr http.ResponseWriter, req *http.Request) {
	http.Error(wr, "Client Error: Headers", http.StatusOK)
}

func server() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHTTP)
	logrus.Info("Server online...")
	logrus.Fatal("ListenAndServe:", http.ListenAndServe(":80", mux))
}
