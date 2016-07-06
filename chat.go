package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

func main() {
	//go run chat.go
	r := gin.Default()
	m := melody.New()
	t := time.Now()
	timeFormat := "2006-01-02 15:04:05"

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
		msgWithTime := []byte(t.Format(timeFormat) + string(msg))

		// Broadcast web socket.
		m.Broadcast(msgWithTime)

		// Open database.
		f, err := os.OpenFile("./chat.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatalln("Error opening file: ", err)
		}

		// Write to database.
		bytes, err := f.WriteString(string(msgWithTime) + "\n")
		if err != nil {
			log.Fatalln("Error writing string: ", err)
		}

		fmt.Printf("Wrote %d bytes to file\n", bytes)
		fmt.Println(string(msgWithTime))

		f.Close()
	})

	r.Run(":5000")
}
