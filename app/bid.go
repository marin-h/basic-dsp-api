package app

type Bid struct {
	Id        string  `json:"-"`
	UserId    string  `json:"-"`
	Timestamp int64   `json:"timestamp"`
	Status    string  `json:"-"`
	Price     float64 `json:"price"`
}

func CreateBid(id string, userId string, price float64, timestamp int64) *Bid {
	return &Bid{id, userId, timestamp, "pending", price}
}
