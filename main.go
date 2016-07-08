package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	// "github.com/njern/gonexmo"
	"./chatty"
	// "./lib"
	"github.com/olahol/melody"
)

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

		// status, err := lib.DelegateSendSMS(msg)
		ps1, err := chatty.HandleChatMessage(s, msg)

		if err != nil {
			// m.Broadcast([]byte("{\"status\": \"" + err.Error() + "\"}"))
			log.Fatalln(err)
		}

		log.Printf("Broadcasting PS1 as: %v", string(ps1))
		fmt.Println()

		m.Broadcast(ps1)
		// m.Broadcast([]byte("I am broadcasting!!!"))

		// FIXME: None of the log.Printf logs are making it into chat.log.

	})

	r.Run(":5000")
}
