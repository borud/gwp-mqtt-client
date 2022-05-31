package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/borud/gwp-mqtt-client/pkg/client"
)

var (
	mqttBrokerURL *url.URL
	mqttBroker    string
	clientID      string
	username      string
	password      string
)

const (
	keepAliveSeconds  = 15
	connectRetryDelay = 5 * time.Second
	connectTimeout    = 5 * time.Second
	commandTopicQoS   = 1
)

func init() {
	clientID = fmt.Sprintf("gwp-mqtt-client-%s", big.NewInt(time.Now().UnixNano()).Text(36))
	flag.StringVar(&mqttBroker, "broker", "mqtt://127.0.0.1:1883", "mqtt broker address (mqtt://<host>:<port>)")
	flag.StringVar(&clientID, "client-id", clientID, "mqtt broker address (mqtt://<host>:<port>)")
	flag.StringVar(&username, "user", "", "MQTT user name")
	flag.StringVar(&password, "password", "", "MQTT password")
	flag.Parse()

	var err error
	mqttBrokerURL, err = url.Parse(mqttBroker)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Printf("broker=[%s] clientID=[%s]", mqttBrokerURL, clientID)

	c, err := client.Create(client.Config{
		ClientID:  clientID,
		BrokerURL: mqttBrokerURL,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	c.Shutdown()
}
