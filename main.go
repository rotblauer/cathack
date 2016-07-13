package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"./config"
	"./controllers"
	"./lib"
	"./models"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

//CORSMiddleware ...
// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://chat.areteh.co:5000")
// 	}
// }

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Open Bolt DB.
	defer models.GetDB().Close()

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

	// Start Snippet.
	snippet := new(controllers.SnippetController)

	r.GET("/hack/b/:bucketId", snippet.All)

	// /hack/s/:snippetId?bucket=snippets
	r.DELETE("/hack/s/:snippetId", snippet.Delete)

	// Start Bucket.
	bucket := new(controllers.BucketController)
	r.GET("/hack/b", bucket.All)

	// Start FS.
	fs := new(controllers.FSController)
	r.PUT("/hack/fs/b/:bucketId", fs.SetBucket)

	// Save all files within HacksRootPath to Bolt db.
	// Controller should respond with index of bolt buckets.
	//
	r.GET("/hack/boltify", func(c *gin.Context) {

		werr := filepath.Walk(config.FSStorePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return nil
			}

			// Only for files.
			if !info.IsDir() {

				// Handle paths.
				cleanPath := filepath.Clean(path) // hacks/snippets/todo/MOAR
				dir := filepath.Dir(cleanPath)    // hacks/snippets/todo
				withinHacksRootDir := strings.Replace(dir, config.FSStorePath+"/", "", 1)
				folders := strings.Split(withinHacksRootDir, "/")
				bucket := folders[0]
				withinBucketDir := strings.Replace(withinHacksRootDir, bucket, "", 1)
				name := withinBucketDir + "/" + info.Name()

				// Get file contents and parse path (if not dir).
				contents, ioerr := ioutil.ReadFile(path)

				if ioerr != nil {
					fmt.Printf("Error reading file: %v\n", ioerr)
				} else {
					fmt.Printf("cleanPath: %v\n", cleanPath)
					fmt.Printf("dir: %v\n", dir)
					fmt.Printf("withinHacksRootDir: %v\n", withinHacksRootDir)
					fmt.Printf("bucket: %v\n", bucket)
					fmt.Printf("withinBucketDir: %v\n", withinBucketDir)
					fmt.Printf("name: %v\n", name)
					fmt.Printf("Contents: \n---\n%v\n---\n", string(contents))

					// Get snippet if exists by bucket name and filename.
					// FIXME: ew.
					// var snip models.Snippet

					var snip models.Snippet

					db.View(func(tx *bolt.Tx) error {
						snip = models.GetSnippetByName(bucket, name, tx)
						return nil
					})

					if snip == (models.Snippet{}) {

						// Snippify.
						snip.Name = name
						snip.BucketName = bucket
						newId := lib.RandSeq(6)
						snip.Id = newId

					}

					// Make updates.
					snip.Content = string(contents)
					snip.Language = lib.GetLanguageModeByExtension(name)
					snip.TimeStamp = int(time.Now().UTC().Unix() * 1000)

					snipJSONBytes, _ := json.Marshal(snip)

					// Save snippet to given bucket.
					dberr := db.Update(func(tx *bolt.Tx) error {

						return models.SetSnippet(snip.Id, snipJSONBytes, bucket, tx)
					})
					if dberr != nil {
						fmt.Printf("Error saving file snippet to bolt: %v\n", dberr)
					}
				}
			}
			return nil
		})

		if werr != nil {
			fmt.Println("Impossible.")
		}

		var buckets models.SnippetBuckets

		indexerr := db.View(func(tx *bolt.Tx) error {
			tx.ForEach(func(name []byte, b *bolt.Bucket) error {
				buckets = append(buckets, models.SnippetBucket{Name: string(name)})
				return nil
			})
			return nil
		})
		if indexerr != nil {
			c.JSON(500, indexerr)
		} else {
			c.JSON(200, buckets)
		}

	})

	r.Run(":5000")
}
