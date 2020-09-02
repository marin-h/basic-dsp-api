package test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	cmd "github.com/marin-h/simple-dsp/cmd"
)

var bidData cmd.BidData

func TestHandleBid(t *testing.T) {

	err, rr := PostBidRequest(BidJson("3cb627bf", 0.5, "100.123.230.3", "a"))

	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if err := json.NewDecoder(rr.Body).Decode(&bidData); err != nil {
		t.Errorf("Couldn't parse response to BidData struct")
	}
}

func TestWinNotice(t *testing.T) {

	err, rr := PostWinNotice(bidData.Imp.Nurl, NoticeJson(0.5, time.Now().Unix()))

	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
