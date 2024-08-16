package shared

import (
	"sync"
)

// Define a global Mutex or RWMutex that can be used across packages
var Mu sync.Mutex

// Alternatively, if you need to support concurrent reads:
// var RWmu sync.RWMutex
