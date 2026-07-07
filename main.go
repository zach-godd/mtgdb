package main

import (
	"log"
	"net/http"

	"mtgdb/api"
	"mtgdb/data"
)

func main() {
	a := api.NewAPI(data.ScryfallClient{})

	mux := http.NewServeMux()
	mux.HandleFunc("/api/bulk-data", a.HandleFetchBulkData)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
