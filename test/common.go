package test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/marin-h/simple-dsp/cmd"
)

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

func NoticeJson(price float64, timestamp int64) []byte {
	return []byte(fmt.Sprintf(`{"timestamp": %d, "price": %e }`, timestamp, price))
}

func BidJson(id string, bidfloor float64, ip string, userId string) []byte {
	return []byte(fmt.Sprintf(`{
		"id": "%s",
		"imp": {
			"bidfloor": %e
		},
		"device": {
			"ip": "%s"
		},
		"user": {
			"id": "%s"
		}
	}`, id, bidfloor, ip, userId))
}
