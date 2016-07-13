package controllers

import "github.com/gin-gonic/gin"

type SnippetController struct{}

func (ctrl SnippetController) Delete(c *gin.Context) {
	snippetId := c.Param("snippetid") // func (c *Context) Param(key string) string
	bucketId := c.Query("bucketid")
	if len(bucketId) == 0 {
		c.JSON(400, "BucketId not present.")
	}
	err := snippetModel.Delete(bucketId, snippetId)
	if err != nil {
		c.JSON(400, "No snippet found with snippetId: "+snippetId)
	} else {
		snippets, _ := snippetModel.All(bucketId)
		// j, _ := json.Marshal(snippets)
		// h.Broadcast(o)
		if err == nil {
			c.JSON(200, gin.H{
				"snippetId": snippetId,
				"bucketId":  bucketId,
				"snippets":  snippets,
			})
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
