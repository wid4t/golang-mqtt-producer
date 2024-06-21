package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
)

var (
	mqttBroker   = "mqtts://mosquitto-service.server.svc.cluster.local:8883"
	mqttClientID = "go-mqtt-producer"
)

func main() {

	app := fiber.New()

	mqttOpts := MQTT.NewClientOptions()
	mqttOpts.AddBroker(mqttBroker)
	mqttOpts.SetClientID(mqttClientID)
	mqttOpts.SetUsername(os.Getenv("username"))
	mqttOpts.SetPassword(os.Getenv("password"))

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	mqttOpts.SetTLSConfig(tlsConfig)

	mqttClient := MQTT.NewClient(mqttOpts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	defer mqttClient.Disconnect(250)

	app.Get("/api/mqtt/send", func(c *fiber.Ctx) error {

		message := c.Query("message")
		token := mqttClient.Publish("topic", 1, false, message)
		token.Wait()

		if token.Error() != nil {
			return c.SendString("Failed to publish MQTT message.")
		}

		return c.SendString("Message sent: " + message)
	})

	log.Fatal(app.Listen(":8082"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down...")

	time.Sleep(2 * time.Second)
}
