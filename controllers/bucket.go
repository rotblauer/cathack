package controllers

import (
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
