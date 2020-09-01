package test

import (
	"fmt"
	"testing"
	"time"

	cmd "github.com/marin-h/simple-dsp/cmd"
)

func TestUserIsCappedWith1MinuteLimit(t *testing.T) {

	var i int8
	for i = 0; i < cmd.Dsp.MaxImpressionsPerMinute-1; i++ {
		TestHandleBid(t)
		TestWinNotice(t)
	}
	fmt.Printf("%+v\n", cmd.Dsp.Registry["f345nf0k"])

	if !cmd.Dsp.FrequencyCapped("f345nf0k", time.Now()) {
		t.Errorf("User should by capped by now!")
	}
}
