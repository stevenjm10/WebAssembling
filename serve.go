//go:build !js

// Minimal static file server for the web/ folder.
// Browsers require .wasm files to be served over HTTP with the
// "application/wasm" content type; opening index.html with file:// does not work.
package main

import (
	"log"
	"mime"
	"net/http"
	"os"
)

func main() {
	_ = mime.AddExtensionType(".wasm", "application/wasm")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Serving ./web at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, http.FileServer(http.Dir("web"))))
}
