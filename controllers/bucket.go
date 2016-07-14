package controllers

import (
	"../models"
	"github.com/gin-gonic/gin"
)

type BucketController struct{}

func (m BucketController) All(c *gin.Context) {
	buckets, err := bucketModel.All()
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, buckets)
	}
}

func (m BucketController) Create(c *gin.Context) {
	bName := c.Param("bucketName")
	if len(bName) == 0 {
		c.JSON(400, "Must include a bucket name as parameter.")
	}

	bucket, err := bucketModel.Create(bName)
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, bucket)
	}
}

func (m BucketController) Destroy(c *gin.Context) {
	bId := c.Param("bucketId")
	if len(bId) == 0 {
		c.JSON(400, "Must include a bucketId as parameter.")
	}

	err := bucketModel.Destroy(bId)
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, gin.H{"status": "success"})
	}
}

// http://phalt.co/a-simple-api-in-go/
func (m BucketController) Set(c *gin.Context) {
	var json models.Bucket
	c.Bind(&json)
	// fmt.Printf("json name is: %v, name: %v\n", json.Id, json.Meta.Name)
	e := bucketModel.Set(json)
	if e != nil {
		c.JSON(500, e)
	} else {
		c.JSON(201, json)
	}
}
