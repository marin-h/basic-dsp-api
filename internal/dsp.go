package main

import (
	"errors"
	"math/rand"
	"time"
)

type DSP struct {
	Budget                    float64 // Reset to 0 at 00:00 using lambda from aws -> https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
	Registry                  map[string]ImpressionRegistry
	Bids                      map[string]Bid // Will need to implement some value-cleaning
	MaxImpressionsPerMinute   int8
	MaxImpressionsPer3Minutes int8
}

func (dsp *DSP) notEnough(money float64) bool {
	return money > dsp.Budget
}

func (dsp *DSP) spend(money float64) error {
	if dsp.notEnough(money) {
		return errors.New("Budget exceeded")
	} else {
		dsp.Budget -= money
	}
	return nil
}

func (dsp *DSP) getBid(userId string, bidFloor float64) (error, *Bid) {

	bid := &Bid{}

	price := bidFloor + rand.Float64()
	if dsp.notEnough(price) {
		return errors.New("Out of budget"), bid
	}

	now := time.Now()
	if dsp.frequencyCapped(userId, now) {
		return errors.New("User frequency is capped"), bid
	}

	timestamp := now.Unix()
	bid = createBid(getBidId(), userId, price, timestamp)

	return nil, bid
}

func (dsp *DSP) setup(dailyBudget float64, limitPerMinute int8, limitPer3Minute int8) {
	dsp.Budget = dailyBudget
	dsp.MaxImpressionsPer3Minutes = limitPer3Minute
	dsp.MaxImpressionsPerMinute = limitPerMinute
}

func (dsp *DSP) frequencyCapped(userId string, now time.Time) bool {

	var count1MinuteImpressions int8
	var count3MinutesImpressions int8

	timeStampAMinuteAgo := now.Add(-1 * time.Minute).Unix()
	timeStamp3MinuteAgo := now.Add(-3 * time.Minute).Unix()

	currentImpression := dsp.Registry[userId].end

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
	return false
}
