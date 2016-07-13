package controllers

import (
	"./models"
	"github.com/gin-gonic/gin"
)

type BucketController struct{}

var bucketModel = new(models.BucketModel)

func (m bucketModel) All(c *gin.Context) {
	buckets, err := bucketModel.All()
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, buckets)
	}
}
