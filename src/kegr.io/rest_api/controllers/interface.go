package controllers

import "github.com/gin-gonic/gin"

// IController is a controller's interface
type IController interface {
	Register(router *gin.Engine)
}
