package nats

import (
	"encoding/json"
	"strings"

	"github.com/mishudark/triper"
	nats "github.com/nats-io/nats.go"
)

// Client nats
type Client struct {
	nc *nats.Conn
}

// NewClient returns the basic client to access to nats
func NewClient(urls string, useTLS bool) (*Client, error) {
	opts := nats.DefaultOptions
	opts.Secure = useTLS
	opts.Servers = strings.Split(urls, ",")

	for i, s := range opts.Servers {
		opts.Servers[i] = strings.Trim(s, " ")
	}

	nc, err := opts.Connect()
	if err != nil {
		return nil, err
	}

	return &Client{
		nc: nc,
	}, nil
}

// Publish a event
func (c *Client) Publish(event triper.Event, bucket, subset string) error {
	blob, err := json.Marshal(event)
	if err != nil {
		return err
	}

	subj := bucket + "." + subset
	c.nc.Publish(subj, blob)
	c.nc.Flush()

	return c.nc.LastError()
}
