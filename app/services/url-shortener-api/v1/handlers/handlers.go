package handlers

import (
	"github.com/MinaMamdouh2/URL-Shortener/app/services/url-shortener-api/v1/handlers/hackgrp"
	v1 "github.com/MinaMamdouh2/URL-Shortener/business/web/v1"
	"github.com/gin-gonic/gin"
)

type Routes struct{}

// Add implements the RouterAdder interface
func (Routes) Add(router *gin.Engine, cfg v1.APIMuxConfig) {
	hackgrp.Routes(router)
}
