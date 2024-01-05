package main

import (
	"fmt"
	"os"
	"time"

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

func main() {
	// Reading the YAML configuration file
	configFile, err := os.ReadFile("config.yaml")
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
	opts.SetClientID("go_mqtt_test_client")
	opts.SetUsername(config.MQTT.Username)
	opts.SetPassword(config.MQTT.Password)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connected to MQTT Broker")

	// Publish a message to trigger the activate function
	payload := "activate" // The payload that your main application expects
	token := client.Publish(config.MQTT.Topic, 0, false, payload)
	token.Wait()
	fmt.Printf("Published message '%s' to topic '%s'\n", payload, config.MQTT.Topic)

	// Giving some time for the message to be delivered before exiting
	time.Sleep(2 * time.Second)

	client.Disconnect(250)
	fmt.Println("Disconnected from MQTT Broker")
}
