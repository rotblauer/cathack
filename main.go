package main

import (
	"./chatty"
	"./lib"
	"./web"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"log"
)

// const (
// 	snippetsBucketName := []byte("snippets")
// )

func main() {
	gin.SetMode(gin.ReleaseMode) // DebugMode
	r := gin.Default()
	m := melody.New()
	h := melody.New()

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("hack.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("snippets"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	// r.StaticFile("/chat.txt", "./chat.txt")
	r.GET("/", web.GetChat)
	r.GET("/r/chat", web.GetChatData)
	r.GET("/ws", func(c *gin.Context) {
		log.Printf("getChatWS")
		fmt.Println()
		m.HandleRequest(c.Writer, c.Request)
	})

	r.GET("/hack", web.GetHack)
	r.GET("/hack/ws", func(c *gin.Context) {
		fmt.Println("Got hack/ws request.")
		h.HandleRequest(c.Writer, c.Request)
	})

	h.HandleConnect(func(s *melody.Session) {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("snippets"))
			v := b.Get([]byte("testSnip"))
			s.Write(v)
			return nil
		})

	})

	h.HandleMessage(func(s *melody.Session, hackery []byte) {
		fmt.Printf("HackHandleMessage: %v", string(hackery))
		fmt.Println()
		h.BroadcastOthers(hackery, s)
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("snippets"))
			err := b.Put([]byte("testSnip"), hackery)
			if err != nil {
				return fmt.Errorf("putting to bucket: %s", err)
			}
			return err
		})
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
