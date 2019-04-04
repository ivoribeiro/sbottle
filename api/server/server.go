package server //import "sbottle/api/server"

import (
	"fmt"
	"net/http"
)

func serveAPI() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	http.ListenAndServe(":80", nil)
}
