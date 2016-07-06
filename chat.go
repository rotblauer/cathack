package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"

	"./lib"
)

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
		t := time.Now()
		timeFormat := "2006-01-02 15:04:05"

		// IP
		ip, err := lib.GetClientIPHelper(s.Request)
		if err != nil {
			log.Fatalln("Error getting client IP: ", err)
		}

		msgWithTime := []byte(t.Format(timeFormat) + "\n" + lib.BootsEncoded(ip) + string(msg))

		// Broadcast web socket. 
		// @msgWithTime []byte
		m.Broadcast(msgWithTime)

		// Open database.
		f, err := os.OpenFile("./chat.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatalln("Error opening file: ", err)
		}

		// Write to database.
		msgWithTimeString := string(msgWithTime)
		bytes, err := f.WriteString(msgWithTimeString + "\n")
		if err != nil {
			log.Fatalln("Error writing string: ", err) // Will this out to same place as fmt? ie &>chat.log
		}

		fmt.Printf("Wrote %d bytes to file\n", bytes)
		fmt.Println(msgWithTimeString)

		f.Close()
	})

	r.Run(":5000")
}

