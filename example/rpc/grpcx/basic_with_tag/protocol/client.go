package protocol

import (
	"google.golang.org/grpc"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

const (
	// AppID is the application ID for the protobuf service.
	AppID = "demo"
)

// Client is the client for protobuf.
type Client struct {
	conn *grpc.ClientConn
}

// NewClient creates and returns a new client.
func NewClient(options ...grpc.DialOption) (*Client, error) {
	conn, err := grpcx.Client.NewGrpcClientConn(AppID, options...)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

// Echo is the client for protobuf.Echo.
func (c *Client) Echo() EchoClient {
	return NewEchoClient(c.conn)
}
