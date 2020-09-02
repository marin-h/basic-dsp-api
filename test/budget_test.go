package test

import (
	"net/http"
	"testing"

	cmd "github.com/marin-h/simple-dsp/cmd"
	"github.com/marin-h/simple-dsp/utils"
)

func TestHandleBidOutOfBudget(t *testing.T) {

	cmd.Dsp.Spend(9.7)

	err, rr := PostBidRequest(utils.BidJson("3cb627bf", 0.5, "100.123.230.3", "a"))

	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
	cmd.Dsp.Budget = 10
}
