package main

import (
	"fmt"
	"io"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	fmt.Println("===== WEBHOOK RECEIVED =====")
	fmt.Println(string(body))
	fmt.Println("============================")

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/my/webhook", handler)

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}