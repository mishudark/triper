package nats

import (
	"encoding/json"
	"strings"

	"github.com/mishudark/triper"
	nats "github.com/nats-io/nats.go"
)

// Client nats
type Client struct {
	Options nats.Options
}

// NewClient returns the basic client to access to nats
func NewClient(urls string, useTLS bool) (*Client, error) {
	opts := nats.DefaultOptions
	opts.Secure = useTLS
	opts.Servers = strings.Split(urls, ",")

	for i, s := range opts.Servers {
		opts.Servers[i] = strings.Trim(s, " ")
	}

	return &Client{
		opts,
	}, nil
}

// Publish a event
func (c *Client) Publish(event triper.Event, bucket, subset string) error {
	nc, err := c.Options.Connect()
	if err != nil {
		return err
	}

	defer nc.Close()

	blob, err := json.Marshal(event)
	if err != nil {
		return err
	}

	subj := bucket + "." + subset
	nc.Publish(subj, blob)
	nc.Flush()

	return nc.LastError()
}
