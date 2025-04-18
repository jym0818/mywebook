package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jym/mywebook/internal/events/article"
)

type App struct {
	consumers []article.Consumer
	web       *gin.Engine
}
