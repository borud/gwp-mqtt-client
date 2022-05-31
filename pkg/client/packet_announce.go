package client

import (
	"context"
	"fmt"
	"log"

	"github.com/borud/gwp-mqtt-client/pkg/addr"
	"github.com/eclipse/paho.golang/paho"
	gwp "go.buf.build/protocolbuffers/go/autro/gwp/v1"
)

// HandleAnnounce handles incoming (GW) announcements.
func (c *Client) HandleAnnounce(packet *gwp.Packet, ann *gwp.Announcement) {
	gwAddr, err := addr.ParseGWAddr(packet.From)
	if err != nil {
		log.Printf("ANN> error parsing GW address: %v", err)
		return
	}

	_, ok := c.gateways[gwAddr.ID]
	if ok {
		log.Printf("ANN> already known gateway [%s]", gwAddr.ID)
		return
	}

	log.Printf("ANN> new gateway [%s]", gwAddr.ID)

	//data topic
	sampleTopic := fmt.Sprintf(sampleTopicFmt, gwAddr.ID)

	_, err = c.cm.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			sampleTopic: {QoS: 0},
		}})
	if err != nil {
		log.Printf("ANN> unable to subscribe to [%s] for GW [%s]", sampleTopic, gwAddr.ID)
		return
	}
	log.Printf("ANN> subscribing to data from [%s]", gwAddr.ID)
	c.gateways[gwAddr.ID] = true
}
