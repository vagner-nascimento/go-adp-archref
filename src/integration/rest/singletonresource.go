package rest

import (
	"sync"
)

// TODO generalize it to use in all connections
type singletonResource struct {
	once     sync.Once
	resource interface{}
}
