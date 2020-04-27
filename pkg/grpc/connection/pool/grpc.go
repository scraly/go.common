package pool

import "google.golang.org/grpc"

// Connection is a gRPC connection wrapper
type Connection struct {
	Conn *grpc.ClientConn
	c    *grpcPool
}

// Close and release the connection to the pool
func (p Connection) Close() error {
	return p.c.put(p.Conn)
}
