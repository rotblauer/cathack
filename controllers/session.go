package controllers

import (
	"encoding/json"
	"fmt"
	"log"

	"./config"
	"./models"

	"github.com/olahol/melody"
)

var snippetModel = new(models.SnippetModel)
var bucketModel = new(models.BucketModel)

// Chat.
var m *Melody

// Hack.
var h *Melody

func init() {
	m = melody.New()
	h = melody.New()

	// Set high limits.
	m.Config.MaxMessageSize = 1024 * 1000
	h.Config.MaxMessageSize = 1024 * 1000

	// Chat.
	//
	m.HandleMessage(func(s *melody.Session, msg []byte) {

		fmt.Printf("Got WS message: %v", string(msg))
		fmt.Println()

		// is typing
		if string(msg) == "***" {
			m.BroadcastOthers([]byte("***"), s)

			// is not typing
		} else if string(msg) == "!***" {
			m.BroadcastOthers([]byte("!***"), s)

			// sent message
		} else {
			ps1, err := catchat.HandleChatMessage(s, msg)
			if err != nil {
				log.Fatalln(err)
			}
			m.Broadcast(ps1)
		}

		// Now check for @SMS.
		sms, err := catchat.DelegateSendSMS(msg)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("SMS: %v", string(sms))
		fmt.Println()

	})

	// Hack.
	//

	// Send all of the snippets when a user connects.
	// h.HandleConnect(func(s *melody.Session) {
	// 	s.Write([]byte("Connected."))
	// })
	h.HandleConnect(func(s *melody.Session) {
		snippets, err := snippetModel.All(config.DefaultBucketName)
		j, _ := json.Marshal(snips)
		s.Write(o)
	})

	h.HandleMessage(func(s *melody.Session, msg []byte) {
		h.BroadcastOthers(msg, s)
		var snip snippetModel.Snippet
		json.Unmarshal(msg, &snip)
		snippetModel.Set(snip)
	})

	// Handle error for hackery.
	h.HandleError(func(s *melody.Session, err error) {
		fmt.Printf("Melody error: %v", err)
	})

}
