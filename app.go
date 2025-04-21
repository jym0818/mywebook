package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jym/mywebook/internal/events/article"
	"github.com/robfig/cron/v3"
)

type App struct {
	consumers []article.Consumer
	web       *gin.Engine
	cron      *cron.Cron
}
