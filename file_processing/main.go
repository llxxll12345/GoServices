package main

import (
	"filemaster/handlers"
	"filemaster/processor"
	"fmt"
	"net/http"
)

func main() {
	processor.StartJobs()
	http.HandleFunc("/upload", handlers.Upload)
	http.HandleFunc("/status", handlers.Status)
	http.HandleFunc("/download", handlers.Download)

	fmt.Println("Serving at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(fmt.Errorf("Server error: %s", err))
	}
}
