package utils

import (
	"fmt"
)

func NoticeJson(price float64, timestamp int64) []byte {
	return []byte(fmt.Sprintf(`{"timestamp": %d, "price": %f }`, timestamp, price))
}

func BidJson(id string, bidfloor float64, ip string, userId string) []byte {
	return []byte(fmt.Sprintf(`{
		"id": "%s",
		"imp": {
			"bidfloor": %f
		},
		"device": {
			"ip": "%s"
		},
		"user": {
			"id": "%s"
		}
	}`, id, bidfloor, ip, userId))
}
