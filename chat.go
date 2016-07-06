package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	// "strconv"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"

	"./lib"
)

// func formatPS1()

func main() {
	//go run chat.go
	r := gin.Default()
	m := melody.New()

	// Serves file,
	r.StaticFile("/chat.txt", "./chat.txt")

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {

		// Message with timestamp.
		// timeString := strconv.FormatInt(time.Now().UTC().Unix(), 16)
		// timeFormat := "2006-01-02 15:04:05"
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

		bs := ""
		bs += timeString 
		bs += ","
		bs += geoip[0] //lat,lon
		bs += "," + geoip[1] //tz
		bs += "," + geoip[2] //subdiv
		bs += ","
		bs += lib.BootsEncoded(ip)
		bs += string(msg)

		ps1 := []byte(bs)

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

