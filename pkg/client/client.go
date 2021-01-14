package client

import (
	"github.com/xh3b4sd/tracer"
	"google.golang.org/grpc"

	"github.com/venturemark/apigengo/pkg/pbf/audience"
)

type Config struct {
}

type Client struct {
	connection *grpc.ClientConn

	audience audience.APIClient
}

func New(c Config) (*Client, error) {
	var err error

	var con *grpc.ClientConn
	{
		con, err = grpc.Dial("127.0.0.1:7777", grpc.WithInsecure())
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var aud audience.APIClient
	{
		aud = audience.NewAPIClient(con)
	}

	cli := &Client{
		connection: con,

		audience: aud,
	}

	return cli, nil
}

func (c *Client) Connection() *grpc.ClientConn {
	return c.connection
}

func (c *Client) Audience() audience.APIClient {
	return c.audience
}
