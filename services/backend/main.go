package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	// Получаем версию из переменной окружения (понадобится для Istio)
	version := os.Getenv("VERSION")
	if version == "" {
		version = "v1"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response from Backend (%s)\n", version)
	})

	fmt.Printf("Backend %s is running on port 8080...\n", version)
	http.ListenAndServe(":8080", nil)
}
