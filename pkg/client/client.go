package client

import (
	"github.com/venturemark/apigengo/pkg/pbf/audience"
	"github.com/venturemark/apigengo/pkg/pbf/message"
	"github.com/venturemark/apigengo/pkg/pbf/texupd"
	"github.com/venturemark/apigengo/pkg/pbf/timeline"
	"github.com/venturemark/apigengo/pkg/pbf/update"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/pkg/client"
	"github.com/xh3b4sd/tracer"
	"google.golang.org/grpc"
)

type Config struct {
	Address string
}

type Client struct {
	grpc   *grpc.ClientConn
	redigo redigo.Interface

	audience audience.APIClient
	message  message.APIClient
	texupd   texupd.APIClient
	timeline timeline.APIClient
	update   update.APIClient
}

func New(c Config) (*Client, error) {
	if c.Address == "" {
		c.Address = "127.0.0.1:7777"
	}

	var err error

	var con *grpc.ClientConn
	{
		con, err = grpc.Dial(c.Address, Credential())
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var red redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
		}

		red, err = client.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var aud audience.APIClient
	{
		aud = audience.NewAPIClient(con)
	}

	var mes message.APIClient
	{
		mes = message.NewAPIClient(con)
	}

	var tex texupd.APIClient
	{
		tex = texupd.NewAPIClient(con)
	}

	var tim timeline.APIClient
	{
		tim = timeline.NewAPIClient(con)
	}

	var upd update.APIClient
	{
		upd = update.NewAPIClient(con)
	}

	cli := &Client{
		grpc:   con,
		redigo: red,

		audience: aud,
		message:  mes,
		texupd:   tex,
		timeline: tim,
		update:   upd,
	}

	return cli, nil
}

func (c *Client) Grpc() *grpc.ClientConn {
	return c.grpc
}

func (c *Client) Redigo() redigo.Interface {
	return c.redigo
}

func (c *Client) Audience() audience.APIClient {
	return c.audience
}

func (c *Client) Message() message.APIClient {
	return c.message
}

func (c *Client) TexUpd() texupd.APIClient {
	return c.texupd
}

func (c *Client) Timeline() timeline.APIClient {
	return c.timeline
}

func (c *Client) Update() update.APIClient {
	return c.update
}
