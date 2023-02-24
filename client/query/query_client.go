package query

import (
	"context"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	queryTimeout = 15 * time.Second
)

// QueryClient is an object that can be used to send queries to the Ojo node
type QueryClient struct {
	grpcEndpoint string
	grpcConn     *grpc.ClientConn
}

// NewQueryClient returns a new instance of the QueryClient
func NewQueryClient(grpcEndpoint string) (*QueryClient, error) {
	qc := &QueryClient{grpcEndpoint: grpcEndpoint}
	err := qc.dialGrpcConn()
	if err != nil {
		return nil, err
	}
	return qc, nil
}

func (qc *QueryClient) dialGrpcConn() (err error) {
	qc.grpcConn, err = grpc.Dial(
		qc.grpcEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialerFunc),
	)
	return err
}

func dialerFunc(_ context.Context, addr string) (net.Conn, error) {
	return connect(addr)
}

// connect dials the given address and returns a net.Conn.
// The protoAddr argument should be prefixed with the protocol,
// eg. "tcp://127.0.0.1:8080" or "unix:///tmp/test.sock".
func connect(protoAddr string) (net.Conn, error) {
	proto, address := protocolAndAddress(protoAddr)
	conn, err := net.Dial(proto, address)
	return conn, err
}

// protocolAndAddress splits an address into the protocol and address components.
// For instance, "tcp://127.0.0.1:8080" will be split into "tcp" and "127.0.0.1:8080".
// If the address has no protocol prefix, the default is "tcp".
func protocolAndAddress(listenAddr string) (string, string) {
	protocol, address := "tcp", listenAddr

	parts := strings.SplitN(address, "://", 2)
	if len(parts) == 2 {
		protocol, address = parts[0], parts[1]
	}

	return protocol, address
}
