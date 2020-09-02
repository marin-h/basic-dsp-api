package app

import (
	"sync"

	"github.com/google/uuid"
)

var Mutex *sync.Mutex

func init() {
	Mutex = &sync.Mutex{}
}

func UUID() string {
	return uuid.New().String()
}
