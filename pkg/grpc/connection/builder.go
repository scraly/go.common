/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package connection

import (
	"context"
	"math"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defaultWindowSize     = 65535
	initialWindowSize     = defaultWindowSize * 32 // for an RPC
	initialConnWindowSize = initialWindowSize * 16 // for a connection
)

// MakeGrpcClient creates new gRPC client connection on given addr
// and wraps it with given credentials
func MakeGrpcClient(ctx context.Context, addr string, creds credentials.TransportCredentials, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	var secureOpt = grpc.WithInsecure()
	if creds != nil {
		secureOpt = grpc.WithTransportCredentials(creds)
	}

	var extraOpts = append(opts,
		secureOpt,
		grpc.WithInitialWindowSize(initialWindowSize),
		grpc.WithInitialConnWindowSize(initialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)),
	)

	cc, err := grpc.DialContext(ctx, addr, extraOpts...)
	if err != nil {
		return nil, err
	}
	return cc, err
}

// MakeGrpcServer creates new gRPC server
func MakeGrpcServer(creds credentials.TransportCredentials, extraopts ...grpc.ServerOption) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.InitialWindowSize(initialWindowSize),
		grpc.InitialConnWindowSize(initialConnWindowSize),
		grpc.MaxConcurrentStreams(math.MaxInt32),
	}
	if creds != nil {
		opts = append(opts, grpc.Creds(creds))
	}
	srv := grpc.NewServer(append(opts, extraopts...)...)
	return srv
}

// ExtractMethod from call
func ExtractMethod(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	return parts[len(parts)-1]
}
