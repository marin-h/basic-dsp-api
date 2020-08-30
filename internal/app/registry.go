package app

type ImpressionRegistry struct {
	length int
	start  *Impression
	end    *Impression
}

type Impression struct {
	timestamp int64       // Unix timestamp
	previous  *Impression // link to the previous Impression
}

func (f *ImpressionRegistry) Append(newImpression *Impression) {
	if f.length == 0 {
		f.start = newImpression
		f.end = newImpression
	} else {
		lastImpression := f.end
		newImpression.previous = lastImpression
		f.end = newImpression
	}
	f.length++
}
