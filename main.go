package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:6600",
	}
	http.HandleFunc("/", index)
	server.ListenAndServe()
}
