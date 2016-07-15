package controllers

import (
	"encoding/json"
	"fmt"

	"../models"

	"github.com/gin-gonic/gin"
)

type FSController struct{}

var fsModel = new(models.FSModel)

// TODO: Move to model.
// func (ctrl FSController) WriteBucketToDir(c *gin.Context) {
// 	bucketId := c.Param("bucketId")
// 	bucket := bucketModel.One([]byte(bucketId))

// 	var err error

// 	// DANGERZONE. Delete dir we're about to write.
// 	err = fsModel.DeleteDir(filepath.Join(config.FSStorePath, bucket.Meta.Name))
// 	if err != nil {
// 		fmt.Printf("Error cleaning bucket path: %v", err)
// 	}

// 	snippets, _ := snippetModel.All(bucketId)

// 	for _, snippet := range snippets {
// 		err = fsModel.WriteFile(bucket, snippet) // per bucketName, snippetName, content
// 		if err != nil {
// 			break
// 		}
// 	}

// 	if err != nil {
// 		c.JSON(500, "Internal server error: "+err.Error())
// 	} else {
// 		c.JSON(200, bucket)
// 	}
// }

func (ctrl FSController) WriteFile(c *gin.Context) {
	snippetId := c.Param("snippetId")
	bucketId := c.Query("bucketId")
	if len(snippetId) == 0 || len(bucketId) == 0 {
		c.JSON(400, "Your query must be of the format: /hack/fs/s/<snippetId>?bucketId=<bucketId>")
	}

	err := fsModel.WriteFile(bucketId, snippetId)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "data": "error writing file"})
	} else {
		c.JSON(200, gin.H{"status": "success"})
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

	path := c.Query("path")
	fmt.Printf("got FS query path: %v\n", path)

	var p string
	json.Unmarshal([]byte(path), &p)

	b, s, e := fsModel.SnippetizeFile(p)
	if e != nil {
		c.JSON(500, e)
	} else {
		c.JSON(200, gin.H{"b": b, "s": s})
	}
}

func (ctrl FSController) SnippetizeMany(c *gin.Context) {
	path := c.Query("path")
	fmt.Printf("got FS query path: %v\n", path)

	var p string
	json.Unmarshal([]byte(path), &p)

	bs, ss, e := fsModel.SnippetizeDir(path)
	if e != nil {
		c.JSON(500, e)
	} else {
		c.JSON(200, gin.H{"b": bs, "s": ss})
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
