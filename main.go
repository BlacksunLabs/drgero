package main

// Copyright 2018 Blacksun Research Labs

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
