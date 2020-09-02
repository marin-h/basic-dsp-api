package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cmd "github.com/marin-h/simple-dsp/cmd"
	"github.com/marin-h/simple-dsp/utils"
)

var bidData cmd.BidData

func TestHandleBid(t *testing.T) {

	err, rr := PostBidRequest(utils.BidJson("3cb627bf", 0.5, "100.123.230.3", "a"))

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

	err, rr := PostWinNotice(bidData.Imp.Nurl, utils.NoticeJson(0.5, time.Now().Unix()))

	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func PostWinNotice(nurl string, jsonStr []byte) (error, *httptest.ResponseRecorder) {

	req, err := http.NewRequest("POST", nurl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err, httptest.NewRecorder()
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cmd.HandleNotice)
	handler.ServeHTTP(rr, req)
	return nil, rr
}

func PostBidRequest(jsonStr []byte) (error, *httptest.ResponseRecorder) {

	req, err := http.NewRequest("POST", "/bid", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err, httptest.NewRecorder()
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cmd.HandleBid)
	handler.ServeHTTP(rr, req)

	return nil, rr
}
