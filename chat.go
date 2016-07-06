package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

// f, err := os.OpenFile(filename, os.O_APPEND, 0666)

// n, err := f.WriteString(text)

// f.Close()

func main() {
	r := gin.Default()
	m := melody.New()

	r.StaticFile("/test.txt", "./test.txt")

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
		f, err := os.OpenFile("./test.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatalln("Error opening file: ", err)
		}
		bytes, err := f.WriteString(string(msg) + "\n")
		if err != nil {
			log.Fatalln("Error writing string: ", err)
		}
		fmt.Printf("Wrote %d bytes to file\n", bytes)
		fmt.Println(string(msg))
		f.Close()
	})

	r.Run(":5000")
}
