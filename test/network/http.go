package network

import (
	"net/http"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}

func MakeHttpServer() {
	http.HandleFunc("/hww", myHandler)
	http.ListenAndServe(":8090", nil)
}
