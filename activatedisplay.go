package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
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

var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	if string(msg.Payload()) == "activate" {
		cmd := exec.Command(
			"xdotool", "search", "--onlyvisible", "--class", "chromium", "windowfocus",
			"key", "ctrl+r",
		)
		cmd.Env = append(os.Environ(), "DISPLAY=:0")
		var stdoutBuf, stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error executing command:", err)
			fmt.Println("Stderr: ", stderrBuf.String())
		} else {
			fmt.Println("Browser refreshed")
		}
		cmd = exec.Command("xset", "dpms", "force", "on")
		cmd.Env = append(os.Environ(), "DISPLAY=:0")
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error executing command:", err)
			fmt.Println("Stderr: ", stderrBuf.String())
		} else {
			fmt.Println("Display activated")
		}
	}
}

var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	fmt.Println("Connected")
	// Publishing a message upon connection
	token := client.Publish(config.MQTT.Topic, 0, false, "Display ready to be activated")
	token.Wait()
	fmt.Println("Published message to topic:", config.MQTT.Topic)
	token = client.Subscribe(config.MQTT.Topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", config.MQTT.Topic)
}

var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

var config Config

func main() {
	// Reading the YAML configuration file
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

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
	// Exponential backoff configuration
	maxRetries := 500
	retryCount := 0
	backoff := 1 * time.Second // Initial backoff duration
	maxBackoff := 32 * time.Second
	backoffFactor := 2

	fmt.Printf("Connect to MQTT server: %v with username: %v\n", config.MQTT.Server, config.MQTT.Username)

	for {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			fmt.Printf("Failed to connect: %v\n", token.Error())

			if retryCount >= maxRetries {
				fmt.Println("Max retries reached, stopping.")
				break
			}

			fmt.Printf("Retrying in %v...\n", backoff)
			time.Sleep(backoff)
			retryCount++
			backoff *= time.Duration(backoffFactor)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		} else {
			fmt.Println("Connected to MQTT broker.")
			break // Successfully connected
		}
	}

	select {}
}
