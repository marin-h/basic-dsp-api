package main

type ImpressionRegistry struct {
	length int
	start  *Impression
	end    *Impression
}

type Impression struct {
	timestamp int64       // Unix timestamp
	previous  *Impression // link to the previous Impression
}
