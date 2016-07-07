package main

import (
	// "fmt"
	"log"
	"net/http"
	// "strconv"
	// "time"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	// "github.com/njern/gonexmo"
	"github.com/olahol/melody"

	"./chatty"
	"./lib"
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
}
func getChatWS(c *gin.Context) {
	m := melody.New()

	m.HandleRequest(c.Writer, c.Request)
	m.HandleMessage(func(s *melody.Session, msg []byte) {

		ps1, err := chatty.HandleChatMessage(s, msg)

		if err != nil {
			m.Broadcast([]byte(err.Error()))
			log.Fatalln(err)
		}

		// Broadcast message with metadata on successful handling.
		// @ps1 []byte
		m.Broadcast(ps1)

		messageText := string(msg)

		phoneBook := make(map[string]string)
		phoneBook["john"] = "2182606849"
		phoneBook["isaac"] = "2183494908"

		re, err := regexp.Compile(`@(\w+)`) // FIXME: this should capture only the name, not the @ part. it captures @name. don't know why.
		if err != nil {
			log.Printf("Error compiling regex: %v", err)
		}

		if re.MatchString(messageText) {

			log.Printf("Regex matches. Sending smss.")

			// get @names
			names := re.FindAllString(messageText, 3) // limit to first 3 matches (from left -> right)

			// Remove @name's if we want.
			// messageText = re.ReplaceAllString(messageText, "")

			// send smss
			for _, name := range names {

				phoneNumber := phoneBook[strings.Replace(string(name), "@", "", 1)]

				log.Printf("Name: %v, Number: %v", name, phoneNumber)

				if len(phoneNumber) > 0 {

					err := lib.SendSMS(phoneNumber, messageText)

					if err != nil {
						log.Printf("Textbelt error: %v", err)
						m.Broadcast([]byte("{\"status\": \"" + err.Error() + "\"}"))
					} else {
						log.Printf("Sent SMS to %v with content: %v", phoneNumber, messageText)
						m.Broadcast([]byte("{\"status\": \"" + "Whoosh!" + "\"}"))
					}
				}
			}
		} else {
			log.Printf("Regex doesn't match any @names.")
		}

	})
}

func getHack(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "hack.html")
}

func main() {
	//go run chat.go
	r := gin.Default()

	r.StaticFile("/chat.txt", "./chat.txt")
	r.GET("/", getChat)
	r.GET("/ws", getChatWS)

	r.Run(":5000")
}
