package main

import (
	"./chatty"
	"./lib"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type ChatMessage struct {
	Time           string
	UnixNano       string
	Message        string
	Ip             string
	BootsIP        string
	Lat            string
	Lon            string
	City           string
	Subdiv         string
	CountryIsoCode string
	Tz             string
}

func getChat(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "index.html")
	log.Printf("Getting chat.")
	fmt.Println()
}

func getChatData(c *gin.Context) {
	// func ReadFile(filename string) ([]byte, error)
	fileContents, err := ioutil.ReadFile("./chat.txt")
	if err != nil {
		fmt.Printf("Error ioutiling chat.txt: %v", err)
	}

	messageStrings := strings.Split(string(fileContents), "\n")

	var collection []ChatMessage

	for _, messageString := range messageStrings {
		bytes := []byte(messageString)
		var cm ChatMessage
		json.Unmarshal(bytes, &cm)
		collection = append(collection, cm)
	}

	collectionBytes, err := json.Marshal(collection) // []byte

	c.JSON(200, gin.H{
		"status": "200 OK",
		"data":   string(collectionBytes), // again, the hanging commas are strangely necessary
	})
}

func getHack(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "hack.html")
}

func main() {
	gin.SetMode(gin.ReleaseMode) // DebugMode
	r := gin.Default()
	m := melody.New()
	h := melody.New()

	// r.StaticFile("/chat.txt", "./chat.txt")
	r.GET("/", getChat)
	r.GET("/r/chat", getChatData)
	r.GET("/ws", func(c *gin.Context) {
		log.Printf("getChatWS")
		fmt.Println()
		m.HandleRequest(c.Writer, c.Request)
	})

	r.GET("/hack", getHack)
	r.GET("/hack/ws", func(c *gin.Context) {
		fmt.Println("Got hack/ws request.")
		h.HandleRequest(c.Writer, c.Request)
	})

	h.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Printf("HackHandleMessage: %v", string(msg))
		fmt.Println()
		h.BroadcastOthers(msg, s)
	})

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
			ps1, err := chatty.HandleChatMessage(s, msg)
			if err != nil {
				log.Fatalln(err)
			}
			m.Broadcast(ps1)
		}

		// Now check for @SMS.
		sms, err := lib.DelegateSendSMS(msg)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("SMS: %v", string(sms))
		fmt.Println()

	})

	r.Run(":5000")
}
