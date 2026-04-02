package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// В Kubernetes мы будем обращаться к бэкенду по его DNS-имени сервиса
		resp, err := http.Get("http://backend-service:8080")
		if err != nil {
			http.Error(w, "Failed to call backend: "+err.Error(), 500)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "Frontend received: %s", string(body))
	})

	fmt.Println("Frontend is running on port 8081...")
	http.ListenAndServe(":8081", nil)
}
