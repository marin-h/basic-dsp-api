package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/bid", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			auction := AuctionRequest{}
			if err := json.NewDecoder(r.Body).Decode(&auction); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Println(auction)
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	})
	log.Println("Listening in port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// JSON Structs
type AuctionRequest struct {
	id     string
	imp    Bidfloor
	device Device
	user   User
}

type Device struct {
	ip string
}

type Bidfloor struct {
	bidfloor float64
}

type User struct {
	id string
}
