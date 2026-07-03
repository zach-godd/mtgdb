package main

import (
	"log"
	"net/http"

	"mtgdb/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/bulk-data", api.HandleFetchBulkData)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
