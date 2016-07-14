package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"./catchat"
	"./controllers"
	"./models"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

//CORSMiddleware ...
// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://chat.areteh.co:5000")
// 	}
// }

//

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

	r.GET("/r/chat", controllers.GetChatData) // Get chat.txt
	r.GET("/hack/b", bucket.All)              // Get all buckets

	r.POST("/hack/b/:bucketName", bucket.Create)  // TODO
	r.DELETE("/hack/b/:bucketId", bucket.Destroy) // TODO
	r.PUT("/hack/b/:bucketId", bucket.Set)        // TODO

	r.GET("/hack/b/:bucketId", snippet.All) // Get all snippets for a given bucket
	r.GET("/hack/s", snippet.UberAll)
	r.DELETE("/hack/s/:snippetId", snippet.Delete) // Delete snippet @ /hack/s/:snippetId?bucket=snippets

	r.GET("/hack/fs", fs.Walk)

	r.PUT("/hack/fs/b/:bucketId", fs.SetBucket)

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

	// IS NOT SAVING TO BOLT.
	h.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Printf("Handling hack message: %v\n", string(msg))
		snip := models.Snippet{}
		json.Unmarshal(msg, &snip)
		snip.TimeStamp = int(time.Now().UTC().UnixNano() / 1000000)
		err := snippetModel.Set(snip)
		if err != nil {
			h.Broadcast([]byte("'ERROR':" + err.Error()))
		}
		h.BroadcastOthers(msg, s)
		// h.Broadcast(j) // send with timestamp update
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
			ps1, err := catchat.HandleChatMessage(s, msg)
			if err != nil {
				log.Fatalln(err)
			}
			m.Broadcast(ps1)
		}

		// Now check for @SMS.
		_, err := catchat.DelegateSendSMS(msg)
		if err != nil {
			log.Fatalln(err)
		}
	})

	// Save all files within HacksRootPath to Bolt db.
	// Controller should respond with index of bolt buckets.
	//
	// r.GET("/hack/boltify", func(c *gin.Context) {

	// 	werr := filepath.Walk(config.FSStorePath, func(path string, info os.FileInfo, err error) error {
	// 		if err != nil {
	// 			fmt.Printf("Error: %v\n", err)
	// 			return nil
	// 		}

	// 		// Only for files.
	// 		if !info.IsDir() {

	// 			// Handle paths.
	// 			cleanPath := filepath.Clean(path) // hacks/snippets/todo/MOAR
	// 			dir := filepath.Dir(cleanPath)    // hacks/snippets/todo
	// 			withinHacksRootDir := strings.Replace(dir, config.FSStorePath+"/", "", 1)
	// 			folders := strings.Split(withinHacksRootDir, "/")
	// 			bucket := folders[0]
	// 			withinBucketDir := strings.Replace(withinHacksRootDir, bucket, "", 1)
	// 			name := withinBucketDir + "/" + info.Name()

	// 			// Get file contents and parse path (if not dir).
	// 			contents, ioerr := ioutil.ReadFile(path)

	// 			if ioerr != nil {
	// 				fmt.Printf("Error reading file: %v\n", ioerr)
	// 			} else {
	// 				fmt.Printf("cleanPath: %v\n", cleanPath)
	// 				fmt.Printf("dir: %v\n", dir)
	// 				fmt.Printf("withinHacksRootDir: %v\n", withinHacksRootDir)
	// 				fmt.Printf("bucket: %v\n", bucket)
	// 				fmt.Printf("withinBucketDir: %v\n", withinBucketDir)
	// 				fmt.Printf("name: %v\n", name)
	// 				fmt.Printf("Contents: \n---\n%v\n---\n", string(contents))

	// 				// Get snippet if exists by bucket name and filename.
	// 				// FIXME: ew.
	// 				// var snip models.Snippet

	// 				var snip models.Snippet

	// 				db.View(func(tx *bolt.Tx) error {
	// 					snip = models.GetSnippetByName(bucket, name, tx)
	// 					return nil
	// 				})

	// 				if snip == (models.Snippet{}) {

	// 					// Snippify.
	// 					snip.Name = name
	// 					snip.BucketName = bucket
	// 					newId := lib.RandSeq(6)
	// 					snip.Id = newId

	// 				}

	// 				// Make updates.
	// 				snip.Content = string(contents)
	// 				snip.Language = lib.GetLanguageModeByExtension(name)
	// 				snip.TimeStamp = int(time.Now().UTC().Unix() * 1000)

	// 				snipJSONBytes, _ := json.Marshal(snip)

	// 				// Save snippet to given bucket.
	// 				dberr := db.Update(func(tx *bolt.Tx) error {

	// 					return models.SetSnippet(snip.Id, snipJSONBytes, bucket, tx)
	// 				})
	// 				if dberr != nil {
	// 					fmt.Printf("Error saving file snippet to bolt: %v\n", dberr)
	// 				}
	// 			}
	// 		}
	// 		return nil
	// 	})

	// 	if werr != nil {
	// 		fmt.Println("Impossible.")
	// 	}

	// 	var buckets models.SnippetBuckets

	// 	indexerr := db.View(func(tx *bolt.Tx) error {
	// 		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
	// 			buckets = append(buckets, models.SnippetBucket{Name: string(name)})
	// 			return nil
	// 		})
	// 		return nil
	// 	})
	// 	if indexerr != nil {
	// 		c.JSON(500, indexerr)
	// 	} else {
	// 		c.JSON(200, buckets)
	// 	}

	// })

	r.Run(":5000")
}
