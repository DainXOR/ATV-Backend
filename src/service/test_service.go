package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type testType struct{}

var Test testType

func (testType) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get test",
	})
}

func (testType) Post(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "post test",
	})
}

func (testType) Put(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "put test",
	})
}

func (testType) Delete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "delete test",
	})
}

func (testType) Patch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "patch test",
	})
}
