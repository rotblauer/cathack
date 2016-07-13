package controllers

import (
	"encoding/json"

	"./models"

	"github.com/gin-gonic/gin"
)

type SnippetController struct{}

var snippetModel = new(models.SnippetModel)

func (ctrl SnippetController) Delete(c *gin.Context) {
	snippetID := c.Param("snippetid") // func (c *Context) Param(key string) string
	bucketID := c.Param("bucketid")
	if len(bucketid) == 0 {
		c.JSON(400, "BucketID not present.")
	}
	err := snippetModel.Delete(bucketID, snippetID)
	if err != nil {
		c.JSON(400, "No snippet found with snippetID: "+snippetID)
	} else {
		snippets, _ := snippetModel.All(bucketID)
		j, _ := json.Marshal(snippets)
		h.Broadcast(o)
		if err == nil {
			c.JSON(200, "Deleted snippet: "+snippetID)
		} else {
			c.JSON(500, "Internal server error.")
		}
	}
}

func (ctrl SnippetController) All(c *gin.Context) {
	bucketId := c.Param("bucketid")
	snippets, err := snippetModel.All(bucketId)
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, snippets)
	}
}
