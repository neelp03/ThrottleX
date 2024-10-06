// ratelimiter/mutex.go

package ratelimiter

import (
	"sync"
	"time"
)

// keyMutex holds a mutex and the last access time.
type keyMutex struct {
	mu         *sync.Mutex
	lastAccess time.Time
}
