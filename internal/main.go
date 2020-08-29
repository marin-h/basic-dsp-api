package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

	fs := http.FileServer(http.Dir("./doc"))
	http.Handle("/doc/", http.StripPrefix("/doc/", fs))

	http.HandleFunc("/bid", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			auction := AuctionData{}
			if err := json.NewDecoder(r.Body).Decode(&auction); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Println("Bid posted")
			fmt.Printf("%+v\n", auction)
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotImplemented)
		return
	})
	log.Println("Listening in port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// JSON Structs
type AuctionData struct {
	Id     string       `json:"id"`
	Imp    BidfloorData `json:"imp"`
	Device DeviceData   `json:"device"`
	User   UserData     `json:"user"`
}

type DeviceData struct {
	Ip string `json:"ip"`
}

type BidfloorData struct {
	Bidfloor float64 `json:"bidfloor"`
}

type UserData struct {
	Id string `json:"id"`
}
