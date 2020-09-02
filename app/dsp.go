package app

import (
	"errors"
	"math/rand"
	"time"
)

type DSP struct {
	Budget                    float64 // Reset to 0 at 00:00 using lambda from aws -> https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
	Registry                  map[string]ImpressionRegistry
	Bids                      map[string]Bid // Clear values based on pending status and timestamp using expiration criteria
	MaxImpressionsPerMinute   int64
	MaxImpressionsPer3Minutes int64
}

func (dsp *DSP) notEnough(money float64) bool {
	return money > dsp.Budget
}

func (dsp *DSP) Spend(money float64) error {
	if dsp.notEnough(money) {
		return errors.New("Budget exceeded")
	} else {
		dsp.Budget -= money
	}
	return nil
}

func (dsp *DSP) GetBid(userId string, bidFloor float64) (error, *Bid) {

	bid := &Bid{}

	price := bidFloor + rand.Float64()
	if dsp.notEnough(price) {
		return errors.New("Out of budget"), bid
	}

	now := time.Now()
	if dsp.FrequencyCapped(userId, now) {
		return errors.New("User frequency is capped"), bid
	}

	timestamp := now.Unix()
	bid = CreateBid(UUID(), userId, price, timestamp)

	return nil, bid
}

func (dsp *DSP) Setup(dailyBudget float64, limitPerMinute int64, limitPer3Minute int64) {
	rand.Seed(time.Now().UnixNano())
	dsp.Budget = dailyBudget
	dsp.MaxImpressionsPerMinute = limitPerMinute
	dsp.MaxImpressionsPer3Minutes = limitPer3Minute
	dsp.Bids = make(map[string]Bid)
	dsp.Registry = make(map[string]ImpressionRegistry)
}

func (dsp *DSP) FrequencyCapped(userId string, now time.Time) bool {

	var count1MinuteImpressions int64
	var count3MinutesImpressions int64
	timeStampAMinuteAgo := now.Add(-1 * time.Minute).Unix()
	timeStamp3MinuteAgo := now.Add(-3 * time.Minute).Unix()

	if userRegistry, ok := dsp.Registry[userId]; ok {

		currentImpression := userRegistry.end

		for currentImpression.timestamp >= timeStamp3MinuteAgo {

			if currentImpression.timestamp >= timeStamp3MinuteAgo {
				count3MinutesImpressions += 1
				if count3MinutesImpressions == dsp.MaxImpressionsPer3Minutes {
					return true
				}

				if currentImpression.timestamp >= timeStampAMinuteAgo {
					count1MinuteImpressions += 1
					if count1MinuteImpressions == dsp.MaxImpressionsPerMinute {
						return true
					}
				}
			}
			currentImpression = currentImpression.previous
		}
	}
	return false
}

func (dsp *DSP) RegisterBid(bid Bid) {
	Mutex.Lock()
	dsp.Bids[bid.Id] = bid
	Mutex.Unlock()
}

func (dsp *DSP) RegisterImpression(bid Bid) {
	Mutex.Lock()
	registry := dsp.Registry[bid.UserId]
	registry.Append(&Impression{bid.Timestamp, &Impression{}})
	dsp.Registry[bid.UserId] = registry // check pointers..
	Mutex.Unlock()
}

func (dsp *DSP) UpdateBid(id string, price float64, timestamp int64) {
	// todo: check bid id
	bid := dsp.Bids[id]
	bid.Price = price
	bid.Timestamp = timestamp
	bid.Status = "won"
}
