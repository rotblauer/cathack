package controllers

import (
	"fmt"
	"log"
	"net/http"

	"../models"
	"github.com/gin-gonic/gin"
)

// GetChat simply renders HTML page.
func GetChat(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "chat.html")
	log.Printf("Getting chat.")
	fmt.Println()
}

// GetChatData is called by AJAX on /chat HTML inline JS page load.
func GetChatData(c *gin.Context) {
	msgs, err := models.AllChatMsgs() // msgs are type []ChatMessageForm
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, msgs)
	}
}

// GetHack simply renders associated /hack HTML page.
func GetHack(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "hack.html")
}
