package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/barmband"
	"gitlab.hs-flensburg.de/flar3845/barmband/bandcommand/messaging"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

const MqttBroker = "localhost"
const MqttPort = "1883"
const ChallengeTopic = "barmband/challenge"

const SetupTopic = "barmband/setup"

func makeConnectionString(host string, port string) string {
	return fmt.Sprintf("tcp://%s:%s", host, port)
}

func connectMqtt(host string, port string) (mqtt.Client, error) {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	connectionString := makeConnectionString(host, port)

	opts := mqtt.NewClientOptions().AddBroker(connectionString).SetClientID("bandcommander")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(2 * time.Second)

	c := mqtt.NewClient(opts)
	token := c.Connect()
	token.Wait()

	return c, token.Error()
}

func main() {

	client, err := connectMqtt(MqttBroker, MqttPort)
	if err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %s", err)
	}

	bc := bandcommand.New(func(pair barmband.Pair) {
		firstS := fmt.Sprintf("%X", pair.First)
		secondS := fmt.Sprintf("%X", pair.Second)

		client.Publish(ChallengeTopic, 0, false, fmt.Sprintf("New pair %s %s %s", firstS, secondS, pair.Color))
	})

	messageHandler := mqttMessageHandler(bc)

	token := client.Subscribe(ChallengeTopic, 0, messageHandler)
	token.Wait()

	if token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %s", token.Error())
	}

	token = client.Subscribe(SetupTopic, 0, messageHandler)
	token.Wait()

	if token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %s", token.Error())
	}

	select {}
}

// mqttMessageHandler parses messages and send it to the BandCommand
func mqttMessageHandler(bc bandcommand.BandCommand) func(client mqtt.Client, message mqtt.Message) {

	return func(client mqtt.Client, message mqtt.Message) {
		fmt.Println(message)
		messageString := string(message.Payload())
		msg, err := messaging.ParseMessage(messageString)

		if err != nil {
			log.Printf("Failed to parse message '%s': %s\n", messageString, err)
		} else {
			fmt.Printf("Got message: %v\n", msg)
			bc.HandleMessage(msg)
		}
	}
}
