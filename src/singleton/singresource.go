package singleton

import (
	"sync"
)

// TODO: review it, is used only on http clients
type SingResource struct {
	Once     sync.Once
	Resource interface{}
}
