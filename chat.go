package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/njern/gonexmo"
	"github.com/olahol/melody"
	j "github.com/ricardolonga/jsongo"

	"./lib"
)

//SMS text sender, nexmo to test...need a sign up with keys
func sms(number string, messageToSend string, key string, secret string) {
	nexmoClient, _ := nexmo.NewClientFromAPI(key, secret)
	// https://github.com/njern/gonexmo
	// Send an SMS
	// See https://docs.nexmo.com/index.php/sms-api/send-message for details.
	message := &nexmo.SMSMessage{
		From:            "12529178592",
		To:              number,
		Type:            nexmo.Text,
		Text:            messageToSend,
		ClientReference: "gonexmo-test " + strconv.FormatInt(time.Now().Unix(), 10),
		Class:           nexmo.Standard,
	}

	messageResponse, err := nexmoClient.SMS.Send(message)
	if err != nil {
		log.Fatalln("Error getting sending sms: ", err)
	}
	fmt.Println(messageResponse)
}

func main() {
	//go run chat.go
	r := gin.Default()
	m := melody.New()

	// Serves file,
	r.StaticFile("/chat.txt", "./chat.txt")

	// These are _really_ slow. WTF.
	// r.Static("/assets", "./assets")
	// r.StaticFile("app.js", "./app.js")
	// r.GET("/app.js", func(c *gin.Context) {
	// 	http.ServeFile(c.Writer, c.Request, "app.js")
	// })

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {

		// Message with timestamp.
		timeUnixNano := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		timeString := time.Now().UTC().String()

		// IP
		ip, err := lib.GetClientIPHelper(s.Request)
		if err != nil {
			log.Fatalln("Error getting client IP: ", err)
		}

		geoip, err := lib.GetGeoFromIP(ip)
		if err != nil {
			log.Fatalln("Error getting Geo IP.", err)
		}

		data := j.Object().Put("time", timeString). // btw hanging dots are no go
								Put("unixNano", timeUnixNano).
								Put("message", string(msg)).
								Put("ip", ip).
								Put("bootsIP", lib.BootsEncoded(ip)).
								Put("lat", geoip["lat"]).
								Put("lon", geoip["lon"]).
								Put("city", geoip["city"]).
								Put("subdiv", geoip["subdiv"]).
								Put("countryIsoCode", geoip["countryIsoCode"]).
								Put("tz", geoip["tz"])

		dataIndentedString := data.String()
		//sms("12182606849", dataIndentedString, "c330fe3b", "d69e9ca6c8245f6a")

		ps1 := []byte(dataIndentedString)
		// Broadcast web socket.
		// @ps1 []byte
		m.Broadcast(ps1)

		// Open database.
		f, err := os.OpenFile("./chat.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatalln("Error opening file: ", err)
		}

		// Write to database.
		ps1String := string(ps1)
		bytes, err := f.WriteString(ps1String + "\n")
		if err != nil {
			log.Fatalln("Error writing string: ", err) // Will this out to same place as fmt? ie &>chat.log
		}

		fmt.Printf("Wrote %d bytes to file\n", bytes)
		fmt.Println(ps1String)

		f.Close()
	})

	r.Run(":5000")
}
