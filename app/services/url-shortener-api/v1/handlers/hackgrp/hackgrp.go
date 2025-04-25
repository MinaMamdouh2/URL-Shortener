package hackgrp

import (
	"github.com/gin-gonic/gin"
)

// Hack handles the /hack route using Gin context.
func Hack(c *gin.Context) {
	c.JSON(200, gin.H{
		"Status": "Ok",
	})
}
