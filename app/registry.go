package app

type ImpressionRegistry struct {
	Length int64
	start  *Impression
	end    *Impression
}

type Impression struct {
	timestamp int64
	previous  *Impression
}

func (f *ImpressionRegistry) Append(newImpression *Impression) {
	if f.Length == 0 {
		f.start = newImpression
		f.end = newImpression
	} else {
		lastImpression := f.end
		newImpression.previous = lastImpression
		f.end = newImpression
	}
	f.Length++
}
