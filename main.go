package main

import (
	"encoding/json"
	"fmt"

	"github.com/BlacksunLabs/drgero/event"
	"github.com/BlacksunLabs/drgero/mq"
	"github.com/gin-gonic/gin"
)

var m = new(mq.Client)

func eventPOST(c *gin.Context) {
	event := new(event.Event)

	body, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Unable to parse JSON request body",
		})
	}

	msg := json.RawMessage(string(body))
	host := c.Request.Host
	ua := c.Request.UserAgent()

	event.Message = string(msg)
	event.Host = host
	event.UserAgent = ua

	eventBody, err := json.Marshal(event)
	if err != nil {
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("failed to marshal event into JSON: %v", err),
		})
	}
	err = m.PublishJSONToFanoutExchange(eventBody, "events")
	if err != nil {
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("failed to publish event to events exchange: %v", err),
		})
	}
}

func main() {
	err := m.Connect("amqp://guest:guest@localhost:5672")
	if err != nil {
		fmt.Printf("%v", err)
	}

	router := gin.Default()

	router.POST("/event", eventPOST)

	router.Run()
}
