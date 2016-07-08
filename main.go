package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	// "github.com/njern/gonexmo"
	"./chatty"
	"./lib"
	// "encoding/json"
	"github.com/olahol/melody"
)

// -------------------------------------------------
// // sms("12182606849", "DDF", "c330fe3b", "d69e9ca6c8245f6a")

// //SMS text sender, nexmo to test...need a sign up with keys
// func sms(number string, messageToSend string, key string, secret string) {
// 	nexmoClient, _ := nexmo.NewClientFromAPI(key, secret)
// 	// https://github.com/njern/gonexmo
// 	// Send an SMS
// 	// See https://docs.nexmo.com/index.php/sms-api/send-message for details.
// 	message := &nexmo.SMSMessage{
// 		From:            "12529178592",
// 		To:              number,
// 		Type:            nexmo.Text,
// 		Text:            messageToSend,
// 		ClientReference: "gonexmo-test " + strconv.FormatInt(time.Now().Unix(), 10),
// 		Class:           nexmo.Standard,
// 	}

// 	messageResponse, err := nexmoClient.SMS.Send(message)
// 	if err != nil {
// 		log.Fatalln("Error getting sending sms: ", err)
// 	}
// 	fmt.Println(messageResponse)
// }
// -------------------------------------------------

// -------------------------------------------------
// type MessageResponder struct {
// 	Action string // message, sms, etc..
// 	Status string // 200 OK,
// 	Body   string // text, error, confirmation, notfound
// }
// Action: MESSAGE
// Status: 200 OK
// Body: "Hey there."
//
// Action: SMS
// Status: 200 OK
// Body: "Whoosh!"
//
// Action: SMS
// Status: NUMBER NOT FOUND
// Body: "@trump isn't in the phone book."
//
// Action: SMS
// Status: 500 BAD
// Body: "Failed to send SMS."
// -------------------------------------------------

func getChat(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "index.html")
	log.Printf("Getting chat.")
	fmt.Println()
}

func getHack(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "hack.html")
}

func main() {
	gin.SetMode(gin.DebugMode) // ReleaseMode
	r := gin.Default()
	m := melody.New()

	r.StaticFile("/chat.txt", "./chat.txt")
	r.GET("/", getChat)
	r.GET("/ws", func(c *gin.Context) {
		log.Printf("getChatWS")
		fmt.Println()
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {

		fmt.Printf("Got WS message: %v", string(msg))
		fmt.Println()

		ps1, err := chatty.HandleChatMessage(s, msg)
		m.Broadcast(ps1)

		// var messageRes MessageResponder
		// messageRes.Action = "Message"

		// if err != nil {
		// 	messageRes.Status = "ERR"
		// 	messageRes.Body = "There was an error handling your message."
		// 	log.Fatalln(err)
		// } else {
		// 	messageRes.Status = "200 OK"
		// 	messageRes.Body = string(ps1)
		// }
		// mr, err := json.Marshal(messageRes)
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		// fmt.Printf("BROADCASTING: %v", string(mr))
		// m.Broadcast(mr)

		// Now check for @SMS.
		sms, err := lib.DelegateSendSMS(msg)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("SMS: %v", string(sms))

		// Error sending SMS.
		// if err != nil {
		// 	log.Fatalln(err)
		// 	smsRes := MessageResponder{
		// 		Action: "SMS",
		// 		Status: "ERR",
		// 		Body: err.Error(),
		// 	}
		// 	smsr, _ := json.Marshal(smsRes)
		// 	m.Broadcast(smsr)
		// } else {
		// 	switch {
		// 	case string(sms) == "SENDOK":
		// 		smsRes := MessageResponder{
		// 			Action: "SMS",
		// 			Status: "200 OK",
		// 			Body: "Whoosh!",
		// 		}
		// 		smsr, _ = json.Marshal(smsRes)
		// 		m.Broadcast(smsr)
		// 	case string(sms) == "NOFIND":
		// 		smsRes := MessageResponder{
		// 			Action: "SMS",
		// 			Status: "200 OK",
		// 			Body: "Phone number not found.",
		// 		}
		// 		smsr, _ = json.Marshal(smsRes)
		// 		m.Broadcast(smsr)
		// 		...
		// 	}
		// }
	})

	r.Run(":5000")
}
