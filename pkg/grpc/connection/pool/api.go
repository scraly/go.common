package pool

import (
	"errors"

	"google.golang.org/grpc"
)

// Pool defines a generic pool contract
type Pool interface {
	Get() (*Connection, error)
	Close() error
	Len() int
}

// Factory is the GRPC connection factory
type Factory func() (*grpc.ClientConn, error)

var (
	// ErrClosed is raised when the connection is closed
	ErrClosed = errors.New("pool is closed")
)
