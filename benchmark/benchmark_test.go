package benchmark

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/marin-h/simple-dsp/app"
	"github.com/marin-h/simple-dsp/cmd"
	"github.com/marin-h/simple-dsp/utils"
)

var BidRequests []*http.Request
var WinRequests []*http.Request

func init() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Setting up data...")
	for i := 0; i <= 500000; i++ {
		auctionId := app.UUID()
		bidId := app.UUID()
		userId := app.UUID()
		timestamp := time.Now().Add(time.Duration(rand.Int()-10) * time.Second).Unix()

		bidPayload := utils.BidJson(auctionId, 0.5, "100.123.230.3", userId)
		bidReq, _ := http.NewRequest("POST", "/bid", bytes.NewBuffer(bidPayload))
		bidReq.Header.Set("Content-Type", "application/json")
		BidRequests = append(BidRequests, bidReq)

		cmd.Dsp.RegisterBid(*app.CreateBid(bidId, userId, 0.5, timestamp))

		winPayload := utils.NoticeJson(0.00001, timestamp)
		winReq, _ := http.NewRequest("POST", "/winningnotice?bidid="+bidId, bytes.NewBuffer(winPayload))
		winReq.Header.Set("Content-Type", "application/json")
		winReq.Close = true
		WinRequests = append(WinRequests, winReq)
	}
	fmt.Println("All set!")
}

func Benchmark_BidRequestDifferentUsers(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			http.HandlerFunc(cmd.HandleBid).ServeHTTP(httptest.NewRecorder(), BidRequests[i])
		}
	})
}

func Benchmark_BidRequestSameUser(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			http.HandlerFunc(cmd.HandleBid).ServeHTTP(httptest.NewRecorder(), BidRequests[0])
		}
	})
}

func Benchmark_WinRequest(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			rr := httptest.NewRecorder()
			http.HandlerFunc(cmd.HandleNotice).ServeHTTP(rr, WinRequests[i])
			rr.Result()
		}
	})
}
