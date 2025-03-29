package handlers

import "github.com/gin-gonic/gin"

type HandlerFunc func(*gin.Context, *Context)
