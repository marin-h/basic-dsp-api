package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	app "simple-dsp/internal/app"
)

var dsp app.DSP

func main() {

	// initialize dsp
	dsp = app.DSP{}
	dsp.Setup(10, 5, 10)
	fmt.Println("DSP setup")
	fmt.Printf("%+v\n", dsp)

	// swagger
	fs := http.FileServer(http.Dir("./doc"))
	http.Handle("/doc/", http.StripPrefix("/doc/", fs))

	// dsp endpoints
	http.HandleFunc("/bid", HandleBid)
	http.HandleFunc("/winningnotice", HandleNotice)

	log.Println("Listening in port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

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

type BidData struct {
	Id    string         `json:"id"`
	BidId string         `json:"bidid"`
	Imp   ImpressionData `json:"bid"`
}

type ImpressionData struct {
	Price     float64 `json:"price,omitempty"`
	Timestamp int64   `json:"timestamp,omitempty"`
	Nurl      string  `json:"nurl,omitempty"`
}

func HandleBid(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		auction := AuctionData{}
		if err := json.NewDecoder(r.Body).Decode(&auction); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Auction request received")
		fmt.Printf("%+v\n", auction)

		err, bid := dsp.GetBid(auction.User.Id, auction.Imp.Bidfloor)
		if err == nil {
			dsp.RegisterBid(*bid)
			responseBody := BidData{auction.Id, bid.Id, ImpressionData{Price: bid.Price, Nurl: "/winningnotice?bidid=" + bid.Id}}
			json.NewEncoder(w).Encode(responseBody)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	return
}

func HandleNotice(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		id := r.URL.Query().Get("bidid")

		if id == "" {
			http.Error(w, "Request malformed", http.StatusBadRequest)
			return
		}

		if bid, ok := dsp.Bids[id]; ok {

			impressionData := ImpressionData{}
			if err := json.NewDecoder(r.Body).Decode(&impressionData); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			fmt.Println("Win notice received")
			fmt.Printf("%+v\n", impressionData)

			err := dsp.Spend(impressionData.Price)

			if err != nil {
				http.Error(w, err.Error(), http.StatusPreconditionFailed)
				return
			}
			dsp.RegisterImpression(bid)
			dsp.UpdateBid(bid.Id, impressionData.Price, impressionData.Timestamp)

			w.WriteHeader(http.StatusOK)
			return
		} else {
			http.Error(w, "Bid not found", http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	return
}
