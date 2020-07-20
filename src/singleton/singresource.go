package singleton

import (
	"sync"
)

type SingResource struct {
	Once     sync.Once
	Resource interface{}
}
