package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"./chatty"
	"./controllers"
	"./lib"
	"./models"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

const (
	hacksRootPath         string = "./hacks/"
	hacksDBPath           string = "hack.db"
	placeHolderBucketName string = "snippets"
)

func main() {
	gin.SetMode(gin.ReleaseMode) // DebugMode
	r := gin.Default()
	m := melody.New()
	h := melody.New()
	// Overclock: set this to 1000KB = 1MB
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
		_, err := tx.CreateBucketIfNotExists([]byte(placeHolderBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	r.StaticFS("/assets", http.Dir("assets"))
	r.GET("/", controllers.GetChat)
	r.GET("/r/chat", controllers.GetChatData)
	r.GET("/ws", func(c *gin.Context) {
		log.Printf("getChatWS")
		fmt.Println()
		m.HandleRequest(c.Writer, c.Request)
	})

	r.GET("/hack", controllers.GetHack)
	r.GET("/hack/ws", func(c *gin.Context) {
		fmt.Println("Got hack/ws request.")
		h.HandleRequest(c.Writer, c.Request)
	})

	// Save bucket as files on server.
	r.GET("/hack/repofy/:bucketName", func(c *gin.Context) {

		var err error

		bucketName := c.Param("bucketName") //string

		// clean it out (in case file names have changed)
		// FIXME: danger.
		err = os.RemoveAll(hacksRootPath + bucketName)
		if err != nil {
			fmt.Printf("Error cleaning bucket path: %v", err)
		}

		err = db.View(func(tx *bolt.Tx) error {
			return models.WriteBucketToFileSys(hacksRootPath, bucketName, tx)
		})

		if err == nil {
			c.JSON(200, "Saved bucket.")
		} else {
			c.JSON(500, "Internal server error."+err.Error())
		}
	})

	//////////

	// Get Snippet.
	r.GET("/hack/s/:snippetId", func(c *gin.Context) {

	})

	// Get all snippets for single bucket.
	r.GET("/hack/b/:bucketId", func(c *gin.Context) {
		bid := c.Param("bucketId")
		var snippets models.Snippets

		err := db.View(func(tx *bolt.Tx) error {
			snippets, _ = models.IndexSnippets(bid, tx)
			return nil
		})
		if err != nil {
			c.JSON(500, err)
		} else {
			c.JSON(200, snippets)
		}
	})

	// Get all buckets.
	// Returns list of bucket names.
	r.GET("/hack/b", func(c *gin.Context) {

		// Buckets slice struct.
		// Will return once full.
		var buckets models.SnippetBuckets

		err := db.View(func(tx *bolt.Tx) error {
			tx.ForEach(func(name []byte, b *bolt.Bucket) error {
				buckets = append(buckets, models.SnippetBucket{Name: string(name)})
				return nil
			})
			return nil
		})
		if err != nil {
			c.JSON(500, err)
		} else {
			c.JSON(200, buckets)
		}
	})

	// DELETE SNIPPET.
	// /hack/s/:snippetId?bucket=snippets
	r.DELETE("/hack/s/:snippetId", func(c *gin.Context) {
		snippetId := c.Param("snippetId") // func (c *Context) Param(key string) string
		// bucketId := c.DefaultQuery("bucket", "snippets")

		// remove given snippet by snippetId
		err := db.Update(func(tx *bolt.Tx) error {
			return models.DeleteSnippet(snippetId, placeHolderBucketName, tx)
		})

		if err != nil {
			c.JSON(400, "No snippet found with snippetId: "+snippetId)
		} else {

			// broadcast new index
			err := db.View(func(tx *bolt.Tx) error {
				snippets, _ := models.IndexSnippets(placeHolderBucketName, tx)
				o, _ := json.Marshal(snippets) // JSON-ified array of snippets
				h.Broadcast(o)
				return nil
			})

			if err == nil {
				c.JSON(200, "Deleted snippet: "+snippetId)
			} else {
				c.JSON(500, "Internal server error.")
			}
		}
	})

	// Send all of the snippets when a user connects.
	// h.HandleConnect(func(s *melody.Session) {
	// 	s.Write([]byte("Connected."))
	// })
	h.HandleConnect(func(s *melody.Session) {
		db.View(func(tx *bolt.Tx) error {
			snips, err := models.IndexSnippets(placeHolderBucketName, tx)
			o, err := json.Marshal(snips) // JSON-ified array of snippets
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

		db.Update(func(tx *bolt.Tx) error {
			return models.SetSnippet(models.SnipFromJSON(hackery).Id, hackery, placeHolderBucketName, tx)
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
