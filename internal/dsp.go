package main

type DSP struct {
	Budget                    float64 // Reset to 0 at 00:00 using lambda from aws -> https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
	Registry                  map[string]ImpressionRegistry
	Bids                      map[string]Bid // Will need to implement some value-cleaning
	MaxImpressionsPerMinute   int8
	MaxImpressionsPer3Minutes int8
}
