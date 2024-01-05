package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v2"
)

type Config struct {
	MQTT struct {
		Server   string `yaml:"server"`
		Topic    string `yaml:"topic"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"mqtt"`
}

var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	if string(msg.Payload()) == "activate" {
		cmd := exec.Command("xset", "dpms", "force", "on")
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error executing command:", err)
		} else {
			fmt.Println("Display activated")
		}
	}
}

var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	fmt.Println("Connected")
}

var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	// Reading the YAML configuration file
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	// Setting up MQTT client options
	opts := MQTT.NewClientOptions()
	opts.AddBroker(config.MQTT.Server)
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername(config.MQTT.Username)
	opts.SetPassword(config.MQTT.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := client.Subscribe(config.MQTT.Topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", config.MQTT.Topic)

	select {}
}
