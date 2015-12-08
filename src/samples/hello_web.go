package main

import (
	"io"
	"net/http"
)

func Hello(writer http.ResponseWriter, request *http.Request) {
	io.WriteString(writer, "hello web !")
}

func main() {
	http.HandleFunc("/", Hello)
	http.ListenAndServe(":8080", nil)
}
