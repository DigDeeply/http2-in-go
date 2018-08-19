package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
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

	// 用于push的 handler
	http.HandleFunc("/crt", func(w http.ResponseWriter, r *http.Request) {
		tpl := template.Must(template.ParseFiles("server.crt"))
		tpl.Execute(w, nil)
	})

	// 请求该Path会触发Push
	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		pusher, ok := w.(http.Pusher)
		if !ok {
			log.Println("not support server push")
		} else {
			err := pusher.Push("/crt", nil)
			if err != nil {
				log.Printf("Failed for server push: %v", err)
			}
		}
		w.WriteHeader(http.StatusOK)
	})

	// 服务端定时自己push内容
	http.HandleFunc("/autoPush", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-custom-header", "custom")
		w.WriteHeader(http.StatusNoContent)

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		pusher, ok := w.(http.Pusher)
		if ok {
			for {
				select {
				case <-time.Tick(5 * time.Second):
					err := pusher.Push("/crt", nil)
					if err != nil {
						log.Printf("Failed for server push: %v", err)
					}
				}
			}
		}
	})

	log.Println("start listen on 8080...")
	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil))
}
