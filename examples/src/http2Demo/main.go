package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/header", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-custom-header", "custom header")
		w.WriteHeader(http.StatusNoContent)

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		select {}
	})

	log.Println("start listen on 8080...")
	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil))
}
