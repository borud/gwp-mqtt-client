package client

import (
	"context"
	"log"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	gwp "go.buf.build/protocolbuffers/go/autro/gwp/v1"
	"google.golang.org/protobuf/proto"
)

// OnPublish is called whenever we receive a message from the broker.
func (c *Client) OnPublish(p *paho.Publish) {
	// unmarshal he payload as a GWP packet.
	var deserialized gwp.Packet
	err := proto.Unmarshal(p.Payload, &deserialized)
	if err != nil {
		log.Printf("unable to unmarshal packet: %v", err)
		return
	}

	packet := &deserialized

	switch payload := packet.Payload.(type) {
	case *gwp.Packet_Announcement:
		c.HandleAnnounce(packet, payload.Announcement)
	case *gwp.Packet_Sample:
		c.HandleSample(packet, payload.Sample)

	default:
		log.Printf("unhandled message type:%v", packet)
	}
}

// OnConnectionUp handles when we first connect to the MQTT broker. It will populate
// the subscriptions with the default subscription(s).
func (c *Client) OnConnectionUp(cm *autopaho.ConnectionManager, _ *paho.Connack) {
	_, err := cm.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			announcementTopicPrefix: {QoS: 0},
		},
	})
	if err != nil {
		log.Printf("failed to subscribe to announcement topic")
	}
}
