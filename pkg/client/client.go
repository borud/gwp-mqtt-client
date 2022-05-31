package client

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

// Client is is a client that understands GWP over MQTT
type Client struct {
	cm       *autopaho.ConnectionManager
	config   Config
	gateways map[string]bool
}

// Config for the Client
type Config struct {
	ClientID  string
	BrokerURL *url.URL
	Username  string
	Password  string
}

const (
	announcementTopicPrefix = "hb/announce/gw/#"
	commandTopicFmt         = "hb/gw/%s/cmd/#"  // expects the GWid
	sampleTopicFmt          = "hb/gw/%s/data/#" // expects the GWid
	keepAliveSeconds        = 15
	connectRetryDelay       = 5 * time.Second
	connectTimeout          = 5 * time.Second
	commandTopicQoS         = 1
)

// Create client.  Creates client and initializes connection. This function
// does not block.
func Create(c Config) (*Client, error) {
	client := &Client{
		config:   c,
		gateways: map[string]bool{},
	}

	clientConfig := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{c.BrokerURL},
		ConnectTimeout:    connectTimeout,
		KeepAlive:         keepAliveSeconds,
		ConnectRetryDelay: connectRetryDelay,
		ClientConfig: paho.ClientConfig{
			ClientID: "my-client-id",
			Router:   paho.NewSingleHandlerRouter(client.OnPublish),
		},
		OnConnectionUp: client.OnConnectionUp,
	}

	// create connection to MQTT broker. This is a managed connection so it will
	// take care of reconnecting if the connection dies.
	var err error
	client.cm, err = autopaho.NewConnection(context.Background(), clientConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating new MQTT connection to [%s]: %v", c.BrokerURL.String(), err)
	}

	return client, nil
}

// Shutdown client
func (c *Client) Shutdown() {
	c.cm.Disconnect(context.Background())
	<-c.cm.Done() // block until client is really done
}
