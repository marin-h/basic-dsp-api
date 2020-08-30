package main

import (
	"github.com/google/uuid"
)

type Bid struct {
	Id        string  `json:"-"`
	UserId    string  `json:"-"`
	Timestamp int64   `json:"timestamp"`
	Status    string  `json:"-"`
	Price     float64 `json:"price"`
}

func createBid(id string, userId string, price float64, timestamp int64) *Bid {
	return &Bid{id, userId, timestamp, "pending", price}
}

func getBidId() string {
	return uuid.New().String()
}
