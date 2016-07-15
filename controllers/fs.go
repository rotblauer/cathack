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

// TODO: Move to model.
func (ctrl FSController) WriteBucketToDir(c *gin.Context) {
	bucketId := c.Param("bucketId")
	bucket := bucketModel.One([]byte(bucketId))

	var err error

	// DANGERZONE. Delete dir we're about to write.
	err = fsModel.DeleteDir(filepath.Join(config.FSStorePath, bucket.Meta.Name))
	if err != nil {
		fmt.Printf("Error cleaning bucket path: %v", err)
	}

	snippets, _ := snippetModel.All(bucketId)

	for _, snippet := range snippets {
		err = fsModel.WriteFile(bucket, snippet) // per bucketName, snippetName, content
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

// Gets array of filepaths (and info!) within HacksRootsDir.
func (ctrl FSController) Walk(c *gin.Context) {
	filepaths, err := fsModel.CollectDirPaths()
	if err != nil {
		c.JSON(500, "error walking fspath")
	} else {
		c.JSON(200, filepaths)
	}
}

func (ctrl FSController) SnippetizeOne(c *gin.Context) {
	var json string
	c.Bind(&json)
	b, s, e := fsModel.SnippetizeFile(json)
	if e != nil {
		c.JSON(500, e)
	} else {
		c.JSON(200, gin.H{"bucket": b, "snippet": s})
	}
}

func (ctrl FSController) SnippetizeMany(c *gin.Context) {
	var json string
	c.Bind(&json)
	bs, ss, e := fsModel.SnippetizeDir(json)
	if e != nil {
		c.JSON(500, e)
	} else {
		c.JSON(200, gin.H{"buckets": bs, "snippets": ss})
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
