package pool

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/scraly/go.common/pkg/log"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
)

type grpcPool struct {
	mu      sync.Mutex
	conns   chan io.Closer
	factory Factory
}

// NewGrpcPool returns a gRPC connection pool
func NewGrpcPool(initialCap, maxCap int, factory Factory) (Pool, error) {
	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}

	c := &grpcPool{
		conns:   make(chan io.Closer, maxCap),
		factory: factory,
	}

	// Prefill connection pool
	for i := 0; i < initialCap; i++ {
		conn, err := factory()
		if err != nil {
			log.SafeClose(c, "Unable to connection")
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- conn
	}

	return c, nil
}

// -----------------------------------------------------------------------------

func (c *grpcPool) Get() (*Connection, error) {
	conns := c.getConns()
	if conns == nil {
		return nil, ErrClosed
	}

	// wrap our connections with out custom grpc.ClientConn implementation (wrapConn
	// method) that puts the connection back to the pool if it's closed.
	select {
	case conn := <-conns:
		if conn == nil {
			return nil, ErrClosed
		}

		return c.wrapConn(conn.(*grpc.ClientConn)), nil
	default:
		conn, err := c.factory()
		if err != nil {
			return nil, err
		}

		return c.wrapConn(conn), nil
	}
}

func (c *grpcPool) Close() error {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.mu.Unlock()

	if conns == nil {
		return nil
	}

	var err error
	close(conns)
	for conn := range conns {
		err = multierr.Append(err, conn.Close())
	}

	return err
}

func (c *grpcPool) Len() int {
	return len(c.getConns())
}

// -----------------------------------------------------------------------------

func (c *grpcPool) wrapConn(conn *grpc.ClientConn) *Connection {
	return &Connection{Conn: conn, c: c}
}

func (c *grpcPool) getConns() chan io.Closer {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()
	return conns
}

func (c *grpcPool) put(conn io.Closer) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conns == nil {
		// pool is closed, close passed connection
		return conn.Close()
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case c.conns <- conn:
		return nil
	default:
		// pool is full, close passed connection
		return conn.Close()
	}
}
