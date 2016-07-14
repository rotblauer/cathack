package controllers

import "github.com/gin-gonic/gin"

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
