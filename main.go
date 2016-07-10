package main

import (
	"./chatty"
	"./lib"
	"./web"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	io "io/ioutil"
	"log"
	"os"
	fp "path/filepath"
)

type Snippet struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Language  string `json:"language"`
	Content   string `json:"content"`
	TimeStamp int    `json:"timestamp"`
	Meta      string `json:"meta"`
}

func main() {
	gin.SetMode(gin.ReleaseMode) // DebugMode
	r := gin.Default()
	m := melody.New()
	h := melody.New()
	// Overclock: set this to 100KB = 1MB
	// https://sourcegraph.com/github.com/olahol/melody/-/info/GoPackage/github.com/olahol/melody/-/New
	m.Config.MaxMessageSize = 1024 * 1000
	h.Config.MaxMessageSize = 1024 * 1000 // (default was 512). suckas.

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

	// Save bucket as files on server.
	r.GET("/hack/repofy/:bucketId", func(c *gin.Context) {

		bucketid := c.Param("bucketId") //string
		hacksBucketRootPath := "./hacks/" + bucketid + "/"

		// clean it out (in case file names have changed)
		// FIXME: danger.
		rerr := os.RemoveAll(hacksBucketRootPath)
		if rerr != nil {
			fmt.Printf("Error cleaning bucket path: %v", rerr)
		}

		// broadcast new index
		err := db.View(func(tx *bolt.Tx) (viewerr error) {

			b := tx.Bucket([]byte(bucketid))
			c := b.Cursor()

			for snipkey, snipval := c.First(); snipkey != nil; snipkey, snipval = c.Next() {
				var snip Snippet
				json.Unmarshal(snipval, &snip)

				// snip.Name -> /some/where/./in//here/boots.go
				cleanFullName := fp.Clean(snip.Name)

				filepath := fp.Dir(cleanFullName)
				if filepath == "." {
					filepath = ""
				}
				// make filepath
				// returns nil if exists
				viewerr = os.MkdirAll(hacksBucketRootPath+filepath, 0777)                                                   //rw
				viewerr = io.WriteFile(hacksBucketRootPath+filepath+"/"+fp.Base(cleanFullName), []byte(snip.Content), 0666) //rw, truncates before write
			}

			return viewerr // hopefully nil
		})

		if err == nil {
			c.JSON(200, "Saved bucket.")
		} else {
			c.JSON(500, "Internal server error."+err.Error())
		}
	})

	r.DELETE("/hack/s/:snippetId", func(c *gin.Context) {

		id := c.Param("snippetId") // func (c *Context) Param(key string) string

		// remove given snippet by id
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("snippets"))

			// remove from os
			// get snippet so can access Name
			var snip Snippet
			v := b.Get([]byte(id))
			json.Unmarshal(v, &snip)

			path := "./hacks/snippets/" + snip.Name
			fmt.Printf("Removing file at path: %v", path)
			err := os.Remove(path)
			if err != nil {
				fmt.Printf("Error removing file: %v", err)
				return err
			} else {
				// remove from bucket
				err := b.Delete([]byte(id))

				if err != nil {
					fmt.Printf("Error deleting from bucket: %v", err)
					return err
				}
			}
			return nil
		})

		if err != nil {
			c.JSON(400, "No snippet found with id: "+id)
		} else {

			// broadcast new index
			err := db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("snippets"))
				// v := b.Get([]byte("testSnip"))

				// iterate through snippets
				c := b.Cursor()
				var snippets []Snippet // Array of Go snippet structs

				for snipkey, snipval := c.First(); snipkey != nil; snipkey, snipval = c.Next() {
					var snip Snippet
					json.Unmarshal(snipval, &snip)
					snippets = append(snippets, snip)
				}

				o, err := json.Marshal(snippets) // JSON-ified array of snippets
				if err != nil {
					fmt.Printf("Error marshaling snippet array: %v", err)
					return err
				}

				h.Broadcast(o)
				return nil
			})

			if err == nil {
				c.JSON(200, "Deleted snippet: "+id)
			} else {
				c.JSON(500, "Internal server error.")
			}
		}

	})

	// Send all of the snippets when a user connects.
	h.HandleConnect(func(s *melody.Session) {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("snippets"))
			// v := b.Get([]byte("testSnip"))

			// iterate through snippets
			c := b.Cursor()
			var snippets []Snippet // Array of Go snippet structs

			for snipkey, snipval := c.First(); snipkey != nil; snipkey, snipval = c.Next() {
				var snip Snippet
				json.Unmarshal(snipval, &snip)
				snippets = append(snippets, snip)
			}

			o, err := json.Marshal(snippets) // JSON-ified array of snippets
			if err != nil {
				fmt.Printf("Error marshaling snippet array: %v", err)
			}

			s.Write(o)
			return nil
		})

	})

	h.HandleMessage(func(s *melody.Session, hackery []byte) {
		fmt.Printf("HackHandleMessage: %v", string(hackery))
		fmt.Println()
		h.BroadcastOthers(hackery, s)

		var snip Snippet

		json.Unmarshal(hackery, &snip)
		id := snip.Id // find snippet name and update db by id

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("snippets"))
			err := b.Put([]byte(id), hackery)

			if err != nil {
				return fmt.Errorf("putting to bucket: %s", err)
			}
			return err
		})
	})

	// Handle error for hackery.
	h.HandleError(func(s *melody.Session, err error) {
		fmt.Printf("Melody error: %v", err)
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
