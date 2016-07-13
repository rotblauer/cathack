package controllers

import (
	"fmt"
	"path/filepath"

	"../config"
	"../models"

	"github.com/gin-gonic/gin"
)

type FSController struct{}

var fsModel = new(models.FSModel)

func (ctrl FSController) SetBucket(c *gin.Context) {
	bucketId := c.Param("bucketId")
	bucket := bucketModel.One([]byte(bucketId))

	var err error
	// var snippets = *snippetModel.Snippet

	err = fsModel.DeleteDir(filepath.Join(config.FSStorePath, bucket.Meta.Name))
	if err != nil {
		fmt.Printf("Error cleaning bucket path: %v", err)
	}

	snippets, _ := snippetModel.All(bucketId)
	// if err != nil {
	// 	fmt.Printf("Error gettings snippets for bucketId: %v. Error:", bucketId, err)
	// }

	for _, snippet := range snippets {
		err = fsModel.SetFile(snippet)
		if err != nil {
			break
		}
	}

	if err != nil {
		c.JSON(500, "Internal server error: "+err.Error())
	} else {
		c.JSON(200, bucket)
	}
}

// func (ctrl FSController) SetSnippet(c *gin.Context) {

// }

// // Does not overwrite any existing buckets.
// func (ctrl FSController) GetGently(c *gin.Context) {

// }

// // Overwrite all buckets. ie --force
// func (ctrl) FSController) GetForce(c *gin.Context) {

// }
