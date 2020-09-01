package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	app "github.com/marin-h/simple-dsp/app"
)

var Dsp app.DSP

func init() {
	// initialize Dsp
	Dsp = app.DSP{}
	Dsp.Setup(10, 5, 10)
	fmt.Println("Dsp setup")
	fmt.Printf("%+v\n", Dsp)
}

func Run() {

	// swagger
	fs := http.FileServer(http.Dir("./doc"))
	http.Handle("/doc/", http.StripPrefix("/doc/", fs))

	// Dsp endpoints
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
	Price float64 `json:"price,omitempty"`
	Nurl  string  `json:"nurl,omitempty"`
}

type WinNotice struct {
	Timestamp int64   `json:"timestamp"`
	Price     float64 `json:"price"`
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

		err, bid := Dsp.GetBid(auction.User.Id, auction.Imp.Bidfloor)
		if err == nil {
			Dsp.RegisterBid(*bid)
			responseBody := BidData{auction.Id, bid.Id, ImpressionData{Price: bid.Price, Nurl: r.Host + "/winningnotice?bidid=" + bid.Id}}
			json.NewEncoder(w).Encode(responseBody)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			fmt.Println("Error:", err.Error())
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
		}
		notice := WinNotice{}
		if err := json.NewDecoder(r.Body).Decode(&notice); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println("Error:", err.Error())
			return
		}

		fmt.Println("Win Notice request received")
		fmt.Printf("%+v\n", notice)

		if bid, ok := Dsp.Bids[id]; ok {

			err := Dsp.Spend(notice.Price)

			if err != nil {
				http.Error(w, err.Error(), http.StatusPreconditionFailed)
				fmt.Println("Error:", err.Error())
				return
			}
			fmt.Printf("%+v\n", bid)
			Dsp.RegisterImpression(bid)
			Dsp.UpdateBid(bid.Id, notice.Price, notice.Timestamp)

			w.WriteHeader(http.StatusOK)
			return
		} else {
			http.Error(w, "Bid not found", http.StatusNotFound)
			fmt.Println("Error: Bid not found")
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	return
}
