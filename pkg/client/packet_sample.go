package client

import (
	"log"
	"time"

	gwp "go.buf.build/protocolbuffers/go/autro/gwp/v1"
)

func (c *Client) HandleSample(packet *gwp.Packet, sample *gwp.Sample) {
	diff := time.Since(time.UnixMilli(int64(sample.Timestamp)))
	log.Printf("from=[%s] timestamp=[%d] timediff=[%s] sample=[%v]", packet.From, sample.Timestamp, diff, sample)
}
