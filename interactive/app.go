package main

import (
	"github.com/jym/mywebook/pkg/grpcx"
	"github.com/jym/mywebook/pkg/saramax"
)

type App struct {
	//在这里 所有需要main函数启动和关闭都在这里
	server    *grpcx.Server
	consumers []saramax.Consumer
}
