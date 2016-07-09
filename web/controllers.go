package web

import (
	"../chatty"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetChat(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "index.html")
	log.Printf("Getting chat.")
	fmt.Println()
}

func GetChatData(c *gin.Context) {
	// func ReadFile(filename string) ([]byte, error)
	fileContents, err := ioutil.ReadFile("./chat.txt")
	if err != nil {
		fmt.Printf("Error ioutiling chat.txt: %v", err)
	}

	messageStrings := strings.Split(string(fileContents), "\n")

	var collection []chatty.ChatMessageAs

	for _, messageString := range messageStrings {
		bytes := []byte(messageString)
		var cm chatty.ChatMessageAs
		json.Unmarshal(bytes, &cm)
		collection = append(collection, cm)
	}

	collectionBytes, err := json.Marshal(collection) // []byte

	c.JSON(200, gin.H{
		"status": "200 OK",
		"data":   string(collectionBytes), // again, the hanging commas are strangely necessary
	})
}

func GetHack(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "hack.html")
}
