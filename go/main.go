package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NikitaLitvishko/Genesis-school-API-service/go/services"
	"github.com/gorilla/mux"

)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/rate", services.GetRate).Methods("GET")
	router.HandleFunc("/api/subscribe", services.SubscribeEmail).Methods("POST")
	router.HandleFunc("/api/sendEmails", services.SendEmails).Methods("POST")

	port := ":3000"
	fmt.Println("Server is listening on port", port)
	log.Fatal(http.ListenAndServe(port, router))
}
