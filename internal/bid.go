package main

type Bid struct {
	id        string `json:"-"`
	userId    string `json:"-"`
	timestamp int64  `json:"-"`
	price     float64
	nurl      string
}
