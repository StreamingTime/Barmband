package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"time"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

const MqttBroker = "test.mosquitto.org"
const MqttPort = "1883"
const MqttTopic = "barmband/test"

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

	if token := client.Subscribe(MqttTopic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %s", token.Error())
	}

	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := client.Publish(MqttTopic, 0, false, text)
		token.Wait()
	}

	select {}
}
