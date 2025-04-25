package hackgrp

import (
	"github.com/gin-gonic/gin"
)

// Routes adds specific routes for this group using Gin.
// Each group of handlers can define routes related to its domain
// and bind them to the router that is passed in.
func Routes(router *gin.Engine) {
	router.GET("/hack", Hack)
}
