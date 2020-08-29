package main

import (
	"errors"
	"math/rand"
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

	price := rand.Float64()
	if dsp.notEnough(price) {
		return errors.New("Out of budget"), bid
	}

	return nil, bid
}

func (dsp *DSP) setup(dailyBudget float64, limitPerMinute int8, limitPer3Minute int8) {
	dsp.Budget = dailyBudget
	dsp.MaxImpressionsPer3Minutes = limitPer3Minute
	dsp.MaxImpressionsPerMinute = limitPerMinute
}
