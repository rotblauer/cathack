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
		msgWithTime = time.Now() + msg
		
		m.Broadcast(msgWithTime)
		
		f, err := os.OpenFile("./chat.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatalln("Error opening file: ", err)
		}
		
		bytes, err := f.WriteString(time.Now() + ' ' + string(msg) + "\n")
		if err != nil {
			log.Fatalln("Error writing string: ", err)
		}
		
		fmt.Printf("Wrote %d bytes to file\n", bytes)
		fmt.Println(string(msg))
		
		f.Close()
	})

	r.Run(":5000")
}
