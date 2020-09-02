package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cmd "github.com/marin-h/simple-dsp/cmd"
)

var bidData cmd.BidData

func TestHandleBid(t *testing.T) {

	var jsonStr = []byte(`{
			"id": "3cb627bf",
			"imp": {
			"bidfloor": 0.5
			},
			"device": {
			"ip": "100.123.230.3"
			},
			"user": {
			"id": "f345nf0k"
			}
		}`)

	req, err := http.NewRequest("POST", "/bid", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cmd.HandleBid)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if err := json.NewDecoder(rr.Body).Decode(&bidData); err != nil {
		t.Errorf("Couldn't parse response to BidData struct")
	}
}

func TestWinNotice(t *testing.T) {

	PostWinNotice(t, 0)
}

func PostWinNotice(t *testing.T, timestamp int64) {

	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}

	var jsonStr = []byte(fmt.Sprintf(`{
		"timestamp": %d,
		"price": 0.5
	}`, timestamp))

	req, err := http.NewRequest("POST", bidData.Imp.Nurl, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cmd.HandleNotice)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
