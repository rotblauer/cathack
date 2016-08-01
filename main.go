package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"./config"
	"./controllers"
	"./models"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	m := melody.New()
	h := melody.New()
	m.Config.MaxMessageSize = 1024 * 1000
	h.Config.MaxMessageSize = 1024 * 1000

	// Open Bolt DB.
	defer models.GetDB().Close()

	// Start shit.
	snippet := new(controllers.SnippetController)
	bucket := new(controllers.BucketController)
	fs := new(controllers.FSController)

	snippetModel := new(models.SnippetModel)

	// Routes.
	//

	r.StaticFS("/assets", http.Dir("assets")) // Static assets.

	r.GET("/", controllers.GetChat)     // index.html
	r.GET("/hack", controllers.GetHack) // hack.html

	r.GET("/ws", func(c *gin.Context) { // Chat WS.
		log.Printf("getChatWS")
		fmt.Println()
		m.HandleRequest(c.Writer, c.Request)
	})
	r.GET("/hack/ws", func(c *gin.Context) { // Hack WS.
		fmt.Println("Got hack/ws request.")
		h.HandleRequest(c.Writer, c.Request)
	})

	r.GET("/r/chat", controllers.GetChatData) // Get chat.txt database.
	r.GET("/hack/b", bucket.All)              // Get all buckets

	r.POST("/hack/b/:bucketName", bucket.Create)  // Create bucket.
	r.DELETE("/hack/b/:bucketId", bucket.Destroy) // Delete bucket.
	r.PUT("/hack/b/:bucketId", bucket.Set)        // Update bucket (rename).

	r.GET("/hack/b/:bucketId", snippet.All)        // Get all snippets for a given bucket
	r.GET("/hack/s", snippet.UberAll)              // Get all snippets ever.
	r.DELETE("/hack/s/:snippetId", snippet.Delete) // Delete snippet @ /hack/s/:snippetId?bucket=snippets

	r.GET("/hack/fs", fs.Walk) // Get all available filepaths.
	// r.PUT("/hack/fs/b/:bucketId", fs.WriteBucketToDir) // Write bucket to directory.

	// Should write contents of file at filePath to associated bucket-snippets (by file's /basepath).
	r.GET("/hack/fs/s", fs.SnippetizeOne)          // Accepts json-ified string as param. Writes one file to bolt snippet.
	r.GET("/hack/fs/b", fs.SnippetizeMany)         // Write each file in path to bolt by n-1 dir name as bucket.
	r.POST("/hack/fs/s/:snippetId", fs.WriteFile)  // Write snippet to file by name. @ /hack/s/:snippetId?bucket=snippets
	r.POST("/hack/fs/b/:bucketId", fs.WriteBucket) // Write bucket's snippets to files by name. No query params required.

	// Websockets.
	//
	// Hack.
	h.HandleConnect(func(s *melody.Session) {
		// snippets, err := snippetModel.All(config.DefaultBucketName)
		// if err != nil {
		// 	s.Write([]byte(err.Error()))
		// }
		// j, _ := json.Marshal(snippets)
		// s.Write(j)
		s.Write([]byte("Connected."))
	})

	h.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Printf("Handling hack message: %v\n", string(msg))

		snippetRequest := models.SnippetChangedRequest{}
		json.Unmarshal(msg, &snippetRequest)

		snippetRequest.Snippet.TimeStamp = int(time.Now().UTC().UnixNano() / 1000000)
		err := snippetModel.Set(snippetRequest.Snippet)
		if err != nil {
			h.Broadcast([]byte("error: " + err.Error()))
		}

		// The change will be peeled off by Angular controller in incoming WS.
		h.BroadcastOthers(msg, s)

	})
	h.HandleError(func(s *melody.Session, err error) {
		fmt.Printf("Melody error: %v", err)
	})

	// Chat.
	m.HandleMessage(func(s *melody.Session, msg []byte) {

		// is typing
		if string(msg) == "***" {
			m.BroadcastOthers([]byte("***"), s)

			// is not typing
		} else if string(msg) == "!***" {
			m.BroadcastOthers([]byte("!***"), s)

			// sent message
		} else {
			ps1, err := models.SaveChatMsg(s, msg)
			if err != nil {
				log.Fatalln(err)
			}
			m.Broadcast(ps1)
		}

		// // Now check for @SMS.
		// _, err := catchat.DelegateSendSMS(msg)
		// if err != nil {
		// 	log.Fatalln(err)
		// }

	})

	r.Run(config.MakeThisMyPort)
}
