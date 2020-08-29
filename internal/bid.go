package main

type Bid struct {
	Id        string `json:"-"`
	UserId    string `json:"-"`
	Timestamp int64  `json:"-"`
	Price     float64
	Nurl      string
}
