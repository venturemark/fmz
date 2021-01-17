package client

import (
	"github.com/venturemark/apigengo/pkg/pbf/audience"
	"github.com/venturemark/apigengo/pkg/pbf/timeline"
	"github.com/xh3b4sd/tracer"
	"google.golang.org/grpc"
)

type Config struct {
	Address string
}

type Client struct {
	connection *grpc.ClientConn

	audience audience.APIClient
	timeline timeline.APIClient
}

func New(c Config) (*Client, error) {
	if c.Address == "" {
		c.Address = "127.0.0.1:7777"
	}

	var err error

	var con *grpc.ClientConn
	{
		con, err = grpc.Dial(c.Address, grpc.WithInsecure())
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var aud audience.APIClient
	{
		aud = audience.NewAPIClient(con)
	}

	var tim timeline.APIClient
	{
		tim = timeline.NewAPIClient(con)
	}

	cli := &Client{
		connection: con,

		audience: aud,
		timeline: tim,
	}

	return cli, nil
}

func (c *Client) Connection() *grpc.ClientConn {
	return c.connection
}

func (c *Client) Audience() audience.APIClient {
	return c.audience
}

func (c *Client) Timeline() timeline.APIClient {
	return c.timeline
}
